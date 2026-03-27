package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anomalyco/bedrock-timeline/internal/journal"
	"github.com/anomalyco/bedrock-timeline/internal/parser"
	"github.com/anomalyco/bedrock-timeline/internal/store"
	"github.com/anomalyco/bedrock-timeline/internal/ws"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/*
var static embed.FS

func main() {
	dbPath := getEnv("DB_PATH", "data/bedrock.db")
	port := getEnv("PORT", "3000")
	journalUnit := getEnv("JOURNAL_UNIT", "bedrock")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := store.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	playerStore := store.NewPlayerStore(db)
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	http.HandleFunc("/api/events", makeEventsHandler(playerStore))
	http.HandleFunc("/api/players", makePlayersHandler(playerStore))
	http.HandleFunc("/api/online", makeOnlineHandler(playerStore))
	http.HandleFunc("/api/stats", makeStatsHandler(playerStore))
	http.HandleFunc("/api/playtime", makePlaytimeHandler(playerStore))
	http.HandleFunc("/api/sessions", makeSessionsHandler(playerStore))

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("Failed to parse static files: %v", err)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	log.Printf("Starting bedrock-timeline on port %s", port)

	jr := journal.NewReader(journalUnit)

	log.Println("Loading historical journal entries...")
	historicalLines, err := jr.ReadHistory(30)
	if err != nil {
		log.Printf("Warning: Failed to read history: %v", err)
	} else {
		count := 0
		for _, line := range historicalLines {
			event, ok := parser.ParseLine(line)
			if !ok {
				continue
			}
			if err := playerStore.SaveEvent(event); err != nil {
				log.Printf("Failed to save historical event: %v", err)
				continue
			}
			count++
		}
		log.Printf("Loaded %d historical events", count)
	}

	if err := jr.Start(); err != nil {
		log.Printf("Warning: Failed to start journal reader: %v", err)
	}

	go func() {
		for line := range jr.Lines() {
			event, ok := parser.ParseLine(line)
			if !ok {
				continue
			}
			if err := playerStore.SaveEvent(event); err != nil {
				log.Printf("Failed to save event: %v", err)
				continue
			}
			hub.Broadcast(event)
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func makeEventsHandler(store *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		player := r.URL.Query().Get("player")
		eventType := r.URL.Query().Get("type")

		events, err := store.GetEvents(r.Context(), player, eventType, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(events)
	}
}

func makePlayersHandler(playerStore *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		players, err := playerStore.GetPlayers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(players)
	}
}

func makeStatsHandler(playerStore *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stats, err := playerStore.GetStats(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(stats)
	}
}

func makeOnlineHandler(playerStore *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		online, err := playerStore.GetOnlinePlayers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(online)
	}
}

func makePlaytimeHandler(playerStore *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		player := r.URL.Query().Get("player")
		if player != "" {
			playtime, err := playerStore.GetPlayerPlaytime(r.Context(), player)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(playtime)
			return
		}

		playtime, err := playerStore.GetAllPlaytime(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(playtime)
	}
}

func makeSessionsHandler(playerStore *store.PlayerStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		player := r.URL.Query().Get("player")
		startDate := r.URL.Query().Get("start")
		if startDate == "" {
			startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		}

		var sessions []byte
		var err error

		if player == "" {
			sessions, err = playerStore.GetAllPlayerSessions(r.Context(), startDate)
		} else {
			sessions, err = playerStore.GetPlayerSessions(r.Context(), player, startDate)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(sessions)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
