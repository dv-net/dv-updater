package app

import (
	"context"
	"errors" //nolint:goimports
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"
	"time" //nolint:goimports
)

func initTickers(ctx context.Context, s *service.Services, l logger.Logger, conf *config.AutoUpdateConfig) error {
	if s.PackageManager != nil {
		go autoUpdatePackages(ctx, s, l)

		go func() {
			ticker := time.NewTicker(time.Second * 10)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := SelfUpdate(ctx, conf, s, l); err != nil {
						l.Error("self update ticker failed", err)
					}
				}
			}
		}()
	}

	return nil
}

func autoUpdatePackages(ctx context.Context, s *service.Services, l logger.Logger) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.PackageManager.UpdateRepository(ctx); err != nil {
				l.Error("Repository update failed: %v", err)
			}
		}
	}
}

func SelfUpdate(ctx context.Context, conf *config.AutoUpdateConfig, s *service.Services, l logger.Logger) error {
	if !conf.Enabled {
		return errors.New("auto-update is disabled")
	}

	updates, err := s.PackageManager.CheckForUpdates(ctx, service.DVUpdaterServiceName)
	if err != nil {
		l.Error("self update new version check", err)
		return err
	}

	if updates.AvailableVersion != "" && updates.InstalledVersion != updates.AvailableVersion {
		if err = s.PackageManager.UpgradePackage(ctx, service.DVUpdaterServiceName); err != nil {
			l.Error("self update upgrade failed", err)
			return err
		}
	}

	return nil
}
