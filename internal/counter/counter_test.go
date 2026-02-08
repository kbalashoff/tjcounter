package counter

import "testing"

func TestCounterLifecycle(t *testing.T) {
	c := New()
	if got := c.Get(); got != 0 {
		t.Fatalf("initial value = %d, want 0", got)
	}

	if got := c.Increment(); got != 1 {
		t.Fatalf("increment value = %d, want 1", got)
	}

	if got := c.Decrement(); got != 0 {
		t.Fatalf("decrement value = %d, want 0", got)
	}

	c.Increment()
	c.Increment()
	if got := c.Reset(); got != 0 {
		t.Fatalf("reset value = %d, want 0", got)
	}
}
