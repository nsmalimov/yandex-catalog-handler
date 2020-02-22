CREATE TABLE operate_log (
    id SERIAL NOT NULL PRIMARY KEY,
    cause VARCHAR(255),
    results JSON,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);