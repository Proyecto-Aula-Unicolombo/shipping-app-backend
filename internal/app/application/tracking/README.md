# Módulo de Tracking y Paradas

Este módulo maneja el rastreo en tiempo real de paquetes para destinatarios y el registro de paradas durante las entregas.

## Casos de Uso Implementados

### 1. Rastrear Paquete para Destinatario (TRACK PACKAGE)

#### Por Número de Paquete (Público)
**Endpoint**: `GET /api/v1/tracking/package?num_package=PKG-2025-001&receiver_id=1`

**Descripción**: Permite a los destinatarios rastrear sus paquetes usando el número de paquete. Opcionalmente valida que el destinatario tenga acceso al paquete.

**Query Parameters**:
- `num_package` (requerido): Número único del paquete
- `receiver_id` (opcional): ID del destinatario para validar acceso

**Respuesta exitosa** (200):
```json
{
  "data": {
    "package_id": 1,
    "num_package": "PKG-2025-001",
    "status": "En camino",
    "origin": "Calle 100 #45-67, Bogotá",
    "destination": "Carrera 7 #85-12, Bogotá",
    "estimated_delivery": "2025-11-16T15:00:00Z",
    "current_location": {
      "latitude": 4.6758,
      "longitude": -74.0498,
      "timestamp": "2025-11-16T14:30:00Z"
    },
    "receiver_name": "Juan Pérez García",
    "receiver_phone": "+57 300 123 4567",
    "is_fragile": false,
    "weight": 2.5,
    "tracking_history": [
      {
        "latitude": 4.6758,
        "longitude": -74.0498,
        "timestamp": "2025-11-16T14:30:00Z"
      },
      {
        "latitude": 4.6695,
        "longitude": -74.0532,
        "timestamp": "2025-11-16T14:00:00Z"
      }
    ]
  }
}
```

#### Por ID de Paquete (Interno)
**Endpoint**: `GET /api/v1/tracking/package/:packageId`

**Descripción**: Rastreo interno usando ID del paquete.

**Respuesta**: Igual que el anterior

**Errores posibles**:
- `400 Bad Request`: Parámetros inválidos
- `403 Forbidden`: El destinatario no tiene acceso a este paquete
- `404 Not Found`: Paquete no encontrado

---

### 2. Registrar Parada (REGISTER STOP)
**Endpoint**: `POST /api/v1/stops/register`

**Descripción**: Permite a los conductores registrar paradas durante la ruta de entrega, incluyendo paradas programadas, incidentes, recogidas y entregas.

**Request Body**:
```json
{
  "order_id": 1,
  "latitude": 4.6758,
  "longitude": -74.0498,
  "type_stop": "Parada",
  "description": "Parada para almuerzo",
  "evidence": "https://example.com/photo.jpg"
}
```

**Tipos de Parada Válidos**:
- `"Parada"`: Parada programada o descanso
- `"Incidente"`: Problema durante la ruta (tráfico, accidente, etc.)
- `"Recogida"`: Recogida de paquete en origen
- `"Entrega"`: Entrega de paquete en destino

**Respuesta exitosa** (201):
```json
{
  "message": "stop registered successfully",
  "data": {
    "id": 1,
    "order_id": 1,
    "type_stop": "Parada",
    "timestamp": "2025-11-16T14:30:00Z"
  }
}
```

**Errores posibles**:
- `400 Bad Request`: Input inválido, coordenadas fuera de rango, o tipo de parada inválido
- `404 Not Found`: Orden no encontrada o no está en progreso
- `409 Conflict`: La orden no está en estado válido para registrar paradas

---

## Reglas de Negocio

### Rastreo de Paquetes
1. El paquete debe existir en el sistema
2. Si se proporciona `receiver_id`, se valida que el destinatario tenga acceso
3. Solo se muestra tracking si el paquete tiene una orden asignada
4. La ubicación actual es la más reciente del historial
5. El historial se ordena de más reciente a más antiguo
6. Se incluye información del destinatario para contacto

### Registro de Paradas
1. La orden debe existir y estar en estado "En camino" o "Pendiente"
2. Las coordenadas deben ser válidas (lat: -90 a 90, lng: -180 a 180)
3. El tipo de parada debe ser uno de los 4 tipos válidos
4. La descripción y evidencia son opcionales
5. Cada parada registra automáticamente un track en el sistema
6. Las paradas se almacenan con timestamp automático

---

## Flujo de Trabajo

### Rastreo para Destinatario:
```
1. Destinatario recibe número de paquete (ej: PKG-2025-001)
2. Ingresa número en app o web
3. Sistema busca paquete y valida acceso
4. Sistema obtiene orden asignada al paquete
5. Sistema recupera historial de ubicaciones
6. Se muestra mapa con ruta y ubicación actual
7. Se muestra información de entrega estimada
8. Actualizaciones en tiempo real cada pocos segundos
```

