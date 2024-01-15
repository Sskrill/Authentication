package psql

import (
	"database/sql"

	"github.com/GOLANG-NINJA/crud-app/internal/domain"
)

type Sessions struct {
	DB *sql.DB
}

func NewSession(db *sql.DB) *Sessions {
	return &Sessions{DB: db}
}

func (s *Sessions) Create(session domain.Session) error {
	_, err := s.DB.Exec("INSERT INTO session (id,user_id) VALUES ($1,$2)", session.Id, session.UserId)
	return err
}
func (s *Sessions) GetId(id string) error {
	var session domain.Session
	err := s.DB.QueryRow("SELECT session_id,id,user_id FROM session WHERE id=$1", id).Scan(&session.SessionId, &session.Id, &session.UserId)
	return err
}
func (s *Sessions) Delete(id string) error {
	_, err := s.DB.Exec("DELETE FROM session WHERE id=$1", id)
	return err
}
