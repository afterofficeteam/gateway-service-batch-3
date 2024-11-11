package users

import (
	"database/sql"
	"fmt"
	model "gateway-service/models"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

type UserRepository interface {
	UserRegister(req model.Users) (*uuid.UUID, error)
	GetUserDetail(req model.Users) (*model.Users, error)
}

func (s *store) UserRegister(req model.Users) (*uuid.UUID, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userID uuid.UUID
	queryStatement := `
		INSERT INTO users(
			email,
		    username,
			role,
 			address,
 			category_preferences,
			password,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			now()
		) RETURNING id
	`

	if err := tx.QueryRow(
		queryStatement,
		req.Email,
		req.Username,
		req.Role,
		req.Address,
		pq.Array(req.CategoryPreferences),
		req.Password,
	).Scan(&userID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &userID, nil
}

func (s *store) GetUserDetail(req model.Users) (*model.Users, error) {
	queryStatement := `
		SELECT
			*
		FROM
		    users
	`

	var queryConditions []string
	if req.Email != "" {
		queryConditions = append(queryConditions, fmt.Sprintf("email = '%s'", req.Email))
	}

	if req.Id != uuid.Nil {
		queryConditions = append(queryConditions, fmt.Sprintf("id = '%v'", req.Id))
	}

	if req.Username != "" {
		queryConditions = append(queryConditions, fmt.Sprintf("username = '%s'", req.Username))
	}

	if len(queryConditions) > 0 {
		queryStatement += " WHERE " + strings.Join(queryConditions, " AND ")
	}

	queryStatement += `
		ORDER BY created_at DESC limit 1
	`

	var response model.Users
	rows, err := s.db.Query(queryStatement)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&response.Id,
			&response.Email,
			&response.Username,
			&response.Role,
			&response.Address,
			pq.Array(&response.CategoryPreferences),
			&response.CreatedAt,
			&response.UpdatedAt,
			&response.DeletedAt,
			&response.Password,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("no user found")
			}
			return nil, fmt.Errorf("failed to fetch user data")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterate over user: %v", err)
	}

	return &response, nil
}
