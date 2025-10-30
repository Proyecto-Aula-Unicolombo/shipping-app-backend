CREATE INDEX IF NOT EXISTS idx_addresspackages_route ON addresspackages(origin, destination);
CREATE INDEX IF NOT EXISTS idx_senders_email ON senders(email);
CREATE INDEX IF NOT EXISTS idx_senders_document ON senders(document);
CREATE INDEX IF NOT EXISTS idx_receivers_email ON receivers(email);
