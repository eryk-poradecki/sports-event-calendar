package team

import (
	"database/sql"
	"fmt"
)

func GetAll(db *sql.DB) ([]Team, error) {
	teams := make([]Team, 0)
	const query = `SELECT id, name, slug, _country_id, _sport_id, website_url FROM teams ORDER BY name ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get teams: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var team Team
		var websiteURL sql.NullString
		if err := rows.Scan(&team.ID, &team.Name, &team.Slug, &team.CountryID, &team.SportID, &websiteURL); err != nil {
			return nil, fmt.Errorf("could not scan team row: %w", err)
		}
		if websiteURL.Valid {
			team.WebsiteURL = new(websiteURL.String)
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating team rows: %w", err)
	}
	return teams, nil
}

func GetByID(db *sql.DB, id uint64) (Team, error) {
	var team Team
	var websiteURL sql.NullString
	const query = `SELECT id, name, slug, _country_id, _sport_id, website_url FROM teams WHERE id = $1`

	row := db.QueryRow(query, id)
	if err := row.Scan(&team.ID, &team.Name, &team.Slug, &team.CountryID, &team.SportID, &websiteURL); err != nil {
		return team, fmt.Errorf("could not get team by id: %w", err)
	}
	if websiteURL.Valid {
		team.WebsiteURL = new(websiteURL.String)
	}
	return team, nil
}
