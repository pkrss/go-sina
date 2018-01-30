cd /d $GOROOT/src/github.com/pkrss/go-sina/so/quote
go build -buildmode=c-archive -o quote.a
gcc -shared -pthread -o libgoSinaQuote.so c/dummy.c quote.a
/bin/cp -f libgoSinaQuote.so $GOROOT/bin/