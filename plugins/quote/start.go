package quote

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Start(p func(channel string, data string)) (ret func(string)) {

	if runtime.GOOS == "windows" {
		ret = startDirectMode(p)
		// Build()
		// ret = startDllMode(p)
	} else {
		Build()
		ret = startPluginMode(p)
	}
	return
}

func Build() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatalln("GOPATH environment is not set? please sure it!")
		return
	}

	if !(strings.HasSuffix(gopath, "/") || strings.HasSuffix(gopath, "\\")) {
		gopath += "/"
	}

	target := gopath + "bin/libgoSinaQuote.so"

	if runtime.GOOS == "windows" {
		gopath = strings.Replace(gopath, "/", "\\", -1)
		target = strings.Replace(target, "/", "\\", -1)
	}

	_, err := os.Stat(target)
	if err == nil {
		log.Println("skip, found so in:" + target)
		return
	}

	log.Println("build, not found so in:" + target)

	if runtime.GOOS == "windows" {
		buildSo(gopath, target)
	} else {
		buildPlugin(gopath, target)
	}
}

func buildSo(gopath string, target string) {
	cmd := gopath + "src/github.com/pkrss/go-sina/so/quote/build"
	// m := "sh"
	if runtime.GOOS == "windows" {
		cmd += ".bat"
		// m = "cmd"
	} else {
		cmd += ".sh"
	}

	c := exec.Command(cmd)

	if e := c.Run(); e != nil {
		log.Println(e.Error())
	}
}
func buildPlugin(gopath string, target string) {
	cmd := "cd " + gopath + "src/github.com/pkrss/go-sina/so/quote"
	cmd += " && go build -buildmode=plugin -o libgoSinaQuote.so"
	cmd += " && "
	if runtime.GOOS == "windows" {
		cmd += "copy /y"
	} else {
		cmd += "/bin/cp -f"
	}
	cmd += " libgoSinaQuote.so " + target

	c := exec.Command("sh", cmd)

	if e := c.Run(); e != nil {
		log.Println(e.Error())
	}
}
