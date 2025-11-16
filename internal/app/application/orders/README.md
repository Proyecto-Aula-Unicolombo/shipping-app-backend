# Módulo de Órdenes - Documentación

## Descripción General

El módulo de órdenes gestiona el ciclo completo de vida de las órdenes de envío en el sistema de logística. Permite crear, asignar, actualizar, consultar y eliminar órdenes, así como gestionar la relación entre órdenes, paquetes, conductores y vehículos.

## Arquitectura

El módulo sigue la arquitectura hexagonal (Clean Architecture) del proyecto:

```
domain/
├── entities/order.go           # Entidad Order
└── ports/repository/
    └── orderRepository.go      # Interface OrderRepository

application/orders/
├── createOrderUseCase.go       # Crear orden
├── listOrdersUseCase.go        # Listar todas las órdenes
├── getOrderUseCase.go          # Consultar orden por ID
├── assignOrderUseCase.go       # Asignar conductor y vehículo
├── updateOrderStatusUseCase.go # Actualizar estado
├── deleteOrderUseCase.go       # Eliminar orden
└── listOrdersByDriverUseCase.go # Listar órdenes por conductor

infrastructure/
├── adapters/
│   └── orderRepositoryAdapter.go  # Implementación PostgreSQL
└── fiber/
    ├── handlers/orders/
    │   ├── orderHandler.go        # Handlers HTTP
    │   └── handleError.go         # Manejo de errores
    └── routers/
        └── setOrderRouter.go      # Configuración de rutas
```

## Casos de Uso Implementados

### 1. Crear Orden (CREAR ORDEN)
**Endpoint**: `POST /api/v1/orders`

**Input**:
```json
{
  "observation": "Entregar antes de las 5pm",
  "driver_id": 1,
  "vehicle_id": 201,
  "package_ids": [1, 2, 3]
}
```

**Validaciones**:
- Driver ID y Vehicle ID son requeridos
- Al menos un paquete debe ser proporcionado
- El conductor debe existir en la base de datos
- El vehículo debe existir en la base de datos
- Los paquetes deben estar disponibles (sin orden asignada)

**Respuesta exitosa** (201):
```json
{
  "message": "order created successfully",
  "id": 1,
  "status": "Pendiente",
  "create_at": "2024-03-15T10:30:00Z"
}
```

**Errores posibles**:
- `400 Bad Request`: Input inválido o sin paquetes
- `404 Not Found`: Conductor o vehículo no encontrado
- `409 Conflict`: Paquetes no disponibles
- `500 Internal Server Error`: Error del servidor

### 2. Listar Todas las Órdenes (LISTAR TODAS LAS ORDENES)
**Endpoint**: `GET /api/v1/orders?page=1&limit=10`

**Parámetros de query**:
- `page` (opcional): Número de página (default: 1)
- `limit` (opcional): Elementos por página (default: 10)

**Respuesta exitosa** (200):
```json
{
  "data": [
    {
      "id": 1,
      "create_at": "2024-03-15T10:30:00Z",
      "assigned_at": "2024-03-15T10:35:00Z",
      "observation": "Entregar antes de las 5pm",
      "status": "En camino",
      "driver_id": 1,
      "vehicle_id": 201
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 45
  }
}
```

### 3. Consultar Orden por ID (CONSULTAR ORDEN)
**Endpoint**: `GET /api/v1/orders/:id`

**Respuesta exitosa** (200):
```json
{
  "id": 1,
  "create_at": "2024-03-15T10:30:00Z",
  "assigned_at": "2024-03-15T10:35:00Z",
  "observation": "Entregar antes de las 5pm",
  "status": "En camino",
  "driver_id": 1,
  "vehicle_id": 201
}
```

**Errores posibles**:
- `400 Bad Request`: ID inválido
- `404 Not Found`: Orden no encontrada

### 4. Asignar Conductor y Vehículo (ASIGNAR ORDEN)
**Endpoint**: `PUT /api/v1/orders/:id/assign`

