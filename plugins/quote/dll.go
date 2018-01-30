package quote

/*
#include <stdlib.h>

typedef  void (*CFuncOnPublishChannel)(char*,char*);

extern void goCFuncOnPublishChannel(char* p1, char* p2);
*/
import "C"

import (
	"log"
	"syscall"
	"unsafe"
)

var funOnPublishChannel func(channel string, data string)

//export goCFuncOnPublishChannel
func goCFuncOnPublishChannel(p1 *C.char, p2 *C.char) {
	if funOnPublishChannel == nil {
		return
	}

	funOnPublishChannel(C.GoString(p1), C.GoString(p2))
}

var pdll *syscall.DLL

// not worked, becaused redis.PublishChannel can not passed in .so
func startDllMode(p func(channel string, data string)) (ret func(string)) {
	funOnPublishChannel = p

	var err error

	pdll, err = syscall.LoadDLL("libgoSinaQuote.so")
	if err != nil {
		log.Print(err)
		return
	}

	notifyStkChanged, err := pdll.FindProc("SoNotifyStkChangedC")
	if err != nil {
		log.Print(err.Error())
		return
	}
	ret = func(s string) {
		p := C.CString(s)
		notifyStkChanged.Call(uintptr(unsafe.Pointer(p)))
		C.free(unsafe.Pointer(p))
	}

	setOnPublishChannel, err := pdll.FindProc("SoSetOnPublishChannelC")
	if err != nil {
		log.Print(err.Error())
		return
	}
	setOnPublishChannel.Call(uintptr(unsafe.Pointer(C.goCFuncOnPublishChannel)))

	go asyncStart()

	return
}

func asyncStart() {
	startAsyncFetchData, err := pdll.FindProc("SoStart")
	if err != nil {
		log.Print(err.Error())
		return
	}
	startAsyncFetchData.Call()
}
