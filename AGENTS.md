# Agent Rules

When creating a new sketch in this repository:

1. Create a new directory under `internal/sketches/` for the sketch.
2. Put the sketch implementation and its test file in that directory.
3. Export a constructor from the sketch package.
4. Register the constructor in `internal/sketches/sketches.go`.
5. Keep the user-facing sketch name unchanged unless the task explicitly says otherwise.

Keep new sketch packages focused and local. Do not add new sketch source files directly under `internal/sketches/` at the package root.
