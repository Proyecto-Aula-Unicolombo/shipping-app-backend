-- Seed data for shipping application database

-- Insert Users
INSERT INTO users (name, lastname, email, password, role) VALUES
('Ana María', 'González', 'ana.gonzalez@logistica.com', '$2a$10$example.hash.admin', 'admin'),
('Roberto', 'Mendoza', 'roberto.mendoza@logistica.com', '$2a$10$example.hash.admin2', 'admin'),
('Carlos', 'Ramírez', 'carlos.ramirez@logistica.com', '$2a$10$example.hash.driver1', 'driver'),
('Sofía', 'García', 'sofia.garcia@logistica.com', '$2a$10$example.hash.driver2', 'driver'),
('Diego', 'Martínez', 'diego.martinez@logistica.com', '$2a$10$example.hash.driver3', 'driver'),
('Luisa', 'Hernández', 'luisa.hernandez@logistica.com', '$2a$10$example.hash.driver4', 'driver'),
('Miguel', 'Torres', 'miguel.torres@logistica.com', '$2a$10$example.hash.driver5', 'driver'),
('María Elena', 'Rodríguez', 'maria.rodriguez@logistica.com', '$2a$10$example.hash.coord1', 'coord'),
('Alejandro', 'Vargas', 'alejandro.vargas@logistica.com', '$2a$10$example.hash.coord2', 'coord'),
('Carmen', 'Jiménez', 'carmen.jimenez@logistica.com', '$2a$10$example.hash.coord3', 'coord');

-- Insert Drivers (references users 3-7)
INSERT INTO drivers (phonenumber, license, is_active, iduser) VALUES
('+57 300 123 4567', 'COL-2023-001234', false, 3),
('+57 301 234 5678', 'COL-2023-005678', false, 4),
('+57 302 345 6789', 'COL-2023-009012', false, 5),
('+57 303 456 7890', 'COL-2023-003456', false, 6),
('+57 304 567 8901', 'COL-2023-007890', false, 7);

-- Insert Vehicles
INSERT INTO vehicles (plate, brand, model, color, vehicletype) VALUES
('ABC123', 'Toyota', 'Hiace', 'Blanco', 'Furgoneta'),
('DEF456', 'Hyundai', 'H350', 'Gris', 'Furgón'),
('GHI789', 'Ford', 'Transit', 'Azul', 'Camión'),
('JKL012', 'Chevrolet', 'N300', 'Rojo', 'Furgoneta'),
('MNO345', 'Nissan', 'Urvan', 'Blanco', 'Furgoneta'),
('PQR678', 'Mercedes-Benz', 'Sprinter', 'Negro', 'Furgón'),
('STU901', 'Renault', 'Master', 'Blanco', 'Camión'),
('VWX234', 'Iveco', 'Daily', 'Amarillo', 'Camión'),
('YZA567', 'Fiat', 'Ducato', 'Gris', 'Furgoneta');

-- Insert Senders (companies/businesses)
INSERT INTO senders (name, document, address, phonenumber, email, api_key, is_active) VALUES
('ABC Logística S.A.S', '900123456-1', 'Calle 50 #45-30, Medellín', '+57 4 444 1111', 'contacto@abclogistica.com', 'api_key_abc_2024_secure_token_001', true),
('Comercial Del Caribe', '900234567-2', 'Carrera 10 #20-15, Barranquilla', '+57 5 555 2222', 'ventas@comercialcaribe.com', 'api_key_cdc_2024_secure_token_002', true),
('Distribuidora Cartagena', '900345678-3', 'Avenida Pedro de Heredia #30-40, Cartagena', '+57 5 666 3333', 'info@distcartagena.com', 'api_key_dc_2024_secure_token_003', true),
('Mercado Central LTDA', '900456789-4', 'Calle 72 #10-34, Bogotá', '+57 1 777 4444', 'pedidos@mercadocentral.com', 'api_key_mc_2024_secure_token_004', true),
('Importadora Andina', '900567890-5', 'Carrera 43A #1-50, Medellín', '+57 4 888 5555', 'ordenes@impandina.com', NULL, true),
('Textiles del Norte', '900678901-6', 'Calle 30 #25-10, Cúcuta', '+57 7 999 6666', 'envios@textilesn.com', NULL, true),
('Electrónica Express', '900789012-7', 'Avenida El Dorado #68-90, Bogotá', '+57 1 111 7777', 'logistica@electroexpress.com', 'api_key_ee_2024_secure_token_007', true);

