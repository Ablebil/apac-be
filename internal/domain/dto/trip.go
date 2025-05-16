package dto

type TripSummaryResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Destination string `json:"destination"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Travelers   int    `json:"travelers"`
}
