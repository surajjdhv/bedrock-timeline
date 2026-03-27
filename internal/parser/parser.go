package parser

import (
	"regexp"
	"strings"
	"time"
)

type PlayerEvent struct {
	PlayerName string    `json:"player_name"`
	XUID       string    `json:"xuid"`
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	RawLine    string    `json:"raw_line,omitempty"`
}

var (
	connectPattern = regexp.MustCompile(`(?i)player\s+connected:\s*(.+?)\s*,\s*xuid:\s*([0-9]+)`)
	disconnectPattern = regexp.MustCompile(`(?i)player\s+disconnected:\s*(.+?)\s*,\s*xuid:\s*([0-9]+)`)
	joinPattern = regexp.MustCompile(`(?i)(\S+)\s+joined the game`)
	leavePattern = regexp.MustCompile(`(?i)(\S+)\s+left the game`)
	timestampPattern = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})`)
)

func ParseLine(line string) (PlayerEvent, bool) {
	var event PlayerEvent

	loc := time.Now().Location()
	event.Timestamp = time.Now().In(loc)

	if matches := timestampPattern.FindStringSubmatch(line); len(matches) > 1 {
		if t, err := time.ParseInLocation("2006-01-02T15:04:05", matches[1], loc); err == nil {
			event.Timestamp = t
		}
	}

	if matches := connectPattern.FindStringSubmatch(line); len(matches) == 3 {
		event.PlayerName = strings.TrimSpace(matches[1])
		event.XUID = matches[2]
		event.EventType = "join"
		return event, true
	}

	if matches := disconnectPattern.FindStringSubmatch(line); len(matches) == 3 {
		event.PlayerName = strings.TrimSpace(matches[1])
		event.XUID = matches[2]
		event.EventType = "leave"
		return event, true
	}

	if matches := joinPattern.FindStringSubmatch(line); len(matches) == 2 {
		event.PlayerName = strings.TrimSpace(matches[1])
		event.EventType = "join"
		return event, true
	}

	if matches := leavePattern.FindStringSubmatch(line); len(matches) == 2 {
		event.PlayerName = strings.TrimSpace(matches[1])
		event.EventType = "leave"
		return event, true
	}

	return event, false
}