package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	info, _ := os.Stdin.Stat()

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: cat <file> | cowsay")
		return
	}

	var lines []string
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}

	const COW = `         \  ^__^
          \ (oo)\_______
	    (__)\       )\/\
	        ||----w |
	        ||     ||
		`

	lines = tabsToSpaces(lines)
	maxWidth := calcMaxWidth(lines)
	messages := normalizeLines(lines, maxWidth)
	balloon := buildBalloon(messages, maxWidth)

	fmt.Println(balloon)
	fmt.Println(COW)
	fmt.Println()
}

func tabsToSpaces(lines []string) []string {
	var ret []string
	for _, line := range lines {
		line = strings.Replace(line, "\t", "    ", -1)
		ret = append(ret, line)
	}
	return ret
}

func calcMaxWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		lineLength := utf8.RuneCountInString(line)
		if maxWidth < lineLength {
			maxWidth = lineLength
		}
	}
	return maxWidth
}

func normalizeLines(lines []string, maxWidth int) []string {
	var ret []string
	for _, line := range lines {
		normalizedLine := line + strings.Repeat(" ", maxWidth-utf8.RuneCountInString(line))
		ret = append(ret, normalizedLine)
	}
	return ret
}

func buildBalloon(lines []string, maxWidth int) string {
	borders := []string{"/", "\\", "\\", "/", "|", "<", ">"}
	numLines := len(lines)
	var ret []string

	top := " " + strings.Repeat("_", maxWidth+2)
	bottom := " " + strings.Repeat("-", maxWidth+2)

	ret = append(ret, top)
	if numLines == 1 {
		// < message >
		line := fmt.Sprintf("%s %s %s", borders[5], lines[0], borders[6])
		ret = append(ret, line)
	} else {
		// / message \
		// | message |
		// \ message /
		line := fmt.Sprintf("%s %s %s", borders[0], lines[0], borders[1])
		ret = append(ret, line)
		for i := 1; i < numLines-1; i++ {
			line = fmt.Sprintf("%s %s %s", borders[4], lines[i], borders[4])
			ret = append(ret, line)
		}
		line = fmt.Sprintf("%s %s %s", borders[2], lines[numLines-1], borders[3])
		ret = append(ret, line)
	}

	ret = append(ret, bottom)
	return strings.Join(ret, "\n")
}
