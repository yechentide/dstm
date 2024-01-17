package env

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"strings"

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
	return []string{"lib32gcc-s1", "lua5.3", "curl", "tmux"}
}

func (d *DebianHelper) IsInstalled(packages []string) (map[string]bool, error) {
	installed := make(map[string]bool, len(packages))

	tmpFilePath := "/tmp/is_installed.txt"
	listInstalledCmd := "dpkg-query -l | awk '{print $2}' > " + tmpFilePath
	err := exec.Command("bash", "-c", listInstalledCmd).Run()
	if err != nil {
		return installed, err
	}

	for _, pkg := range packages {
		output, err := shell.ExecuteAndGetOutput("grep", "^"+pkg+"$", tmpFilePath)
		if err != nil {
			return installed, errors.New(output)
		}
		installed[pkg] = output == pkg
	}
	return installed, nil
}

func (d *DebianHelper) InstallPackages(packages []string) error {
	isRoot := d.IsRoot()
	if !isRoot && !d.HasSudoPermmission() {
		slog.Error("You must have sudo permission to install packages")
		os.Exit(1)
	}
	cmd := "apt install -y " + strings.Join(packages, " ")
	if isRoot {
		slog.Debug("Running command: " + cmd)
		return exec.Command("bash", "-c", cmd).Run()
	} else {
		slog.Debug("Running command: sudo " + cmd)
		return exec.Command("sudo", "bash", "-c", cmd).Run()
	}
}

func (d *DebianHelper) InstallAllRequired() error {
	pkgs := d.Dependencies()
	slog.Debug("Installing dependencies: " + strings.Join(pkgs, " "))
	return d.InstallPackages(pkgs)
}
