package service

import (
	"mysql/config"
	"mysql/helper"
	"mysql/response"

	"gorm.io/gorm"
)

type RecieveService interface {
	GetRecieve(id int) ([]response.RecieveResponse, error)
}

type recieveservice struct {
	db *gorm.DB
}

func NewRecieveService() RecieveService {
	return &recieveservice{
		db: config.DB,
	}
}

func (s *recieveservice) GetRecieve(id int) ([]response.RecieveResponse, error) {
	var recieve []response.RecieveResponse
	query := s.db.Table("recieves r").
		Select(`
			r.id AS id,
			r.code AS code,
			b.id AS branch_id,
			b.name AS branch_name,
			r.receive_date AS receive_date,
			r.total_receive AS total_receive,
			c.id AS currency_id,
			c.code AS currency_code,
			r.note AS note,
			r.receive_by AS receive_by,
			e.name_kh AS recieve_by_name
		`).
		Joins("LEFT JOIN branches b ON b.id = r.branch_id").
		Joins("LEFT JOIN currencies c ON c.id = r.currency_id").
		Joins("LEFT JOIN users u ON u.id = r.receive_by").
		Joins("LEFT JOIN employees e ON e.id = u.employee_id").
		Joins("LEFT JOIN loans l ON l.id = r.loan_id").
		Where("l.id = ?", id)

	if err := query.Scan(&recieve).Error; err != nil {
		return nil, err
	}

	for i := range recieve {
		recieve[i].RecieveDate = helper.FormatDate(recieve[i].RecieveDate)
	}

	if len(recieve) > 0 {
		recieveIDs := make([]int, len(recieve))
		for i, recieves := range recieve {
			recieveIDs[i] = recieves.ID

		}

		var recievedetail []response.RecieveDetailResponse
		if err := s.db.Table("recieve_details rd").
			Select(`
				rd.id AS recieve_detail_id,
				rd.receive_id AS receive_id,
				rd.principal AS principal,
				rd.rate AS rate,
				rd.income AS income,
			p.payroll_date AS payroll_date,
			pt.name AS payroll_type
			`).
			Joins("LEFT JOIN recieves r ON r.id = rd.receive_id").
			Joins("LEFT JOIN payrolls p ON p.id = r.payroll_id").
			Joins("LEFT JOIN pay_roll_types pt ON pt.id = p.payroll_type_id").
			Where("rd.receive_id IN ?", recieveIDs).
			Scan(&recievedetail).Error; err != nil {
			return nil, err
		}
		recievedetailmap := make(map[int][]response.RecieveDetailResponse)
		for j := range recievedetail {
			recievedetail[j].PayrollDate = helper.FormatDate(recievedetail[j].PayrollDate)
			recievedetailmap[recievedetail[j].RecieveID] = append(recievedetailmap[recievedetail[j].RecieveID], recievedetail[j])
		}

		for i := range recieve {
			recieve[i].RecieveDetailResponse = recievedetailmap[recieve[i].ID]
			if recieve[i].RecieveDetailResponse == nil {
				recieve[i].RecieveDetailResponse = []response.RecieveDetailResponse{}
			}
		}
	}
	return recieve, nil
}
