CREATE TABLE IF NOT EXISTS shipments (
    id               TEXT PRIMARY KEY,
    reference_number TEXT NOT NULL UNIQUE,
    origin           TEXT NOT NULL,
    destination      TEXT NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending',
    driver_name      TEXT NOT NULL,
    driver_unit      TEXT NOT NULL,
    amount           NUMERIC(10,2) NOT NULL DEFAULT 0,
    driver_revenue   NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL,
    updated_at       TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS shipment_events (
    id          TEXT PRIMARY KEY,
    shipment_id TEXT NOT NULL REFERENCES shipments(id),
    status      TEXT NOT NULL,
    note        TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_shipment_events_shipment_id 
    ON shipment_events(shipment_id);