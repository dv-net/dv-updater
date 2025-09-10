package systeminfo

type Service struct {
	appVersion string
	appCommit  string
}

type InfoResponse struct {
	AppVersion string `json:"app_version"`
	AppCommit  string `json:"app_commit"`
}

func NewService(appVersion, appCommit string) *Service {
	return &Service{
		appVersion: appVersion,
		appCommit:  appCommit,
	}
}

func (o *Service) GetAppVersion() string {
	return o.appVersion
}

func (o *Service) GetAppCommit() string {
	return o.appCommit
}

func (o *Service) GetSystemInfo() *InfoResponse {
	return &InfoResponse{
		AppVersion: o.appVersion,
		AppCommit:  o.appCommit,
	}
}
