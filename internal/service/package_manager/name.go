package package_manager

import (
	"errors"
)

// ErrInvalidPackageName is returned for any package name other than the dv packages the updater manages.
var ErrInvalidPackageName = errors.New("invalid package name")

// allowedPackages are the only package names passed to sudo yum/apt. Any other value could be read by the
// package manager as an option ("-o..."), a URL or a local file ("x.rpm"), which runs its scriptlets as root.
var allowedPackages = map[string]struct{}{
	"dv-merchant":   {},
	"dv-processing": {},
	"dv-updater":    {},
}

func validatePackageName(name string) error {
	if _, ok := allowedPackages[name]; !ok {
		return ErrInvalidPackageName
	}
	return nil
}
