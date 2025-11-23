package delivery

import (
	"shipping-app/internal/app/application/delivery"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ReportDeliveryRequest struct {
	PackageID         uint    `json:"package_id"`
	Observation       *string `json:"observation"`
	SignatureReceived *string `json:"signature_received"`
	PhotoDelivery     string  `json:"photo_delivery"`
}

type ReportIncidentRequest struct {
	PackageID          uint    `json:"package_id"`
	ReasonCancellation string  `json:"reason_cancellation"`
	Observation        *string `json:"observation"`
	PhotoEvidence      string  `json:"photo_evidence"`
	Status             string  `json:"status"`
}

type DeliveryHandler struct {
	reportDeliveryUC   *delivery.ReportDeliveryUseCase
	reportIncidentUC   *delivery.ReportIncidentUseCase
	getPackageReportUC *delivery.GetPackageReportUseCase
}

func NewDeliveryHandler(
	reportDeliveryUC *delivery.ReportDeliveryUseCase,
	reportIncidentUC *delivery.ReportIncidentUseCase,
	getPackageReportUC *delivery.GetPackageReportUseCase,
) *DeliveryHandler {
	return &DeliveryHandler{
		reportDeliveryUC:   reportDeliveryUC,
		reportIncidentUC:   reportIncidentUC,
		getPackageReportUC: getPackageReportUC,
	}
}

func (h *DeliveryHandler) ReportDelivery(ctx fiber.Ctx) error {
	var req ReportDeliveryRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	input := delivery.ReportDeliveryInput{
		PackageID:         req.PackageID,
		Observation:       req.Observation,
		SignatureReceived: req.SignatureReceived,
		PhotoDelivery:     req.PhotoDelivery,
	}

	output, err := h.reportDeliveryUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "delivery reported successfully",
		"data":    output,
	})
}

func (h *DeliveryHandler) ReportIncident(ctx fiber.Ctx) error {
	var req ReportIncidentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	input := delivery.ReportIncidentInput{
		PackageID:          req.PackageID,
		ReasonCancellation: req.ReasonCancellation,
		Observation:        req.Observation,
		PhotoEvidence:      req.PhotoEvidence,
		Status:             req.Status,
	}

	output, err := h.reportIncidentUC.Execute(ctx.Context(), input)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "incident reported successfully",
		"data":    output,
	})
}

func (h *DeliveryHandler) GetPackageReportByPackageID(ctx fiber.Ctx) error {
	packageIDStr := ctx.Params("packageId")
	packageID, err := strconv.ParseUint(packageIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_package_id",
			"message": "Package ID must be a valid number",
		})
	}

	report, err := h.getPackageReportUC.ExecuteByPackageID(ctx.Context(), uint(packageID))
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data": report,
	})
}

func (h *DeliveryHandler) GetPackageReportByReportID(ctx fiber.Ctx) error {
	reportIDStr := ctx.Params("reportId")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_report_id",
			"message": "Report ID must be a valid number",
		})
	}

	report, err := h.getPackageReportUC.ExecuteByReportID(ctx.Context(), uint(reportID))
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data": report,
	})
}

func (h *DeliveryHandler) handleError(ctx fiber.Ctx, err error) error {
	switch err {
	case delivery.ErrInvalidDeliveryInput, delivery.ErrInvalidIncidentInput:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_input",
			"message": err.Error(),
		})
	case delivery.ErrPackageAlreadyDelivered, delivery.ErrPackageHasIncident:
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "conflict",
			"message": err.Error(),
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "An error occurred while processing the request",
		})
	}
}
