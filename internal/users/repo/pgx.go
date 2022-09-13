package repo

import (
	"context"
	"database/sql"

	"github.com/Ollub/user_service/internal/users"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) GetAll(ctx context.Context) ([]*users.User, error) {
	items := []*users.User{}
	rows, err := repo.DB.QueryContext(ctx, "SELECT id, first_name, last_name, email, version FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := &users.User{}
		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Ver)
		if err != nil {
			return nil, err
		}
		items = append(items, u)
	}

	return items, nil
}

func (repo *RepoPgx) GetByID(ctx context.Context, id uint32) (*users.User, error) {
	u := &users.User{}

	err := repo.DB.
		QueryRowContext(ctx, `SELECT id, first_name, last_name, email, version, password FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Ver, &u.PassHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *RepoPgx) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	u := &users.User{}

	err := repo.DB.
		QueryRowContext(ctx, `SELECT id, first_name, last_name, email, version, password FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Ver, &u.PassHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *RepoPgx) Add(ctx context.Context, u *users.User) (int64, error) {
	var lastInsertId int64
	err := repo.DB.QueryRowContext(
		ctx,
		`INSERT INTO users (first_name, last_name, email, version, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Ver,
		u.PassHash,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

//
func (repo *RepoPgx) Update(ctx context.Context, u *users.User) (int64, error) {
	result, err := repo.DB.ExecContext(
		ctx,
		`UPDATE users SET`+
			`"first_name" = $1`+
			`,"last_name" = $2`+
			`,"email" = $3`+
			`,"version" = $4`+
			`WHERE id = $5`,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Ver,
		u.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
