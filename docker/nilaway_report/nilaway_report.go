package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Lines     int       `json:"lines"`
	Rating    float64   `json:"rating"`
}

type History struct {
	Records []Record `json:"records"`
	// WorstLines is the largest issue-line count ever seen (or the supplied
	// baseline, whichever is bigger). It's what "0/10" is anchored to.
	WorstLines int `json:"worst_lines"`
}

func main() {
	historyPath := flag.String("history", ".nilaway_history.json", "path to history JSON file")
	passthrough := flag.Bool("passthrough", true, "echo the original nilaway output before printing the report")
	baselineWorst := flag.Int("baseline-worst", 1000, "issue-line count treated as the worst case (0/10) if history hasn't recorded anything worse")
	flag.Parse()

	lines := countStdinLines(*passthrough)

	hist, _ := loadHistory(*historyPath)
	var prev *Record
	if n := len(hist.Records); n > 0 {
		prev = &hist.Records[n-1]
	}

	worst := *baselineWorst
	if hist.WorstLines > worst {
		worst = hist.WorstLines
	}
	if lines > worst {
		worst = lines
	}

	rating := computeRating(lines, worst)

	printReport(lines, rating, prev, worst)

	hist.WorstLines = worst
	hist.Records = append(hist.Records, Record{
		Timestamp: time.Now(),
		Lines:     lines,
		Rating:    rating,
	})
	if len(hist.Records) > 200 {
		hist.Records = hist.Records[len(hist.Records)-200:]
	}
	if err := saveHistory(*historyPath, hist); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not save history: %v\n", err)
	}
}

func countStdinLines(passthrough bool) int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var out *bufio.Writer
	if passthrough {
		out = bufio.NewWriter(os.Stdout)
		defer out.Flush()
	}

	count := 0
	for scanner.Scan() {
		if passthrough {
			out.WriteString(scanner.Text())
			out.WriteByte('\n')
		}
		count++
	}
	if passthrough {
		out.Flush() // make sure the original report is visible before we print the summary
	}
	return count
}

func computeRating(issueLines, worst int) float64 {
	if worst <= 0 {
		return 10
	}
	r := 10 * (1 - float64(issueLines)/float64(worst))
	return math.Max(0, math.Min(10, r))
}

func printReport(lines int, rating float64, prev *Record, worst int) {
	fmt.Printf("%d lines of nilaway issues left.\n", lines)

	if prev == nil {
		fmt.Printf("Your code has been rated at %.2f/10 (no previous run to compare, worst case anchor: %d lines)\n", rating, worst)
		return
	}

	diffLines := prev.Lines - lines
	switch {
	case diffLines > 0:
		fmt.Printf("You improved your code by %d lines in comparison to previous run. %s\n", diffLines, encouragement(diffLines))
	case diffLines < 0:
		fmt.Printf("Your code got %d lines worse in comparison to previous run. %s\n", -diffLines, discouragement(-diffLines))
	default:
		fmt.Println("No change since the previous run.")
	}

	diffRating := rating - prev.Rating
	fmt.Printf("Your code has been rated at %.2f/10 (previous run: %.2f/10, %+.2f)\n", rating, prev.Rating, diffRating)

	if lines == worst && (diffLines < 0) {
		fmt.Println("New worst case recorded — the 0/10 line has been redrawn to fit it.")
	}

	if isNewBest(prev, lines) {
		fmt.Println("New best score!")
	}
}

func encouragement(n int) string {
	switch {
	case n >= 100:
		return "Huge cleanup!"
	case n >= 20:
		return "Nice work!"
	default:
		return "Keep it up!"
	}
}

func discouragement(n int) string {
	switch {
	case n >= 20:
		return "Might be worth a look."
	default:
		return "Small regression, no big deal."
	}
}

func isNewBest(prev *Record, lines int) bool {
	return prev != nil && lines < prev.Lines
}

func loadHistory(path string) (History, error) {
	var h History
	data, err := os.ReadFile(path)
	if err != nil {
		return h, err
	}
	if err := json.Unmarshal(data, &h); err != nil {
		return History{}, err
	}
	return h, nil
}

func saveHistory(path string, h History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
