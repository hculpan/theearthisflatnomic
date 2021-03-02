package web

import "time"

// TemplateData is the minimum data structure
// to send to all pages
type TemplateData struct {
	Cookies         map[string]string
	ErrorText       string
	Data            interface{}
	MessageText     string
	DestinationURL  string
	UserDisplayName string
	IsTurn          bool
	EndOfTurn       time.Time

	OriginalURL     string
	OriginalURLName string
}

// DisplayEndOfTurn returns the EndOfTurn as a
// readable string
func (t TemplateData) DisplayEndOfTurn() string {
	return t.EndOfTurn.Format("Monday, Jan 2, 2006 3:04pm EST")
}
