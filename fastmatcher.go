package hlconf2017

import (
	"regexp"
	"strings"
)

// FastPatternMatcher implements high-performance Graphite metric filtering
type FastPatternMatcher struct {
	AllowedPatterns []string
    CachePatterns map[string]*regexp.Regexp
}

// InitPatterns accepts allowed patterns in Graphite format, e.g.
//   metric.name.single
//   metric.name.*
//   metric.name.wild*card
//   metric.name.{one,two}.maybe.longer
func (p *FastPatternMatcher) InitPatterns(allowedPatterns []string) {
	p.AllowedPatterns = allowedPatterns
    p.CachePatterns = make(map[string]*regexp.Regexp)

    for _, pattern := range p.AllowedPatterns {
        patternParts := strings.Split(pattern, ".")

        for _, part := range patternParts {
            regexPart := "^" + part + "$"
            regexPart = strings.Replace(regexPart, "*", ".*", -1)
            regexPart = strings.Replace(regexPart, "{", "(", -1)
            regexPart = strings.Replace(regexPart, "}", ")", -1)
            regexPart = strings.Replace(regexPart, ",", "|", -1)
            regex := regexp.MustCompile(regexPart)

            p.CachePatterns[part] = regex
        }
    }
}

// DetectMatchingPatterns returns a list of allowed patterns that match given metric
func (p *FastPatternMatcher) DetectMatchingPatterns(metricName string) (matchingPatterns []string) {
	metricParts := strings.Split(metricName, ".")

NEXTPATTERN:
	for _, pattern := range p.AllowedPatterns {
		patternParts := strings.Split(pattern, ".")
		if len(patternParts) != len(metricParts) {
			continue NEXTPATTERN
		}
		for i, part := range patternParts {
            if part == metricParts[i] {
                continue
            }

            var regex *regexp.Regexp
            regex = p.CachePatterns[part]

			if !regex.MatchString(metricParts[i]) {
				continue NEXTPATTERN
			}
		}
		matchingPatterns = append(matchingPatterns, pattern)
	}

	return
}
