package orders

import (
	"strconv"

	"shipping-app/internal/app/application/orders"
	"shipping-app/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type CreateOrderRequest struct {
	Observation *string `json:"observation"`
	DriverID    *uint   `json:"driver_id"`
	VehicleID   *uint   `json:"vehicle_id"`
	PackageIDs  []uint  `json:"package_ids"`
	TypeService string  `json:"type_service"`
}

type AssignOrderRequest struct {
	DriverID  uint `json:"driver_id"`
	VehicleID uint `json:"vehicle_id"`
}

type UpdateStatusRequest struct {
	Status      string  `json:"status"`
	Observation *string `json:"observation"`
}

type OrderHandler struct {
	createUC         *orders.CreateOrderUseCase
	listUC           *orders.ListOrdersUseCase
	getUC            *orders.GetOrderUseCase
	assignUC         *orders.AssignOrderUseCase
	updateStatusUC   *orders.UpdateOrderStatusUseCase
	deleteUC         *orders.DeleteOrderUseCase
	listByDriverUC   *orders.ListOrdersByDriverUseCase
	listUnassignedUC *orders.ListOrdersUnassignedUseCase
}

func NewOrderHandler(
	createUC *orders.CreateOrderUseCase,
	listUC *orders.ListOrdersUseCase,
	getUC *orders.GetOrderUseCase,
	assignUC *orders.AssignOrderUseCase,
	updateStatusUC *orders.UpdateOrderStatusUseCase,
	deleteUC *orders.DeleteOrderUseCase,
	listByDriverUC *orders.ListOrdersByDriverUseCase,
	listUnassignedUC *orders.ListOrdersUnassignedUseCase,
) *OrderHandler {
	return &OrderHandler{
		createUC:         createUC,
		listUC:           listUC,
		getUC:            getUC,
		assignUC:         assignUC,
		updateStatusUC:   updateStatusUC,
		deleteUC:         deleteUC,
		listByDriverUC:   listByDriverUC,
		listUnassignedUC: listUnassignedUC,
	}
}

func (h *OrderHandler) CreateOrder(ctx fiber.Ctx) error {
	var req CreateOrderRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Failed to parse request body",
		})
	}

	input := orders.CreateOrderInput{
		Observation: req.Observation,
		DriverID:    req.DriverID,
		VehicleID:   req.VehicleID,
		PackageIDs:  req.PackageIDs,
		TypeService: req.TypeService,
	}

	output, err := h.createUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleErrorCreate(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "order created successfully",
		"id":        output.ID,
		"status":    output.Status,
		"create_at": output.CreateAt,
	})
}

func (h *OrderHandler) ListOrders(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)
	typeService := ctx.Query("type_service")
	status := ctx.Query("status")
	orderIDStr := ctx.Query("order_id")
	var orderID uint
	if orderIDStr != "" {
		parsedID, err := strconv.ParseUint(orderIDStr, 10, 32)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_order_id",
				"message": "Order ID must be a valid number",
			})
		}
		orderID = uint(parsedID)
	}

	input := orders.ListOrdersInput{
		Limit:       params.Limit,
		Offset:      params.Offset,
		TypeService: typeService,
		Status:      status,
		OrderID:     orderID,
	}

	ordersList, total, err := h.listUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	if ordersList == nil {
		ordersList = []*orders.ListOrdersByDriverOutput{}
	}

	response := utils.NewPaginationResponse(ordersList, int(total), params.Page, params.Limit)
	return ctx.JSON(response)
}

func (h *OrderHandler) GetOrder(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_id",
			"message": "ID must be a valid number",
		})
	}

	order, err := h.getUC.Execute(ctx.Context(), uint(id))
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(order)
}

func (h *OrderHandler) AssignOrder(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_id",
			"message": "ID must be a valid number",
		})
	}

	var req AssignOrderRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Failed to parse request body",
		})
	}

	input := orders.AssignOrderInput{
		OrderID:   uint(id),
		DriverID:  req.DriverID,
		VehicleID: req.VehicleID,
	}

	if err := h.assignUC.Execute(ctx.Context(), input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"message": "order assigned successfully",
	})
}

func (h *OrderHandler) UpdateStatus(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_id",
			"message": "ID must be a valid number",
		})
	}

	var req UpdateStatusRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Failed to parse request body",
		})
	}

	input := orders.UpdateOrderStatusInput{
		OrderID:     uint(id),
		Status:      req.Status,
		Observation: req.Observation,
	}

	if err := h.updateStatusUC.Execute(ctx.Context(), input); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"message": "order status updated successfully",
	})
}

func (h *OrderHandler) DeleteOrder(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_id",
			"message": "ID must be a valid number",
		})
	}

	if err := h.deleteUC.Execute(ctx.Context(), uint(id)); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (h *OrderHandler) ListOrdersByDriver(ctx fiber.Ctx) error {
	driverIDStr := ctx.Params("driverId")
	driverID, err := strconv.ParseUint(driverIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_driver_id",
			"message": "Driver ID must be a valid number",
		})
	}

	params := utils.GetPaginationParams(ctx)

	input := orders.ListOrdersByDriverInput{
		DriverID: uint(driverID),
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	ordersList, total, err := h.listByDriverUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	response := utils.NewPaginationResponse(ordersList, int(total), params.Page, params.Limit)
	return ctx.JSON(response)
}

func (h *OrderHandler) ListOrdersUnassigned(ctx fiber.Ctx) error {
	params := utils.GetPaginationParams(ctx)
	idStr := ctx.Query("id")
	id := uint64(0)
	err := error(nil)
	if idStr != "" {
		id, err = strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_id",
				"message": "ID must be a valid number",
			})
		}

	}
	input := orders.ListOrdersUnassignedInput{
		Limit:  params.Limit,
		Offset: params.Offset,
		ID:     uint(id),
	}

	ordersList, total, err := h.listUnassignedUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	response := utils.NewPaginationResponse(ordersList, int(total), params.Page, params.Limit)
	return ctx.JSON(response)
}
