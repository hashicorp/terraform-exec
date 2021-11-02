package testutil

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}
	return log.New(ioutil.Discard, "", 0)
}
