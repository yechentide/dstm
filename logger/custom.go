package logger

import (
	"fmt"

	"github.com/yechentide/dstm/global"
)

func OK(msg string) {
	prefix := ""
	if showPrefix {
		prefix = "[OK]"
	}
	if showColor {
		prefix = global.OkStyle.Render(prefix)
	}
	fmt.Printf("%s %s", prefix, msg)
}
