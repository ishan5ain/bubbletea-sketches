package main

import "testing"

func TestSelectSketchDefaultsToHelloWorld(t *testing.T) {
	model, err := selectSketch(nil)
	if err != nil {
		t.Fatalf("select default sketch: %v", err)
	}

	if got := model.View().Content; got == "" {
		t.Fatal("default sketch returned an empty view")
	}
}

func TestSelectSketchByName(t *testing.T) {
	model, err := selectSketch([]string{"hello-world"})
	if err != nil {
		t.Fatalf("select named sketch: %v", err)
	}

	if got := model.View().Content; got == "" {
		t.Fatal("named sketch returned an empty view")
	}
}

func TestSelectSketchStyledHelloWorld(t *testing.T) {
	model, err := selectSketch([]string{"styled-hello-world"})
	if err != nil {
		t.Fatalf("select styled sketch: %v", err)
	}

	if got := model.View().Content; got == "" {
		t.Fatal("styled sketch returned an empty view")
	}
}

func TestSelectSketchDirectorySelector(t *testing.T) {
	model, err := selectSketch([]string{"directory-selector"})
	if err != nil {
		t.Fatalf("select directory-selector sketch: %v", err)
	}

	if got := model.View().Content; got == "" {
		t.Fatal("directory-selector sketch returned an empty view")
	}
}

func TestSelectSketchFlexibleKeyValuePairList(t *testing.T) {
	model, err := selectSketch([]string{"flexible-key-value-pair-list"})
	if err != nil {
		t.Fatalf("select flexible-key-value-pair-list sketch: %v", err)
	}

	if got := model.View().Content; got == "" {
		t.Fatal("flexible-key-value-pair-list sketch returned an empty view")
	}
}

func TestSelectSketchRejectsUnknownSketch(t *testing.T) {
	if _, err := selectSketch([]string{"missing"}); err == nil {
		t.Fatal("expected an error for an unknown sketch")
	}
}
