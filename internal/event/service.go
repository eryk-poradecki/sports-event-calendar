package event

import (
	"database/sql"
	"fmt"
)

func CreateEvent(db *sql.DB, event *Event) error {
	if event.HomeTeamID == event.AwayTeamID {
		return fmt.Errorf("home and away teams must be different")
	}
	if event.StartTime.IsZero() {
		return fmt.Errorf("start_time is required")
	}
	if event.IsNeutralVenue && event.VenueID == nil {
		return fmt.Errorf("neutral venue events must have a venue")
	}
	if event.HomeScore != nil && *event.HomeScore < 0 {
		return fmt.Errorf("home_score cannot be negative")
	}
	if event.AwayScore != nil && *event.AwayScore < 0 {
		return fmt.Errorf("away_score cannot be negative")
	}
	if event.Status == Scheduled {
		if event.HomeScore != nil || event.AwayScore != nil {
			return fmt.Errorf("scheduled events cannot have scores")
		}
	}
	if event.Status == Cancelled {
		if event.HomeScore != nil || event.AwayScore != nil {
			return fmt.Errorf("cancelled events cannot have scores")
		}
	}

	return Create(db, event)
}
