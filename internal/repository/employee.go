package repos

import (
	"database/sql"

	"github.com/Sskrill/Authentication.git/internal/domain"
)

type Employees struct {
	DB *sql.DB
}

func NewEmployees(db *sql.DB) *Employees {
	return &Employees{DB: db}
}

func (e *Employees) Create(empl domain.Employee) error {
	_, err := e.DB.Exec("INSERT INTO employee (name,age,job) VALUES ($1,$2,$3)", empl.Name, empl.Age, empl.Job)
	return err
}

func (e *Employees) Update(id int, empl domain.UpdateEmployee) error {
	var count int
	err := e.DB.QueryRow("SELECT COUNT(*) FROM employee WHERE id=$1", id).Scan(&count)
	if err != nil || count == 0 {
		return domain.ErrEmplNotFound
	}

	_, err = e.DB.Exec("UPDATE employee SET name=$1,age=$2,job=$3", empl.Name, empl.Age, empl.Job)

	if err != nil {
		return err
	}
	return nil
}

func (e *Employees) Get(id int) (domain.Employee, error) {
	var empl domain.Employee
	err := e.DB.QueryRow("SELECT id,name,age,job FROM employee WHERE id=$1", id).Scan(&empl.Id, &empl.Name, &empl.Age, &empl.Job)
	if err == sql.ErrNoRows {
		return empl, domain.ErrEmplNotFound
	}
	return empl, err
}
func (e *Employees) Delete(id int) error {
	_, err := e.DB.Exec("DELETE FROM employee WHERE id=$1", id)
	return err
}
func (e *Employees) GetAll() ([]domain.Employee, error) {
	rows, err := e.DB.Query("SELECT id,name,age,job FROM employee")
	if err != nil {
		return nil, err
	}

	empls := make([]domain.Employee, 0)
	var empl domain.Employee
	for rows.Next() {
		err := rows.Scan(&empl.Id, &empl.Name, &empl.Age, &empl.Job)
		if err != nil {
			return nil, err
		}
		empls = append(empls, empl)
	}
	return empls, rows.Err()
}
