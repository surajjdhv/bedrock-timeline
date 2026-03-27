package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anomalyco/bedrock-timeline/internal/parser"
)

const schema = `
CREATE TABLE IF NOT EXISTS player_events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	player_name TEXT NOT NULL,
	xuid TEXT,
	event_type TEXT NOT NULL,
	timestamp DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_player_name ON player_events(player_name);
CREATE INDEX IF NOT EXISTS idx_timestamp ON player_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_event_type ON player_events(event_type);
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_event ON player_events(player_name, event_type, timestamp);
`

type PlayerStore struct {
	db *sql.DB
}

func NewPlayerStore(db *sql.DB) *PlayerStore {
	return &PlayerStore{db: db}
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}

func (s *PlayerStore) SaveEvent(event parser.PlayerEvent) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO player_events (player_name, xuid, event_type, timestamp) VALUES (?, ?, ?, ?)`,
		event.PlayerName, event.XUID, event.EventType, event.Timestamp,
	)
	return err
}

func (s *PlayerStore) GetEvents(ctx context.Context, player, eventType string, limit int) ([]byte, error) {
	query := `SELECT player_name, xuid, event_type, timestamp FROM player_events ORDER BY timestamp DESC LIMIT ?`
	args := []interface{}{limit}

	if player != "" {
		query = `SELECT player_name, xuid, event_type, timestamp FROM player_events WHERE player_name = ? ORDER BY timestamp DESC LIMIT ?`
		args = []interface{}{player, limit}
	}
	if eventType != "" && player != "" {
		query = `SELECT player_name, xuid, event_type, timestamp FROM player_events WHERE player_name = ? AND event_type = ? ORDER BY timestamp DESC LIMIT ?`
		args = []interface{}{player, eventType, limit}
	} else if eventType != "" {
		query = `SELECT player_name, xuid, event_type, timestamp FROM player_events WHERE event_type = ? ORDER BY timestamp DESC LIMIT ?`
		args = []interface{}{eventType, limit}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var playerName, xuid, eventType string
		var timestamp time.Time
		if err := rows.Scan(&playerName, &xuid, &eventType, &timestamp); err != nil {
			return nil, err
		}
		events = append(events, map[string]interface{}{
			"player_name": playerName,
			"xuid":        xuid,
			"event_type":  eventType,
			"timestamp":   timestamp.Format(time.RFC3339),
		})
	}

	if events == nil {
		events = []map[string]interface{}{}
	}

	return json.Marshal(events)
}

func (s *PlayerStore) GetPlayers(ctx context.Context) ([]byte, error) {
	query := `
		WITH event_pairs AS (
			SELECT 
				player_name,
				event_type,
				timestamp,
				LEAD(timestamp) OVER (PARTITION BY player_name ORDER BY timestamp) as next_timestamp,
				LEAD(event_type) OVER (PARTITION BY player_name ORDER BY timestamp) as next_event_type
			FROM player_events
		),
		playtime AS (
			SELECT 
				player_name,
				SUM(
					CASE 
						WHEN event_type = 'join' AND next_event_type = 'leave' AND next_timestamp IS NOT NULL
						THEN CAST((julianday(next_timestamp) - julianday(timestamp)) * 86400 AS INTEGER)
						ELSE 0 
					END
				) as total_seconds
			FROM event_pairs
			WHERE event_type = 'join'
			GROUP BY player_name
		)
		SELECT 
			pe.player_name,
			MAX(CASE WHEN pe.event_type = 'join' THEN pe.timestamp END) as last_join,
			MAX(CASE WHEN pe.event_type = 'leave' THEN pe.timestamp END) as last_leave,
			COALESCE(pt.total_seconds, 0) as total_playtime
		FROM player_events pe
		LEFT JOIN playtime pt ON pe.player_name = pt.player_name
		GROUP BY pe.player_name
		ORDER BY total_playtime DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []map[string]interface{}
	for rows.Next() {
		var playerName string
		var lastJoin, lastLeave sql.NullString
		var totalPlaytime int
		if err := rows.Scan(&playerName, &lastJoin, &lastLeave, &totalPlaytime); err != nil {
			return nil, err
		}
		player := map[string]interface{}{
			"name":           playerName,
			"total_playtime": totalPlaytime,
		}
		if lastJoin.Valid {
			player["last_join"] = lastJoin.String
		}
		if lastLeave.Valid {
			player["last_leave"] = lastLeave.String
		}
		players = append(players, player)
	}

	return json.Marshal(players)
}

