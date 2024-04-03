-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS devices (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR(255) NOT NULL,
        description TEXT,
        type VARCHAR(255) NOT NULL,
        location VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT now (),
        updated_at TIMESTAMP DEFAULT now ()
    );

CREATE TRIGGER update_devices_updated_at BEFORE
UPDATE ON devices FOR EACH ROW EXECUTE PROCEDURE trigger_set_updated ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;

-- +goose StatementEnd