-- Insert Receivers (end customers)
INSERT INTO receivers (name, lastname, phonenumber, email) VALUES
('Juan Carlos', 'Pérez López', '+57 310 111 2233', 'juan.perez@email.com'),
('Laura Marcela', 'Gómez Ruiz', '+57 311 222 3344', 'laura.gomez@email.com'),
('Andrés Felipe', 'Moreno Castro', '+57 312 333 4455', 'andres.moreno@email.com'),
('Carolina', 'Sánchez Díaz', '+57 313 444 5566', 'carolina.sanchez@email.com'),
('Roberto', 'Silva Vargas', '+57 314 555 6677', 'roberto.silva@email.com'),
('Patricia', 'Ramírez Torres', '+57 315 666 7788', 'patricia.ramirez@email.com'),
('Fernando', 'López Martínez', '+57 316 777 8899', 'fernando.lopez@email.com'),
('Diana', 'Rodríguez Ospina', '+57 317 888 9900', 'diana.rodriguez@email.com'),
('Camilo', 'Hernández Ríos', '+57 318 999 0011', 'camilo.hernandez@email.com'),
('Valentina', 'Castro Mejía', '+57 319 000 1122', 'valentina.castro@email.com');

-- Insert Commercial Informations
INSERT INTO comercialinformations (cost_sending, is_paid) VALUES
(25000.00, true),
(35000.00, true),
(45000.00, false),
(30000.00, true),
(50000.00, true),
(28000.00, false),
(40000.00, true),
(32000.00, true),
(38000.00, true),
(42000.00, false),
(27000.00, true),
(33000.00, true),
(48000.00, true),
(29000.00, false),
(36000.00, true);


-- Insert Address Packages
INSERT INTO addresspackages (origin, destination, delivery_instructions) VALUES
('Calle 50 #45-30, Medellín', 'Carrera 70 #52-10, Medellín', 'Tocar el timbre dos veces'),
('Carrera 10 #20-15, Barranquilla', 'Calle 84 #51-32, Barranquilla', 'Entregar en portería'),
('Avenida Pedro de Heredia #30-40, Cartagena', 'Calle 30 #8-45, Cartagena', 'Llamar antes de llegar'),
('Calle 72 #10-34, Bogotá', 'Carrera 15 #88-20, Bogotá', 'Oficina 301, tercer piso'),
('Carrera 43A #1-50, Medellín', 'Calle 10 #32-15, Medellín', 'Casa blanca con reja verde'),
('Calle 30 #25-10, Cúcuta', 'Avenida 0 #11-50, Cúcuta', 'Dejar con el vigilante'),
('Avenida El Dorado #68-90, Bogotá', 'Calle 127 #15-20, Bogotá', 'Tocar apartamento 402'),
('Calle 50 #45-30, Medellín', 'Carrera 80 #30-25, Medellín', 'Entregar solo al destinatario'),
('Carrera 10 #20-15, Barranquilla', 'Calle 72 #46-83, Barranquilla', 'Horario de 8am a 6pm'),
('Avenida Pedro de Heredia #30-40, Cartagena', 'Calle 25 #5-60, Cartagena', 'Casa esquinera amarilla'),
('Calle 72 #10-34, Bogotá', 'Carrera 7 #32-16, Bogotá', 'Local comercial'),
('Carrera 43A #1-50, Medellín', 'Calle 33 #70-20, Medellín', 'Edificio Torre del Parque'),
('Calle 30 #25-10, Cúcuta', 'Avenida 6 #14-35, Cúcuta', 'Conjunto residencial Los Pinos'),
('Avenida El Dorado #68-90, Bogotá', 'Calle 100 #18-30, Bogotá', 'Centro Comercial, local 205'),
('Calle 50 #45-30, Medellín', 'Carrera 65 #48-15, Medellín', 'Llamar al llegar');

