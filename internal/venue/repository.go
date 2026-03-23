package venue

import (
	"database/sql"
	"fmt"
)

func GetAll(db *sql.DB) ([]Venue, error) {
	venues := make([]Venue, 0)
	const query = `SELECT id, name, city, _country_id, address, capacity, website_url FROM venues ORDER BY name ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get venues: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var venue Venue
		var address sql.NullString
		var capacity sql.NullInt64
		var websiteURL sql.NullString
		if err := rows.Scan(&venue.ID, &venue.Name, &venue.City, &venue.CountryID, &address, &capacity, &websiteURL); err != nil {
			return nil, fmt.Errorf("could not scan venue row: %w", err)
		}
		if address.Valid {
			venue.Address = &address.String
		}
		if capacity.Valid {
			venue.Capacity = new(int(capacity.Int64))
		}
		if websiteURL.Valid {
			venue.WebsiteURL = &websiteURL.String
		}
		venues = append(venues, venue)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating venue rows: %w", err)
	}
	return venues, nil
}
