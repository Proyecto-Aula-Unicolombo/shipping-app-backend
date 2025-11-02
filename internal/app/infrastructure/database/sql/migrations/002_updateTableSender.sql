ALTER TABLE senders 
ADD COLUMN IF NOT EXISTS api_key VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

CREATE INDEX IF NOT EXISTS idx_senders_api_key ON senders(api_key) WHERE api_key IS NOT NULL;

UPDATE senders SET is_active = true WHERE api_key IS NULL;
