package web

// TemplateData is the minimum data structure
// to send to all pages
type TemplateData struct {
	Cookies         map[string]string
	ErrorText       string
	Data            interface{}
	MessageText     string
	DestinationURL  string
	UserDisplayName string

	OriginalURL     string
	OriginalURLName string
}
