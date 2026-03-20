CREATE TABLE teams (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(150) UNIQUE NOT NULL,
    _country_id BIGINT NOT NULL,
    _sport_id BIGINT NOT NULL,
    website_url VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT teams_pk PRIMARY KEY (id),
    CONSTRAINT teams_country_fk FOREIGN KEY (_country_id) REFERENCES countries(id),
    CONSTRAINT teams_sport_fk FOREIGN KEY (_sport_id) REFERENCES sports(id),
    CONSTRAINT teams_valid_timestamps_check CHECK ( created_at <= updated_at )
);