**Input**:
```json
{
  "driver_id": 2,
  "vehicle_id": 202
}
```

**Comportamiento**:
- Solo órdenes en estado "Pendiente" pueden ser asignadas
- Cambia automáticamente el estado a "En camino"
- Actualiza `assigned_at` con timestamp actual
- Valida que el conductor y vehículo existan

**Respuesta exitosa** (200):
```json
{
  "message": "order assigned successfully"
}
```

**Errores posibles**:
- `404 Not Found`: Orden, conductor o vehículo no encontrado
- `409 Conflict`: Orden ya asignada

### 5. Actualizar Estado (ACTUALIZAR ORDEN)
**Endpoint**: `PUT /api/v1/orders/:id/status`

**Input**:
```json
{
  "status": "Entregado",
  "observation": "Entregado sin novedad"
}
```

**Estados válidos**:
- `Pendiente`: Orden creada, esperando asignación
- `En camino`: Conductor en ruta hacia destino
- `Entregado`: Orden completada exitosamente
- `Cancelado`: Orden cancelada

**Respuesta exitosa** (200):
```json
{
  "message": "order status updated successfully"
}
```

**Errores posibles**:
- `400 Bad Request`: Estado inválido
- `404 Not Found`: Orden no encontrada

### 6. Eliminar Orden (ELIMINAR ORDEN)
**Endpoint**: `DELETE /api/v1/orders/:id`

**Restricciones**:
- Solo órdenes en estado "Pendiente" pueden ser eliminadas
- Las órdenes en proceso o entregadas no pueden eliminarse

**Respuesta exitosa** (204 No Content)

**Errores posibles**:
- `404 Not Found`: Orden no encontrada
- `409 Conflict`: No se puede eliminar orden en ese estado

### 7. Listar Órdenes por Conductor (LISTAR ORDENES ASIGNADAS)
**Endpoint**: `GET /api/v1/orders/driver/:driverId?page=1&limit=10`

**Parámetros**:
- `driverId`: ID del conductor (en la URL)
- `page` (opcional): Número de página
- `limit` (opcional): Elementos por página

**Respuesta exitosa** (200):
```json
{
  "data": [
    {
      "id": 1,
      "create_at": "2024-03-15T10:30:00Z",
      "assigned_at": "2024-03-15T10:35:00Z",
      "status": "En camino",
      "driver_id": 1,
      "vehicle_id": 201
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5
  }
}
```

## Base de Datos

### Tabla `orders`
```sql
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_at TIMESTAMP,
    observation TEXT,
    status VARCHAR(50) NOT NULL,
    iddriver INTEGER NOT NULL,
    idvehicle INTEGER NOT NULL,
    FOREIGN KEY (iddriver) REFERENCES drivers(id),
    FOREIGN KEY (idvehicle) REFERENCES vehicles(id)
);
```

### Relaciones
- **orders → drivers**: Cada orden está asignada a un conductor
- **orders → vehicles**: Cada orden usa un vehículo específico
- **packages → orders**: Múltiples paquetes pueden pertenecer a una orden (FK: `idorder`)

## Flujos de Trabajo

### Flujo de Creación de Orden
```
1. Cliente/Coordinador crea orden
2. Sistema valida conductor y vehículo
3. Sistema valida disponibilidad de paquetes
4. Sistema crea orden en transacción
5. Sistema asigna paquetes a la orden
6. Sistema confirma transacción
7. Orden queda en estado "Pendiente"
```

### Flujo de Asignación
```
1. Coordinador asigna conductor y vehículo
2. Sistema valida que orden esté "Pendiente"
3. Sistema valida existencia de conductor y vehículo
4. Sistema actualiza orden con asignación
5. Sistema cambia estado a "En camino"
6. Conductor recibe notificación (futuro)
```

