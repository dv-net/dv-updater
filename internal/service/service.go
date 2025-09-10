package service

import (
	"errors"
	"fmt"

	"github.com/dv-net/dv-updater/internal/distro"
	"github.com/dv-net/dv-updater/internal/service/package_manager"
	systeminfo "github.com/dv-net/dv-updater/internal/service/system_info"
	"github.com/dv-net/dv-updater/pkg/logger"
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
