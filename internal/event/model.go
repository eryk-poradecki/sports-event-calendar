package event

import "time"

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
