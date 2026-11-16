package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const flowerWidth = 11

// Each flower is drawn in a fixed 11 by 9 cell so different varieties align in
// the garden while keeping their petals, stems, and leaves easy to recognize.
var flowerTemplates = [][]string{
	// Daisy: a round, open flower with a sunny center.
	{
		`    .-.    `,
		` .-(   )-. `,
		`(  (o o)  )`,
		` '-(  ^)-' `,
		`     |     `,
		`   _/|     `,
		`  /  |     `,
		`     |     `,
		`    / \    `,
	},
	// Tulip: a cup shaped bloom with pointed petals.
	{
		`    .-.    `,
		`   /   \   `,
		`  /     \  `,
		`  \     /  `,
		`   \___/   `,
		`     |     `,
		`   _/|     `,
		`  /  |     `,
		`     |     `,
	},
	// Rose: layered petals curling around a spiral center.
	{
		`   .---.   `,
		` .' .-. '. `,
		`/  ( @ )  \`,
		`|   '-'   |`,
		` \       / `,
		`  '.___.'  `,
		`     |     `,
		`   _/|\_   `,
		`     |     `,
	},
	// Sunflower: broad petals surrounding a dark seed head.
	{
		` \  |  /  `,
		`  \ | /   `,
		`--(###)-- `,
		`  / | \   `,
		` /  |  \  `,
		`    |     `,
		`  _/|     `,
		` /  |     `,
		`    |     `,
	},
	// Lily: open petals and long stamens.
	{
		`  _  |  _ `,
		` ( \ | / )`,
		`  \ \|/ / `,
		` ---(o)---`,
		`  / /|\ \ `,
		`    /|     `,
		`   / |     `,
		`     |     `,
		`    / \    `,
	},
	// Poppy: softly folded petals around a dark center.
	{
		`   .---.   `,
		` .'     '. `,
		`/  .---.  \`,
		`|  (###)  |`,
		` \  '---' / `,
		`  '.___.'  `,
		`     |     `,
		`   _/|\_   `,
		`     |     `,
	},
}

var petalColors = []lipgloss.Color{
	lipgloss.Color("#FFF4D6"), // daisy cream
	lipgloss.Color("#FF5C8A"), // tulip pink
	lipgloss.Color("#E84A5F"), // rose red
	lipgloss.Color("#FFD447"), // sunflower gold
	lipgloss.Color("#FFF0F5"), // lily white
	lipgloss.Color("#FF765E"), // poppy coral
}

var stemColor = lipgloss.Color("#5DBB63")

func renderFlower(t thought, growFrame int) []string {
	tmpl := flowerTemplates[t.FlowerType%len(flowerTemplates)]
	color := petalColors[t.ColorIndex%len(petalColors)]

	petalStyle := lipgloss.NewStyle().Foreground(color)
	stemStyle := lipgloss.NewStyle().Foreground(stemColor)
	blank := strings.Repeat(" ", flowerWidth)

	showFrom := 0
	if growFrame >= 0 && growFrame < len(tmpl)-1 {
		showFrom = len(tmpl) - 1 - growFrame
	}

	stemStart := 5

	lines := make([]string, len(tmpl))
	for i, line := range tmpl {
		if growFrame >= 0 && i < showFrom {
			lines[i] = blank
			continue
		}
		if i >= stemStart {
			lines[i] = stemStyle.Render(line)
		} else {
			lines[i] = petalStyle.Render(line)
		}
	}

	return lines
}
