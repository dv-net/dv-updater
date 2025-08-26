package request

type UpdatePackageRequest struct {
	Name string `json:"name" validate:"required,oneof=dv-processing dv-merchant"`
}
