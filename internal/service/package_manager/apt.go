package package_manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/dv-net/dv-updater/pkg/logger" //nolint:goimports
	"github.com/dv-net/dv-updater/pkg/retry"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time" //nolint:goimports
)

type AptManager struct {
	logger     logger.Logger
	binaryPath string
}

const repo = "sources.list.d/dvnet.list"
const flagForceUpdate = "--force-confold"

var _ PackageManager = (*AptManager)(nil)

func NewAptManager(l logger.Logger) (*AptManager, error) {
	binaryPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get binary path: %w", err)
	}

	return &AptManager{
		logger:     l,
		binaryPath: binaryPath,
	}, nil
}

func (a *AptManager) GetInstalledPackage(ctx context.Context, packageName string) (Package, error) {
	out, err := exec.CommandContext(ctx, "sudo", "apt", "list", "--installed", packageName).Output()
	if err != nil {
		a.logger.Error("Failed to get installed packages %s", err)
		return Package{}, err
	}

	return a.parseAptOutput(out, packageName)
}

func (a *AptManager) CheckForUpdates(ctx context.Context, packageName string) (Package, error) {
	out, err := exec.CommandContext(ctx, "sudo", "apt", "list", "--upgradable", packageName).Output()
	if err != nil {
		a.logger.Error("Failed to check for updates: %s", err)
		return Package{}, ErrNothingToUpdate
	}

	return a.parseAptOutput(out, packageName)
}

func (a *AptManager) UpgradePackage(ctx context.Context, packageName string) error {
	a.logger.Error("Attempting to upgrade package", nil, "pkg", packageName)
	err := a.runAptCommandWithSpinLock(ctx, "install", "-o", "Dpkg::Options::="+flagForceUpdate, "-y", "--only-upgrade", packageName)
	if err != nil {
		a.logger.Error("Failed to upgrade package", err, "pkg", packageName)
		pkg, checkErr := a.CheckForUpdates(ctx, packageName)
		if checkErr != nil || pkg.NeedForUpdate {
			a.logger.Error("Package still needs update after dpkg configure", nil, "pkg", packageName, "checkErr", checkErr)
			return fmt.Errorf("failed to update package %s: still needs update", packageName)
		}
		a.logger.Info("Package updated successfully", "pkg", packageName)
	} else {
		a.logger.Info("Package updated successfully", "pkg", packageName)
	}

	if packageName == "dv-updater" {
		if err := a.selfUpdate(ctx); err != nil {
			a.logger.Error("Failed to self-update", err)
			return fmt.Errorf("self-update failed: %w", err)
		}
	}

	a.logger.Info("UpgradePackage completed", "pkg", packageName)
	return nil
}

func (a *AptManager) UpdateRepository(ctx context.Context) error {
	a.logger.Info("start Updating repository")
	out, err := exec.CommandContext(ctx, "sudo", "apt", "update", "-o", "Dir::Etc::sourcelist="+repo).CombinedOutput()
	if err != nil {
		a.logger.Error("Failed to update package: %v", err, "out", string(out))
		return errors.New("failed to update package list")
	}
	a.logger.Info("Package list updated successfully")
	a.logger.Debug("Output: %s", string(out))

	return nil
}

func (a *AptManager) parseAptOutput(out []byte, packageName string) (Package, error) {
	lines := strings.Split(string(out), "\n")

	var installedVersion, availableVersion string

	regexPattern := fmt.Sprintf(`^%s/(?:unknown(?:,now)?)?\s*(\S+)\s+amd64\s+\[(?:installed|upgradable from:\s*(\S+))\]`, regexp.QuoteMeta(packageName))
	re := regexp.MustCompile(regexPattern)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			if len(matches) > 2 && matches[2] != "" {
				availableVersion = matches[1]
				installedVersion = matches[2]
			} else {
				installedVersion = matches[1]
				availableVersion = ""
			}
			break
		}
	}

	if installedVersion == "" && availableVersion == "" {
		return Package{}, fmt.Errorf("package %s not found or no version information in provided output", packageName)
	}

	return Package{
		Name:             packageName,
		InstalledVersion: installedVersion,
		AvailableVersion: availableVersion,
		NeedForUpdate:    installedVersion != availableVersion,
	}, nil
}

func (a *AptManager) runAptCommandWithSpinLock(ctx context.Context, args ...string) error {
	return retry.New(
		retry.WithPolicy(retry.PolicyLinear),
		retry.WithDelay(5*time.Second),
		retry.WithMaxAttempts(5),
	).Do(func() error {
		if a.isDpkgLocked(ctx) {
			a.logger.Error("dpkg lock detected, retrying in 5s", nil)
			return retry.ErrRetry
		}

		cmdArgs := append([]string{"sudo", "apt"}, args...)
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
		out, err := cmd.CombinedOutput()
		if err != nil {
			if a.isLockError(err) {
				out, errDpkg := exec.CommandContext(ctx, "sudo", "dpkg", "--configure", "-a").CombinedOutput()
				if errDpkg != nil {
					a.logger.Error("Failed to configure dpkg", errDpkg, "pkg", "packageName", "out", string(out))
					return fmt.Errorf("failed to update package %s: apt error: %w, dpkg error: %w", "packageName", err, errDpkg)
				}

				a.logger.Debug("dpkg configured", "pkg", "packageName", "out", string(out))
				a.logger.Debug("dpkg lock detected during command, retrying in 5s")
				return retry.ErrRetry
			}

			return fmt.Errorf("apt command failed: %w, output: %s", err, string(out))
		}

		a.logger.Error("apt command succeeded", nil, "output", string(out))
		return nil
	})
}

func (a *AptManager) isDpkgLocked(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "fuser", "/var/lib/dpkg/lock-frontend")
	output, err := cmd.CombinedOutput()
	return err == nil && len(output) > 0
}

func (a *AptManager) isLockError(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode() == 100
	}
	return false
}

func (a *AptManager) selfUpdate(ctx context.Context) error {
	a.logger.Debug("checking binary in use for self-update")

	cleanPath := filepath.Clean(a.binaryPath)
	if info, err := os.Stat(cleanPath); err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		a.logger.Error("invalid binary path", err, "path", cleanPath)
		return fmt.Errorf("invalid or non-executable binary path: %s", cleanPath)
	}

	if err := retry.New(
		retry.WithPolicy(retry.PolicyLinear),
		retry.WithDelay(2*time.Second),
		retry.WithMaxAttempts(5),
	).Do(func() error {
		cmd := exec.CommandContext(ctx, "fuser", cleanPath)
		if output, err := cmd.CombinedOutput(); err == nil && len(output) > 0 {
			a.logger.Debug("binary is locked. retry", "out", string(output))
			return retry.ErrRetry
		}

		return nil
	}); err != nil {
		return err
	}

	a.logger.Error("systemd restart initialization...", nil)
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		a.logger.Error("failed to send SIGTERM", err)
	}

	a.logger.Info("process terminating")
	time.Sleep(3 * time.Second)
	os.Exit(0)

	return nil
}
