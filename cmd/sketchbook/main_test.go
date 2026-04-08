package main

import "testing"

func TestSelectSketchDefaultsToHelloWorld(t *testing.T) {
	model, err := selectSketch(nil)
	if err != nil {
		t.Fatalf("select default sketch: %v", err)
	}

	if got := model.View(); got == "" {
		t.Fatal("default sketch returned an empty view")
	}
}

func TestSelectSketchByName(t *testing.T) {
	model, err := selectSketch([]string{"hello-world"})
	if err != nil {
		t.Fatalf("select named sketch: %v", err)
	}

	if got := model.View(); got == "" {
		t.Fatal("named sketch returned an empty view")
	}
}

func TestSelectSketchRejectsUnknownSketch(t *testing.T) {
	if _, err := selectSketch([]string{"missing"}); err == nil {
		t.Fatal("expected an error for an unknown sketch")
	}
}
