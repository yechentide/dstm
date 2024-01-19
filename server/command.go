package server

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yechentide/dstm/shell"
)

func GetPlayerNumber(sessionName string) (int, error) {
	err := shell.SendMessageToTmuxSession(sessionName, "c_getnumplayers()", true)
	if err != nil {
		return -1, err
	}
	time.Sleep(100 * time.Millisecond)
	log, err := shell.CaptureTmuxSessionOutput(sessionName, true)
	if err != nil {
		return -1, err
	}
	lastNewLineIdx := strings.LastIndexByte(log, '\n')
	cmdOutput := log[lastNewLineIdx+1:]

	pattern := `^\[\d+:\d{2}:\d{2}\]:\s(\d+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(cmdOutput)
	if matches == nil || len(matches) != 2 {
		return 0, errors.New("cannot get player number")
	}
	return strconv.Atoi(matches[1])
}
