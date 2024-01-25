package psql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
	"github.com/GOLANG-NINJA/crud-app/internal/service"
	"github.com/Sskrill/gRpc-log/proto/audit"
	"github.com/sirupsen/logrus"
)

var (
	ErrEmplNotFound        = errors.New("empl not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type Employees struct {
	db          *sql.DB
	auditClient service.Audit
}

func NewEmpls(db *sql.DB, audit service.Audit) *Employees {
	return &Employees{db: db, auditClient: audit}
}

func (r *Employees) Create(ctx context.Context, empl domain.Employee) error {
	_, err := r.db.Exec("INSERT INTO employee (name, age, job) values ($1, $2, $3)",
		empl.Name, empl.Age, empl.Job)
	if err != nil {
		return err
	}
	id, err := r.getEmplId(empl)
	if err != nil {
		return err
	}
	err = r.auditClient.SenLogReq(ctx, audit.LogItem{Action: audit.ACTION_CREATE, Entity: audit.ENTITY_EMPLOYEE, EntityID: int64(id), Timestamp: time.Now()})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignUp",
		}).Error("failed to send log request:", err)
	}
	return err
}

func (r *Employees) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	var empl domain.Employee
	err := r.db.QueryRow("SELECT id, name, age, job FROM employee WHERE id=$1", id).
		Scan(&empl.Id, &empl.Name, &empl.Age, &empl.Job)
	if err == sql.ErrNoRows {
		return empl, ErrEmplNotFound
	}
	err = r.auditClient.SenLogReq(ctx, audit.LogItem{Action: audit.ACTION_GET, Entity: audit.ENTITY_EMPLOYEE, EntityID: int64(id), Timestamp: time.Now()})

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
	if err != nil {
		return err
	}
	err = r.auditClient.SenLogReq(ctx, audit.LogItem{Action: audit.ACTION_DELETE, Entity: audit.ENTITY_EMPLOYEE, EntityID: int64(id), Timestamp: time.Now()})
	return err
}

func (r *Employees) Update(ctx context.Context, id int64, empl domain.UpdateEmployee) error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM employee WHERE id=$1", id).Scan(&count)
	if err != nil || count == 0 {
		return ErrEmplNotFound
	}

	_, err = r.db.Exec("UPDATE employee SET name=$1,age=$2,job=$3 WHERE id=$4", empl.Name, empl.Age, empl.Job, id)
	if err != nil {
		return err
	}
	err = r.auditClient.SenLogReq(ctx, audit.LogItem{Action: audit.ACTION_UPDATE, Entity: audit.ENTITY_EMPLOYEE, EntityID: int64(id), Timestamp: time.Now()})
	return err
}

func (r *Employees) getEmplId(employee domain.Employee) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM employee WHERE name=$1 AND age =$2 AND job=$3", employee.Name, employee.Age, employee.Job).
		Scan(&id)
	if err == sql.ErrNoRows {
		return id, ErrEmplNotFound
	}
	return id, nil
}
