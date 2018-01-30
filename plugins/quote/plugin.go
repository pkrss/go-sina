package quote

import (
	"log"
	"plugin"
)

func startPluginMode(p func(channel string, data string)) (ret func(string)) {
	//打开动态库
	pdll, err := plugin.Open("plugins/bin/sina.so")
	if err != nil {
		log.Print(err)
		return
	}

	notifyStkChanged, err := pdll.Lookup("NotifyStkChanged")
	if err != nil {
		log.Print(err.Error())
		return
	}
	ret = notifyStkChanged.(func(string))

	setOnPublishChannel, err := pdll.Lookup("SetOnPublishChannel")
	if err != nil {
		log.Print(err.Error())
		return
	}
	setOnPublishChannel.(func(interface{}))(p)

	startAsyncFetchData, err := pdll.Lookup("StartAsyncFetchData")
	if err != nil {
		log.Print(err.Error())
		return
	}
	go startAsyncFetchData.(func())()

	return
}
