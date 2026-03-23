package venue

type Venue struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	CountryID  uint64  `json:"country_id"`
	Address    *string `json:"address"`
	Capacity   *int    `json:"capacity"`
	WebsiteURL *string `json:"website_url"`
}
