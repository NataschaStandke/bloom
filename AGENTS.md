# Repository guide

## Project

This is a small Go terminal UI that turns each submitted thought into a flower in a persistent garden. It uses Bubble Tea for the application loop, Bubbles for text input, and Lip Gloss for terminal styling.

## Structure

- `main.go` loads saved thoughts and starts the Bubble Tea program.
- `model.go` owns the screens, keyboard handling, garden layout, and growth animation.
- `flower.go` defines the flower drawings, colors, and flower renderer.
- `thought.go` defines the persisted thought data and JSON storage.
- Thoughts are stored at `os.UserConfigDir()/thought-garden/thoughts.json`.

## Behavior to preserve

- The splash screen shows the existing garden. `i` or Enter opens the garden; `q` quits.
- In the garden, `n`, `i`, or Enter opens thought input; `x` asks for wipe confirmation; `q` quits; Esc returns to the splash screen.
- In thought input, Enter adds a non-empty thought, saves it, and returns to the garden; Esc cancels and returns to the garden.
- On the wipe confirmation screen, `y` clears and saves an empty garden; any other key cancels.
- The garden displays each thought beneath its flower, keeps the newest flowers visible in a small terminal, and animates new growth.
- Ctrl+C quits from any screen.
- Keep the JSON field names in `thought.go` compatible with existing saved data unless a migration is included.

## Go workflow

- Format changed Go files with `gofmt`.
- Build with `go build -buildvcs=false ./...`. Disabling VCS stamping keeps builds working in checkouts without Git metadata.
- Keep dependencies focused; update `go.mod` and `go.sum` when dependencies change.
- Prefer small changes that preserve the existing Bubble Tea and Lip Gloss patterns.
