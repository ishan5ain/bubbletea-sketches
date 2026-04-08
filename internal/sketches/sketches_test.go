package sketches

import "testing"

func TestGetReturnsRegisteredSketch(t *testing.T) {
	factory, ok := Get("hello-world")
	if !ok {
		t.Fatal("expected hello-world to be registered")
	}

	if got := factory().View().Content; got == "" {
		t.Fatal("registered sketch returned an empty view")
	}
}

func TestNamesIncludesDefaultSketch(t *testing.T) {
	names := Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 sketches, got %d", len(names))
	}

	if names[0] != DefaultName() {
		t.Fatalf("expected %q, got %q", DefaultName(), names[0])
	}

	if names[1] != "styled-hello-world" {
		t.Fatalf("expected %q, got %q", "styled-hello-world", names[1])
	}
}

func TestGetReturnsStyledHelloWorld(t *testing.T) {
	factory, ok := Get("styled-hello-world")
	if !ok {
		t.Fatal("expected styled-hello-world to be registered")
	}

	if got := factory().View().Content; got == "" {
		t.Fatal("styled-hello-world returned an empty view")
	}
}
