package req

type AspPayoutAuditJob struct {
	Id          int    `json:"id" validate:"required"`
	Status      string `json:"status" validate:"required"`
	OperationID int    `json:"operation_id" validate:"required"`
}
