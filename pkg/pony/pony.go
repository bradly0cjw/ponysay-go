package pony

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"ponysay-go/pkg/color"
)

var balloonTagRegex = regexp.MustCompile(`\$balloon[0-9a-zA-Z,]*\$`)

// Pony represents a parsed pony artwork file.
type Pony struct {
	Name          string
	Metadata      map[string]string
	RawInfo       string
	HasInfo       bool
	BalloonTop    int
	BalloonBottom int
	BodyLines     []string
}

// FormatInfo formats metadata text according to ponysay's formatting rules.
// Tag lines with ALL-CAPS keys (A-Z and space) are formatted as "\x1b[1mKEY\x1b[21m: value".
// Non-tag lines (comments) are preserved after the tag section.
func FormatInfo(info string) string {
	lines := strings.Split(info, "\n")
	var tags []string
	var comments []string

	for _, line := range lines {
		sep := strings.Index(line, ":")
		if sep > 0 {
			key := line[:sep]
			test := key
			for _, c := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ " {
				test = strings.ReplaceAll(test, string(c), "")
			}
			if len(test) == 0 && len(strings.TrimSpace(key)) > 0 {
				value := strings.TrimSpace(line[sep+1:])
				formattedLine := fmt.Sprintf("\x1b[1m%s\x1b[22m: %s", strings.TrimSpace(key), value)
				tags = append(tags, formattedLine)
				continue
			}
		}
		comments = append(comments, line)
	}

	commentStr := strings.Join(comments, "\n")
	commentStr = strings.TrimLeft(commentStr, "\n")

	res := strings.Join(tags, "\n")
	if len(tags) > 0 && len(commentStr) > 0 {
		res += "\n\n" + commentStr
	} else if len(commentStr) > 0 {
		res += commentStr
	}
	return res
}

// ParsePony parses the raw content of a .pony file.
func ParsePony(name, rawContent string) (*Pony, error) {
	p := &Pony{
		Name:     name,
		Metadata: make(map[string]string),
	}

	content := strings.ReplaceAll(rawContent, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	bodyStartIndex := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "$$$" {
		endIdx := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "$$$" {
				endIdx = i
				break
			}
		}

		if endIdx != -1 {
			p.HasInfo = true
			p.RawInfo = strings.Join(lines[1:endIdx], "\n")
			for _, metaLine := range lines[1:endIdx] {
				idx := strings.Index(metaLine, ":")
				if idx > 0 {
					key := strings.TrimSpace(metaLine[:idx])
					val := strings.TrimSpace(metaLine[idx+1:])
					p.Metadata[key] = val
				}
			}
			bodyStartIndex = endIdx + 1
		}
	}

	p.BodyLines = lines[bodyStartIndex:]

	if topStr, ok := p.Metadata["BALLOON TOP"]; ok {
		if topVal, err := strconv.Atoi(topStr); err == nil {
			p.BalloonTop = topVal
		}
	}

	if botStr, ok := p.Metadata["BALLOON BOTTOM"]; ok {
		if botVal, err := strconv.Atoi(botStr); err == nil {
			p.BalloonBottom = botVal
		}
	}

	return p, nil
}

// RenderPonyWithBalloon combines the formatted balloon lines with the pony art lines.
func (p *Pony) RenderPonyWithBalloon(balloonLines []string, linkChar, linkColor string) string {
	if len(balloonLines) == 0 {
		return p.RenderPonyOnly()
	}

	coloredLink := color.ApplyColor(linkChar, linkColor)

	output := make([]string, 0, len(p.BodyLines)+len(balloonLines))
	balloonInserted := false

	for _, line := range p.BodyLines {
		if strings.Contains(line, "$balloon") && balloonTagRegex.MatchString(line) {
			if !balloonInserted {
				output = append(output, balloonLines...)
				balloonInserted = true
			}
			continue
		}

		processedLine := line
		if strings.Contains(line, "$\\$") {
			processedLine = strings.ReplaceAll(line, "$\\$", coloredLink)
		}
		output = append(output, processedLine)
	}

	if !balloonInserted {
		newOutput := make([]string, 0, len(balloonLines)+len(output))
		newOutput = append(newOutput, balloonLines...)
		newOutput = append(newOutput, output...)
		output = newOutput
	}

	return strings.Join(output, "\n")
}

// RenderPonyOnly returns pony artwork only, slicing off balloon stem lines.
func (p *Pony) RenderPonyOnly() string {
	lines := p.BodyLines
	topCut := 0

	if len(lines) > 0 {
		if strings.Contains(lines[0], "$balloon") && balloonTagRegex.MatchString(lines[0]) {
			topCut = 1 + p.BalloonTop
		} else if p.BalloonTop > 0 {
			topCut = p.BalloonTop
		}
	}

	if topCut > len(lines) {
		topCut = len(lines)
	}

	lines = lines[topCut:]

	if p.BalloonBottom > 0 && p.BalloonBottom <= len(lines) {
		lines = lines[:len(lines)-p.BalloonBottom]
	}

	output := make([]string, 0, len(lines))
	for _, line := range lines {
		clean := line
		if strings.Contains(clean, "$balloon") {
			clean = balloonTagRegex.ReplaceAllString(clean, "")
		}
		if strings.Contains(clean, "$\\$") {
			clean = strings.ReplaceAll(clean, "$\\$", "")
		}
		output = append(output, clean)
	}

	return strings.Join(output, "\n")
}
