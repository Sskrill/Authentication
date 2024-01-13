package domain

type Employee struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

type UpdateEmployee struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}