### Registro de Parada por Conductor:
```
1. Conductor hace una parada durante la ruta
2. App captura ubicación GPS automáticamente
3. Conductor selecciona tipo de parada
4. Opcionalmente agrega descripción/foto
5. Envía registro con POST /stops/register
6. Sistema valida orden y estado
7. Sistema registra parada en deliverystops
8. Sistema actualiza tracking en tracks
9. Ubicación se sincroniza para rastreo en tiempo real
```

---

## Datos de la Base de Datos

### Tabla `tracks` (Tracking de Órdenes):
```sql
CREATE TABLE tracks (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    idorder INTEGER NOT NULL,
    FOREIGN KEY (idorder) REFERENCES orders(id)
);
```

### Tabla `deliverystops` (Paradas de Entrega):
```sql
CREATE TABLE deliverystops (
    id SERIAL PRIMARY KEY,
    stoplocation GEOGRAPHY(POINT, 4326) NOT NULL,
    typestop VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    evidence TEXT,
    idorder INTEGER NOT NULL,
    FOREIGN KEY (idorder) REFERENCES orders(id)
);
```

**Índices Geoespaciales**:
- `idx_tracks_location`: Índice GIST para búsquedas geográficas en tracks
- `idx_deliverystops_location`: Índice GIST para búsquedas geográficas en paradas

---

## Ejemplos de Uso

### Usando cURL

#### Rastrear Paquete por Número:
```bash
curl "http://localhost:8080/api/v1/tracking/package?num_package=PKG-2025-001"
```

#### Rastrear con Validación de Destinatario:
```bash
curl "http://localhost:8080/api/v1/tracking/package?num_package=PKG-2025-001&receiver_id=1"
```

#### Rastrear por ID Interno:
```bash
curl "http://localhost:8080/api/v1/tracking/package/1"
```

#### Registrar Parada:
```bash
curl -X POST http://localhost:8080/api/v1/stops/register \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "latitude": 4.6758,
    "longitude": -74.0498,
    "type_stop": "Parada",
    "description": "Descanso programado",
    "evidence": "https://s3.amazonaws.com/stops/photo789.jpg"
  }'
```

#### Registrar Incidente en Ruta:
```bash
curl -X POST http://localhost:8080/api/v1/stops/register \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "latitude": 4.6695,
    "longitude": -74.0532,
    "type_stop": "Incidente",
    "description": "Congestión vehicular en la vía",
    "evidence": null
  }'
```

---

## Integraciones

### Con Otros Módulos:
- **Paquetes**: Obtiene información completa del paquete para rastreo
- **Órdenes**: Valida estado de orden antes de registrar paradas
- **PostGIS**: Usa tipos geográficos para ubicaciones precisas
- **WebSockets**: Para actualizaciones en tiempo real (opcional)

### Tecnologías Geoespaciales:
- **PostGIS**: Extension de PostgreSQL para datos geográficos
- **SRID 4326**: Sistema de coordenadas WGS84 (estándar GPS)
- **ST_MakePoint**: Crea puntos geográficos desde coordenadas
- **ST_AsBinary**: Convierte geometrías a formato WKB para Go
- **go-geom**: Librería Go para manejo de geometrías

---

## Características Avanzadas

### Optimizaciones de Rendimiento:
- Índices geoespaciales GIST para búsquedas rápidas
- Caché de ubicaciones recientes
- Límite de historial para evitar cargas grandes
- Paginación de paradas por orden

### Seguridad y Privacidad:
- Validación de acceso por destinatario
- Solo se muestran datos relevantes al destinatario
- No se exponen IDs internos en tracking público
- Validación de coordenadas para evitar datos inválidos

### Extensiones Futuras:
- [ ] Notificaciones push cuando el conductor está cerca
- [ ] Predicción de tiempo de llegada con ML
- [ ] Ruta optimizada sugerida para conductor
- [ ] Heatmap de zonas con más incidentes
- [ ] Integración con Google Maps/Waze
- [ ] Geofencing para alertas automáticas
- [ ] Histórico de rutas completadas
- [ ] Analytics de eficiencia de conductores

---

## Casos de Uso del Rastreo

### Para Destinatarios:
1. Ver ubicación actual del paquete en mapa
2. Conocer tiempo estimado de llegada
3. Ver historial de la ruta recorrida
4. Contactar al destinatario si es necesario
5. Recibir notificaciones de proximidad

### Para Administradores:
1. Monitorear todas las entregas en tiempo real
2. Identificar retrasos o desvíos de ruta
3. Analizar eficiencia de rutas
4. Detectar patrones de incidentes
5. Generar reportes de desempeño

### Para Conductores:
1. Registrar paradas durante la ruta
2. Documentar incidentes con evidencia
3. Marcar recogidas y entregas
4. Mantener historial de movimientos
5. Justificar tiempos de entrega

---

**Versión**: 1.0.0  
**Última actualización**: Noviembre 2025  
**Autor**: Equipo de Desarrollo Logística Express
