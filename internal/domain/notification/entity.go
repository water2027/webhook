package notification

type Message struct {
	EventID string `json:"eventId"`
	Source  string `json:"source"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

