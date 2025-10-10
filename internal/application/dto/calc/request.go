package calc

type SendFormToCalcRequest struct {
	FormID uint `json:"formID" validate:"required"`
	CalcID uint `json:"calcID" validate:"required"`
}
