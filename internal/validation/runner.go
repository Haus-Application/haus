// Package validation runs Buster- and GOB-style Q&A against each KB device
// entry and scores the answers against the frontmatter ground truth. Output is
// a JSON report per device plus an aggregate summary.
//
// GOB is on stage tonight. The illusion is 100 devices describing themselves
// back to us, with no props, no tools, no live state — just the docs.
package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/coalson/haus/internal/kb"
)

// Score is the per-device score breakdown.
type Score struct {
	Buster    int `json:"buster"`
	BusterMax int `json:"buster_max"`
	GOB       int `json:"gob"`
	GOBMax    int `json:"gob_max"`
	TotalPct  int `json:"total_pct"`
}

// QA is a single question + graded answer.
type QA struct {
	Question string   `json:"q"`
	Answer   string   `json:"a"`
	Score    int      `json:"score"`
	Max      int      `json:"max"`
	Notes    []string `json:"notes"`
}

// DeviceReport is the per-device validation output.
type DeviceReport struct {
	Slug              string    `json:"slug"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`
	IntegrationKey    string    `json:"integration_key"`
	IntegrationStatus string    `json:"integration_status"`
	RanAt             time.Time `json:"ran_at"`
	Score             Score     `json:"score"`
	Buster            []QA      `json:"buster"`
	GOB               []QA      `json:"gob"`
	Gaps              []string  `json:"gaps"`
}

// Summary aggregates all reports.
type Summary struct {
	RanAt        time.Time      `json:"ran_at"`
	TotalDevices int            `json:"total_devices"`
	AvgScore     int            `json:"avg_score"`
	ByCategory   map[string]int `json:"by_category"`
	Failing      []string       `json:"failing"`
	Passing      int            `json:"passing"`
	Warning      int            `json:"warning"`
}

// Progress is emitted per completed device.
type Progress struct {
	JobID  string `json:"job_id,omitempty"`
	Slug   string `json:"slug"`
	Done   int    `json:"done"`
	Total  int    `json:"total"`
	Score  int    `json:"score"`
	Status string `json:"status"` // running | scoring | done | failed
	Error  string `json:"error,omitempty"`
}

// Runner wires a catalog + Claude client + output dir together.
type Runner struct {
	Catalog     *kb.Catalog
	Client      anthropic.Client
	Model       string
	Concurrency int
	OutDir      string
}

// BusterQuestions are the 7 technical questions asked of every device.
var BusterQuestions = []string{
	"Introduce yourself. What's your manufacturer, model number, and what do you do?",
	"How would Haus discover you on the network? Tell me the OUI/MAC prefix, mDNS service type, and default ports.",
	"What protocol and port do you use? Is it TLS? Self-signed cert? What encoding?",
	"How do I authenticate with you? Walk me through initial pairing or setup.",
	"What can I control or query about you? List your main capabilities.",
	"Show me one concrete example — the exact HTTP request or command to turn you on, query state, or do your main thing.",
	"What's one quirk, gotcha, or rate limit I should know about?",
}

// GOBQuestions are the 4 UX/UI questions.
var GOBQuestions = []string{
	"If I'm building your control page in the app, what 3-4 things must the user see right away when the page loads?",
	"Map each of your capabilities to a UI widget — toggle, slider, color picker, number input, button group, etc. Be specific.",
	"Describe yourself in one sentence for a non-technical user. No jargon.",
	"If the app can't reach you (you're offline), what should the app show instead? What's the graceful degradation?",
}

// NewRunner builds a Runner with sane defaults.
func NewRunner(catalog *kb.Catalog, apiKey, model, outDir string, concurrency int) *Runner {
	if model == "" {
		model = "claude-haiku-4-5-20251001"
	}
	if concurrency <= 0 {
		concurrency = 8
	}
	if outDir == "" {
		outDir = "validation"
	}
	var opts []option.RequestOption
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	return &Runner{
		Catalog:     catalog,
		Client:      anthropic.NewClient(opts...),
		Model:       model,
		Concurrency: concurrency,
		OutDir:      outDir,
	}
}

// RunOne validates a single device by slug and writes its JSON report.
func (r *Runner) RunOne(ctx context.Context, slug string) (*DeviceReport, error) {
	dev, ok := r.Catalog.BySlug[slug]
	if !ok {
		return nil, fmt.Errorf("unknown device slug %q", slug)
	}
	return r.runDevice(ctx, dev)
}

// RunAll validates every device in the catalog. Calls `progress` after each
// completion (or failure) if non-nil.
func (r *Runner) RunAll(ctx context.Context, progress func(Progress)) (*Summary, error) {
	if r.Catalog == nil || len(r.Catalog.All) == 0 {
		return nil, errors.New("empty catalog")
	}

	if err := os.MkdirAll(filepath.Join(r.OutDir, "devices"), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir outdir: %w", err)
	}

	total := len(r.Catalog.All)
	sem := make(chan struct{}, r.Concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var done int
	reports := make([]*DeviceReport, 0, total)

	for _, dev := range r.Catalog.All {
		d := dev
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			rep, err := r.runDevice(ctx, d)

			mu.Lock()
			done++
			currentDone := done
			mu.Unlock()

			if err != nil {
				log.Printf("[valid] %s: FAILED: %v", d.Slug, err)
				if progress != nil {
					progress(Progress{Slug: d.Slug, Done: currentDone, Total: total, Status: "failed", Error: err.Error()})
				}
				return
			}
			mu.Lock()
			reports = append(reports, rep)
			mu.Unlock()
			log.Printf("[valid] %s: %d%% (B %d/%d G %d/%d)", d.Slug, rep.Score.TotalPct, rep.Score.Buster, rep.Score.BusterMax, rep.Score.GOB, rep.Score.GOBMax)
			if progress != nil {
				progress(Progress{Slug: d.Slug, Done: currentDone, Total: total, Score: rep.Score.TotalPct, Status: "done"})
			}
		}()
	}
	wg.Wait()

	sum := buildSummary(reports)
	if err := writeJSON(filepath.Join(r.OutDir, "summary.json"), sum); err != nil {
		return sum, fmt.Errorf("write summary: %w", err)
	}

	return sum, nil
}

