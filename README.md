# EPL Stats (Go + Gin + GORM)

API-first EPL statistics website with authentication, personalization, live league table, player stats, matches, and an admin panel.

## Tech Stack
- Go (Gin)
- GORM (SQLite for dev, PostgreSQL for prod)
- HTML + CSS + JS (templates + static)

## Project Structure
- cmd/server/main.go: Server entry point
- internal/config: Env config
- internal/database: DB connection
- internal/models: GORM models (Users, Teams, Players, PlayerStats, Matches)
- internal/migrations: AutoMigrate + seed
- internal/middleware: JWT auth, admin guard
- internal/services: Business logic
- internal/handlers: REST API handlers
- web/templates + web/static: Frontend

## Environment
- DB_DRIVER=sqlite (or postgres)
- DB_DSN (sqlite default: file:epl.db?cache=shared&_journal_mode=WAL)
- JWT_SECRET=<set a strong secret>
- ADMIN_EMAIL=admin@epl.local

## Setup
1. Ensure Go is installed.
2. From the project folder:
   - go mod tidy
   - go run ./cmd/server
3. Visit http://localhost:8080
   - / renders league table
   - /profile for login and favorite team

## Database
GORM models:
- Users: id, name, email (unique), password_hash, role, favorite_team_id
- Teams: name (unique), short_name, colors, points, matches_played, goal_diff
- Players: name, team_id, position
- PlayerStats: player_id, season, goals, assists, clean_sheets, minutes_played
- Matches: home_team_id, away_team_id, scores, date, stadium, status

AutoMigrate runs at startup and seeds Top-6 teams and sample matches.

## API Endpoints
- POST /api/auth/register {name,email,password}
- POST /api/auth/login {email,password} -> {token,user}
- GET /api/teams
- GET /api/players?teamId=
- GET /api/matches
- GET /api/table
- Auth required:
  - GET /api/profile/me (Bearer token)
  - POST /api/profile/favorite {teamId}
- Admin required:
  - POST /api/admin/teams (Team JSON)
  - POST /api/admin/players (Player JSON)
  - POST /api/admin/matches/:id/result {home,away,status}

## Example Requests
Register:
curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d "{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"password\":\"pass\"}"

Login:
curl -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d "{\"email\":\"alice@example.com\",\"password\":\"pass\"}"

Get Table:
curl http://localhost:8080/api/table

Set Favorite:
curl -X POST http://localhost:8080/api/profile/favorite -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d "{\"teamId\":1}"

Admin Update Match:
curl -X POST http://localhost:8080/api/admin/matches/1/result -H "Authorization: Bearer <ADMIN_TOKEN>" -H "Content-Type: application/json" -d "{\"home\":2,\"away\":1,\"status\":\"finished\"}"

## Main Logic
- Auth: bcrypt password hashing; JWT for stateless sessions; role in claims for admin guard.
- Personalization: favorite team persisted; theme colors applied from team selection.
- League Table: computed from Teams; cached in-memory for 30s; frontend polls every 10s.
- Player Stats: list by team; sortable client-side; stats preloaded.
- Matches: list upcoming and finished; admin can set results; table can be updated accordingly if you update points and GD in Teams.

## Notes
- For PostgreSQL set DB_DRIVER=postgres and DB_DSN to your connection string.
- Ensure JWT_SECRET is changed for production.
