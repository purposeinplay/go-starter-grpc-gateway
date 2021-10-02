CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    PRIMARY KEY (id),
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    email         VARCHAR(255)  UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);