### Flujo de Entrega
```
1. Conductor actualiza estado a "Entregado"
2. Sistema registra observaciones de entrega
3. Sistema actualiza timestamps
4. Sistema notifica al cliente (futuro)
5. Orden finalizada
```

## Dependencias

### Repositorios requeridos
- `OrderRepository`: Operaciones sobre órdenes
- `DriverRepository`: Validación de conductores
- `VehicleRepository`: Validación de vehículos
- `PackageRepository`: Gestión de paquetes
- `TxProvider`: Manejo de transacciones

### Servicios externos (futuro)
- Servicio de notificaciones
- Servicio de geolocalización
- Servicio de tracking en tiempo real

## Códigos de Error

| Código | Descripción |
|--------|-------------|
| `invalid_input` | Datos de entrada inválidos |
| `driver_not_found` | Conductor no encontrado |
| `vehicle_not_found` | Vehículo no encontrado |
| `no_packages` | No se proporcionaron paquetes |
| `package_not_available` | Paquetes no disponibles |
| `order_not_found` | Orden no encontrada |
| `order_already_assigned` | Orden ya asignada |
| `invalid_status` | Estado inválido |
| `cannot_delete_order` | No se puede eliminar la orden |
| `internal_server_error` | Error interno del servidor |

## Ejemplos de Uso

### Ejemplo completo: Crear y asignar orden

```bash
# 1. Crear orden
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "observation": "Entregar en recepción",
    "driver_id": 1,
    "vehicle_id": 201,
    "package_ids": [1, 2, 3]
  }'

# Respuesta:
# {
#   "message": "order created successfully",
#   "id": 1,
#   "status": "Pendiente",
#   "create_at": "2024-03-15T10:30:00Z"
# }

# 2. Consultar orden creada
curl http://localhost:8080/api/v1/orders/1

# 3. Actualizar estado a "En camino"
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "En camino"
  }'

# 4. Listar órdenes del conductor
curl http://localhost:8080/api/v1/orders/driver/1

# 5. Actualizar a "Entregado"
curl -X PUT http://localhost:8080/api/v1/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Entregado",
    "observation": "Entregado sin novedad"
  }'
```

## Testing

### Pruebas unitarias recomendadas
```bash
# Ejecutar tests del módulo
go test ./internal/app/application/orders/...
go test ./internal/app/infrastructure/adapters/...
```

### Escenarios de prueba
1. ✅ Crear orden con paquetes válidos
2. ✅ Crear orden sin conductor → Error
3. ✅ Crear orden sin vehículo → Error
4. ✅ Crear orden con paquetes no disponibles → Error
5. ✅ Asignar orden pendiente → Éxito
6. ✅ Asignar orden ya asignada → Error
7. ✅ Actualizar estado con estado válido → Éxito
8. ✅ Actualizar estado con estado inválido → Error
9. ✅ Eliminar orden pendiente → Éxito
10. ✅ Eliminar orden en camino → Error

## Mejoras Futuras

### Features planificados
- [ ] Notificaciones push para conductores
- [ ] Tracking en tiempo real con GPS
- [ ] Historial de cambios de estado
- [ ] Reprogramación de órdenes
- [ ] Optimización de rutas
- [ ] Asignación automática de conductores
- [ ] Métricas de desempeño por conductor
- [ ] Reportes de entregas

### Optimizaciones técnicas
- [ ] Cache de órdenes activas
- [ ] Índices en base de datos para búsquedas
- [ ] Paginación con cursor en lugar de offset
- [ ] Eventos de dominio para notificaciones
- [ ] Validación de horarios de entrega

## Soporte

Para preguntas o problemas:
1. Revisar logs del servidor en `/var/log/shipping-app/`
2. Verificar estado de la base de datos
3. Consultar documentación de casos de uso específicos
4. Revisar tests unitarios para ejemplos de uso

---

**Versión**: 1.0.0  
**Última actualización**: Noviembre 2025  
**Autor**: Equipo de Desarrollo Logística Express
