package event

import (
	"database/sql"
	"time"
)

type Status string

const (
	Scheduled Status = "scheduled"
	Finished  Status = "finished"
	Cancelled Status = "cancelled"
)

type Event struct {
	ID             uint64    `json:"id"`
	SportID        uint64    `json:"_sport_id"`
	CompetitionID  *uint64   `json:"_competition_id"`
	VenueID        *uint64   `json:"_venue_id"`
	HomeTeamID     uint64    `json:"_home_team_id"`
	AwayTeamID     uint64    `json:"_away_team_id"`
	StartTime      time.Time `json:"start_time"`
	Status         Status    `json:"status"`
	HomeScore      *int      `json:"home_score"`
	AwayScore      *int      `json:"away_score"`
	Description    *string   `json:"description"`
	IsNeutralVenue bool      `json:"is_neutral_venue"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type EventListItem struct {
	ID              uint64    `json:"id"`
	SportName       string    `json:"sport_name"`
	CompetitionName string    `json:"competition_name"`
	VenueName       string    `json:"venue_name"`
	HomeTeamName    string    `json:"home_team_name"`
	AwayTeamName    string    `json:"away_team_name"`
	StartTime       time.Time `json:"start_time"`
	Status          Status    `json:"status"`
}

type EventDetails struct {
	ID                  uint64    `json:"id"`
	SportName           string    `json:"sport_name"`
	CompetitionName     string    `json:"competition_name"`
	VenueName           string    `json:"venue_name"`
	HomeTeamName        string    `json:"home_team_name"`
	AwayTeamName        string    `json:"away_team_name"`
	HomeTeamCountryName string    `json:"home_team_country_name"`
	AwayTeamCountryName string    `json:"away_team_country_name"`
	StartTime           time.Time `json:"start_time"`
	Status              Status    `json:"status"`
	HomeScore           *int      `json:"home_score"`
	AwayScore           *int      `json:"away_score"`
	Description         *string   `json:"description"`
	IsNeutralVenue      bool      `json:"is_neutral_venue"`
	HomeTeamURL         *string   `json:"home_team_url"`
	AwayTeamURL         *string   `json:"away_team_url"`
	VenueURL            *string   `json:"venue_url"`
}

type eventDetailsNullableFields struct {
	CompetitionName sql.NullString
	VenueName       sql.NullString
	HomeScore       sql.NullInt64
	AwayScore       sql.NullInt64
	Description     sql.NullString
	HomeTeamURL     sql.NullString
	AwayTeamURL     sql.NullString
	VenueURL        sql.NullString
}

type PaginatedEventsResponse struct {
	Items      []EventListItem `json:"items"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
}
