# Z&A United API Documentation

## Public Endpoints

### GET /api/matchtracker
- Returns real-time match data (score, commentary, player stats)

### GET /api/table
- Returns the current league table

### GET /api/teams
- Returns all EPL teams

### GET /api/players?teamId=
- Returns players (optionally filtered by team)

### GET /api/matches
- Returns all matches

### GET /api/threads
- Returns all match threads

### POST /api/threads/comment
- Adds a comment to a match thread
- Body: `{ "threadId": int, "user": string, "message": string }`

### GET /api/stats
- Returns top scorers and team standings

### GET /api/historical
- Returns historical EPL season data

## Authenticated Endpoints

### POST /api/auth/register
- Register a new user
- Body: `{ "name": string, "email": string, "password": string }`

### POST /api/auth/login
- Login and receive JWT token
- Body: `{ "email": string, "password": string }`

### GET /api/profile/me
- Returns user profile (requires Bearer token)

### POST /api/profile/favorite
- Set favorite team (requires Bearer token)
- Body: `{ "teamId": int }`

### GET /api/feed
- Returns personalized news for user's favorite team (requires Bearer token)

## Admin Endpoints (require admin role)

### POST /api/admin/teams
- Add or update a team

### POST /api/admin/players
- Add or update a player

### POST /api/admin/matches/:id/result
- Update match result
- Body: `{ "home": int, "away": int, "status": string }`
