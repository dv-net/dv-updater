package package_manager

import "testing"

func TestValidatePackageName(t *testing.T) {
	valid := []string{"dv-merchant", "dv-processing", "dv-updater"}
	for _, name := range valid {
		if err := validatePackageName(name); err != nil {
			t.Errorf("validatePackageName(%q) = %v, want nil", name, err)
		}
	}

	invalid := []string{
		"",
		"-oAPT::Install::Post-Invoke-Success::=id",
		"--setopt=gpgcheck=0",
		"http://10.0.1.23:8899/dvm99.rpm",
		"https://d.uguu.se/x.rpm",
		"/tmp/evil.rpm",
		"./evil.rpm",
		"dv-merchant evil",
		"dv-merchant;id",
		"DV-Merchant",
		"dv-merchant.rpm",
		"dv.deb",
		"dv-merchant ",
		"libfoo2.1",
	}
	for _, name := range invalid {
		if err := validatePackageName(name); err == nil {
			t.Errorf("validatePackageName(%q) = nil, want error", name)
		}
	}
}
