package models

type Message struct {
	Content string `json:"content"`
	Sender  int    `json:"sender"`
}

type NotificationBody struct {
	Channel int     `json:"channel"`
	Message Message `json:"message"`
}

type SmsBody struct {
	Route          string `json:"route"`
	SenderID       string `json:"sender_id"`
	Message        string `json:"message"`
	VariableValues string `json:"variables_values"`
	ScheduledTime  string `json:"schedule_time"`
	Numbers        string `json:"numbers"`
}
