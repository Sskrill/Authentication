package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
)

var (
	ErrEmplNotFound        = errors.New("empl not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type Employees struct {
	db *sql.DB
}

func NewEmpls(db *sql.DB) *Employees {
	return &Employees{db}
}

func (r *Employees) Create(ctx context.Context, empl domain.Employee) error {
	_, err := r.db.Exec("INSERT INTO employee (name, age, job) values ($1, $2, $3)",
		empl.Name, empl.Age, empl.Job)

	return err
}

func (r *Employees) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	var empl domain.Employee
	err := r.db.QueryRow("SELECT id, name, age, job FROM employee WHERE id=$1", id).
		Scan(&empl.Id, &empl.Name, &empl.Age, &empl.Job)
	if err == sql.ErrNoRows {
		return empl, ErrEmplNotFound
	}

	return empl, err
}

func (r *Employees) GetAll(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.Query("SELECT id, name, age, job FROM employee")
	if err != nil {
		return nil, err
	}

	empls := make([]domain.Employee, 0)
	for rows.Next() {
		var empl domain.Employee
		if err := rows.Scan(&empl.Id, &empl.Name, &empl.Age, &empl.Job); err != nil {
			return nil, err
		}

		empls = append(empls, empl)
	}

	return empls, rows.Err()
}

func (r *Employees) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec("DELETE FROM employee WHERE id=$1", id)

	return err
}

func (e *Employees) Update(ctx context.Context, id int64, empl domain.UpdateEmployee) error {
	var count int
	err := e.db.QueryRow("SELECT COUNT(*) FROM employee WHERE id=$1", id).Scan(&count)
	if err != nil || count == 0 {
		return ErrEmplNotFound
	}

	_, err = e.db.Exec("UPDATE employee SET name=$1,age=$2,job=$3", empl.Name, empl.Age, empl.Job)

	return err
}
