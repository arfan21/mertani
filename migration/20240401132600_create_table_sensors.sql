-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS sensors (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        device_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        type VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now (),
        CONSTRAINT fk_sensors_device_id FOREIGN KEY (device_id) REFERENCES devices (id) ON DELETE CASCADE
    );

CREATE TRIGGER update_sensors_updated_at BEFORE
UPDATE ON sensors FOR EACH ROW EXECUTE PROCEDURE trigger_set_updated ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sensors;

-- +goose StatementEnd