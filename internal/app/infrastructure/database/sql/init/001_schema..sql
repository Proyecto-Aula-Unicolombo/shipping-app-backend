CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    phonenumber VARCHAR(20) NOT NULL,
    license VARCHAR(50) NOT NULL UNIQUE,
    iduser INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT false,
    FOREIGN KEY (iduser) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    plate VARCHAR(20) NOT NULL UNIQUE,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    color VARCHAR(30) NOT NULL,
    vehicletype VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS senders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    document VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    phonenumber VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    api_key VARCHAR(255) UNIQUE,
    is_active BOOLEAN DEFAULT true
);


CREATE TABLE IF NOT EXISTS receivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    phonenumber VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS comercialinformations (
    id SERIAL PRIMARY KEY,
    cost_sending DECIMAL(10,2) NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE
);


CREATE TABLE IF NOT EXISTS addresspackages (
    id SERIAL PRIMARY KEY,
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    delivery_instructions TEXT
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_at TIMESTAMP,
    observation TEXT,
    status VARCHAR(50) NOT NULL,
    typeservice VARCHAR(100) NOT NULL,
    iddriver INTEGER,
    idvehicle INTEGER,
    FOREIGN KEY (iddriver) REFERENCES drivers(id),
    FOREIGN KEY (idvehicle) REFERENCES vehicles(id)
);

CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    numpackage VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    descriptioncontent TEXT,
    weight DECIMAL(8,2),
    dimension VARCHAR(100),
    declared_value DECIMAL(10,2),
    type_package VARCHAR(50),
    is_fragile BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    idaddresspackage INTEGER NOT NULL,
    idcomercialinformation INTEGER NOT NULL,
    idsender INTEGER NOT NULL,
    idreceivers INTEGER NOT NULL,
    idorder INTEGER,
    FOREIGN KEY (idaddresspackage) REFERENCES addresspackages(id),
    FOREIGN KEY (idcomercialinformation) REFERENCES comercialinformations(id),
    FOREIGN KEY (idsender) REFERENCES senders(id),
    FOREIGN KEY (idreceivers) REFERENCES receivers(id),
    FOREIGN KEY (idorder) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS informationdeliveries (
    id SERIAL PRIMARY KEY,
    observations TEXT,
    signature_received TEXT,
    photo_delivery TEXT,
    reason_cancellation TEXT,
    idpackage INTEGER NOT NULL,
    FOREIGN KEY (idpackage) REFERENCES packages(id)
);

CREATE TABLE IF NOT EXISTS tracks (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    idorder INTEGER NOT NULL,
    FOREIGN KEY (idorder) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS deliverystops (
    id SERIAL PRIMARY KEY,
    stoplocation GEOGRAPHY(POINT, 4326) NOT NULL,
    typestop VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    evidence TEXT,
    idorder INTEGER NOT NULL,
    FOREIGN KEY (idorder) REFERENCES orders(id)
);

CREATE INDEX idx_packages_numpackage ON packages(numpackage);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_drivers_license ON drivers(license);
CREATE INDEX idx_vehicles_plate ON vehicles(plate);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_tracks_timestamp ON tracks(timestamp);
CREATE INDEX idx_tracks_location ON tracks USING GIST(location);
CREATE INDEX idx_deliverystops_location ON deliverystops USING GIST(stoplocation);
CREATE INDEX IF NOT EXISTS idx_addresspackages_route ON addresspackages(origin, destination);
CREATE INDEX IF NOT EXISTS idx_senders_email ON senders(email);
CREATE INDEX IF NOT EXISTS idx_senders_document ON senders(document);
CREATE INDEX IF NOT EXISTS idx_receivers_email ON receivers(email);
CREATE INDEX IF NOT EXISTS idx_senders_api_key ON senders(api_key) WHERE api_key IS NOT NULL;



-- trigger to update status of package when status of order is updated

CREATE OR REPLACE FUNCTION update_status_order_and_package()
RETURNS TRIGGER AS $$
DECLARE  
	status_order VARCHAR(50);
BEGIN	
	
	SELECT status INTO status_order
	FROM orders
	WHERE id = NEW.idorder;

	IF status_order = 'asignada' THEN
		
		UPDATE orders
		SET status = 'en camino'
		WHERE id = NEW.idorder;
	
		UPDATE packages 
		SET status = 'en camino',
			updated_at = CURRENT_TIMESTAMP
		WHERE idorder = NEW.idorder 
			AND status = 'asignado';

	END IF;

		RETURN NEW;
	 
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION check_order_completion()
RETURNS TRIGGER AS $$
DECLARE
    total_paquetes INTEGER;
    paquetes_entregados INTEGER;
    paquetes_en_camino INTEGER;
    paquetes_cancelados INTEGER;
    paquetes_incidente INTEGER;
    estado_final VARCHAR(50);
BEGIN
    -- Ejecutar si el paquete cambió a un estado final (entregado, cancelado o incidente)
    IF NEW.status IN ('entregado', 'cancelado', 'incidente') 
       AND OLD.status = 'en camino' THEN
        
        -- Contar todos los paquetes de la orden
        SELECT COUNT(*) INTO total_paquetes
        FROM packages
        WHERE idorder = NEW.idorder;
        
        -- Contar paquetes por estado
        SELECT 
            COUNT(*) FILTER (WHERE status = 'entregado') AS entregados,
            COUNT(*) FILTER (WHERE status = 'en camino') AS en_camino,
            COUNT(*) FILTER (WHERE status = 'cancelado') AS cancelados,
            COUNT(*) FILTER (WHERE status = 'incidente') AS incidentes
        INTO 
            paquetes_entregados,
            paquetes_en_camino,
            paquetes_cancelados,
            paquetes_incidente
        FROM packages
        WHERE idorder = NEW.idorder;
        
        -- Solo actualizar la orden si NO quedan paquetes "en camino"
        IF paquetes_en_camino = 0 THEN
            
            -- Decidir el estado final según lo que haya pasado
            IF paquetes_entregados = total_paquetes THEN
                -- Todos entregados exitosamente
                estado_final := 'completada';
                
            ELSIF paquetes_cancelados = total_paquetes THEN
                estado_final := 'cancelada';
                
            ELSIF paquetes_incidente = total_paquetes THEN
                estado_final := 'incidente';
                
            ELSIF paquetes_entregados > 0 AND (paquetes_cancelados > 0 OR paquetes_incidente > 0) THEN
                -- Mix: algunos entregados y otros con problemas
                estado_final := 'parcialmente_completada';
                
            ELSE
                estado_final := 'finalizada';
            END IF;
            
            UPDATE orders
            SET status = estado_final
            WHERE id = NEW.idorder;
            
        END IF;
        
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;




CREATE TRIGGER trigger_check_order_completion
AFTER UPDATE ON packages
FOR EACH ROW
EXECUTE FUNCTION check_order_completion();

CREATE TRIGGER trigger_update_status_order_and_package
AFTER INSERT ON tracks 
FOR EACH ROW
EXECUTE FUNCTION update_status_order_and_package();

