package web

import (
	"fmt"
	"net/http"
	"strings"
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
