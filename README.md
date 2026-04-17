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
- `directory-selector`: terminal-style directory prompt with `tab` autocomplete against the real filesystem
- `flexible-key-value-pair-list`: stacked parameter editor with required rows, optional-key autocomplete, and live preview
- `styled-hello-world`: Lip Gloss v2 style explorer with paged examples for color, spacing, borders, alignment, and more

## Adding a sketch

1. Add a new sketch directory under `internal/sketches/` with the model code and tests.
2. Export a constructor from that package.
3. Register the constructor in `internal/sketches/sketches.go`.
4. Run it through `go run ./cmd/sketchbook <sketch-name>`.
