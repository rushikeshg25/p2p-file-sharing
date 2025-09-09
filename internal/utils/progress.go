package utils

import (
	"fmt"
	"strings"
	"time"
)

type ProgressBar struct {
	totalBytes   int64
	currentBytes int64
	startTime    time.Time
	label        string
	width        int
	done         bool
	lastRender   time.Time
	minInterval  time.Duration
}

func NewProgressBar(totalBytes int64, label string) *ProgressBar {
	if label == "" {
		label = "Progress"
	}
	return &ProgressBar{
		totalBytes:   totalBytes,
		currentBytes: 0,
		startTime:    time.Now(),
		label:        label,
		width:        24,
		lastRender:   time.Time{},
		minInterval:  100 * time.Millisecond,
	}
}

func (p *ProgressBar) Add(n int64) {
	if p.done {
		return
	}
	p.currentBytes += n
	if p.currentBytes > p.totalBytes {
		p.currentBytes = p.totalBytes
	}
	p.render(false)
}

func (p *ProgressBar) Finish() {
	if p.done {
		return
	}
	p.done = true
	p.currentBytes = p.totalBytes
	p.render(true)
	fmt.Print("\n")
}

func (p *ProgressBar) render(force bool) {
	if !force && !p.lastRender.IsZero() && time.Since(p.lastRender) < p.minInterval {
		return
	}
	percentage := 0.0
	if p.totalBytes > 0 {
		percentage = float64(p.currentBytes) / float64(p.totalBytes)
	}
	filled := int(percentage * float64(p.width))
	if filled > p.width {
		filled = p.width
	}
	bar := strings.Repeat("=", filled)
	if filled < p.width {
		bar += ">"
		bar += strings.Repeat(".", p.width-filled-1)
	}

	elapsed := time.Since(p.startTime)
	rate := float64(0)
	if elapsed > 0 {
		rate = float64(p.currentBytes) / elapsed.Seconds()
	}
	eta := "--:--"
	if p.totalBytes > 0 && rate > 0 {
		remaining := float64(p.totalBytes-p.currentBytes) / rate
		eta = formatDuration(time.Duration(remaining) * time.Second)
	}

	status := fmt.Sprintf("\r\x1b[2K%s %3.0f%% [%s] %s/%s %s/s ETA %s",
		p.label,
		percentage*100,
		bar,
		formatBytes(p.currentBytes),
		formatBytes(p.totalBytes),
		formatBytes(int64(rate)),
		eta,
	)
	fmt.Print(status)
	p.lastRender = time.Now()
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	value := float64(b) / float64(div)
	suffix := "KMGTPE"[exp : exp+1]
	return fmt.Sprintf("%.1f %ciB", value, suffix[0])
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	seconds := int(d.Seconds())
	mins := seconds / 60
	secs := seconds % 60
	if mins > 99 {
		return "99:59"
	}
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