func (r *Runner) runDevice(ctx context.Context, d *kb.Device) (*DeviceReport, error) {
	ranAt := time.Now().UTC()

	busterAnswers := make([]string, len(BusterQuestions))
	for i, q := range BusterQuestions {
		ans, err := r.ask(ctx, d, q)
		if err != nil {
			return nil, fmt.Errorf("buster q%d: %w", i+1, err)
		}
		busterAnswers[i] = ans
	}
	gobAnswers := make([]string, len(GOBQuestions))
	for i, q := range GOBQuestions {
		ans, err := r.ask(ctx, d, q)
		if err != nil {
			return nil, fmt.Errorf("gob q%d: %w", i+1, err)
		}
		gobAnswers[i] = ans
	}

	buster := gradeBuster(busterAnswers, d)
	gob := gradeGOB(gobAnswers, d)

	bScore, bMax := sumQA(buster)
	gScore, gMax := sumQA(gob)
	total := bMax + gMax
	pct := 0
	if total > 0 {
		pct = (bScore + gScore) * 100 / total
	}

	gaps := detectGaps(buster, gob)

	rep := &DeviceReport{
		Slug:              d.Slug,
		Name:              d.Name,
		Category:          d.Category,
		IntegrationKey:    d.IntegrationKey,
		IntegrationStatus: d.IntegrationStatus,
		RanAt:             ranAt,
		Score:             Score{Buster: bScore, BusterMax: bMax, GOB: gScore, GOBMax: gMax, TotalPct: pct},
		Buster:            buster,
		GOB:               gob,
		Gaps:              gaps,
	}

	// Write per-device report
	if r.OutDir != "" {
		if err := os.MkdirAll(filepath.Join(r.OutDir, "devices"), 0o755); err == nil {
			_ = writeJSON(filepath.Join(r.OutDir, "devices", d.Slug+".json"), rep)
		}
	}
	return rep, nil
}

// ask issues one Claude call: system = device persona + KB body; user = question.
func (r *Runner) ask(ctx context.Context, d *kb.Device, question string) (string, error) {
	system := fmt.Sprintf(
		"You ARE the device \"%s\" (%s %s). Answer as the device, in first person.\nYour entire knowledge is the documentation below — use it, don't make things up.\nBe concise. Lead with the answer. No preamble.\n\n# Your Documentation\n%s",
		d.Name, d.Manufacturer, d.Model, d.Body,
	)

	resp, err := r.Client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(r.Model),
		MaxTokens: 400,
		System:    []anthropic.TextBlockParam{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(question)),
		},
	})
	if err != nil {
		return "", err
	}

	var out string
	for _, b := range resp.Content {
		if b.Type == "text" {
			out += b.Text
		}
	}
	return out, nil
}

func sumQA(qs []QA) (score, max int) {
	for _, q := range qs {
		score += q.Score
		max += q.Max
	}
	return
}

func detectGaps(buster, gob []QA) []string {
	var gaps []string
	for i, q := range buster {
		if q.Max > 0 && q.Score*2 < q.Max { // scored less than half
			gaps = append(gaps, fmt.Sprintf("Buster Q%d scored %d/%d", i+1, q.Score, q.Max))
		}
	}
	for i, q := range gob {
		if q.Max > 0 && q.Score*2 < q.Max {
			gaps = append(gaps, fmt.Sprintf("GOB Q%d scored %d/%d", i+1, q.Score, q.Max))
		}
	}
	if len(gaps) > 3 {
		gaps = gaps[:3]
	}
	return gaps
}

func buildSummary(reports []*DeviceReport) *Summary {
	sum := &Summary{
		RanAt:        time.Now().UTC(),
		TotalDevices: len(reports),
		ByCategory:   map[string]int{},
	}
	if len(reports) == 0 {
		return sum
	}

	totalScore := 0
	categoryTotals := map[string]int{}
	categoryCounts := map[string]int{}
	var failing []string
	for _, r := range reports {
		totalScore += r.Score.TotalPct
		categoryTotals[r.Category] += r.Score.TotalPct
		categoryCounts[r.Category]++
		if r.Score.TotalPct < 70 {
			failing = append(failing, r.Slug)
		} else if r.Score.TotalPct < 85 {
			sum.Warning++
		} else {
			sum.Passing++
		}
	}
	sort.Strings(failing)
	sum.Failing = failing
	sum.AvgScore = totalScore / len(reports)
	for cat, tot := range categoryTotals {
		sum.ByCategory[cat] = tot / categoryCounts[cat]
	}
	return sum
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
