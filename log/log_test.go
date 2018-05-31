/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */
package log

import (
	"fmt"
	"testing"
)

func Test_log(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic error %v\n", err)
		}
	}()
	NewLogger("debug", false)
	Debug("debug log")
	Info("info log")
	Fatal("fatal log")
	Error("error log")
	Warn("warn log")
	Panic("panic log")
}
