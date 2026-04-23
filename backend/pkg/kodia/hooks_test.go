package kodia

import (
	"sync"
	"testing"
)

func TestHookManager(t *testing.T) {
	m := NewHookManager()
	
	t.Run("it dispatches events to listeners", func(t *testing.T) {
		called := false
		m.Listen("test.event", func(data any) {
			called = true
			if data.(string) != "payload" {
				t.Errorf("expected payload, got %v", data)
			}
		})
		
		m.Dispatch("test.event", "payload")
		
		if !called {
			t.Error("callback was not called")
		}
	})
	
	t.Run("it handles multiple listeners", func(t *testing.T) {
		count := 0
		m.Listen("multi.event", func(data any) { count++ })
		m.Listen("multi.event", func(data any) { count++ })
		
		m.Dispatch("multi.event", nil)
		
		if count != 2 {
			t.Errorf("expected 2 calls, got %d", count)
		}
	})
	
	t.Run("it is thread-safe", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 1000
		
		wg.Add(iterations)
		for i := 0; i < iterations; i++ {
			go func() {
				defer wg.Done()
				m.Listen("concurrent.event", func(data any) {})
				m.Dispatch("concurrent.event", nil)
			}()
		}
		
		wg.Wait()
	})
}
