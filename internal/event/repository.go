package event

import (
	"database/sql"
	"errors"
	"fmt"
)

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

	var ev Event

	var competitionID sql.NullInt64
	var venueID sql.NullInt64
	var homeScore sql.NullInt64
	var awayScore sql.NullInt64
	var description sql.NullString

	err := db.QueryRow(query, id).Scan(
		&ev.ID,
		&ev.SportID,
		&competitionID,
		&venueID,
		&ev.HomeTeamID,
		&ev.AwayTeamID,
		&ev.StartTime,
		&ev.Status,
		&homeScore,
		&awayScore,
		&description,
		&ev.IsNeutralVenue,
		&ev.CreatedAt,
		&ev.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event with id %d: %w", id, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("could not get event by id: %w", err)
	}

	if competitionID.Valid {
		v := uint64(competitionID.Int64)
		ev.CompetitionID = &v
	}

	if venueID.Valid {
		v := uint64(venueID.Int64)
		ev.VenueID = &v
	}

	if homeScore.Valid {
		v := int(homeScore.Int64)
		ev.HomeScore = &v
	}

	if awayScore.Valid {
		v := int(awayScore.Int64)
		ev.AwayScore = &v
	}

	if description.Valid {
		v := description.String
		ev.Description = &v
	}

	return &ev, nil
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
		return events, fmt.Errorf("could not get events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ev Event

		var competitionID sql.NullInt64
		var venueID sql.NullInt64
		var homeScore sql.NullInt64
		var awayScore sql.NullInt64
		var description sql.NullString

		err := rows.Scan(
			&ev.ID,
			&ev.SportID,
			&competitionID,
			&venueID,
			&ev.HomeTeamID,
			&ev.AwayTeamID,
			&ev.StartTime,
			&ev.Status,
			&homeScore,
			&awayScore,
			&description,
			&ev.IsNeutralVenue,
			&ev.CreatedAt,
			&ev.UpdatedAt,
		)
		if err != nil {
			return events, fmt.Errorf("could not scan event row: %w", err)
		}

		if competitionID.Valid {
			v := uint64(competitionID.Int64)
			ev.CompetitionID = &v
		}

		if venueID.Valid {
			v := uint64(venueID.Int64)
			ev.VenueID = &v
		}

		if homeScore.Valid {
			v := int(homeScore.Int64)
			ev.HomeScore = &v
		}

		if awayScore.Valid {
			v := int(awayScore.Int64)
			ev.AwayScore = &v
		}

		if description.Valid {
			v := description.String
			ev.Description = &v
		}

		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return events, fmt.Errorf("error while iterating event rows: %w", err)
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
