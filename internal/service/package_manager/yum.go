package package_manager

import (
	"context"
	"errors"
	"github.com/dv-net/dv-updater/pkg/logger"
	"os/exec"
	"strings"
)

type YumManager struct {
	logger logger.Logger
}

var _ PackageManager = (*YumManager)(nil)

func NewYumManager(log logger.Logger) *YumManager {
	return &YumManager{logger: log}
}

func (y *YumManager) GetInstalledPackage(ctx context.Context, packageName string) (Package, error) {
	// sudo yum list installed
	out, err := exec.CommandContext(ctx, "sudo", "yum", "list", "installed", packageName).Output()
	if err != nil {
		y.logger.Error("Failed to get installed package: %v", err)
		return Package{}, errors.New("package not found")
	}

	return y.parseYumOutput(out, packageName)
}

func (y *YumManager) CheckForUpdates(ctx context.Context, packageName string) (Package, error) {
	// sudo yum --repo=dvnet list --refresh
	out, err := exec.CommandContext(ctx, "sudo", "yum", "--repo=dvnet", "list", "--refresh", packageName).Output()
	if err != nil {
		y.logger.Error("Failed to check for updates: %v", err)
		return Package{}, ErrNothingToUpdate
	}

	return y.parseYumOutput(out, packageName)
}

func (y *YumManager) UpgradePackage(ctx context.Context, packageName string) error {
	y.logger.Info("start Updating repository")
	out, err := exec.CommandContext(ctx, "sudo", "yum", "--repo", "dvnet", "update", "-y", packageName).Output()
	if err != nil {
		y.logger.Error("Failed to update package: %s", err)
		return errors.New("failed to update package")
	}

	y.logger.Info("Package %s updated successfully", packageName)
	y.logger.Debug("Output: %s", string(out))

	return nil
}

func (y *YumManager) UpdateRepository(ctx context.Context) error {
	// sudo yum --repo dvnet list available --refresh"
	out, err := exec.CommandContext(ctx, "sudo", "yum", "--repo", "dvnet", "list", "available", "--refresh").Output()

	if err != nil {
		y.logger.Error("Failed to update package: %v", err)
		return errors.New("failed to update package list")
	}

	y.logger.Info("Package list updated successfully")
	y.logger.Debug("Output: %s", string(out))
	return nil
}

func (y *YumManager) SearchPackage(ctx context.Context, packageName string) ([]string, error) {
	out, err := exec.CommandContext(ctx, "sudo", "yum", "search", packageName).Output()
	if err != nil {
		y.logger.Error("Failed to search for package: %v", err)
		return nil, err
	}

	// Разделяем вывод на строки
	lines := strings.Split(string(out), "\n")
	var results []string

	// Фильтруем строки, оставляя только названия пакетов
	for _, line := range lines {
		if strings.HasPrefix(line, " ") && strings.Contains(line, ".") {
			results = append(results, strings.Fields(line)[0])
		}
	}

	return results, nil
}

func (y *YumManager) parseYumOutput(out []byte, packageName string) (Package, error) {
	lines := strings.Split(string(out), "\n")

	var installedVersion, availableVersion string
	var parsingInstalled bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Installed Packages") {
			parsingInstalled = true
			continue
		} else if strings.HasPrefix(line, "Available Packages") {
			parsingInstalled = false
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		name := strings.Split(parts[0], ".")[0]
		version := parts[1]

		if name == packageName {
			if parsingInstalled {
				installedVersion = version
			} else {
				availableVersion = version
			}
		}
	}

	if installedVersion == "" && availableVersion == "" {
		return Package{}, errors.New("package not found")
	}

	needForUpdate := false

	if availableVersion == "" {
		availableVersion = installedVersion
	}

	if installedVersion != availableVersion {
		needForUpdate = true
	}

	return Package{
		Name:             packageName,
		InstalledVersion: installedVersion,
		AvailableVersion: availableVersion,
		NeedForUpdate:    needForUpdate,
	}, nil
}
