package entity

type Event struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}
