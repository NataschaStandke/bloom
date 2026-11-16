# Bloom

Bloom is a small terminal garden: write down a thought and watch a flower grow for it. Thoughts are saved locally, so the garden is still there the next time you open the app.

## Requirements

- Go 1.26.4 or newer
- A terminal that supports ANSI color

## Run

From the repository directory:

```sh
go run -buildvcs=false .
```

The `-buildvcs=false` flag also lets this run from a directory that has not been initialized as a Git repository. To build the program instead:

```sh
go build -buildvcs=false -o bloom .
./bloom
```

## Controls

| Screen | Key | Action |
| --- | --- | --- |
| Splash | `i` or `Enter` | Open the garden |
| Splash | `q` | Quit |
| Garden | `n`, `i`, or `Enter` | Write a new thought |
| Garden | `x` | Open the wipe confirmation |
| Garden | `q` | Quit |
| Garden | `Esc` | Return to the splash screen |
| Input | `Enter` | Plant the thought and return to the garden |
| Input | `Esc` | Cancel and return to the garden |
| Wipe confirmation | `y` | Clear and save an empty garden |
| Wipe confirmation | Any other key | Cancel |
| Any screen | `Ctrl+C` | Quit |

Each thought is paired with a randomly selected flower and petal color. The garden shows the thought beneath its flower and animates newly planted flowers.

## Saved garden

Thoughts are stored as JSON at:

```text
<user config directory>/thought-garden/thoughts.json
```

The user config directory is provided by Go's `os.UserConfigDir()` and varies by operating system. Use the app's wipe option to clear the saved garden.

## Development

Format Go files and build the project with:

```sh
gofmt -w *.go
go build -buildvcs=false ./...
```

The terminal interface uses [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).
