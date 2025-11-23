# Sistema de Rastreo en Tiempo Real

## 📡 Arquitectura del Sistema

El sistema de rastreo GPS implementa una arquitectura basada en WebSockets para comunicación bidireccional en tiempo real entre conductores, administradores y clientes.

### Componentes Principales

1. **Backend (Go + Fiber + WebSocket)**
2. **Base de datos PostgreSQL con PostGIS** (almacenamiento de coordenadas GPS)
3. **Frontend Admin/Coordinador** (React - pendiente)
4. **Frontend Cliente** (React - pendiente)
5. **App Conductor** (PWA/React - pendiente)

---

## 🔌 API Endpoints Disponibles

### WebSocket
```
WS: /api/v1/ws
```
Conexión WebSocket para recibir actualizaciones en tiempo real.

### HTTP Endpoints

#### 1. Registrar Ubicación (Conductor)
```http
POST /api/v1/tracks
Authorization: Bearer <jwt_token>

{
  "order_id": 123,
  "latitude": 6.2442,
  "longitude": -75.5812
}
```

**Respuesta:**
```json
{
  "message": "Track created successfully"
}
```

**Comportamiento:**
- Guarda la ubicación en la base de datos
- Envía actualización en tiempo real vía WebSocket a:
  - Todos los administradores
  - Clientes siguiendo esa orden específica

---

#### 2. Obtener Historial de Ubicaciones de una Orden
```http
GET /api/v1/tracks/order/:orderId?limit=50
Authorization: Bearer <jwt_token>
```

**Respuesta:**
```json
{
  "data": {
    "order_id": 123,
    "status": "En camino",
    "tracks": [
      {
        "track_id": 1,
        "order_id": 123,
        "latitude": 6.2442,
        "longitude": -75.5812,
        "timestamp": "2024-01-20T15:30:00Z"
      }
    ]
  }
}
```

**Parámetros opcionales:**
- `limit`: Cantidad máxima de tracks a retornar (por defecto: todos)

---

#### 3. Obtener Ubicaciones de Todos los Conductores Activos
```http
GET /api/v1/tracks/active-drivers
Authorization: Bearer <jwt_token>
```

**Respuesta:**
```json
{
  "data": {
    "total_drivers": 5,
    "drivers": [
      {
        "driver_id": 1,
        "driver_name": "Carlos Ramírez López",
        "phone_number": "+57 300 123 4567",
        "order_id": 123,
        "order_status": "En camino",
        "latitude": 6.2442,
        "longitude": -75.5812,
        "last_update": "2024-01-20T15:30:00Z",
        "vehicle_id": 1,
        "vehicle_plate": "ABC123"
      }
    ]
  }
}
```

**Uso:**
- Para dashboard de administrador/coordinador
- Muestra solo conductores con órdenes activas (no entregadas ni canceladas)
- Muestra la última ubicación reportada de cada conductor

---

## 🔄 Flujo de Datos en Tiempo Real

### 1. Conductor Envía Ubicación

```
[App Conductor] 
    ↓ POST /api/v1/tracks (cada 15-30 seg)
[Backend] → Guarda en DB
    ↓
[Hub WebSocket] → BroadcastToOrder(orderID)
    ↓
[Admin Dashboard] ← Actualiza mapa
[Cliente Web]     ← Actualiza tracking
```

### 2. Admin/Coordinador Ve Todos los Conductores

```
[Admin Dashboard]
    ↓ GET /api/v1/tracks/active-drivers
[Backend] → Consulta última ubicación de cada conductor activo
    ↓
[Admin Dashboard] ← Muestra mapa con marcadores
```

### 3. Cliente Final Rastrea su Pedido

```
[Cliente Web]
    ↓ WS: /api/v1/ws (conecta con orderID)
    ↓ GET /api/v1/tracks/order/123
[Backend] → Retorna historial + stream en tiempo real
    ↓
[Cliente Web] ← Muestra mapa + ETA actualizado
```

---

## 📨 Mensajes WebSocket

### Estructura de Mensajes

```json
{
  "type": "track_update",
  "payload": {
    "order_id": 123,
    "track_id": 456,
    "latitude": 6.2442,
    "longitude": -75.5812,
    "timestamp": "2024-01-20T15:30:00Z"
  }
}
```

