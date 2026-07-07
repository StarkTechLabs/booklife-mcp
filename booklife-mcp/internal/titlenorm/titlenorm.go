package titlenorm

import (
	"regexp"
	"strings"
)

// Result holds the normalized title and any series metadata extracted from the raw string.
type Result struct {
	Cleaned      string // clean, searchable title
	Series       string // series name if detected, e.g. "Don Tillman"
	SeriesNumber string // series position if detected, e.g. "2"
	Raw          string // original unmodified input
}

var (
	smartQuoteReplacer = strings.NewReplacer(
		"‘", "'", "’", "'",
		"“", `"`, "”", `"`,
		"ʼ", "'",
	)

	// Matches a full parenthetical including preceding whitespace.
	parenRe = regexp.MustCompile(`\s*\(([^)]*)\)`)

	// Matches "Name Book N" or "Name #N" at end of a parenthetical's content.
	bookNumRe = regexp.MustCompile(`(?i)^(.*?)\s+(?:Book\s+|#)(\d+)\s*$`)

	// Parentheticals matching these words are marketing copy, not series names.
	marketingRe = regexp.MustCompile(`(?i)\b(club|pick|selection|tie.in|edition|winner|finalist|award)\b`)

	// Em/en dash used as subtitle separator (must have surrounding spaces).
	dashSubtitle = regexp.MustCompile(`\s+[–—]\s+.+$`)

	// Spaced hyphen as subtitle separator, followed by uppercase word.
	spacedDash = regexp.MustCompile(`\s+-\s+[A-Z].+$`)

	// "A Novel" / "A Novella" suffix.
	novelSuffix = regexp.MustCompile(`(?i)[,\s]+a\s+novella?\s*$`)

	// "A <adjectives> Romance/Thriller/etc" genre-tagline suffix.
	genreSuffix = regexp.MustCompile(`(?i)[,\s]+a\s+\w+(?:\s+\w+){0,6}\s+(romance|thriller|mystery|romantasy|adventure|suspense|fantasy)\s*$`)

	// Trailing exclamation marks from marketing copy.
	trailingBang = regexp.MustCompile(`[!]+$`)

	// Trailing separators left after other stripping.
	trailingPunct = regexp.MustCompile(`[\s:,\-–—]+$`)
)

// Normalize strips marketing copy and subtitle text from a raw book title. It returns a
// Result with a clean string suitable for API lookups, plus any series metadata extracted
// from parentheticals. If the cleaned result would be too short, Raw is returned as Cleaned.
//
// Callers should search with Cleaned first, then retry with Raw on a miss.
func Normalize(raw string) Result {
	if raw == "" {
		return Result{Raw: raw}
	}

	s := smartQuoteReplacer.Replace(raw)

	// Extract series metadata from the full string before any truncation destroys it.
	var series, seriesNum string
	for _, m := range parenRe.FindAllStringSubmatch(s, -1) {
		content := strings.TrimSpace(m[1])
		if content == "" || marketingRe.MatchString(content) {
			continue
		}
		if nm := bookNumRe.FindStringSubmatch(content); nm != nil {
			name := strings.TrimSpace(nm[1])
			// Skip series names that themselves contain colons — too ambiguous to extract cleanly.
			if !strings.Contains(name, ":") {
				series = name
				seriesNum = nm[2]
			}
		} else if !strings.Contains(content, ":") {
			series = content
		}
		break // use first valid parenthetical only
	}

	// Remove all parentheticals.
	s = parenRe.ReplaceAllString(s, "")

	// Truncate at first colon if it leaves a meaningful prefix.
	if idx := strings.Index(s, ":"); idx > 0 {
		pre := strings.TrimSpace(s[:idx])
		if len([]rune(pre)) >= 2 {
			s = pre
		}
	}

	// Truncate at em/en dash subtitle separators and spaced hyphens.
	s = dashSubtitle.ReplaceAllString(s, "")
	s = spacedDash.ReplaceAllString(s, "")

	// Strip trailing junk phrases until no more match.
	for {
		prev := s
		s = trailingBang.ReplaceAllString(s, "")
		s = novelSuffix.ReplaceAllString(s, "")
		s = genreSuffix.ReplaceAllString(s, "")
		if s == prev {
			break
		}
	}

	// Normalize whitespace and strip trailing separator characters.
	s = strings.Join(strings.Fields(s), " ")
	s = trailingPunct.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)

	if len([]rune(s)) < 2 {
		s = raw
	}

	return Result{
		Cleaned:      s,
		Series:       series,
		SeriesNumber: seriesNum,
		Raw:          raw,
	}
}
