package titlenorm

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		raw          string
		wantCleaned  string
		wantSeries   string
		wantSeriesNo string
	}{
		// Colon truncation + series extraction from parenthetical
		{
			raw:          "The Rosie Effect: A Novel (Don Tillman Book 2)",
			wantCleaned:  "The Rosie Effect",
			wantSeries:   "Don Tillman",
			wantSeriesNo: "2",
		},
		// Double colon (first wins)
		{
			raw:         "Sex and Vanity: A GMA Book Club Pick: A Novel",
			wantCleaned: "Sex and Vanity",
		},
		// Marketing parenthetical discarded (not extracted as series)
		{
			raw:         "The Sweetness of Water (Oprah’s Book Club): A Novel",
			wantCleaned: "The Sweetness of Water",
		},
		// Colon truncation removes trailing "!" in subtitle
		{
			raw:         "Twelve Days of Christmas: The bestselling Christmas read to devour in one sitting!",
			wantCleaned: "Twelve Days of Christmas",
		},
		// En-dash subtitle separator
		{
			raw:         "The Wrong Family: A Domestic Thriller – A Twisted Psychological Mystery About the Cracks in a Perfect Life",
			wantCleaned: "The Wrong Family",
		},
		// Series extraction across colon boundary — parens come after subtitle text
		{
			raw:          "Veiled Court: A Slow Burn Romantasy (Legacy of Avalon Book 1)",
			wantCleaned:  "Veiled Court",
			wantSeries:   "Legacy of Avalon",
			wantSeriesNo: "1",
		},
		// Series without a number
		{
			raw:         "Tempting Levi (Cade Brothers)",
			wantCleaned: "Tempting Levi",
			wantSeries:  "Cade Brothers",
		},
		// Complex series name with colon — not extracted; colon truncation does the rest
		{
			raw:         "His to Claim: A Forced Proximity Mafia Romance (Merciless Mercy: Sovarin Bratva Series Book 1)",
			wantCleaned: "His to Claim",
		},
		// Nothing to strip — passes through unchanged
		{
			raw:         "My Emergency Contact is a Wolf Shifter",
			wantCleaned: "My Emergency Contact is a Wolf Shifter",
		},
		// "Why Choose Romance" suffix stripped by genreSuffix after colon truncation
		{
			raw:         "Can I Join You?: A Why Choose Romance",
			wantCleaned: "Can I Join You?",
		},
		// Hyphen inside a word must not be split
		{
			raw:         "Spider-Man",
			wantCleaned: "Spider-Man",
		},
		// Empty input
		{
			raw:         "",
			wantCleaned: "",
		},
		// No subtitle junk at all
		{
			raw:         "Project Hail Mary",
			wantCleaned: "Project Hail Mary",
		},
		// Smart quotes normalized
		{
			raw:         "“The Book”: A Novel",
			wantCleaned: `"The Book"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got := Normalize(tc.raw)

			if got.Cleaned != tc.wantCleaned {
				t.Errorf("Cleaned:\n  got  %q\n  want %q", got.Cleaned, tc.wantCleaned)
			}
			if got.Series != tc.wantSeries {
				t.Errorf("Series:\n  got  %q\n  want %q", got.Series, tc.wantSeries)
			}
			if got.SeriesNumber != tc.wantSeriesNo {
				t.Errorf("SeriesNumber:\n  got  %q\n  want %q", got.SeriesNumber, tc.wantSeriesNo)
			}
			if got.Raw != tc.raw {
				t.Errorf("Raw should equal input:\n  got  %q\n  want %q", got.Raw, tc.raw)
			}
		})
	}
}