-- Insert Orders
INSERT INTO orders (create_at, assigned_at, observation, status, typeservice, iddriver, idvehicle) VALUES
('2025-11-16 08:00:00', '2025-11-16 08:30:00', 'Entrega urgente', 'En camino', 'express delivery', 1, 1),
('2025-11-16 09:00:00', '2025-11-16 09:15:00', 'Frágil, manejar con cuidado', 'En camino', 'standard delivery', 2, 2),
('2025-11-15 10:00:00', '2025-11-15 10:20:00', 'Entrega exitosa', 'Entregado', 'standard delivery', 3, 3),
('2025-11-16 07:30:00', '2025-11-16 08:00:00', 'Cliente solicitó entrega en la mañana', 'En camino', 'express delivery', 4, 4),
('2025-11-14 11:00:00', '2025-11-14 11:15:00', 'Sin novedades', 'Entregado', 'standard delivery', 5, 5),
('2025-11-16 12:00:00', NULL, 'Pendiente de asignación', 'Pendiente', 'standard delivery', 1, 1),
('2025-11-15 14:00:00', '2025-11-15 14:30:00', 'Entregado correctamente', 'Entregado', 'standard delivery', 2, 2),
('2025-11-16 10:00:00', '2025-11-16 10:45:00', 'Ruta con múltiples paradas', 'En camino', 'standard delivery', 3, 3),
('2025-11-13 09:00:00', '2025-11-13 09:30:00', 'Cliente no disponible, cancelado', 'Cancelado', 'express delivery', 4, 4),
('2025-11-14 13:00:00', '2025-11-14 13:20:00', 'Entregado en tiempo récord', 'Entregado', 'express delivery', 5, 5);

