package openai

import (
	"testing"
	"time"
)

func TestStore_AppendAndContext(t *testing.T) {
	s := NewStore(time.Hour, 100)

	s.Append(1, Message{Role: RoleUser, Content: "a"})
	s.Append(1, Message{Role: RoleAssistant, Content: "b"})

	got := s.ContextMessages(1, 10)
	if len(got) != 2 || got[0].Content != "a" || got[1].Content != "b" {
		t.Fatalf("expected [a b], got %+v", got)
	}
}

func TestStore_FIFOEvict(t *testing.T) {
	s := NewStore(time.Hour, 3)

	for i := 0; i < 5; i++ {
		s.Append(1, Message{Content: string(rune('a' + i))})
	}

	got := s.All(1)
	if len(got) != 3 {
		t.Fatalf("expected 3 stored messages, got %d", len(got))
	}
	if got[0].Content != "c" || got[2].Content != "e" {
		t.Fatalf("expected FIFO evict to keep [c d e], got %+v", got)
	}
}

func TestStore_TTLExpiry(t *testing.T) {
	base := time.Now()
	timeNow = func() time.Time { return base }
	defer func() { timeNow = time.Now }()

	s := NewStore(time.Minute, 100)

	s.Append(1, Message{Content: "a"})

	timeNow = func() time.Time { return base.Add(2 * time.Minute) }

	if got := s.ContextMessages(1, 10); got != nil {
		t.Fatalf("expected nil context after TTL expiry, got %+v", got)
	}
	if got := s.All(1); len(got) != 0 {
		t.Fatalf("expected empty history after TTL expiry, got %+v", got)
	}

	// Appending after expiry should start fresh, not carry over old messages.
	s.Append(1, Message{Content: "b"})
	got := s.All(1)
	if len(got) != 1 || got[0].Content != "b" {
		t.Fatalf("expected fresh session with [b], got %+v", got)
	}
}

func TestStore_Clear(t *testing.T) {
	s := NewStore(time.Hour, 100)

	s.Append(1, Message{Content: "a"})
	s.Clear(1)

	if got := s.All(1); len(got) != 0 {
		t.Fatalf("expected empty history after Clear, got %+v", got)
	}
}
