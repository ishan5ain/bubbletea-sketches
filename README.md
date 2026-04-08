# bubbletea-sketches

Small Bubble Tea v2 sketches for exploring terminal UI components and interaction patterns in Go.

## Getting started

Run the default sketch:

```bash
go run ./cmd/sketchbook
```

Run a specific sketch:

```bash
go run ./cmd/sketchbook hello-world
```

## Sketches

- `hello-world`: minimal Bubble Tea program that renders a greeting and exits on `q` or `ctrl+c`
- `styled-hello-world`: Lip Gloss v2 style explorer with paged examples for color, spacing, borders, alignment, and more

## Adding a sketch

1. Add a new model constructor to `internal/sketches`.
2. Register it in the sketch registry.
3. Run it through `go run ./cmd/sketchbook <sketch-name>`.