func (s *PlayerStore) GetStats(ctx context.Context) ([]byte, error) {
	query := `
		SELECT 
			COUNT(DISTINCT player_name) as total_players,
			COUNT(CASE WHEN event_type = 'join' THEN 1 END) as total_joins,
			COUNT(CASE WHEN event_type = 'leave' THEN 1 END) as total_leaves,
			MAX(timestamp) as last_event
		FROM player_events`

	var totalPlayers, totalJoins, totalLeaves int
	var lastEvent sql.NullString
	err := s.db.QueryRowContext(ctx, query).Scan(&totalPlayers, &totalJoins, &totalLeaves, &lastEvent)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_players": totalPlayers,
		"total_joins":   totalJoins,
		"total_leaves":  totalLeaves,
	}
	if lastEvent.Valid {
		stats["last_event"] = lastEvent.String
	}

	return json.Marshal(stats)
}

type DayPlaytime struct {
	Date     string `json:"date"`
	Duration int    `json:"duration_seconds"`
	Hours    string `json:"hours"`
}

func (s *PlayerStore) GetPlayerPlaytime(ctx context.Context, playerName string) ([]byte, error) {
	query := `
		WITH sessions AS (
			SELECT 
				player_name,
				timestamp as join_time,
				LEAD(timestamp) OVER (PARTITION BY player_name ORDER BY timestamp) as leave_time,
				LEAD(event_type) OVER (PARTITION BY player_name ORDER BY timestamp) as next_event
			FROM player_events
			WHERE player_name = ? AND event_type = 'join'
		)
		SELECT 
			DATE(join_time) as play_date,
			SUM(
				CASE 
					WHEN leave_time IS NOT NULL AND next_event = 'leave' 
					THEN CAST((julianday(leave_time) - julianday(join_time)) * 86400 AS INTEGER)
					ELSE 0 
				END
			) as total_seconds
		FROM sessions
		WHERE leave_time IS NOT NULL AND next_event = 'leave'
		GROUP BY DATE(join_time)
		ORDER BY play_date DESC
		LIMIT 365`

	rows, err := s.db.QueryContext(ctx, query, playerName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playtimes []DayPlaytime
	for rows.Next() {
		var date string
		var seconds int
		if err := rows.Scan(&date, &seconds); err != nil {
			return nil, err
		}
		hours := float64(seconds) / 3600
		playtimes = append(playtimes, DayPlaytime{
			Date:     date,
			Duration: seconds,
			Hours:    formatDuration(hours),
		})
	}

	if playtimes == nil {
		playtimes = []DayPlaytime{}
	}

	return json.Marshal(playtimes)
}

func (s *PlayerStore) GetAllPlaytime(ctx context.Context) ([]byte, error) {
	query := `
		WITH sessions AS (
			SELECT 
				player_name,
				timestamp as join_time,
				LEAD(timestamp) OVER (PARTITION BY player_name ORDER BY timestamp) as leave_time,
				LEAD(event_type) OVER (PARTITION BY player_name ORDER BY timestamp) as next_event
			FROM player_events
			WHERE event_type = 'join'
		)
		SELECT 
			player_name,
			DATE(join_time) as play_date,
			SUM(
				CASE 
					WHEN leave_time IS NOT NULL AND next_event = 'leave' 
					THEN CAST((julianday(leave_time) - julianday(join_time)) * 86400 AS INTEGER)
					ELSE 0 
				END
			) as total_seconds
		FROM sessions
		WHERE leave_time IS NOT NULL AND next_event = 'leave'
		GROUP BY player_name, DATE(join_time)
		ORDER BY play_date DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type PlayerDayPlaytime struct {
		PlayerName string `json:"player_name"`
		Date       string `json:"date"`
		Duration   int    `json:"duration_seconds"`
		Hours      string `json:"hours"`
	}

	var playtimes []PlayerDayPlaytime
	for rows.Next() {
		var playerName, date string
		var seconds int
		if err := rows.Scan(&playerName, &date, &seconds); err != nil {
			return nil, err
		}
		hours := float64(seconds) / 3600
		playtimes = append(playtimes, PlayerDayPlaytime{
			PlayerName: playerName,
			Date:       date,
			Duration:   seconds,
			Hours:      formatDuration(hours),
		})
	}

	if playtimes == nil {
		playtimes = []PlayerDayPlaytime{}
	}

	return json.Marshal(playtimes)
}

func formatDuration(hours float64) string {
	if hours < 1 {
		return fmt.Sprintf("%dm", int(hours*60))
	}
	return fmt.Sprintf("%.1fh", hours)
}

func (s *PlayerStore) GetPlayerWeekEvents(ctx context.Context, playerName string, startDate string) ([]byte, error) {
	query := `
		SELECT player_name, xuid, event_type, timestamp
		FROM player_events
		WHERE player_name = ? 
		  AND timestamp >= DATE(?, 'start of day')
		  AND timestamp < DATE(?, '+7 days', 'start of day')
		ORDER BY timestamp ASC`

	rows, err := s.db.QueryContext(ctx, query, playerName, startDate, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var playerName, xuid, eventType string
		var timestamp time.Time
		if err := rows.Scan(&playerName, &xuid, &eventType, &timestamp); err != nil {
			return nil, err
		}
		events = append(events, map[string]interface{}{
			"player_name": playerName,
			"xuid":        xuid,
			"event_type":  eventType,
			"timestamp":   timestamp.Format(time.RFC3339),
		})
	}

	return json.Marshal(events)
}

func (s *PlayerStore) GetPlayerSessions(ctx context.Context, playerName string, startDate string) ([]byte, error) {
	query := `
		WITH events AS (
			SELECT 
				player_name,
				event_type,
				timestamp,
				ROW_NUMBER() OVER (ORDER BY timestamp) as rn
			FROM player_events
			WHERE player_name = ?
			  AND timestamp >= DATE(?, 'start of day')
			  AND timestamp < DATE(?, '+7 days', 'start of day')
		),
		sessions AS (
			SELECT 
				j.timestamp as start_time,
				l.timestamp as end_time,
				CASE 
					WHEN l.timestamp IS NOT NULL 
					THEN CAST((julianday(l.timestamp) - julianday(j.timestamp)) * 86400 AS INTEGER)
					ELSE NULL 
				END as duration_seconds
			FROM events j
			LEFT JOIN events l ON l.rn = j.rn + 1 AND l.event_type = 'leave'
			WHERE j.event_type = 'join'
		)
		SELECT 
			start_time,
			end_time,
			duration_seconds,
			DATE(start_time) as session_date
		FROM sessions
		WHERE start_time IS NOT NULL
		ORDER BY start_time ASC`

	rows, err := s.db.QueryContext(ctx, query, playerName, startDate, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type Session struct {
		StartTime string  `json:"start_time"`
		EndTime   *string `json:"end_time"`
		Duration  *int    `json:"duration_seconds"`
		Date      string  `json:"date"`
	}

	var sessions []Session
	for rows.Next() {
		var startTime sql.NullString
		var endTime sql.NullString
		var duration sql.NullInt64
		var date string
		if err := rows.Scan(&startTime, &endTime, &duration, &date); err != nil {
			return nil, err
		}

		session := Session{
			Date: date,
		}
		if startTime.Valid {
			session.StartTime = startTime.String
		}
		if endTime.Valid {
			t := endTime.String
			session.EndTime = &t
		}
		if duration.Valid {
			d := int(duration.Int64)
			session.Duration = &d
		}
		sessions = append(sessions, session)
	}

	if sessions == nil {
		sessions = []Session{}
	}

	return json.Marshal(sessions)
}

func (s *PlayerStore) GetAllPlayerSessions(ctx context.Context, startDate string) ([]byte, error) {
	query := `
		WITH event_pairs AS (
			SELECT 
				player_name,
				event_type,
				timestamp,
				LEAD(timestamp) OVER (PARTITION BY player_name ORDER BY timestamp) as next_timestamp,
				LEAD(event_type) OVER (PARTITION BY player_name ORDER BY timestamp) as next_event_type
			FROM player_events
			WHERE timestamp >= DATE(?, 'start of day')
			  AND timestamp < DATE(?, '+7 days', 'start of day')
		)
		SELECT 
			player_name,
			timestamp as start_time,
			next_timestamp as end_time,
			CASE 
				WHEN next_event_type = 'leave' AND next_timestamp IS NOT NULL
				THEN CAST((julianday(next_timestamp) - julianday(timestamp)) * 86400 AS INTEGER)
				ELSE NULL 
			END as duration_seconds,
			DATE(timestamp) as session_date
		FROM event_pairs
		WHERE event_type = 'join'
		ORDER BY start_time ASC`

	rows, err := s.db.QueryContext(ctx, query, startDate, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type Session struct {
		PlayerName      string  `json:"player_name"`
		StartTime       string  `json:"start_time"`
		EndTime         *string `json:"end_time"`
		DurationSeconds *int    `json:"duration_seconds"`
		Date            string  `json:"date"`
	}

	var sessions []Session
	for rows.Next() {
		var playerName string
		var startTime sql.NullString
		var endTime sql.NullString
		var duration sql.NullInt64
		var date string
		if err := rows.Scan(&playerName, &startTime, &endTime, &duration, &date); err != nil {
			return nil, err
		}

		session := Session{
			PlayerName: playerName,
			Date:       date,
		}
		if startTime.Valid {
			session.StartTime = startTime.String
		}
		if endTime.Valid {
			t := endTime.String
			session.EndTime = &t
		}
		if duration.Valid {
			d := int(duration.Int64)
			session.DurationSeconds = &d
		}
		sessions = append(sessions, session)
	}

	if sessions == nil {
		sessions = []Session{}
	}

	return json.Marshal(sessions)
}