### Tipos de Mensajes

| Tipo | Descripción | Quién lo recibe |
|------|-------------|-----------------|
| `track_update` | Nueva ubicación de conductor | Admin + Cliente de esa orden |
| `order_status_change` | Cambio de estado de orden | Admin + Cliente de esa orden |
| `driver_assigned` | Conductor asignado a orden | Admin + Cliente |

---

## 🎭 Roles y Permisos WebSocket

### Conexión por Rol

Al conectarse al WebSocket, los clientes deben identificarse:

```javascript
// Ejemplo de conexión desde frontend
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

// Los clientes pueden enviar mensaje inicial para identificarse
ws.send(JSON.stringify({
  type: 'subscribe',
  role: 'client',      // 'admin', 'driver', 'client'
  user_id: 123,
  order_ids: [456, 789] // Solo para clientes
}));
```

### Broadcast Selectivo

El sistema usa estos métodos para enviar mensajes:

1. **`BroadcastToRole(message, role)`** - Solo a un rol específico
2. **`BroadcastToOrder(message, orderID)`** - A admins + clientes de esa orden
3. **`SendToClient(message, clientID)`** - A un cliente específico
4. **`BroadcastJSON(message, nil)`** - A todos (broadcast global)

---

## 🗃️ Estructura de Base de Datos

### Tabla: tracks

```sql
CREATE TABLE tracks (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    location GEOMETRY(Point, 4326) NOT NULL,  -- PostGIS
    idorder INTEGER NOT NULL REFERENCES orders(id),
    CONSTRAINT fk_order FOREIGN KEY (idorder) 
        REFERENCES orders(id) ON DELETE CASCADE
);

-- Índices para optimización
CREATE INDEX idx_tracks_order ON tracks(idorder);
CREATE INDEX idx_tracks_timestamp ON tracks(timestamp DESC);
CREATE INDEX idx_tracks_location ON tracks USING GIST(location);
```

**Notas:**
- `GEOMETRY(Point, 4326)` usa el sistema de coordenadas WGS84 (GPS estándar)
- PostGIS permite consultas espaciales eficientes
- Índice GIST para búsquedas geoespaciales rápidas

---

## 🚀 Implementación en Frontend (Pendiente)

### Para Admin/Coordinador

```typescript
// Hook propuesto para tracking en tiempo real
import { useWebSocket } from '@/hooks/useWebSocket';
import { useQuery } from '@tanstack/react-query';

function AdminTrackingDashboard() {
  const { data: activeDrivers } = useQuery({
    queryKey: ['active-drivers'],
    queryFn: () => fetch('/api/v1/tracks/active-drivers').then(r => r.json()),
    refetchInterval: 30000 // Actualizar cada 30 seg
  });

  const { lastMessage } = useWebSocket('/api/v1/ws', {
    onMessage: (message) => {
      if (message.type === 'track_update') {
        // Actualizar marcador en el mapa
        updateDriverMarker(message.payload);
      }
    }
  });

  return (
    <GoogleMap markers={activeDrivers?.data?.drivers} />
  );
}
```

### Para Cliente Final

```typescript
function ClientOrderTracking({ orderNumber }) {
  const { data: orderTracks } = useQuery({
    queryKey: ['order-tracks', orderNumber],
    queryFn: () => fetch(`/api/v1/tracks/order/${orderNumber}`).then(r => r.json())
  });

  const { lastMessage } = useWebSocket('/api/v1/ws', {
    onConnect: (ws) => {
      // Suscribirse a actualizaciones de esta orden
      ws.send(JSON.stringify({
        type: 'subscribe',
        role: 'client',
        order_ids: [orderNumber]
      }));
    },
    onMessage: (message) => {
      if (message.type === 'track_update' && 
          message.payload.order_id === orderNumber) {
        // Actualizar ubicación del conductor
        updateDriverLocation(message.payload);
      }
    }
  });

  return (
    <GoogleMap 
      driverLocation={lastMessage?.payload}
      route={orderTracks?.data?.tracks}
    />
  );
}
```

### Para Conductor (App Móvil/PWA)

