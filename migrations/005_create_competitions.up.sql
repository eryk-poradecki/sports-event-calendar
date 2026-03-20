CREATE TABLE competitions(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(150) UNIQUE NOT NULL,
    _sport_id BIGINT NOT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    description TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT competitions_pk PRIMARY KEY (id),
    CONSTRAINT competitions_sport_fk FOREIGN KEY (_sport_id) REFERENCES sports(id),
    CONSTRAINT competitions_date_check CHECK ( end_date >= start_date ),
    CONSTRAINT competitions_valid_timestamps_check CHECK ( created_at <= updated_at )
);