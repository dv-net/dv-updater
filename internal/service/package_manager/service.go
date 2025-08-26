package package_manager

import "context"

type PackageManager interface {
	GetInstalledPackage(ctx context.Context, packageName string) (Package, error)
	CheckForUpdates(ctx context.Context, packageName string) (Package, error)
	UpgradePackage(ctx context.Context, packageName string) error
	UpdateRepository(ctx context.Context) error
}

type Package struct {
	Name             string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	AvailableVersion string `json:"available_version"`
	NeedForUpdate    bool   `json:"need_for_update"`
}
