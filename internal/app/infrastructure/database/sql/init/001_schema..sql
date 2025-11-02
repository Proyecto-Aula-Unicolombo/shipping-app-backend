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



CREATE TABLE IF NOT EXISTS statusdelivery (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    date_estimated_delivery TIMESTAMP,
    date_real_delivery TIMESTAMP
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
    iddriver INTEGER NOT NULL,
    idvehicle INTEGER NOT NULL,
    FOREIGN KEY (iddriver) REFERENCES drivers(id),
    FOREIGN KEY (idvehicle) REFERENCES vehicles(id)
);

CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    numpackage VARCHAR(50) NOT NULL UNIQUE,
    startstatus VARCHAR(50) NOT NULL,
    descriptioncontent TEXT,
    weight DECIMAL(8,2),
    dimension VARCHAR(100),
    declared_value DECIMAL(10,2),
    type_package VARCHAR(50),
    is_fragile BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    idaddresspackage INTEGER NOT NULL,
    idstatusdelivery INTEGER NOT NULL,
    idcomercialinformation INTEGER NOT NULL,
    idsender INTEGER NOT NULL,
    idreceivers INTEGER NOT NULL,
    idorder INTEGER,
    FOREIGN KEY (idaddresspackage) REFERENCES addresspackages(id),
    FOREIGN KEY (idstatusdelivery) REFERENCES statusdelivery(id),
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
    location GEOMETRY(POINT, 4326) NOT NULL,
    idorder INTEGER NOT NULL,
    FOREIGN KEY (idorder) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS deliverystops (
    id SERIAL PRIMARY KEY,
    stoplocation GEOMETRY(POINT, 4326) NOT NULL,
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
