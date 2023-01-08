package models

import (
	"os"
	"testing"
)

var TestDir = "../test/modles"

func TestMain(m *testing.M) {
	// _ = os.RemoveAll(TestDir)
	os.MkdirAll(TestDir, 0775)
	m.Run()
}
