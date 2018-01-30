package minute

import (
	"errors"
	"regexp"

	"github.com/pkrss/go-sina/common"
)

var regCtp *regexp.Regexp

func queryMinute(symbol string) (retList []MINUTEDATA, retE error) {

	if symbol == "" {
		retE = errors.New("分线股票不可为空")
		return
	}

	if regCtp == nil {
		regCtp = regexp.MustCompile(common.Regexp_Ctp)
	}
	ss := regCtp.FindStringSubmatch(symbol)
	if len(ss) == 0 {
		retE = errors.New("暂只支持期货K线")
	}

	retList, retE = queryMinuteFuture(symbol)

	return
}

func queryMinuteFuture(symbol string) (retList []MINUTEDATA, retE error) {
	return fetchSinaMinuteFutureDo(symbol)
}