-- Insert Packages
INSERT INTO packages (numpackage, status, descriptioncontent, weight, dimension, declared_value, type_package, is_fragile, idaddresspackage, idcomercialinformation, idsender, idreceivers, idorder) VALUES
('PKG-2025-001', 'En tránsito', 'Documentos legales', 0.50, '30x20x5 cm', 100000.00, 'Documentos', false, 1, 1, 1, 1, 1),
('PKG-2025-002', 'En tránsito', 'Electrodomésticos', 15.00, '60x50x40 cm', 850000.00, 'Electrodomésticos', true, 2, 2, 2, 2, 2),
('PKG-2025-003', 'Entregado', 'Ropa y textiles', 3.50, '40x30x20 cm', 250000.00, 'Textiles', false, 3, 3, 3, 3, 3),
('PKG-2025-004', 'Pendiente', 'Libros educativos', 8.00, '35x25x30 cm', 180000.00, 'Libros', false, 4, 4, 4, 4, 4),
('PKG-2025-005', 'En tránsito', 'Equipos electrónicos', 5.50, '45x35x25 cm', 1200000.00, 'Electrónicos', true, 5, 5, 5, 5, 5),
('PKG-2025-006', 'Entregado', 'Productos de belleza', 2.00, '25x20x15 cm', 150000.00, 'Cosméticos', false, 6, 6, 6, 6, 6),
('PKG-2025-007', 'Pendiente', 'Repuestos automotrices', 12.00, '50x40x30 cm', 450000.00, 'Repuestos', false, 7, 7, 7, 7, 7),
('PKG-2025-008', 'En tránsito', 'Artículos deportivos', 6.50, '55x35x25 cm', 320000.00, 'Deportes', false, 8, 12, 8, 2, 8),
('PKG-2025-009', 'Cancelado', 'Muebles desmontados', 25.00, '120x80x15 cm', 680000.00, 'Muebles', false, 9, 9, 9, 3, 9),
('PKG-2025-004', 'Pendiente', 'Libros educativos', 8.00, '35x25x30 cm', 180000.00, 'Libros', false, 4, 4, 4, 4, 4, 1),
('PKG-2025-005', 'En tránsito', 'Equipos electrónicos', 5.50, '45x35x25 cm', 1200000.00, 'Electrónicos', true, 5, 5, 5, 5, 5, 4),
('PKG-2025-006', 'Entregado', 'Productos de belleza', 2.00, '25x20x15 cm', 150000.00, 'Cosméticos', false, 6, 6, 6, 6, 6, 5),
('PKG-2025-007', 'Pendiente', 'Repuestos automotrices', 12.00, '50x40x30 cm', 450000.00, 'Repuestos', false, 7, 7, 7, 7, 7, 6),
('PKG-2025-008', 'En tránsito', 'Artículos deportivos', 6.50, '55x35x25 cm', 320000.00, 'Deportes', false, 8, 12, 8, 2, 8, 8),
('PKG-2025-009', 'Cancelado', 'Muebles desmontados', 25.00, '120x80x15 cm', 680000.00, 'Muebles', false, 9, 9, 9, 3, 9, 9),
('PKG-2025-010', 'Entregado', 'Instrumentos musicales', 4.50, '90x30x20 cm', 950000.00, 'Instrumentos', true, 10, 10, 4, 10, 10),
('PKG-2025-011', 'Pendiente', 'Alimentos no perecederos', 10.00, '40x35x30 cm', 120000.00, 'Alimentos', false, 11, 11, 11, 6, 1, 6),
('PKG-2025-012', 'En tránsito', 'Piezas de computadora', 3.00, '35x30x20 cm', 780000.00, 'Electrónicos', true, 12, 12, 12, 7, 2, 8),
('PKG-2025-013', 'Entregado', 'Juguetes infantiles', 4.00, '45x40x35 cm', 220000.00, 'Juguetes', false, 13,13, 1, 3, 7),
('PKG-2025-014', 'Pendiente', 'Herramientas de trabajo', 18.00, '60x45x35 cm', 540000.00, 'Herramientas', false, 14, 14, 2, 4, 6),
('PKG-2025-015', 'En tránsito', 'Accesorios de oficina', 5.00, '40x35x25 cm', 190000.00, 'Oficina', false, 15, 15, 5, 5, 1);

-- Insert Information Deliveries (only for delivered and cancelled packages)
INSERT INTO informationdeliveries (observations, signature_received, photo_delivery, reason_cancellation, idpackage) VALUES
('Entrega exitosa, cliente satisfecho', 'data:signature/base64...', 'https://storage.example.com/delivery/pkg003.jpg', NULL, 3),
('Paquete entregado en perfecto estado', 'data:signature/base64...', 'https://storage.example.com/delivery/pkg006.jpg', NULL, 6),
('Cliente canceló por cambio de dirección', NULL, NULL, 'Cliente solicitó cambio de dirección y canceló orden', 9),
('Entrega sin novedades', 'data:signature/base64...', 'https://storage.example.com/delivery/pkg010.jpg', NULL, 10),
('Recibido por portero del edificio', 'data:signature/base64...', 'https://storage.example.com/delivery/pkg013.jpg', NULL, 13);

