package fp

import (
	"fmt"
	"os"
)

func panicError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: "+format, args...)
}
func logWarn(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "WARNING: "+format, args...)
}

func logInfo(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "INFO: "+format, args...)
}

func logDebug(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "DEBUG: "+format, args...)
}
