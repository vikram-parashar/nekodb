package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func humanizeDuration(d time.Duration) string {
	seconds := int(d.Seconds())

	days := seconds / (24 * 3600)
	seconds %= 24 * 3600
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

func (s *Server) clearExpiryRoutine(){
	ticker:= time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			log.Info("Clearing expired keys...")
			s.kvMu.Lock()
			for key, expiration := range s.kv.Expirations {
				if time.Now().After(expiration) {
					delete(s.kv.Strings, key)
					delete(s.kv.Expirations, key)
				}
			}
			s.kvMu.Unlock()
		}
	}
}
