package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"project/internal/cache"
	"project/models"
	"project/storage"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var wsClients = make(map[*websocket.Conn]bool)
var wsMutex sync.Mutex

func broadcastStandings(standings []map[string]interface{}) {
	wsMutex.Lock()
	defer wsMutex.Unlock()
	for c := range wsClients {
		c.WriteJSON(map[string]interface{}{"standings": standings})
	}
}

func wsStandingsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	wsMutex.Lock()
	wsClients[conn] = true
	wsMutex.Unlock()
	// Send initial standings
	standings := getLiveStandings()
	conn.WriteJSON(map[string]interface{}{"standings": standings})
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	wsMutex.Lock()
	delete(wsClients, conn)
	wsMutex.Unlock()
	conn.Close()
}

// Dummy function for now, should call TableService and mark live teams
func getLiveStandings() []map[string]interface{} {
	// Full league table with Manchester United first
	return []map[string]interface{}{
		{"team": "Manchester United", "played": 25, "points": 65, "gd": 40, "live": true},
		{"team": "Arsenal", "played": 25, "points": 60, "gd": 35, "live": false},
		{"team": "Man City", "played": 25, "points": 59, "gd": 33, "live": false},
		{"team": "Liverpool", "played": 25, "points": 58, "gd": 30, "live": false},
		{"team": "Chelsea", "played": 25, "points": 45, "gd": 10, "live": false},
		{"team": "Tottenham", "played": 25, "points": 43, "gd": 8, "live": false},
		{"team": "Newcastle", "played": 25, "points": 41, "gd": 12, "live": false},
		{"team": "Brighton", "played": 25, "points": 39, "gd": 7, "live": false},
		{"team": "Brentford", "played": 25, "points": 36, "gd": 3, "live": false},
		{"team": "Fulham", "played": 25, "points": 35, "gd": 2, "live": false},
		{"team": "Crystal Palace", "played": 25, "points": 30, "gd": -5, "live": false},
		{"team": "Aston Villa", "played": 25, "points": 29, "gd": -7, "live": false},
		{"team": "Leicester", "played": 25, "points": 28, "gd": -8, "live": false},
		{"team": "Wolves", "played": 25, "points": 27, "gd": -10, "live": false},
		{"team": "West Ham", "played": 25, "points": 25, "gd": -12, "live": false},
		{"team": "Leeds", "played": 25, "points": 23, "gd": -15, "live": false},
		{"team": "Everton", "played": 25, "points": 22, "gd": -18, "live": false},
		{"team": "Bournemouth", "played": 25, "points": 21, "gd": -20, "live": false},
		{"team": "Southampton", "played": 25, "points": 18, "gd": -25, "live": false},
	}
}

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

var db = storage.NewStorage()
var taskChan = make(chan string)

func backgroundWorker() {
	for msg := range taskChan {
		time.Sleep(1 * time.Second)
		fmt.Println("Background worker processed:", msg)
	}
}

func main() {
	http.HandleFunc("/api/calendar/", func(w http.ResponseWriter, r *http.Request) {
		// Stub: serve a downloadable .ics file for the favorite team
		w.Header().Set("Content-Type", "text/calendar")
		w.Header().Set("Content-Disposition", "attachment; filename=matches.ics")
		// TODO: Generate real calendar data from DB
		w.Write([]byte("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//EPLHub//EN\nEND:VCALENDAR\n"))
	})
	// Initialize Redis (default: localhost:6379, db 0, no password)
	cache.InitRedis("localhost:6379", "", 0)
	go backgroundWorker()

	http.HandleFunc("/", feedHandler)
	http.HandleFunc("/feed", feedHandler)
	http.HandleFunc("/live", liveMatchesHandler)
	http.HandleFunc("/analytics", analyticsHandler)
	http.HandleFunc("/community", communityHandler)
	http.HandleFunc("/league", leagueTableHandler)
	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/ws/standings", wsStandingsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("web/templates"))))

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users := db.GetAll()
		json.NewEncoder(w).Encode(users)

	case "POST":
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)
		newUser := db.CreateUser(u)
		taskChan <- "New user: " + newUser.Name
		json.NewEncoder(w).Encode(newUser)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	switch r.Method {
	case "GET":
		u, err := db.GetByID(id)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(u)

	case "PUT":
		var u models.User
		json.NewDecoder(r.Body).Decode(&u)
		err := db.Update(id, u)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(u)

	case "DELETE":
		err := db.Delete(id)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
