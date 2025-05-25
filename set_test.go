package expset

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func TestCache(t *testing.T) {
	s := New[string]()
	clk := clock.NewMock()
	s.clock = clk
	s.Start()
	defer s.Stop()
	t.Run(`base`, func(t *testing.T) {
		s.Add("test-key-1", 10*time.Second)
		s.Add("test-key-2", 20*time.Second)
		s.Add("test-key-3", 30*time.Second)
		if !s.Has(`test-key-1`) {
			t.Fatalf("Item missing")
		}
		clk.Add(10 * time.Second)
		if s.Has(`test-key-1`) {
			t.Fatalf("Item not evicted")
		}
		clk.Add(10 * time.Second)
	})
	t.Run(`clear`, func(t *testing.T) {
		s.Clear()
		if s.Len() > 0 {
			t.Fatalf("Set not cleared")
		}
		if s.Has(`test-key-3`) {
			t.Fatalf("Item not cleared")
		}
	})
	t.Run(`duplicate`, func(t *testing.T) {
		s.Add("test-key-1", 10*time.Second)
		s.Add("test-key-1", 20*time.Second)
		s.Add("test-key-2", 30*time.Second)
		if s.Len() > 2 {
			t.Fatalf("Duplicate item inserted")
		}
		clk.Add(10 * time.Second)
		if s.Len() != 2 {
			t.Fatalf("Old TTL not overwritten")
		}
		clk.Add(10 * time.Second)
		if s.Len() != 1 {
			t.Fatalf("New TTL not respected")
		}
		clk.Add(10 * time.Second)
		if s.Len() > 0 {
			t.Fatalf("Item not evicted")
		}
	})
	t.Run(`collision`, func(t *testing.T) {
		s.Clear()
		s.Add("test-key-1", 10*time.Second)
		s.Add("test-key-2", 10*time.Second)
		if s.Len() != 2 {
			t.Fatalf("Incorrect number of items")
		}
		clk.Add(11 * time.Second)
		if s.Len() > 0 {
			t.Fatalf("Items not evicted")
		}
	})
	t.Run(`refresh`, func(t *testing.T) {
		s.Clear()
		s.Add("test-key-1", 10*time.Second)
		if s.Len() != 1 {
			t.Fatalf("Incorrect number of items")
		}
		clk.Add(5 * time.Second)
		s.Refresh("test-key-1")
		clk.Add(5 * time.Second)
		if s.Len() != 1 {
			t.Fatalf("Refresh not respected")
		}
		s.Refresh("test-key-1")
		clk.Add(5 * time.Second)
		if s.Len() != 1 {
			t.Fatalf("Refresh not respected")
		}
		t.Run(`missing`, func(t *testing.T) {
			s.Refresh("test-key-2")
			if s.Len() != 1 {
				t.Fatalf("Wrong number of items")
			}
		})
		t.Run(`collision`, func(t *testing.T) {
			s.Add("test-key-1", 10*time.Second)
			s.Add("test-key-2", 15*time.Second)
			clk.Add(5 * time.Second)
			s.Refresh("test-key-1")
			if s.Len() != 2 {
				t.Fatalf("Wrong number of items")
			}
			clk.Add(10 * time.Second)
			if s.Len() != 1 {
				t.Fatalf("Items not evicted")
			}
			clk.Add(1 * time.Second)
			if s.Len() != 0 {
				t.Fatalf("Items not evicted")
			}
		})
	})
}
