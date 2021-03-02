package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hculpan/theearthisflatnomic/entity"
)

func rootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Requested: %s -> %s\n", req.Method, req.URL.Path)

	if req.URL.Path == "/" {
		http.Redirect(w, req, "/index.html", http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "text/html")

		switch {
		case strings.HasSuffix(req.URL.Path, ".html"):
			f := req.URL.Path[1:len(req.URL.Path)-5] + ".gohtml"
			if err := executeTemplate(f, nil, w, req); err != nil {
				fmt.Println("ERROR:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			fmt.Println("File not found:", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

type ruleDisplay struct {
	RuleNumber int
	RuleText   []string
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Requested index.html: %s \n", req.Method)

	w.Header().Set("Content-Type", "text/html")

	templateData := initializeTemplateData(nil, req)

	if templateData.UserDisplayName != "" {
		rules := entity.FindAllActiveRules()
		ruleDisplays := []ruleDisplay{}
		for _, r := range rules {
			rd := ruleDisplay{
				RuleNumber: r.RuleNumber,
				RuleText:   strings.Split(r.RuleText, "<br>"),
			}
			ruleDisplays = append(ruleDisplays, rd)
		}
		templateData.Data = ruleDisplays
	}

	if err := executeTemplate("index.gohtml", templateData, w, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
