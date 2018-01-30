@echo off

cd /d %GOPATH%\src\github.com\pkrss\go-sina\so\quote

go build -buildmode=c-archive -o quote.a

gcc -shared -pthread -o libgoSinaQuote.so c\dummy.c quote.a -lWinMM -lWS2_32

copy /y libgoSinaQuote.so %GOPATH%\bin\