package hlconf2017

import (
	"regexp"
	"strings"
)

// FastPatternMatcher implements high-performance Graphite metric filtering
type FastPatternMatcher struct {
	AllowedPatterns []string
	AllowedRegexps  []*regexp.Regexp
	PatternPartsNum []int
}

// InitPatterns accepts allowed patterns in Graphite format, e.g.
//   metric.name.single
//   metric.name.*
//   metric.name.wild*card
//   metric.name.{one,two}.maybe.longer
func (p *FastPatternMatcher) InitPatterns(allowedPatterns []string) {
	p.AllowedPatterns = allowedPatterns
	for _, pattern := range p.AllowedPatterns {
		regex_str := strings.Replace("^"+pattern+"$", ".", "\\.", -1)
		regex_str = strings.Replace(regex_str, "*", "[^.]*", -1)
		regex_str = strings.Replace(regex_str, "{", "(", -1)
		regex_str = strings.Replace(regex_str, "}", ")", -1)
		regex_str = strings.Replace(regex_str, ",", "|", -1)

		// fmt.Printf("regex_str: %s\n", regex_str)
		regex := regexp.MustCompile(regex_str)
		p.AllowedRegexps = append(p.AllowedRegexps, regex)
		partsNum := len(strings.Split(pattern, "."))
		p.PatternPartsNum = append(p.PatternPartsNum, partsNum)
	}
}

// DetectMatchingPatterns returns a list of allowed patterns that match given metric
func (p *FastPatternMatcher) DetectMatchingPatterns(metricName string) (matchingPatterns []string) {
	metricParts := strings.Split(metricName, ".")

NEXTPATTERN:
	for i, regex := range p.AllowedRegexps {
		if p.PatternPartsNum[i] != len(metricParts) {
			continue NEXTPATTERN
		}
		if !regex.MatchString(metricName) {
			continue NEXTPATTERN
		}
		matchingPatterns = append(matchingPatterns, p.AllowedPatterns[i])
	}

	return
}
