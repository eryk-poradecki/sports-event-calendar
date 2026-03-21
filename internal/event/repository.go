package event

import (
	"database/sql"
	"errors"
	"fmt"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEvent(scanner rowScanner) (*Event, error) {
	var ev Event
	var evNullable eventNullableFields

	err := scanner.Scan(
		&ev.ID,
		&ev.SportID,
		&evNullable.CompetitionID,
		&evNullable.VenueID,
		&ev.HomeTeamID,
		&ev.AwayTeamID,
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

	if evNullable.CompetitionID.Valid {
		v := uint64(evNullable.CompetitionID.Int64)
		ev.CompetitionID = &v
	}

	if evNullable.VenueID.Valid {
		v := uint64(evNullable.VenueID.Int64)
		ev.VenueID = &v
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

func GetByID(db *sql.DB, id uint64) (*Event, error) {
	const query = `
		SELECT
			id,
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
			is_neutral_venue,
			created_at,
			updated_at
		FROM events
		WHERE id = $1;
	`

	scanner := db.QueryRow(query, id)
	return scanEvent(scanner)
}

func GetAll(db *sql.DB) ([]Event, error) {
	events := make([]Event, 0)

	const query = `
		SELECT
			id,
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
			is_neutral_venue,
			created_at,
			updated_at
		FROM events
		ORDER BY start_time ASC;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ev, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("could not get events: %w", err)
		}

		events = append(events, *ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating event rows: %w", err)
	}

	return events, nil
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
		RETURNING id;
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
	).Scan(&ev.ID)
	if err != nil {
		return fmt.Errorf("could not create event: %w", err)
	}

	return nil
}
