package event

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eryk-poradecki/sports-event-calendar/internal/sport"
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
