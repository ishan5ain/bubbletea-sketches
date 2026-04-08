package sketches

import "testing"

func TestGetReturnsRegisteredSketch(t *testing.T) {
	factory, ok := Get("hello-world")
	if !ok {
		t.Fatal("expected hello-world to be registered")
	}

	if got := factory().View(); got == "" {
		t.Fatal("registered sketch returned an empty view")
	}
}

func TestNamesIncludesDefaultSketch(t *testing.T) {
	names := Names()
	if len(names) != 1 {
		t.Fatalf("expected 1 sketch, got %d", len(names))
	}

	if names[0] != DefaultName() {
		t.Fatalf("expected %q, got %q", DefaultName(), names[0])
	}
}
