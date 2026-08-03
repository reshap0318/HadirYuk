// Package openai is standalone: it must not import models/gorm/gin, only stdlib
// and its own types, so it stays unit-testable without the rest of the app.
package openai

import (
	"sync"
	"time"
)

// timeNow allows tests to control expiry without real sleeps.
var timeNow = time.Now

type session struct {
	mu           sync.Mutex
	messages     []Message
	lastActivity time.Time
}

// Store is an in-memory, per-user chat session store with TTL expiry and a
// FIFO cap on stored messages.
//
// ponytail: single-process sync.Map, not shared across instances. Move to
// Redis if the app ever runs multi-instance (docker-compose today is single instance).
type Store struct {
	sessions  sync.Map // map[uint]*session
	ttl       time.Duration
	maxStored int
}

// NewStore creates a session store and starts its background TTL sweep.
// The sweep goroutine runs for the lifetime of the process (Store lives as
// long as the app does, same as the DB pool).
func NewStore(ttl time.Duration, maxStored int) *Store {
	s := &Store{ttl: ttl, maxStored: maxStored}
	go s.sweepLoop()
	return s
}

func (s *Store) sweepLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.sweep()
	}
}

func (s *Store) sweep() {
	now := timeNow()
	s.sessions.Range(func(key, value any) bool {
		sess := value.(*session)
		sess.mu.Lock()
		expired := now.Sub(sess.lastActivity) > s.ttl
		sess.mu.Unlock()
		if expired {
			s.sessions.Delete(key)
		}
		return true
	})
}

func (s *Store) getOrCreate(userID uint) *session {
	val, _ := s.sessions.LoadOrStore(userID, &session{lastActivity: timeNow()})
	return val.(*session)
}

// Append adds a message to the user's session, evicting the oldest message
// (FIFO) once maxStored is exceeded. An expired session starts fresh.
func (s *Store) Append(userID uint, msg Message) {
	sess := s.getOrCreate(userID)
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if timeNow().Sub(sess.lastActivity) > s.ttl {
		sess.messages = nil
	}

	sess.messages = append(sess.messages, msg)
	if len(sess.messages) > s.maxStored {
		sess.messages = sess.messages[len(sess.messages)-s.maxStored:]
	}
	sess.lastActivity = timeNow()
}

// ContextMessages returns the last n messages for a session, for use as
// OpenAI context. Returns nil for a missing or expired session.
func (s *Store) ContextMessages(userID uint, n int) []Message {
	val, ok := s.sessions.Load(userID)
	if !ok {
		return nil
	}
	sess := val.(*session)
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if timeNow().Sub(sess.lastActivity) > s.ttl {
		return nil
	}
	if len(sess.messages) <= n {
		return append([]Message(nil), sess.messages...)
	}
	return append([]Message(nil), sess.messages[len(sess.messages)-n:]...)
}

// All returns the full stored history (for scrollback), or an empty slice if
// the session is missing or expired.
func (s *Store) All(userID uint) []Message {
	val, ok := s.sessions.Load(userID)
	if !ok {
		return []Message{}
	}
	sess := val.(*session)
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if timeNow().Sub(sess.lastActivity) > s.ttl {
		return []Message{}
	}
	return append([]Message(nil), sess.messages...)
}

// Clear removes a user's session entirely (used by the reset endpoint).
func (s *Store) Clear(userID uint) {
	s.sessions.Delete(userID)
}
