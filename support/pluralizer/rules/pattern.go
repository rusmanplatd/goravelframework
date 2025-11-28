package rules

import (
	"regexp"

	"github.com/rusmanplatd/goravelframework/contracts/support/pluralizer"
)

var _ pluralizer.Pattern = (*Pattern)(nil)

type Pattern struct {
	pattern *regexp.Regexp
}

func NewPattern(pattern string) *Pattern {
	return &Pattern{
		pattern: regexp.MustCompile(pattern),
	}
}

func (r *Pattern) Matches(word string) bool {
	return r.pattern.MatchString(word)
}
