package event

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEventListItem(scanner rowScanner) (*EventListItem, error) {
	var ev EventListItem
	var competitionName sql.NullString
	var venueName sql.NullString

	err := scanner.Scan(
		&ev.ID,
		&ev.SportName,
		&competitionName,
		&venueName,
		&ev.HomeTeamName,
		&ev.AwayTeamName,
		&ev.StartTime,
		&ev.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("could not get event: %w", err)
	}

	if competitionName.Valid {
		v := competitionName.String
		ev.CompetitionName = v
	}

	if venueName.Valid {
		v := venueName.String
		ev.VenueName = v
	}

	return &ev, nil
}

func scanEventDetails(scanner rowScanner) (*EventDetails, error) {
	var ev EventDetails
	var evNullable eventDetailsNullableFields

	err := scanner.Scan(
		&ev.ID,
		&ev.SportName,
		&evNullable.CompetitionName,
		&evNullable.VenueName,
		&ev.HomeTeamName,
		&ev.AwayTeamName,
		&ev.HomeTeamCountryName,
		&ev.AwayTeamCountryName,
		&ev.StartTime,
		&ev.Status,
		&evNullable.HomeScore,
		&evNullable.AwayScore,
		&evNullable.Description,
		&ev.IsNeutralVenue,
		&evNullable.HomeTeamURL,
		&evNullable.AwayTeamURL,
		&evNullable.VenueURL,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("could not get event details: %w", err)
	}

	if evNullable.CompetitionName.Valid {
		v := evNullable.CompetitionName.String
		ev.CompetitionName = v
	}
	if evNullable.VenueName.Valid {
		v := evNullable.VenueName.String
		ev.VenueName = v
	}
	if evNullable.HomeScore.Valid {
		ev.HomeScore = new(int(evNullable.HomeScore.Int64))
	}
	if evNullable.AwayScore.Valid {
		ev.AwayScore = new(int(evNullable.AwayScore.Int64))
	}
	if evNullable.Description.Valid {
		ev.Description = new(evNullable.Description.String)
	}
	if evNullable.HomeTeamURL.Valid {
		ev.HomeTeamURL = new(evNullable.HomeTeamURL.String)
	}
	if evNullable.AwayTeamURL.Valid {
		ev.AwayTeamURL = new(evNullable.AwayTeamURL.String)
	}
	if evNullable.VenueURL.Valid {
		ev.VenueURL = new(evNullable.VenueURL.String)
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
			home_countries.name,
			away_countries.name,
			events.start_time,
			events.status,
			events.home_score,
			events.away_score,
			events.description,
			events.is_neutral_venue,
			home_teams.website_url,
			away_teams.website_url,
			venues.website_url
		FROM events
		JOIN sports ON sports.id = events._sport_id
		LEFT JOIN competitions ON competitions.id = events._competition_id
		LEFT JOIN venues ON venues.id = events._venue_id
		JOIN teams home_teams ON home_teams.id = events._home_team_id
		JOIN teams away_teams ON away_teams.id = events._away_team_id
		JOIN countries home_countries ON home_countries.id = home_teams._country_id
		JOIN countries away_countries ON away_countries.id = away_teams._country_id
		WHERE events.id = $1;
	`

	scanner := db.QueryRow(query, id)
	return scanEventDetails(scanner)
}

func GetAll(db *sql.DB, page, pageSize int, sportID *uint64, dateFrom, dateTo *time.Time) ([]EventListItem, int, error) {
	events := make([]EventListItem, 0)

	baseQuery := `
		SELECT
			events.id,
			sports.name,
			competitions.name,
			venues.name,
			home_teams.name,
			away_teams.name,
			events.start_time,
			events.status
		FROM events
		JOIN sports ON sports.id = events._sport_id
		LEFT JOIN competitions ON competitions.id = events._competition_id
		LEFT JOIN venues ON venues.id = events._venue_id
		JOIN teams home_teams ON home_teams.id = events._home_team_id
		JOIN teams away_teams ON away_teams.id = events._away_team_id
	`

	var conditions []string
	var args []any
	argN := 1

	// Keep separate filter args for the paginated data query and the count query.
	// Currently, the filters are almost identical, but the queries can diverge over time
	// (for example because of different joins or an FK becoming nullable), so separate condition/arg slices are safer.
	var countConditions []string
	var countArgs []any
	countArgN := 1

	if sportID != nil {
		conditions = append(conditions, fmt.Sprintf("events._sport_id = $%d", argN))
		args = append(args, *sportID)
		argN++
		countConditions = append(countConditions, fmt.Sprintf("events._sport_id = $%d", countArgN))
		countArgs = append(countArgs, *sportID)
		countArgN++
	}
	if dateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("events.start_time >= $%d", argN))
		args = append(args, *dateFrom)
		argN++
		countConditions = append(countConditions, fmt.Sprintf("events.start_time >= $%d", countArgN))
		countArgs = append(countArgs, *dateFrom)
		countArgN++
	}
	if dateTo != nil {
		conditions = append(conditions, fmt.Sprintf("events.start_time < $%d", argN))
		args = append(args, *dateTo)
		argN++
		countConditions = append(countConditions, fmt.Sprintf("events.start_time < $%d", countArgN))
		countArgs = append(countArgs, *dateTo)
		countArgN++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY events.start_time ASC LIMIT $%d OFFSET $%d", argN, argN+1)
	args = append(args, pageSize, (page-1)*pageSize)

	countQuery := `
		SELECT COUNT(*)
		FROM events
		JOIN sports ON sports.id = events._sport_id
		JOIN teams home_teams ON home_teams.id = events._home_team_id
		JOIN teams away_teams ON away_teams.id = events._away_team_id
	`

	if len(countConditions) > 0 {
		countQuery += " WHERE " + strings.Join(countConditions, " AND ")
	}

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("could not get events: %w", err)
	}
	defer rows.Close()

	var total int
	err = db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("could not count events: %w", err)
	}

	for rows.Next() {
		ev, err := scanEventListItem(rows)
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
