package distro

type ReleaseDetails = map[string]string

type LinuxDistro struct {
	Name       string         `json:"name"`
	ID         string         `json:"id"`
	Version    string         `json:"version"`
	LsbRelease ReleaseDetails `json:"lsb_release"`
	OsRelease  ReleaseDetails `json:"os_release"`
}
