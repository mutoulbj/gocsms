-- SQL migration
CREATE TABLE charge_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    serial_number VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'UNKNOWN',
    last_heartbeat TIMESTAMPTZ,
    ocpp_protocol VARCHAR(20) NOT NULL DEFAULT '1.6',
    registration_status VARCHAR(20) NOT NULL DEFAULT 'UNKNOWN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    charge_station_id UUID NOT NULL
);

-- Add indexes for performance
CREATE INDEX idx_charge_points_code ON charge_points(code);
CREATE INDEX idx_charge_points_serial_number ON charge_points(serial_number);
CREATE INDEX idx_charge_points_status ON charge_points(status);


CREATE TABLE connectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charge_point_id UUID NOT NULL,
    connector_id VARCHAR(50) NOT NULL,
    standard VARCHAR(20) NOT NULL,
    format VARCHAR(20) NOT NULL,
    power_type VARCHAR(20) NOT NULL,
    max_voltage INTEGER NOT NULL,
    max_amperage INTEGER NOT NULL,
    max_power INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (charge_point_id, connector_id)
);
-- Add indexes for performance
CREATE INDEX idx_connectors_charge_point_id ON connectors(charge_point_id);
CREATE INDEX idx_connectors_connector_id ON connectors(connector_id);
