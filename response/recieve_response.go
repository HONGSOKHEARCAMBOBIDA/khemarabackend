package response

type RecieveResponse struct {
	ID                    int                     `json:"id"`
	Code                  string                  `json:"code"`
	BranchID              int                     `json:"branch_id"`
	BranchName            string                  `json:"branch_name"`
	RecieveDate           string                  `json:"receive_date" gorm:"column:receive_date"`
	TotalRecieve          float64                 `json:"total_receive" gorm:"column:total_receive"`
	CurrencyID            int                     `json:"currency_id"`
	CurrencyCode          string                  `json:"currency_code"`
	Note                  string                  `json:"note"`
	RecieveByID           int                     `json:"receive_by" gorm:"column:receive_by"`
	RecieveByName         string                  `json:"recieve_by_name"`
	RecieveDetailResponse []RecieveDetailResponse `json:"recieve_detaild" gorm:"-"`
}