-- Insert Tracks (GPS tracking for active orders)
INSERT INTO tracks (timestamp, location, idorder) VALUES
-- Order 1 (Driver 1)
('2025-11-16 08:35:00', ST_SetSRID(ST_MakePoint(-75.5812, 6.2442), 4326), 1),
('2025-11-16 09:00:00', ST_SetSRID(ST_MakePoint(-75.5750, 6.2480), 4326), 1),
('2025-11-16 09:30:00', ST_SetSRID(ST_MakePoint(-75.5690, 6.2520), 4326), 1),
('2025-11-16 10:00:00', ST_SetSRID(ST_MakePoint(-75.5630, 6.2560), 4326), 1),
-- Order 2 (Driver 2)
('2025-11-16 09:20:00', ST_SetSRID(ST_MakePoint(-74.7964, 10.9639), 4326), 2),
('2025-11-16 10:00:00', ST_SetSRID(ST_MakePoint(-74.8020, 10.9700), 4326), 2),
('2025-11-16 10:45:00', ST_SetSRID(ST_MakePoint(-74.8080, 10.9760), 4326), 2),
-- Order 4 (Driver 4)
('2025-11-16 08:05:00', ST_SetSRID(ST_MakePoint(-75.1794, 6.2518), 4326), 4),
('2025-11-16 08:40:00', ST_SetSRID(ST_MakePoint(-75.1850, 6.2580), 4326), 4),
('2025-11-16 09:20:00', ST_SetSRID(ST_MakePoint(-75.1910, 6.2640), 4326), 4),
-- Order 8 (Driver 3)
('2025-11-16 10:50:00', ST_SetSRID(ST_MakePoint(-75.5640, 6.2400), 4326), 8),
('2025-11-16 11:30:00', ST_SetSRID(ST_MakePoint(-75.5580, 6.2450), 4326), 8),
('2025-11-16 12:00:00', ST_SetSRID(ST_MakePoint(-75.5520, 6.2500), 4326), 8);

-- Insert Delivery Stops (incidents and stops during delivery)
INSERT INTO deliverystops (stoplocation, typestop, timestamp, description, evidence, idorder) VALUES
-- Order 1
(ST_SetSRID(ST_MakePoint(-75.5700, 6.2500), 4326), 'Parada', '2025-11-16 09:15:00', 'Parada técnica - Verificación de paquete', NULL, 1),
-- Order 2
(ST_SetSRID(ST_MakePoint(-74.8000, 10.9680), 4326), 'Incidente', '2025-11-16 10:30:00', 'Tráfico pesado en la zona', 'https://storage.example.com/incidents/traffic001.jpg', 2),
-- Order 3 (delivered)
(ST_SetSRID(ST_MakePoint(-75.5400, 6.1800), 4326), 'Entrega', '2025-11-15 14:25:00', 'Entrega exitosa al cliente', 'https://storage.example.com/delivery/order003.jpg', 3),
-- Order 4
(ST_SetSRID(ST_MakePoint(-75.1880, 6.2600), 4326), 'Parada', '2025-11-16 09:00:00', 'Parada para combustible', NULL, 4),
-- Order 5 (delivered)
(ST_SetSRID(ST_MakePoint(-75.5200, 6.2200), 4326), 'Entrega', '2025-11-14 13:10:00', 'Paquete entregado correctamente', 'https://storage.example.com/delivery/order005.jpg', 5),
-- Order 8
(ST_SetSRID(ST_MakePoint(-75.5600, 6.2430), 4326), 'Incidente', '2025-11-16 11:15:00', 'Vía cerrada, tomando ruta alterna', 'https://storage.example.com/incidents/road002.jpg', 8),
-- Order 9 (cancelled)
(ST_SetSRID(ST_MakePoint(-75.5350, 6.2350), 4326), 'Incidente', '2025-11-13 10:00:00', 'Cliente no disponible después de múltiples intentos', NULL, 9),
-- Order 10 (delivered)
(ST_SetSRID(ST_MakePoint(-75.5100, 6.2100), 4326), 'Entrega', '2025-11-14 15:50:00', 'Entrega rápida y eficiente', 'https://storage.example.com/delivery/order010.jpg', 10);
