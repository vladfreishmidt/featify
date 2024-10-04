package models

import (
	"database/sql"
	"errors"
	"time"
)

type Project struct {
	ID          int
	Name        string
	Description string
	Created     time.Time
}

type ProjectModel struct {
	DB *sql.DB
}

func (m *ProjectModel) Insert(name string, description string) (int, error) {
	stmt := `INSERT INTO projects (name, description, created)
			VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, name, description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *ProjectModel) Get(id int) (*Project, error) {
	stmt := `SELECT id, name, description, created FROM projects WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	p := &Project{}

	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *ProjectModel) Latest() ([]*Project, error) {
	stmt := `SELECT id, name, description, created FROM projects
		ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		p := &Project{}

		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Created)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (m *ProjectModel) GetAll() ([]*Project, error) {
	stmt := `SELECT id, name, description, created FROM projects
		ORDER BY id DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		p := &Project{}

		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Created)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
