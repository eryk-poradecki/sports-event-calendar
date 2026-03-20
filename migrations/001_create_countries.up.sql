CREATE TABLE countries (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code_alpha2 CHAR(2) NOT NULL UNIQUE ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT countries_pk PRIMARY KEY (id),
    CONSTRAINT countries_valid_timestamps_check CHECK ( created_at <= updated_at )
);