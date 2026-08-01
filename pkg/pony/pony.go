package pony

import (
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
	BalloonTop    int
	BalloonBottom int
	BodyLines     []string
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

	var output []string
	output = append(output, balloonLines...)

	for _, line := range p.BodyLines {
		processedLine := balloonTagRegex.ReplaceAllString(line, "")
		processedLine = strings.ReplaceAll(processedLine, "$\\$", coloredLink)

		if strings.TrimSpace(processedLine) != "" || len(output) > len(balloonLines) {
			output = append(output, processedLine)
		}
	}

	return strings.Join(output, "\n")
}

// RenderPonyOnly returns pony artwork only, slicing off balloon stem lines.
func (p *Pony) RenderPonyOnly() string {
	lines := p.BodyLines
	topCut := 0

	if len(lines) > 0 && balloonTagRegex.MatchString(lines[0]) {
		topCut = 1 + p.BalloonTop
	} else if p.BalloonTop > 0 {
		topCut = p.BalloonTop
	}

	if topCut > len(lines) {
		topCut = len(lines)
	}

	lines = lines[topCut:]

	if p.BalloonBottom > 0 && p.BalloonBottom <= len(lines) {
		lines = lines[:len(lines)-p.BalloonBottom]
	}

	var output []string
	for _, line := range lines {
		clean := balloonTagRegex.ReplaceAllString(line, "")
		clean = strings.ReplaceAll(clean, "$\\$", "")
		output = append(output, clean)
	}

	return strings.Join(output, "\n")
}
