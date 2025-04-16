package logging

import "log"

var Debug bool

// DebugLog prints debug messages if Debug is enabled.
func DebugLog(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}
