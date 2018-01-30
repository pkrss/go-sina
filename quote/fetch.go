package quote

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkrss/gosina/common"

	"github.com/pkrss/gosina/check"

	"github.com/axgle/mahonia"
)

var gbkEnc mahonia.Decoder

var regCtpCffex *regexp.Regexp

func FetchSinaHqDo(list ...string) (ret map[string]map[string]interface{}) {

	url := fmt.Sprintf("http://hq.sinajs.cn/?_=%d000/&list=", time.Now().Unix())
	for idx, stk := range list {
		if idx != 0 {
			url += ","
		}

		if regCtpCffex == nil {
			regCtpCffex = regexp.MustCompile(common.Regexp_Cffex)
		}

		ss := regCtpCffex.FindStringSubmatch(stk)
		if len(ss) > 0 {
			url += "CFF_RE_" + stk
		} else {
			url += stk
		}

	}

	// log.Printf("Fetching %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Println(err2.Error())
		return
	}

	if len(body) == 0 {
		log.Println("fetchSinaHq body len = 0")
		return
	}

	rspStr := string(body)

	if gbkEnc == nil {
		gbkEnc = mahonia.NewDecoder("GBK")
	}

	rspStr = gbkEnc.ConvertString(rspStr)
	if len(rspStr) == 0 {
		log.Println("fetchSinaHq rspStr len = 0")
		return
	}

	sendData := make(map[string]interface{})
	sendData["cat"] = "quote"
	sendData["oper"] = "realtime"

	rspRows := strings.Split(rspStr, "\n")

	if len(rspRows)-1 != len(list) {
		log.Printf("fetchSinaHq rspRows len %d != list len %d", len(rspRows), len(list))
		return
	}

	var bp int
	var ep int

	const hx_str = "hq_str_"
	const hq_str_len = len(hx_str)

	ret = make(map[string]map[string]interface{}, 0)

	for _, rspRow := range rspRows {
		if rspRow == "" {
			continue
		}

		bp = strings.Index(rspRow, hx_str)
		ep = strings.IndexByte(rspRow, '=')

		if ep <= bp {
			log.Println("fetchSinaHq row not found match hq_str_xxx=")
			return
		}

		stk := rspRow[bp+hq_str_len : ep]

		bp = strings.IndexByte(rspRow, '"')
		ep = strings.LastIndexByte(rspRow, '"')

		if ep <= bp {
			log.Println("fetchSinaHq row not found match \"")
			return
		}

		rspStr = rspRow[bp+1 : ep]

		rspRowCols := strings.Split(rspStr, ",")

		// log.Printf("fetchSinaHq prepare parse row: %s", stk)

		var e error

		quote := make(map[string]interface{})
		sendData["data"] = &quote

		if check.IsZh(stk) {
			e = parseZh(quote, stk, rspRowCols)
			if e != nil {
				continue
			}
		} else if check.IsFx(stk) {
			e = parseFx(quote, stk, rspRowCols)
			if e != nil {
				continue
			}
		} else if check.IsCtp(stk) {
			stk, e = parseCtp(quote, stk, rspRowCols)
			if e != nil {
				continue
			}
		} else {
			log.Printf("fetchSinaHq not valid stk %s", stk)
			continue
		}

		ret[stk] = sendData

		// SetCacheHqInfo(stk, sendData)
	}

	return
}
