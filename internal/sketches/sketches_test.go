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
	if len(names) != 3 {
		t.Fatalf("expected 3 sketches, got %d", len(names))
	}

	want := []string{"directory-selector", DefaultName(), "styled-hello-world"}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("expected %q at index %d, got %q", want[i], i, names[i])
		}
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

func TestGetReturnsDirectorySelector(t *testing.T) {
	factory, ok := Get("directory-selector")
	if !ok {
		t.Fatal("expected directory-selector to be registered")
	}

	if got := factory().View().Content; got == "" {
		t.Fatal("directory-selector returned an empty view")
	}
}
