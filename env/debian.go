package env

import (
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/yechentide/dstm/global"
	"github.com/yechentide/dstm/shell"
	"golang.org/x/exp/slices"
)

type DebianHelper struct{}

func (d *DebianHelper) IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		slog.Error("Failed to get current user", err)
		os.Exit(1)
	}
	return currentUser.Username == "root"
}

func (d *DebianHelper) HasSudoPermmission() bool {
	output, err := shell.ExecuteAndGetOutput("groups")
	if err != nil {
		slog.Error("Failed to get groups", err)
		os.Exit(1)
	}
	groups := strings.Split(output, " ")
	return slices.Contains(groups, "sudo")
}

func (d *DebianHelper) Dependencies() []string {
	return []string{"lib32gcc-s1", "lua5.1", "curl", "tmux"}
}

func (d *DebianHelper) IsInstalled(packages []string) (map[string]bool, error) {
	if len(packages) == 0 {
		return map[string]bool{}, nil
	}

	installed := make(map[string]bool, len(packages))

	tmpFilePath := global.FILE_PATH_PKG_INSTALL_CHECK
	listInstalledCmd := "dpkg-query -l | awk '{print $2}' > " + tmpFilePath
	err := exec.Command("bash", "-c", listInstalledCmd).Run()
	if err != nil {
		return installed, err
	}

	for _, pkg := range packages {
		output, err := shell.ExecuteAndGetOutput("grep", "^"+pkg+"$", tmpFilePath)
		if err == nil {
			installed[pkg] = output == pkg
		} else {
			installed[pkg] = false
		}
	}
	return installed, nil
}

func (d *DebianHelper) IsTerminalMultiplexerReady() (bool, error) {
	result, err := d.IsInstalled([]string{"tmux"})
	if err != nil {
		return false, err
	}
	if result["tmux"] {
		return true, nil
	}
	return false, nil
}

func (d *DebianHelper) InstallPackages(packages []string, password string) error {
	if !d.IsRoot() && !d.HasSudoPermmission() {
		slog.Error("You must have sudo permission to install packages")
		os.Exit(1)
	}
	cmdString := "apt install -y " + strings.Join(packages, " ")
	var cmd *exec.Cmd
	if d.IsRoot() {
		slog.Debug("Running command: " + cmdString)
		cmd = exec.Command("bash", "-c", cmdString)
	} else {
		slog.Debug("Running command: sudo -S " + cmdString)
		cmd = exec.Command("bash", "-c", "sudo -S "+cmdString)
		cmd.Stdin = strings.NewReader(password + "\n")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (d *DebianHelper) InstallAllRequired(password string) error {
	pkgs := d.Dependencies()
	slog.Debug("Installing dependencies: " + strings.Join(pkgs, " "))
	return d.InstallPackages(pkgs, password)
}
