package models

type SQSEmailMessage struct {
	To     string `json:"to"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
