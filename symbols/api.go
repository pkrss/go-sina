package symbols

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkrss/go-sina/kline"
	pkFile "github.com/pkrss/go-utils/file"
	"github.com/pkrss/go-utils/profile"
	pkTime "github.com/pkrss/go-utils/time"
)

var openCtpMarketHost = "" // open-ctp-market

func OptInit(host string) {
	openCtpMarketHost = host
}

func GetExchangeList() (retList []Exchange2Items, retE error) {

	saveFileName := profile.ProfileReadString("symbols_save_path")

	lastAccessTime, fileExist := pkFile.FileLastWriteTime(saveFileName)

	if !fileExist {
		pkFile.CreateDir(pkFile.FileDir(saveFileName))
	}

	needFetch := !fileExist

	if !needFetch {
		needFetch = !pkTime.CheckSamePeriod("1d", time.Unix(lastAccessTime, 0))
	}

	if needFetch {
		retList, retE = fetchExchanges()

		if retE == nil {
			if openCtpMarketHost == "" {
				retList = appendCustomExchange2(retList)
			}
			jsonData, err := json.Marshal(retList)
			if err == nil {
				retE = ioutil.WriteFile(saveFileName, jsonData, 0666)
			} else {
				retE = err
			}
		}
	} else {
		fileData, e := ioutil.ReadFile(saveFileName)
		if e != nil {
			retE = e
			return
		}
		retE = json.Unmarshal(fileData, &retList)
		if retE != nil {
			return
		}
	}

	return

}

func appendCustomExchange2(exchangesList []Exchange2Items) []Exchange2Items {
	if exchangesList == nil {
		exchangesList = make([]Exchange2Items, 0)
	}

	var custom Exchange2Items
	custom.Id = Exchange_CustomFutures
	custom.Name = "自选"
	custom.Country = Country_CN
	custom.ExchangeType = ExchangeType_Futures
	custom.KPeriods = kline.GetValidPeriods()

	exchangesList = append([]Exchange2Items{custom}, exchangesList...)

	return exchangesList
}

func GetSymbolList(foundExchange *Exchange2Items) (retList []map[string]string, retE error) {

	exchangeId := foundExchange.Id

	saveFileName := profile.ProfileReadString("symbols_save_path_fmt")

	saveFileName = fmt.Sprintf(saveFileName, exchangeId)

	lastAccessTime, fileExist := pkFile.FileLastWriteTime(saveFileName)

	if !fileExist {
		pkFile.CreateDir(pkFile.FileDir(saveFileName))
	}

	needFetch := !fileExist

	if !needFetch {
		needFetch = !pkTime.CheckSamePeriod("1d", time.Unix(lastAccessTime, 0))
	}

	if needFetch {
		retList, retE = fetchExchangeSymbols(foundExchange)

		if retE == nil {
			jsonData, err := json.Marshal(retList)
			if err == nil {
				retE = ioutil.WriteFile(saveFileName, jsonData, 0666)
			} else {
				retE = err
			}
		}
	} else {
		fileData, e := ioutil.ReadFile(saveFileName)
		if e != nil {
			retE = e
			return
		}
		retE = json.Unmarshal(fileData, &retList)
		if retE != nil {
			return
		}
	}

	return

}
