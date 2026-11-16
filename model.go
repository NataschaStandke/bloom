package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	cellWidth       = flowerWidth + 2
	flowerLineCount = 9
	flowerRowHeight = flowerLineCount + 3
)

type viewState int

const (
	splashView viewState = iota
	gardenScreen
	inputView
	wipeView
)

type growTickMsg struct{}

func growTickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(time.Time) tea.Msg { return growTickMsg{} })
}

type model struct {
	thoughts   []thought
	input      textinput.Model
	state      viewState
	width      int
	height     int
	growingIdx int
	growFrame  int
	err        string
}

func newModel(thoughts []thought) model {
	in := textinput.New()
	in.Placeholder = "A thought to grow your garden…"
	in.CharLimit = 240
	in.Width = 50
	return model{
		thoughts:   thoughts,
		input:      in,
		state:      splashView,
		growingIdx: -1,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case growTickMsg:
		if m.growingIdx >= 0 {
			m.growFrame++
			if m.growFrame > flowerLineCount-1 {
				m.growingIdx = -1
				return m, nil
			}
			return m, growTickCmd()
		}
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.state {
		case splashView:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "i", "enter":
				m.state = gardenScreen
				return m, nil
			default:
				return m, nil
			}
		case gardenScreen:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "n", "i", "enter":
				m.state = inputView
				m.input.Focus()
				return m, textinput.Blink
			case "x":
				if len(m.thoughts) > 0 {
					m.state = wipeView
				}
				return m, nil
			case "esc":
				m.state = splashView
				return m, nil
			}
		case wipeView:
			switch msg.String() {
			case "y":
				m.thoughts = []thought{}
				if err := saveThoughts(m.thoughts); err != nil {
					m.err = err.Error()
				}
				m.growingIdx = -1
				m.state = gardenScreen
				return m, nil
			default:
				m.state = gardenScreen
				return m, nil
			}
		case inputView:
			switch msg.String() {
			case "esc":
				m.state = gardenScreen
				m.input.Blur()
				m.input.Reset()
				return m, nil
			case "enter":
				text := strings.TrimSpace(m.input.Value())
				if text == "" {
					m.state = gardenScreen
					m.input.Blur()
					return m, nil
				}
				m.thoughts = append(m.thoughts, thought{
					Text:       text,
					FlowerType: rand.Intn(len(flowerTemplates)),
					ColorIndex: rand.Intn(len(petalColors)),
					CreatedAt:  time.Now(),
				})
				if err := saveThoughts(m.thoughts); err != nil {
					m.err = err.Error()
				} else {
					m.err = ""
				}
				m.growingIdx = len(m.thoughts) - 1
				m.growFrame = 0
				m.state = gardenScreen
				m.input.Blur()
				m.input.Reset()
				return m, growTickCmd()
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Growing a garden…"
	}
	switch m.state {
	case splashView:
		return m.renderSplash()
	case gardenScreen:
		return m.renderGarden()
	case inputView:
		return m.renderInput()
	case wipeView:
		return m.renderWipe()
	default:
		return ""
	}
}

func (m model) renderSplash() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF69B4")).Render("✿  thought garden  ✿")
	content := title + "\nA little flower for every thought.\n\n" + m.gardenView() + "\n\ni  enter garden     q  quit"
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) renderGarden() string {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true).Width(m.width).Align(lipgloss.Center)
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Width(m.width).Align(lipgloss.Center)
	var b strings.Builder
	b.WriteString(titleStyle.Render("bloom garden") + "\n\n")
	if m.err != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render("Could not save: "+m.err) + "\n")
	}
	b.WriteString(m.gardenView() + "\n\n")
	b.WriteString(footerStyle.Render("[n] New thought  [x] Wipe garden  [q] Quit"))
	return b.String()
}

func (m model) renderInput() string {
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true).Render("Plant a thought")
	return title + "\n" + m.input.View() + "\nPress Enter to plant · Esc to return\n\n" + m.gardenView()
}

func (m model) renderWipe() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		style.Render("Wipe the whole garden?")+"\n\n[y] Yes, clear it     any other key to cancel")
}

func (m model) gardenView() string {
	if len(m.thoughts) == 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Your garden is waiting for its first thought.")
	}
	cols := m.width / cellWidth
	if cols < 1 {
		cols = 1
	}
	visibleRows := (m.height - 8) / flowerRowHeight
	if visibleRows < 1 {
		visibleRows = 1
	}
	first := len(m.thoughts) - visibleRows*cols
	if first < 0 {
		first = 0
	}
	var rows []string
	for start := first; start < len(m.thoughts); start += cols {
		end := min(start+cols, len(m.thoughts))
		cells := make([]string, end-start)
		for i := start; i < end; i++ {
			frame := -1
			if i == m.growingIdx {
				frame = m.growFrame
			}
			flower := strings.Join(renderFlower(m.thoughts[i], frame), "\n")
			caption := strings.Join(captionLines(m.thoughts[i].Text, cellWidth, 2), "\n")
			color := petalColors[m.thoughts[i].ColorIndex%len(petalColors)]
			caption = lipgloss.NewStyle().Foreground(color).Width(cellWidth).Align(lipgloss.Center).Render(caption)
			cells[i-start] = lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center).Render(flower + "\n" + caption)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	if first > 0 {
		rows = append([]string{lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render(fmt.Sprintf("… %d earlier flower(s) …", first))}, rows...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func captionLines(text string, width, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 || width < 1 || maxLines < 1 {
		return nil
	}
	lines := make([]string, 0, maxLines)
	line := ""
	for _, word := range words {
		if lipgloss.Width(word) > width {
			word = truncateWidth(word, width)
		}
		candidate := word
		if line != "" {
			candidate = line + " " + word
		}
		if lipgloss.Width(candidate) <= width {
			line = candidate
			continue
		}
		lines = append(lines, line)
		if len(lines) == maxLines {
			lines[len(lines)-1] = addEllipsis(lines[len(lines)-1], width)
			return lines
		}
		line = word
	}
	if line != "" && len(lines) < maxLines {
		lines = append(lines, line)
	}
	if len(lines) == maxLines && len(strings.Fields(strings.Join(lines, " "))) < len(words) {
		lines[len(lines)-1] = addEllipsis(lines[len(lines)-1], width)
	}
	return lines
}

func truncateWidth(text string, width int) string {
	var b strings.Builder
	for _, r := range text {
		candidate := b.String() + string(r)
		if lipgloss.Width(candidate) > width {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}

func addEllipsis(text string, width int) string {
	if lipgloss.Width(text) < width {
		return text + "…"
	}
	return truncateWidth(text, width-1) + "…"
}
