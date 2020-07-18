package support

import "os"

var dbFlag map[string]bool
var logFilePtr *os.File

func init() {
	dbFlag = make(map[string]bool)
	logFilePtr = os.Stdout
}

func SetSupport(dbf map[string]bool, lfp *os.File) {
	dbFlag = dbf
	logFilePtr = lfp
}
