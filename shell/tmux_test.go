package shell_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yechentide/dstm/shell"
)

func TestTmuxListAndCreateAndKill(t *testing.T) {
	var err error
	var exists bool

	sessions, err := shell.ListTmuxSessions()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(sessions))

	exists, err = shell.HasTmuxSession("test001")
	assert.Nil(t, err)
	assert.False(t, exists)
	exists, err = shell.HasTmuxSession("test002")
	assert.Nil(t, err)
	assert.False(t, exists)

	err = shell.CreateTmuxSession("test001", "")
	assert.Nil(t, err)
	err = shell.CreateTmuxSession("test002", "")
	assert.Nil(t, err)

	exists, err = shell.HasTmuxSession("test001")
	assert.Nil(t, err)
	assert.True(t, exists)
	exists, err = shell.HasTmuxSession("test002")
	assert.Nil(t, err)
	assert.True(t, exists)

	sessions, err = shell.ListTmuxSessions()
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"test001", "test002"}, sessions)

	err = shell.KillTmuxSession("test001")
	assert.Nil(t, err)
	err = shell.KillTmuxSession("test002")
	assert.Nil(t, err)
}

func TestTmuxSendMessageAndCapture(t *testing.T) {
	sessionName := "test001"
	cmd := "uname"
	err := shell.CreateTmuxSession(sessionName, "")
	assert.Nil(t, err)

	err = shell.SendMessageToTmuxSession(sessionName, cmd, true)
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	log, err := shell.CaptureTmuxSessionOutput(sessionName, true)
	assert.Nil(t, err)

	lines := strings.Split(log, "\n")
	assert.Equal(t, 4, len(lines))
	assert.Contains(t, lines[1], cmd)
	assert.Contains(t, lines[1], cmd)
	assert.Equal(t, lines[2], "Linux")

	err = shell.KillTmuxSession(sessionName)
	assert.Nil(t, err)
}
