package usecase

import (
	"context"
	"fmt"
	"user_service/internal/users"
	"user_service/pkg/log"
	"user_service/pkg/utils/password"
)

type Repo interface {
	Add(context.Context, *users.User) (int64, error)
	GetByEmail(context.Context, string) (*users.User, error)
	GetByID(ctx context.Context, id uint32) (*users.User, error)
	GetAll(ctx context.Context) ([]*users.User, error)
	Update(ctx context.Context, u *users.User) (int64, error)
}

type Manager struct {
	repo        Repo
	argonParams *password.ArgonParams
}

func NewManager(repo Repo) *Manager {
	return &Manager{
		repo: repo,
		argonParams: &password.ArgonParams{
			Memory:      64 * 1024, // 64 MB
			Iterations:  3,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

func (m *Manager) Create(ctx context.Context, in *users.UserIn) (*users.User, error) {
	u, err := m.repo.GetByEmail(ctx, in.Email)
	if err != nil {
		log.Clog(ctx).Error("Error while check user exists", log.Fields{"error": err.Error()})
		return nil, fmt.Errorf("create user: %w", err)
	}
	if u != nil {
		log.Clog(ctx).Info("User exist", log.Fields{"user": u})
		return nil, UserExistsError
	}

	pass, err := password.GenerateHash(in.Password, m.argonParams)
	if err != nil {
		log.Clog(ctx).Info("Error during hash generation", log.Fields{"error": err})
		return nil, fmt.Errorf("create user: %w", err)
	}

	user := &users.User{
		LastName:  in.LastName,
		FirstName: in.FirstName,
		Email:     in.Email,
		Ver:       0,
		PassHash:  pass,
	}

	lastId, err := m.repo.Add(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("new store: %w", err)
	}
	user.ID = uint32(lastId)
	log.Clog(ctx).Info("UserId created", log.Fields{"id": user.ID, "email": user.Email})
	return user, nil
}

func (m *Manager) GetUserVersion(ctx context.Context, userId uint32) (int, error) {
	u, err := m.repo.GetByID(ctx, userId)
	if err != nil {
		log.Clog(ctx).Error("Error while retrieving the user", log.Fields{"userId": userId, "error": err.Error()})
		return 0, fmt.Errorf("get user version: %w", err)
	}
	if u == nil {
		return 0, UserNotFoundError
	}
	return u.Ver, nil
}

func (m *Manager) ListUsers(ctx context.Context) ([]*users.User, error) {
	items, err := m.repo.GetAll(ctx)
	if err != nil {
		log.Clog(ctx).Error("Error while listing users", log.Fields{"error": err.Error()})
		return nil, fmt.Errorf("list users: %w", err)
	}
	return items, nil
}

func (m *Manager) PartialUpdate(ctx context.Context, userId uint32, payload *users.UserUpdate) (*users.User, error) {
	u, err := m.repo.GetByID(ctx, userId)
	if err != nil {
		log.Clog(ctx).Error("Error while retrieving the user", log.Fields{"userId": userId, "error": err.Error()})
		return nil, fmt.Errorf("get user version: %w", err)
	}
	if u == nil {
		return nil, UserNotFoundError
	}
	if payload.LastName != "" {
		u.LastName = payload.LastName
	}
	if payload.FirstName != "" {
		u.FirstName = payload.FirstName
	}
	u.Ver++
	_, err = m.repo.Update(ctx, u)
	if err != nil {
		log.Clog(ctx).Error("Error while updating user", log.Fields{"userId": userId, "error": err.Error()})
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func (m *Manager) CheckPassByEmail(ctx context.Context, email, pass string) (*users.User, error) {
	u, err := m.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Clog(ctx).Error("Error while retrieving the user", log.Fields{"userEmail": email, "error": err.Error()})
		return nil, fmt.Errorf("check user password: %w", err)
	}
	if u == nil {
		return nil, UserNotFoundError
	}
	if ok, err := password.VerifyPassword(pass, u.PassHash); !ok || err != nil {
		return nil, BadPasswordError
	}
	return u, nil
}
