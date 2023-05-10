package req

type CheckOutReq struct {
	Sn string `json:"sn" label:"sn" validate:"required,numeric"`
}
