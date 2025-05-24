package dtos

type EmployeeV1Response struct {
	EmployeeID int64  `json:"employee_id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`

	PositionID int64   `json:"position_id"`
	Position   string  `json:"position"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
	StartDate  string  `json:"start_date"`
}
