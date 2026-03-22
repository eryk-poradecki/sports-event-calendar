package team

type Team struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	CountryID  uint64  `json:"country_id"`
	SportID    uint64  `json:"sport_id"`
	WebsiteURL *string `json:"website_url"`
}
