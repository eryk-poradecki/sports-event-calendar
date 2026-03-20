CREATE TABLE venues (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(150) NOT NULL,
    city VARCHAR(100) NOT NULL,
    _country_id BIGINT NOT NULL,
    address VARCHAR(200) NULL,
    capacity INTEGER NULL,
    website_url VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT venues_pk PRIMARY KEY (id),
    CONSTRAINT venues_country_fk FOREIGN KEY (_country_id) REFERENCES countries(id),
    CONSTRAINT venues_capacity_check CHECK ( capacity >= 0 ),
    CONSTRAINT venues_valid_timestamps_check CHECK ( created_at <= updated_at )
);