```typescript
function DriverLocationSender({ orderId }) {
  const [location, setLocation] = useState(null);

  useEffect(() => {
    // Obtener ubicación cada 15 segundos
    const interval = setInterval(() => {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          
          // Enviar al backend
          fetch('/api/v1/tracks', {
            method: 'POST',
            headers: { 
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
              order_id: orderId,
              latitude,
              longitude
            })
          });
          
          setLocation({ latitude, longitude });
        },
        (error) => console.error('GPS error:', error),
        { enableHighAccuracy: true }
      );
    }, 15000);

    return () => clearInterval(interval);
  }, [orderId]);

  return <div>Enviando ubicación: {JSON.stringify(location)}</div>;
}
```

---

## ⚙️ Configuración y Requisitos

### Backend

1. **PostGIS instalado en PostgreSQL**
```bash
CREATE EXTENSION postgis;
```

2. **Variables de entorno**
```env
JWT_SECRET=your-secret-key
DATABASE_URL=postgresql://user:pass@localhost:5432/shipping_db
```

3. **Ejecutar servidor**
```bash
cd shipping-app-backend
go run main.go
# Servidor en http://localhost:8080
# WebSocket en ws://localhost:8080/api/v1/ws
```

### Frontend (Próximamente)

- React 18+
- TanStack Query para cache
- WebSocket API nativa o library (socket.io-client)
- Google Maps API key

---

## 🔐 Seguridad

### Autenticación
- Todos los endpoints requieren JWT token válido
- Token debe incluir `user_id` y `role` en los claims

### WebSocket
- Validación de rol al conectar
- Clientes solo reciben actualizaciones de sus propias órdenes
- Admins reciben todas las actualizaciones

### Rate Limiting (Recomendado)
- Limitar envíos de ubicación a 1 cada 10 segundos por conductor
- Proteger contra spam de conexiones WebSocket

---

## 📊 Optimizaciones Futuras

1. **Clustering de Marcadores** - Agrupar conductores cercanos en el mapa
2. **Predicción de ETA** - Calcular tiempo estimado de llegada con Google Directions API
3. **Geofencing** - Alertas cuando conductor entra/sale de zonas
4. **Persistencia de Conexión** - Reconexión automática de WebSocket
5. **Compresión de Datos** - Reducir payload de mensajes WebSocket
6. **Redis para PubSub** - Escalar a múltiples instancias del servidor

---

## 🧪 Testing

### Probar Endpoint de Registro
```bash
curl -X POST http://localhost:8080/api/v1/tracks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "order_id": 1,
    "latitude": 6.2442,
    "longitude": -75.5812
  }'
```

### Probar WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

---

## 📝 Notas Importantes

1. **Frecuencia de Envío**: Los conductores deben enviar su ubicación cada 15-30 segundos para tracking fluido
2. **Consumo de Batería**: En apps móviles, considerar `enableHighAccuracy: false` para ahorrar batería
3. **Precisión GPS**: La precisión puede variar (5-50 metros) según el dispositivo y condiciones
4. **Offline Handling**: Implementar cola de ubicaciones pendientes cuando no hay conexión

---

## 🎯 Estado Actual del Proyecto

### ✅ Completado (Backend)
- [x] Endpoint para registrar ubicaciones
- [x] Endpoint para obtener historial de ubicaciones
- [x] Endpoint para obtener conductores activos
- [x] Sistema WebSocket con broadcast selectivo
- [x] Integración con PostGIS

### ⏳ Pendiente (Frontend)
- [ ] Dashboard de admin con mapa en tiempo real
- [ ] Página de tracking para cliente
- [ ] App/PWA para conductor con GPS
- [ ] Hook `useWebSocket` reutilizable
- [ ] Integración con Google Maps API
- [ ] Sistema de notificaciones push

---

## 📚 Referencias

- [Fiber WebSocket Docs](https://docs.gofiber.io/api/middleware/websocket/)
- [PostGIS Documentation](https://postgis.net/documentation/)
- [Google Maps JavaScript API](https://developers.google.com/maps/documentation/javascript)
- [Geolocation API](https://developer.mozilla.org/en-US/docs/Web/API/Geolocation_API)
