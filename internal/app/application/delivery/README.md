# Módulo de Entregas e Incidentes

Este módulo maneja el reporte de entregas exitosas e incidentes durante el proceso de entrega de paquetes.

## Casos de Uso Implementados

### 1. Reportar Entrega Exitosa (REPORT DELIVERY)
**Endpoint**: `POST /api/v1/delivery/report`

**Descripción**: Permite a los conductores reportar la entrega exitosa de un paquete con evidencia fotográfica y firma del destinatario.

**Request Body**:
```json
{
  "package_id": 1,
  "observation": "Entrega realizada sin novedad",
  "signature_received": "base64_string_or_url",
  "photo_delivery": "https://example.com/photo.jpg"
}
```

**Respuesta exitosa** (201):
```json
{
  "message": "delivery reported successfully",
  "data": {
    "id": 1,
    "package_id": 1,
    "created_at": "2025-11-16T10:30:00Z"
  }
}
```

**Errores posibles**:
- `400 Bad Request`: Input inválido o foto faltante
- `404 Not Found`: Paquete no encontrado
- `409 Conflict`: El paquete ya tiene información de entrega registrada

---

### 2. Reportar Incidente (REPORT INCIDENT)
**Endpoint**: `POST /api/v1/delivery/incident`

**Descripción**: Permite reportar incidentes durante la entrega (dirección incorrecta, destinatario ausente, paquete dañado, etc.).

**Request Body**:
```json
{
  "package_id": 1,
  "reason_cancellation": "Destinatario ausente",
  "observation": "Se intentó contacto telefónico sin respuesta",
  "photo_evidence": "https://example.com/evidence.jpg"
}
```

**Respuesta exitosa** (201):
```json
{
  "message": "incident reported successfully",
  "data": {
    "id": 1,
    "package_id": 1,
    "created_at": "2025-11-16T10:30:00Z"
  }
}
```

**Errores posibles**:
- `400 Bad Request`: Input inválido, razón o evidencia faltante
- `404 Not Found`: Paquete no encontrado
- `409 Conflict`: El paquete ya tiene un incidente reportado

---

### 3. Obtener Reporte de Paquete (GET PACKAGE REPORT)

#### Por ID de Paquete
**Endpoint**: `GET /api/v1/delivery/package/:packageId/report`

**Descripción**: Obtiene el reporte completo de entrega o incidente de un paquete específico.

**Respuesta exitosa** (200):
```json
{
  "data": {
    "id": 1,
    "package_id": 1,
    "package_num_package": "PKG-2025-001",
    "observation": "Entrega realizada sin novedad",
    "signature_received": "base64_string_or_url",
    "photo_delivery": "https://example.com/photo.jpg",
    "reason_cancellation": null,
    "report_type": "delivery"
  }
}
```

#### Por ID de Reporte
**Endpoint**: `GET /api/v1/delivery/report/:reportId`

**Descripción**: Obtiene el reporte por su ID único.

**Respuesta**: Igual que el anterior

**report_type** puede ser:
- `"delivery"`: Entrega exitosa
- `"incident"`: Incidente reportado

---

## Reglas de Negocio

### Reportar Entrega
1. El paquete debe existir
2. No debe existir un reporte de entrega previo
3. La foto de entrega es obligatoria
4. La firma del destinatario es opcional pero recomendada
5. Se puede agregar observaciones adicionales

### Reportar Incidente
1. El paquete debe existir
2. La razón del incidente es obligatoria
3. La evidencia fotográfica es obligatoria
4. Si ya existe información de entrega, se actualiza con el incidente
5. Se registra en una transacción para garantizar atomicidad

### Tipos de Incidentes Comunes
- Destinatario ausente
- Dirección incorrecta o incompleta
- Paquete dañado
- Condiciones climáticas adversas
- Acceso restringido o denegado
- Destinatario rechaza el paquete
- Problemas de seguridad

---

## Flujo de Trabajo

### Entrega Exitosa:
```
1. Conductor llega al destino
2. Entrega el paquete al destinatario
3. Obtiene firma del destinatario
4. Toma foto del paquete entregado
5. Reporta entrega vía app con POST /delivery/report
6. Sistema registra información en base de datos
7. Estado del paquete se actualiza a "Entregado"
```

### Incidente:
```
1. Conductor encuentra problema en la entrega
2. Documenta el incidente con foto
3. Registra razón específica del problema
4. Reporta incidente vía app con POST /delivery/incident
5. Sistema registra incidente
6. Paquete queda pendiente de resolución
7. Se notifica al coordinador/administrador
```

---

## Datos de la Base de Datos

La información se almacena en la tabla `informationdeliveries`:

```sql
CREATE TABLE informationdeliveries (
    id SERIAL PRIMARY KEY,
    observations TEXT,
    signature_received TEXT,
    photo_delivery TEXT,
    reason_cancellation TEXT,
    idpackage INTEGER NOT NULL,
    FOREIGN KEY (idpackage) REFERENCES packages(id)
);
```

**Campos**:
- `observations`: Comentarios adicionales
- `signature_received`: Firma digital o URL
- `photo_delivery`: URL de la foto de entrega/evidencia
- `reason_cancellation`: Motivo del incidente (si aplica)
- `idpackage`: ID del paquete relacionado

---

## Ejemplos de Uso

### Usando cURL

#### Reportar Entrega:
```bash
curl -X POST http://localhost:8080/api/v1/delivery/report \
  -H "Content-Type: application/json" \
  -d '{
    "package_id": 1,
    "observation": "Entregado en recepción",
    "signature_received": "Juan Pérez",
    "photo_delivery": "https://s3.amazonaws.com/deliveries/photo123.jpg"
  }'
```

#### Reportar Incidente:
```bash
curl -X POST http://localhost:8080/api/v1/delivery/incident \
  -H "Content-Type: application/json" \
  -d '{
    "package_id": 2,
    "reason_cancellation": "Destinatario ausente",
    "observation": "Se dejó aviso en puerta",
    "photo_evidence": "https://s3.amazonaws.com/incidents/photo456.jpg"
  }'
```

#### Consultar Reporte:
```bash
curl http://localhost:8080/api/v1/delivery/package/1/report
```

---

## Integraciones

### Con Otros Módulos:
- **Paquetes**: Valida existencia del paquete antes de crear reporte
- **Transacciones**: Usa transacciones SQL para garantizar consistencia
- **Storage**: Integrar con S3/Cloud Storage para guardar fotos

### Extensiones Futuras:
- [ ] Notificaciones push al destinatario cuando se reporta entrega
- [ ] Envío de email con foto de entrega
- [ ] SMS al destinatario con confirmación
- [ ] Integración con sistema de calificaciones
- [ ] Dashboard de métricas de incidentes
- [ ] Análisis de patrones de incidentes por zona

---

**Versión**: 1.0.0  
**Última actualización**: Noviembre 2025  
**Autor**: Equipo de Desarrollo Logística Express
