package main

/*
#include <stdlib.h>

typedef  void (*CFuncOnPublishChannel)(char*,char*);

static inline void callCFuncOnPublishChannel(CFuncOnPublishChannel f, char* p1, char* p2){
	if(!f || !p1 || !p2)
		return;
	(*f)(p1,p2);
}
*/
import "C"

import (
	"unsafe"

	"github.com/pkrss/gosina/quote"
)

//export SoSetOnPublishChannel
func SoSetOnPublishChannel(f func(string, string)) {
	quote.SetOnPublishChannel(f)
}

var cfuncOnPublishChannel C.CFuncOnPublishChannel

func GoBridgeSetOnPublishChannel(channel string, msg string) {
	if cfuncOnPublishChannel == nil {
		return
	}

	p1 := C.CString(channel)
	p2 := C.CString(msg)
	C.callCFuncOnPublishChannel(cfuncOnPublishChannel, p1, p2)
	C.free(unsafe.Pointer(p1))
	C.free(unsafe.Pointer(p2))
}

//export SoSetOnPublishChannelC
func SoSetOnPublishChannelC(f C.CFuncOnPublishChannel) {
	cfuncOnPublishChannel = f
	quote.SetOnPublishChannel(GoBridgeSetOnPublishChannel)
}

//export SoNotifyStkChanged
func SoNotifyStkChanged(stkStr string) {
	quote.NotifyStkChanged(stkStr)
}

//export SoNotifyStkChangedC
func SoNotifyStkChangedC(stkStr *C.char) {
	quote.NotifyStkChanged(C.GoString(stkStr))
}

//export SoStart
func SoStart() {
	quote.Start()
}

//export SoRunOnce
func SoRunOnce() {
	quote.RunOnce()
}
