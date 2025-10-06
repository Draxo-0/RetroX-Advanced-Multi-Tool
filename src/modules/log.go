package modules

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func getTerminalWidth() int {
	// Try to get actual terminal width, fallback to reasonable default
	if runtime.GOOS == "windows" {
		// For Windows, try to get console width
		cmd := exec.Command("cmd", "/c", "mode con")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Columns:") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						if width, err := strconv.Atoi(parts[1]); err == nil && width > 0 {
							return width
						}
					}
				}
			}
		}
	} else {
		// For Unix-like systems, try to get terminal width
		cmd := exec.Command("tput", "cols")
		output, err := cmd.Output()
		if err == nil {
			if width, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && width > 0 {
				return width
			}
		}
	}
	// Fallback to a reasonable default that works well for most terminals
	return 120
}

func (m *Modules) Menu() {
	opt := map[int]string{
		1:  "Mass DM",
		2:  "Dm Spam",
		3:  "React Verify",
		4:  "Joiner",
		5:  "Leaver",
		6:  "Accept Rules",
		7:  "Raid Channel",
		8:  "Scrape Users",
		9:  "Check Tokens",
		10: "Mass Ping",
		11: "Button Click",
		12: "Friender",
		13: "Token Menu",
		14: "Booster",
		15: "VoiceChat",
		16: "SoundBoard",
		17: "OnBoarding",
		18: "Server Info",
	}
	tkn, _, _ := modules.ReadFile("tokens.txt")

	// Add some vertical spacing at the top
	fmt.Println()

	// Create ASCII art aligned to the left
	asciiArt := `/* +==============================================+ */
/* |                                              | */
/* |                                              | */
/* | ________       _____                ____  __ | */
/* | ___  __ \_____ __  /_______________ __  |/ / | */
/* | __  /_/ /_  _ \_  __/__  ___/_  __ \__    /  | */
/* | _  _, _/ /  __// /_  _  /    / /_/ /_    |   | */
/* | /_/ |_|  \___/ \__/  /_/     \____/ /_/|_|   | */
/* |                                              | */
/* |             Advanced Multi-Tool              | */
/* +==============================================+ */`

	// Print ASCII art aligned to the left
	lines := strings.Split(asciiArt, "\n")
	for _, line := range lines {
		fmt.Println(line)
	}

	// Print token count aligned to the left
	tokenText := "[Tokens: " + strconv.Itoa(len(tkn)) + "]"
	fmt.Printf("\n%s\n\n", tokenText)

	m.PrintMenu(opt)
}

func (m *Modules) StrlogV(text string, data string, s time.Time) {
	e := time.Since(s)
	durationStr := e.String()
	if len(durationStr) > 3 {
		durationStr = durationStr[:3]
	}
	fmt.Printf("[%sms] [✓]%s: %s\n", durationStr, text, data)
}

func (m *Modules) StrlogE(text string, data string, s time.Time) {
	e := time.Since(s)
	durationStr := e.String()
	if len(durationStr) > 3 {
		durationStr = durationStr[:3]
	}
	fmt.Printf("[%sms] [X]%s: %s\n", durationStr, text, data)
}

func (m *Modules) StrlogR(text string, data string, s time.Time) {
	e := time.Since(s)
	durationStr := e.String()
	if len(durationStr) > 3 {
		durationStr = durationStr[:3]
	}
	fmt.Printf("[%sms] [-]%s: %s\n", durationStr, text, data)
}

func (m *Modules) PrintMenu(options map[int]string) {
	// Find the longest menu item to calculate exact width needed
	maxItemLength := 0
	for _, item := range options {
		if len(item) > maxItemLength {
			maxItemLength = len(item)
		}
	}

	var sb strings.Builder
	var count int

	tw := tabwriter.NewWriter(&sb, 0, 0, 4, ' ', 0)
	opts := make([]int, 0, len(options))
	for k := range options {
		opts = append(opts, k)
	}
	sort.Ints(opts)

	num := (len(opts) + 2 - 1) / 2

	for row := 0; row < num; row++ {
		for col := 0; col < 2; col++ {
			i := col*num + row
			if i < len(opts) {
				optnum := strconv.Itoa(opts[i])
				if len(optnum) == 1 {
					optnum = "0" + optnum
				}
				// Use dynamic width based on longest item, aligned to the left
				fmt.Fprintf(tw, "[%s]  %-*s  ", optnum, maxItemLength, options[opts[i]])
				count++
			}
		}
		fmt.Fprintln(tw)
	}
	if count%2 != 0 {
		fmt.Fprint(tw, " ")
	}
	tw.Flush()

	fmt.Println(sb.String())
}

func (m *Modules) Cls() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PrintStruct(data interface{}) {
	fmt.Printf("%+v\n", data)
}
