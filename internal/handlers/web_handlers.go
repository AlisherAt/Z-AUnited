package main

import (
	"html/template"
	"net/http"
)

func feedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, nil)
}

func liveMatchesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/live.html"))
	tmpl.Execute(w, nil)
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/analytics.html"))
	tmpl.Execute(w, nil)
}

func communityHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/community.html"))
	tmpl.Execute(w, nil)
}

func leagueTableHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/league.html"))
	tmpl.Execute(w, nil)
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/account.html"))
	tmpl.Execute(w, nil)
}
