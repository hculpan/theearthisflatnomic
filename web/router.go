package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/hculpan/theearthisflatnomic/utils"
)

var templateList *template.Template

var fileserver = http.FileServer(http.Dir("./resources"))

// SetupRoutes set up the http routes and handlers
func SetupRoutes() {
	LoadTemplates()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/login.html", loginHandler)
	http.HandleFunc("/create_account.html", createAccountHandler)
	http.HandleFunc("/logout.html", logoutHandler)
	http.HandleFunc("/recover_account.html", recoverAccountHandler)
	http.Handle("/resources/", http.StripPrefix("/resources", fileserver))
}

func initializeTemplateData(templateData *TemplateData, req *http.Request) *TemplateData {
	result := templateData

	cookies := map[string]string{}
	for _, v := range req.Cookies() {
		cookies[v.Name] = v.Value
	}

	if result == nil {
		result = &TemplateData{Cookies: cookies}
	} else {
		result.Cookies = cookies
	}

	if token, err := req.Cookie("token"); err == nil {
		if claims, err := utils.DecodeToken(token.Value); err == nil && claims != nil {
			result.UserDisplayName = claims.DisplayName
		}
	}

	return result
}

func executeTemplate(name string, templateData *TemplateData, w http.ResponseWriter, req *http.Request) error {
	templateData = initializeTemplateData(templateData, req)

	if token, err := req.Cookie("token"); err == nil && token != nil {
		if claims, err := utils.DecodeToken(token.Value); err == nil {
			templateData.UserDisplayName = claims.DisplayName
		}
	} else if err.Error() != "http: named cookie not present" {
		fmt.Printf("Error decoding token: %v\n", err)
	}

	return templateList.ExecuteTemplate(w, name, templateData)
}

func errorTemplate(errorText, originalURL, originalURLName string, w http.ResponseWriter, req *http.Request) error {
	templateData := initializeTemplateData(nil, req)

	templateData.ErrorText = errorText
	templateData.OriginalURL = originalURL
	templateData.OriginalURLName = originalURLName

	return templateList.ExecuteTemplate(w, "message_error.gohtml", templateData)
}

func executeTemplateNoUserInfo(name string, data interface{}, errorText string, w http.ResponseWriter, req *http.Request) error {
	cookies := map[string]string{}
	for _, v := range req.Cookies() {
		cookies[v.Name] = v.Value
	}

	templateData := TemplateData{
		ErrorText: errorText,
		Cookies:   cookies,
		Data:      data,
	}

	return templateList.ExecuteTemplate(w, name, templateData)
}

// LoadTemplates load the templates
func LoadTemplates() {
	templateList = template.Must(template.ParseGlob("./templates/*.gohtml"))
}
