CREATE TYPE status AS ENUM ('scheduled', 'finished', 'cancelled');

CREATE TABLE events(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    _sport_id BIGINT NOT NULL,
    _competition_id BIGINT NULL,
    _venue_id BIGINT NULL,
    _home_team_id BIGINT NOT NULL,
    _away_team_id BIGINT NOT NULL,
    start_time timestamptz NOT NULL,
    status status NOT NULL,
    home_score INTEGER NULL,
    away_score INTEGER NULL,
    description TEXT NULL,
    is_neutral_venue BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT events_pk PRIMARY KEY (id),
    CONSTRAINT events_sport_fk FOREIGN KEY (_sport_id) REFERENCES sports(id),
    CONSTRAINT events_competition_fk FOREIGN KEY (_competition_id) REFERENCES competitions(id),
    CONSTRAINT events_venue_fk FOREIGN KEY (_venue_id) REFERENCES venues(id),
    CONSTRAINT events_home_team_fk FOREIGN KEY (_home_team_id) REFERENCES teams(id),
    CONSTRAINT events_away_team_fk FOREIGN KEY (_away_team_id) REFERENCES teams(id),
    CONSTRAINT events_home_score_check CHECK ( home_score >= 0 ),
    CONSTRAINT events_away_score_check CHECK ( away_score >= 0 ),
    CONSTRAINT events_home_not_away_check CHECK ( _home_team_id != _away_team_id ),
    CONSTRAINT events_neutral_venue_check CHECK (is_neutral_venue = FALSE OR _venue_id IS NOT NULL),
    CONSTRAINT events_valid_timestamps_check CHECK ( created_at <= updated_at )
);