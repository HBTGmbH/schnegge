package base

import (
	"io/ioutil"
	"log"
	"os"
)

var Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func DisableLogs() {
	Log = log.New(ioutil.Discard, "", log.Ldate)
}

func EnableLogs() {
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}
