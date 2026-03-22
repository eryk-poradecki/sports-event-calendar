package event

import (
	"database/sql"
	"errors"
	"fmt"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEvent(scanner rowScanner) (*EventDetails, error) {
	var ev EventDetails
	var evNullable eventNullableFields

	err := scanner.Scan(
		&ev.ID,
		&ev.SportName,
		&evNullable.CompetitionName,
		&evNullable.VenueName,
		&ev.HomeTeamName,
		&ev.AwayTeamName,
		&ev.StartTime,
		&ev.Status,
		&evNullable.HomeScore,
		&evNullable.AwayScore,
		&evNullable.Description,
		&ev.IsNeutralVenue,
		&ev.CreatedAt,
		&ev.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("could not get event: %w", err)
	}

	if evNullable.CompetitionName.Valid {
		v := string(evNullable.CompetitionName.String)
		ev.CompetitionName = v
	}

	if evNullable.VenueName.Valid {
		v := string(evNullable.VenueName.String)
		ev.VenueName = v
	}

	if evNullable.HomeScore.Valid {
		v := int(evNullable.HomeScore.Int64)
		ev.HomeScore = &v
	}

	if evNullable.AwayScore.Valid {
		v := int(evNullable.AwayScore.Int64)
		ev.AwayScore = &v
	}

	if evNullable.Description.Valid {
		v := evNullable.Description.String
		ev.Description = &v
	}

	return &ev, nil
}

func GetByID(db *sql.DB, id uint64) (*EventDetails, error) {
	const query = `
		SELECT
			events.id,
			sports.name,
			competitions.name,
			venues.name,
			home_teams.name,
			away_teams.name,
			events.start_time,
			events.status,
			events.home_score,
			events.away_score,
			events.description,
			events.is_neutral_venue,
			events.created_at,
			events.updated_at
		FROM events
		JOIN sports ON sports.id = events._sport_id
		LEFT JOIN competitions ON competitions.id = events._competition_id
		LEFT JOIN venues ON venues.id = events._venue_id
		JOIN teams home_teams ON home_teams.id = events._home_team_id
		JOIN teams away_teams ON away_teams.id = events._away_team_id
		WHERE events.id = $1;
	`

	scanner := db.QueryRow(query, id)
	return scanEvent(scanner)
}

func GetAll(db *sql.DB, page, pageSize int) ([]EventDetails, int, error) {
	events := make([]EventDetails, 0)

	const query = `
		SELECT
			events.id,
			sports.name,
			competitions.name,
			venues.name,
			home_teams.name,
			away_teams.name,
			events.start_time,
			events.status,
			events.home_score,
			events.away_score,
			events.description,
			events.is_neutral_venue,
			events.created_at,
			events.updated_at
		FROM events
		JOIN sports ON sports.id = events._sport_id
		LEFT JOIN competitions ON competitions.id = events._competition_id
		LEFT JOIN venues ON venues.id = events._venue_id
		JOIN teams home_teams ON home_teams.id = events._home_team_id
		JOIN teams away_teams ON away_teams.id = events._away_team_id
		ORDER BY events.start_time ASC
		LIMIT $1 OFFSET $2;
	`

	const queryTotal = `
		SELECT
			COUNT(*)
		FROM events
	`

	rows, err := db.Query(query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("could not get events: %w", err)
	}
	defer rows.Close()
	var total int
	db.QueryRow(queryTotal).Scan(&total)

	for rows.Next() {
		ev, err := scanEvent(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("could not get events: %w", err)
		}

		events = append(events, *ev)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error while iterating event rows: %w", err)
	}

	return events, total, nil
}

func Create(db *sql.DB, ev *Event) error {
	const query = `
		INSERT INTO events (
			_sport_id,
			_competition_id,
			_venue_id,
			_home_team_id,
			_away_team_id,
			start_time,
			status,
			home_score,
			away_score,
			description,
			is_neutral_venue
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at;
	`

	err := db.QueryRow(
		query,
		ev.SportID,
		ev.CompetitionID,
		ev.VenueID,
		ev.HomeTeamID,
		ev.AwayTeamID,
		ev.StartTime,
		ev.Status,
		ev.HomeScore,
		ev.AwayScore,
		ev.Description,
		ev.IsNeutralVenue,
	).Scan(&ev.ID, &ev.CreatedAt, &ev.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return nil
}
