package balloon

import (
	"strings"

	"ponysay-go/pkg/color"
)

// Balloon holds border and link characters for speech/thought bubbles.
type Balloon struct {
	Link       string
	LinkMirror string
	LinkCross  string

	WW, EE           string
	NW, NNW, N, NNE, NE []string
	NEE, E, SEE      string
	SE, SSE, S, SSW, SW []string
	SWW, W, NWW      string

	MinWidth  int
	MinHeight int
}

// ParseBalloon parses a balloon style file string.
func ParseBalloon(content string, isThink bool) *Balloon {
	if content == "" {
		if isThink {
			return &Balloon{
				Link: "o", LinkMirror: "o", LinkCross: "o",
				WW: "( ", EE: " )",
				NW: []string{" _"}, NNW: []string{"_"}, N: []string{"_"}, NNE: []string{"_"}, NE: []string{"_ "},
				NEE: " )", E: " )", SEE: " )",
				SE: []string{"- "}, SSE: []string{"-"}, S: []string{"-"}, SSW: []string{"-"}, SW: []string{" -"},
				SWW: "( ", W: "( ", NWW: "( ",
				MinWidth: 4, MinHeight: 2,
			}
		}
		return &Balloon{
			Link: "\\", LinkMirror: "/", LinkCross: "X",
			WW: "< ", EE: " >",
			NW: []string{" _"}, NNW: []string{"_"}, N: []string{"_"}, NNE: []string{"_"}, NE: []string{"_ "},
			NEE: " \\", E: " |", SEE: " /",
			SE: []string{"- "}, SSE: []string{"-"}, S: []string{"-"}, SSW: []string{"-"}, SW: []string{" -"},
			SWW: "\\ ", W: "| ", NWW: "/ ",
			MinWidth: 4, MinHeight: 2,
		}
	}

	m := make(map[string][]string)
	keys := []string{"\\", "/", "X", "ww", "ee", "nw", "nnw", "n", "nne", "ne", "nee", "e", "see", "se", "sse", "s", "ssw", "sw", "sww", "w", "nww"}
	for _, k := range keys {
		m[k] = []string{}
	}

	lines := strings.Split(content, "\n")
	var lastKey string
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		if line[0] == ':' {
			if lastKey != "" {
				m[lastKey] = append(m[lastKey], line[1:])
			}
		} else {
			idx := strings.Index(line, ":")
			if idx > 0 {
				lastKey = line[:idx]
				m[lastKey] = append(m[lastKey], line[idx+1:])
			}
		}
	}

	getSingle := func(key, def string) string {
		if val, ok := m[key]; ok && len(val) > 0 {
			return val[0]
		}
		return def
	}

	getList := func(key string, def []string) []string {
		if val, ok := m[key]; ok && len(val) > 0 {
			return val
		}
		return def
	}

	b := &Balloon{
		Link:       getSingle("\\", "\\"),
		LinkMirror: getSingle("/", "/"),
		LinkCross:  getSingle("X", "X"),
		WW:         getSingle("ww", "< "),
		EE:         getSingle("ee", " >"),
		NW:         getList("nw", []string{" _"}),
		NNW:        getList("nnw", []string{"_"}),
		N:          getList("n", []string{"_"}),
		NNE:        getList("nne", []string{"_"}),
		NE:         getList("ne", []string{"_ "}),
		NEE:        getSingle("nee", " \\"),
		E:          getSingle("e", " |"),
		SEE:        getSingle("see", " /"),
		SE:         getList("se", []string{"- "}),
		SSE:        getList("sse", []string{"-"}),
		S:          getList("s", []string{"-"}),
		SSW:        getList("ssw", []string{"-"}),
		SW:         getList("sw", []string{" - text"}),
		SWW:        getSingle("sww", "\\ "),
		W:          getSingle("w", "| "),
		NWW:        getSingle("nww", "/ "),
	}

	b.MinWidth = color.DisplayWidth(b.WW) + color.DisplayWidth(b.EE)
	b.MinHeight = len(b.N) + len(b.S)

	return b
}

