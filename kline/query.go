package kline

import (
	"errors"
	"fmt"
	"regexp"
	"sx98/base/sina"
)

var regCtp *regexp.Regexp

func queryKline(symbol string, period string, indicate string) (retList []K_MACDV, retE error) {

	if symbol == "" {
		retE = fmt.Errorf("K线股票不可为空, %s", period)
		return
	}

	if regCtp == nil {
		regCtp = regexp.MustCompile(sina.Regexp_Ctp)
	}
	ss := regCtp.FindStringSubmatch(symbol)
	if len(ss) == 0 {
		retE = errors.New("暂只支持期货K线")
	}

	retList, retE = queryKlineFuture(symbol, period, indicate)

	return
}

func queryKlineFuture(symbol string, period string, indicate string) (retList []K_MACDV, retE error) {

	if indicate != "" {
		retList, _, retE = KlineDataQuery(symbol, period, "", nil)
		if retE != nil {
			return
		}

		

		return
	}

	var fetchPeriod string

	switch period {
	case "5m", "15m", "30m", "60m", "1d":
		fetchPeriod = period
	case "1w", "1mon", "1y":
		fetchPeriod = "1d"
	}

	if fetchPeriod == "" {
		retE = fmt.Errorf("不支持的K线周期:%s", period)
		return
	}

	if fetchPeriod == period {
		retList, retE = fetchSinaKlineFutureDo(symbol, fetchPeriod)
		return
	}

	retList, _, retE = KlineDataQuery(symbol, fetchPeriod, indicate, nil)
	if retE != nil {
		return
	}

	retList = handlePeriodKData(retList, period)

	return
}
