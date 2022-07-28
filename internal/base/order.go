package base

type Order struct {
	Suborder        Suborder `json:"suborder"`
	EmployeeorderId int      `json:"employeeorderId"`
}

type Suborder struct {
	Id              int    `json:"id"`
	Label           string `json:"label"`
	CommentRequired bool   `json:"commentRequired"`
}
