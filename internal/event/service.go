package event

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eryk-poradecki/sports-event-calendar/internal/competition"
	"github.com/eryk-poradecki/sports-event-calendar/internal/sport"
	"github.com/eryk-poradecki/sports-event-calendar/internal/team"
	"github.com/eryk-poradecki/sports-event-calendar/internal/venue"
)

func CreateEvent(db *sql.DB, event *Event) error {
	if event.StartTime.IsZero() {
		return fmt.Errorf("start_time is required")
	}
	if event.IsNeutralVenue && event.VenueID == nil {
		return fmt.Errorf("%w: neutral venue events must have a venue", ErrInvalidEvent)
	}
	if event.HomeScore != nil && *event.HomeScore < 0 {
		return fmt.Errorf("%w: home_score cannot be negative", ErrInvalidEvent)
	}
	if event.AwayScore != nil && *event.AwayScore < 0 {
		return fmt.Errorf("%w: away_score cannot be negative", ErrInvalidEvent)
	}

	switch event.Status {
	case Scheduled, Finished, Cancelled:
	default:
		return fmt.Errorf("%w: invalid event status", ErrInvalidEvent)
	}

	now := time.Now().UTC()

	if event.Status == Scheduled {
		if event.HomeScore != nil || event.AwayScore != nil {
			return fmt.Errorf("%w: scheduled events cannot have scores", ErrInvalidEvent)
		}
		if event.StartTime.Before(now) {
			return fmt.Errorf("%w: cannot schedule an event in the past", ErrInvalidEvent)
		}
	}
	if event.Status == Cancelled {
		if event.HomeScore != nil || event.AwayScore != nil {
			return fmt.Errorf("%w: cancelled events cannot have scores", ErrInvalidEvent)
		}
	}
	if event.Status == Finished {
		if event.StartTime.After(now) {
			return fmt.Errorf("%w: cannot create a finished event in the future", ErrInvalidEvent)
		}
	}

	homeTeam, err := team.GetByID(db, event.HomeTeamID)
	if err != nil {
		return fmt.Errorf("%w: home team not found", ErrInvalidEvent)
	}
	awayTeam, err := team.GetByID(db, event.AwayTeamID)
	if err != nil {
		return fmt.Errorf("%w: away team not found", ErrInvalidEvent)
	}
	if homeTeam.ID == awayTeam.ID {
		return fmt.Errorf("%w: home and away teams must be different", ErrInvalidEvent)
	}

	if homeTeam.SportID != awayTeam.SportID {
		return fmt.Errorf("%w: home and away teams must belong to the same sport", ErrInvalidEvent)
	}

	if homeTeam.SportID != event.SportID || awayTeam.SportID != event.SportID {
		return fmt.Errorf("%w: event sport must match both teams' sport", ErrInvalidEvent)
	}

	if event.CompetitionID != nil {
		eventCompetition, err := competition.GetByID(db, *event.CompetitionID)
		if err != nil {
			return fmt.Errorf("%w: competition not found", ErrInvalidEvent)
		}
		if eventCompetition.SportID != event.SportID {
			return fmt.Errorf("%w: competition sport must match event sport", ErrInvalidEvent)
		}
	}

	if event.VenueID != nil {
		_, err := venue.GetByID(db, *event.VenueID)
		if err != nil {
			return fmt.Errorf("%w: venue not found", ErrInvalidEvent)
		}
	}

	return Create(db, event)
}

func GetAllEvents(db *sql.DB, page, pageSize int, sportSlug, dateFrom, dateTo string) (PaginatedEventsResponse, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 50 {
		pageSize = 50
	}

	if page < 1 {
		page = 1
	}

	var sportID *uint64

	if sportSlug != "" {
		sportBySlug, err := sport.GetBySlug(db, sportSlug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return PaginatedEventsResponse{}, ErrSportNotFound
			}
			return PaginatedEventsResponse{}, err
		}
		sportID = &sportBySlug.ID
	}

	dateFromTime, err := parseDateParam("from", dateFrom)
	if err != nil {
		return PaginatedEventsResponse{}, err
	}

	dateToTime, err := parseDateParam("to", dateTo)
	if err != nil {
		return PaginatedEventsResponse{}, err
	}

	if dateFromTime != nil && dateToTime != nil && dateFromTime.After(*dateToTime) {
		return PaginatedEventsResponse{}, fmt.Errorf("%w: 'from' date cannot be after 'to' date", ErrInvalidDate)
	}

	if dateToTime != nil {
		dateToTime = new(dateToTime.AddDate(0, 0, 1))
	}

	var response PaginatedEventsResponse

	items, total, err := GetAll(db, page, pageSize, sportID, dateFromTime, dateToTime)
	if err != nil {
		return response, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.Items = items
	response.Page = page
	response.PageSize = pageSize
	response.Total = total
	response.TotalPages = totalPages

	return response, nil
}

func parseDateParam(field, value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		if _, ok := errors.AsType[*time.ParseError](err); ok {
			return nil, fmt.Errorf("%w: %q value %q, expected YYYY-MM-DD", ErrInvalidDate, field, value)
		}
		return nil, err
	}
	return &t, nil
}
