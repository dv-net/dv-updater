package service

import (
	"dv-updater/internal/distro"
	"dv-updater/internal/service/package_manager"
	systeminfo "dv-updater/internal/service/system_info"
	"dv-updater/pkg/logger"
	"errors"
	"fmt"
)

const (
	DVUpdaterServiceName    string = "dv-updater"
	DVMerchantServiceName   string = "dv-merchant"
	DVProcessingServiceName string = "dv-processing"
)

func ValidateServiceName(serviceName string) error {
	switch serviceName {
	case DVMerchantServiceName, DVUpdaterServiceName, DVProcessingServiceName:
		return nil
	default:
		return fmt.Errorf("invalid service name")
	}
}

type Services struct {
	PackageManager    package_manager.PackageManager
	SystemInfoService *systeminfo.Service
}

func NewServices(l logger.Logger, dist distro.LinuxDistro, currentAppVersion, currentAppCommitHash string) (*Services, error) {
	var (
		pm  package_manager.PackageManager
		err error
	)
	switch dist.ID {
	case "debian", "ubuntu":
		pm, err = package_manager.NewAptManager(l)
		if err != nil {
			return nil, err
		}
	case "centos", "rhel":
		pm = package_manager.NewYumManager(l)
	default:
		l.Fatal("Unsupported distribution", errors.New("unsupported distro"))
	}

	return &Services{
		PackageManager:    pm,
		SystemInfoService: systeminfo.NewService(currentAppVersion, currentAppCommitHash),
	}, nil
}
