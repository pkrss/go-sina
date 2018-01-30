package check

import (
	"fmt"
	"strings"
	"time"
)

func CheckSamePeriod(period string, d ...time.Time) bool {

	if period == "" {
		return false
	}
	period = strings.ToLower(period)

	var big time.Time
	var small time.Time

	if len(d) < 1 {
		return false
	}

	if len(d) > 1 {
		big = d[1]
	} else {
		big = time.Now()
	}

	small = d[0]

	if small.After(big) {
		t := small
		small = big
		big = t
	}

	var ret bool

	var name string
	var num int
	n, e := fmt.Sscanf(period, "%d%s", &num, &name)
	if e != nil || n < 2 {
		return false
	}

	name2Seconds := map[string]int{"m": 60, "h": 60 * 60, "d": 60 * 60 * 24, "w": 60 * 60 * 24 * 7, "mon": 60 * 60 * 24 * 31, "y": 60 * 60 * 24 * 365}
	seconds, ok := name2Seconds[name]
	if !ok {
		return false
	}

	ret = (big.Unix() - small.Unix()) > int64(seconds)*int64(num)

	if ret {
		return false
	}

	vb := 0
	vs := 0
	switch name {
	case "m":
		vb = big.Minute()
		vs = small.Minute()
	case "h":
		vb = big.Hour()
		vs = small.Hour()
	case "d":
		vb = big.Day()
		vs = small.Day()
	case "w":
		_, vb = big.ISOWeek()
		_, vs = small.ISOWeek()
	case "mon":
		vb = int(big.Month())
		vs = int(small.Month())
	case "y":
		vb = big.Year()
		vs = small.Year()
	}
	if vs > vb {
		return false
	}
	ret = ((vb / num) == (vs / num))

	return ret
}

func main() {
	d1 := time.Date(2017, 1, 1, 0, 0, 0, 0, nil)
	d2 := time.Date(2017, 2, 1, 0, 0, 0, 0, nil)
	period := "1w"
	f := CheckSamePeriod(period, d1, d2)

	fmt.Printf("%v - %v, same %s = %v", d1, d2, period, f)
}
