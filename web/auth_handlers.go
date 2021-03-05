package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hculpan/theearthisflatnomic/entity"
	"github.com/hculpan/theearthisflatnomic/utils"
)

func createAccountHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Create account requested: %s -> %s\n", req.Method, req.URL.Path)

	w.Header().Set("Content-Type", "text/html")

	switch req.Method {
	case "GET":
		if err := executeTemplate("create_account.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		p1 := req.FormValue("inputPassword")
		p2 := req.FormValue("confirmPassword")
		username := req.FormValue("inputEmail")
		fullname := req.FormValue("inputFullName")
		displayname := req.FormValue("inputDisplayName")

		switch {
		case p1 == "":
			if err := executeTemplate("create_account.gohtml", &TemplateData{ErrorText: "Password cannot be empty"}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case p1 != p2:
			if err := executeTemplate("create_account.gohtml", &TemplateData{ErrorText: "Passwords do not match"}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			user, err := entity.AddNewUser(fullname, displayname, username, p1)
			if err != nil {
				errorTemplate(err.Error(), "/create_account.html", "Create an Account", w, req)
				return
			}

			token, err := utils.CreateToken(user.Username, user.FullName, user.DisplayName)
			if err != nil {
				fmt.Printf("ERROR: %+v\n", err)
				if err2 := errorTemplate("Internal server error: "+err.Error(), "/create_account.html", "Create account", w, req); err2 != nil {
					http.Error(w, err2.Error(), http.StatusInternalServerError)
				}
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: time.Now().Add(3 * time.Hour)})
			http.Redirect(w, req, "/index.html", http.StatusSeeOther)
		}
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Login requested: %s -> %s\n", req.Method, req.URL.Path)

	w.Header().Set("Content-Type", "text/html")

	if req.FormValue("inputEmail") != "" {
		username := req.FormValue("inputEmail")
		password := req.FormValue("inputPassword")
		if user, err := entity.Authenticate(username, password); err == nil {
			fmt.Printf("User logged in: %s/%s\n", user.FullName, user.Username)
			token, err := utils.CreateToken(user.Username, user.FullName, user.DisplayName)
			if err != nil {
				fmt.Printf("ERROR: %+v\n", err)
				if err2 := errorTemplate("Internal server error: "+err.Error(), "/login.html", "Login", w, req); err2 != nil {
					http.Error(w, err2.Error(), http.StatusInternalServerError)
				}
			} else {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: time.Now().Add(3 * time.Hour)})
				http.Redirect(w, req, "/index.html", http.StatusSeeOther)
			}
		} else {
			fmt.Printf("Authentication failed: %s\n", username)
			if err2 := errorTemplate("Invalid username/password", "/login.html", "Login", w, req); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		if err := executeTemplate("login.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	})
	if err := executeTemplateNoUserInfo("logout.gohtml", nil, "", w, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func recoverAccountHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Recover account requested: %s -> %s\n", req.Method, req.URL.Path)

	w.Header().Set("Content-Type", "text/html")

	switch req.Method {
	case "GET":
		if err := executeTemplate("recover_account.gohtml", nil, w, req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		username := req.FormValue("inputEmail")

		_, found := entity.FindUserByUsername(username)
		if found {
			if err := executeTemplate("message.gohtml", &TemplateData{ErrorText: "Email with a new password has been sent."}, w, req); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {

		}
	}
}
