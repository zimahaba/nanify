package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var pages = map[string]string{
	"base":    "web/templates/base.html",
	"shrink":  "web/templates/shrink.html",
	"result":  "web/templates/result.html",
	"expired": "web/templates/expired.html",
	"myurls":  "web/templates/myurls.html",
}

func render(w http.ResponseWriter, tmpl string, data any) {
	t := template.Must(template.ParseFiles(pages["base"], tmpl))
	err := t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		fmt.Println("Error when executing template", err)
	}
}

func handleCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("shrink")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "shrink",
			Value: generateKey(),
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
	return cookie
}
