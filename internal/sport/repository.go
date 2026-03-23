package sport

import "database/sql"

func GetAll(db *sql.DB) ([]Sport, error) {
	sports := make([]Sport, 0)
	const query = `SELECT id, name, slug FROM sports ORDER BY name ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sport Sport
		if err := rows.Scan(&sport.ID, &sport.Name, &sport.Slug); err != nil {
			return nil, err
		}
		sports = append(sports, sport)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sports, nil
}

func GetBySlug(db *sql.DB, slug string) (Sport, error) {
	const query = `SELECT id, name, slug FROM sports WHERE slug = $1`

	row := db.QueryRow(query, slug)
	var sport Sport
	if err := row.Scan(&sport.ID, &sport.Name, &sport.Slug); err != nil {
		return sport, err
	}
	return sport, nil
}
