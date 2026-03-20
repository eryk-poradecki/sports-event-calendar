CREATE TABLE sports (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT sports_pk PRIMARY KEY (id),
    CONSTRAINT sports_valid_timestamps_check CHECK ( created_at <= updated_at )
);