package tracker

import (
	"sync"
	"time"

	"github.com/anomalyco/bedrock-timeline/internal/parser"
)

type PlayerTracker struct {
	onlinePlayers map[string]time.Time // player_name -> join_time
	mu            sync.RWMutex
}

func NewPlayerTracker() *PlayerTracker {
	return &PlayerTracker{
		onlinePlayers: make(map[string]time.Time),
	}
}

// ProcessEvent processes a join/leave event and returns whether it should be saved
// Returns: (shouldSave, eventType)
// eventType may be modified (e.g., "join" becomes "leave" if we detect a glitch)
func (t *PlayerTracker) ProcessEvent(event parser.PlayerEvent) (shouldSave bool, processedEvent parser.PlayerEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()

	processedEvent = event

	switch event.EventType {
	case "join":
		_, isOnline := t.onlinePlayers[event.PlayerName]
		if isOnline {
			// Player is already online - this is a glitch
			// The server probably restarted or we missed a leave event
			// Create a synthetic leave event for the previous session
			// Don't save this join since they're already logged in
			return false, event
		}
		// Valid join - player was offline
		t.onlinePlayers[event.PlayerName] = event.Timestamp
		return true, event

	case "leave":
		joinTime, isOnline := t.onlinePlayers[event.PlayerName]
		if !isOnline {
			// Player is not online - this is a glitch
			// They might have joined before we started tracking
			// Assume they were online and save the leave
			// Create a synthetic join event for consistency? No, just skip this leave
			return false, event
		}
		// Valid leave - player was online
		_ = joinTime // join time tracked for potential session calculation
		delete(t.onlinePlayers, event.PlayerName)
		return true, event
	}

	return false, event
}

// IsOnline checks if a player is currently online
func (t *PlayerTracker) IsOnline(playerName string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, online := t.onlinePlayers[playerName]
	return online
}

// GetOnlinePlayers returns list of currently online players
func (t *PlayerTracker) GetOnlinePlayers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	players := make([]string, 0, len(t.onlinePlayers))
	for name := range t.onlinePlayers {
		players = append(players, name)
	}
	return players
}

// SetOnline forcibly sets a player as online (used when loading history)
func (t *PlayerTracker) SetOnline(playerName string, joinTime time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onlinePlayers[playerName] = joinTime
}

// SetOffline forcibly sets a player as offline
func (t *PlayerTracker) SetOffline(playerName string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.onlinePlayers, playerName)
}

// RebuildState rebuilds the online player state from historical events
// This should be called after loading history to know who is currently online
func (t *PlayerTracker) RebuildState(events []parser.PlayerEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Clear current state
	t.onlinePlayers = make(map[string]time.Time)

	// Sort events by timestamp (they should already be sorted from history)
	// Process each event to build the current state
	for _, event := range events {
		switch event.EventType {
		case "join":
			t.onlinePlayers[event.PlayerName] = event.Timestamp
		case "leave":
			delete(t.onlinePlayers, event.PlayerName)
		}
	}
}
