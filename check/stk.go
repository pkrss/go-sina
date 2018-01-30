package check

import (
	"regexp"
	"strings"

	"github.com/pkrss/go-sina/common"
)

var regCtp *regexp.Regexp

func IsCtp(stk string) bool {
	if strings.HasPrefix(stk, "CFF_RE_") {
		return true
	}
	if regCtp == nil {
		regCtp = regexp.MustCompile(common.Regexp_Ctp)
	}
	ss := regCtp.FindStringSubmatch(stk)
	return len(ss) > 0
}

var regFx *regexp.Regexp

func IsFx(stk string) bool {
	if regZh == nil {
		regZh = regexp.MustCompile(`fx_.+`)
	}
	ss := regZh.FindStringSubmatch(stk)
	return len(ss) > 0
}

var regZh *regexp.Regexp

func IsZh(stk string) bool {
	if regZh == nil {
		regZh = regexp.MustCompile(`s[h|z]\d{6}`)
	}
	ss := regZh.FindStringSubmatch(stk)
	return len(ss) > 0
}
