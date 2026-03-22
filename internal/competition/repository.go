package competition

import (
	"database/sql"
	"fmt"
)

func GetAll(db *sql.DB) ([]Competition, error) {
	competitions := make([]Competition, 0)
	const query = `SELECT id, name, slug, _sport_id, start_date, end_date, description FROM competitions ORDER BY start_date ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all competitions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comp Competition
		var description sql.NullString
		if err := rows.Scan(&comp.ID, &comp.Name, &comp.Slug, &comp.SportID, &comp.StartDate, &comp.EndDate, &description); err != nil {
			return nil, fmt.Errorf("could not scan competition row: %w", err)
		}
		if description.Valid {
			comp.Description = new(description.String)
		}
		competitions = append(competitions, comp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not get all competitions: %w", err)
	}
	return competitions, nil
}

func GetByID(db *sql.DB, id uint64) (Competition, error) {
	var competition Competition
	var description sql.NullString
	const query = `SELECT id, name, slug, _sport_id, start_date, end_date, description FROM competitions WHERE id = $1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&competition.ID, &competition.Name, &competition.Slug, &competition.SportID, &competition.StartDate, &competition.EndDate, &description); err != nil {
		return Competition{}, fmt.Errorf("could not get competition by id: %w", err)
	}
	if description.Valid {
		competition.Description = new(description.String)
	}
	return competition, nil
}