// WrapText wraps message lines to specified column width while respecting ANSI color sequences.
func WrapText(msg string, wrapWidth int) []string {
	if wrapWidth <= 0 {
		return strings.Split(msg, "\n")
	}

	rawLines := strings.Split(msg, "\n")
	var wrapped []string

	for _, line := range rawLines {
		if color.DisplayWidth(line) <= wrapWidth {
			wrapped = append(wrapped, line)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		var current string
		for _, w := range words {
			if current == "" {
				current = w
			} else if color.DisplayWidth(current+" "+w) <= wrapWidth {
				current += " " + w
			} else {
				wrapped = append(wrapped, current)
				current = w
			}
		}
		if current != "" {
			wrapped = append(wrapped, current)
		}
	}

	return wrapped
}

// FormatBalloon draws the bubble around the message lines.
func (b *Balloon) FormatBalloon(lines []string, minWidth, minHeight int, balloonColor string) []string {
	maxMsgWidth := 0
	for _, l := range lines {
		w := color.DisplayWidth(l)
		if w > maxMsgWidth {
			maxMsgWidth = w
		}
	}

	contentWidth := maxMsgWidth
	if contentWidth+b.MinWidth < minWidth {
		contentWidth = minWidth - b.MinWidth
	}

	w := b.MinWidth + contentWidth

	var ws, es map[int]string
	numLines := len(lines)
	if numLines > 1 {
		ws = map[int]string{0: b.NWW, numLines - 1: b.SWW}
		es = map[int]string{0: b.NEE, numLines - 1: b.SEE}
		for j := 1; j < numLines-1; j++ {
			ws[j] = b.W
			es[j] = b.E
		}
	} else {
		ws = map[int]string{0: b.WW}
		es = map[int]string{0: b.EE}
	}

	var result []string

	// Top border
	for j := 0; j < len(b.N); j++ {
		nwStr := b.NW[j]
		neStr := b.NE[j]
		nnwStr := b.NNW[j]
		nneStr := b.NNE[j]
		nStr := b.N[j]

		outer := color.DisplayWidth(nwStr) + color.DisplayWidth(neStr)
		inner := color.DisplayWidth(nnwStr) + color.DisplayWidth(nneStr)

		var line string
		if outer+inner <= w {
			repeatCount := w - outer - inner
			if repeatCount < 0 {
				repeatCount = 0
			}
			line = nwStr + nnwStr + strings.Repeat(nStr, repeatCount) + nneStr + neStr
		} else {
			repeatCount := w - outer
			if repeatCount < 0 {
				repeatCount = 0
			}
			line = nwStr + strings.Repeat(nStr, repeatCount) + neStr
		}
		result = append(result, color.ApplyColor(line, balloonColor))
	}

	// Message body
	for j, msgLine := range lines {
		lWidth := color.DisplayWidth(msgLine)
		padding := contentWidth - lWidth
		if padding < 0 {
			padding = 0
		}

		leftEdge := color.ApplyColor(ws[j], balloonColor)
		rightEdge := color.ApplyColor(es[j], balloonColor)

		result = append(result, leftEdge+msgLine+strings.Repeat(" ", padding)+rightEdge)
	}

	// Bottom border
	for j := 0; j < len(b.S); j++ {
		swStr := b.SW[j]
		seStr := b.SE[j]
		sswStr := b.SSW[j]
		sseStr := b.SSE[j]
		sStr := b.S[j]

		outer := color.DisplayWidth(swStr) + color.DisplayWidth(seStr)
		inner := color.DisplayWidth(sswStr) + color.DisplayWidth(sseStr)

		var line string
		if outer+inner <= w {
			repeatCount := w - outer - inner
			if repeatCount < 0 {
				repeatCount = 0
			}
			line = swStr + sswStr + strings.Repeat(sStr, repeatCount) + sseStr + seStr
		} else {
			repeatCount := w - outer
			if repeatCount < 0 {
				repeatCount = 0
			}
			line = swStr + strings.Repeat(sStr, repeatCount) + seStr
		}
		result = append(result, color.ApplyColor(line, balloonColor))
	}

	return result
}
