package repl

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/yechentide/dstm/global"
)

func clearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printError(message string) {
	output := "[ERROR] " + message
	fmt.Println(global.ErrorStyle.Render(output))
}

func printWarn(message string) {
	output := "[WARN] " + message
	fmt.Println(global.WarnStyle.Render(output))
}

func printInfo(message string) {
	output := "[INFO] " + message
	fmt.Println(global.InfoStyle.Render(output))
}
