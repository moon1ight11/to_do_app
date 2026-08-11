package usersservice

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/crypto/bcrypt"

	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/usersrepos"
)

func TestAddUser_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockUserRepo{
		CheckUserNameFunc: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
		CheckUserEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateUserFunc: func(ctx context.Context, name string, hashPass string, email string) (uuid.UUID, error) {
			return userId, nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	user := models.UserAuth{
		Name:  "test",
		Email: "test@test.com",
		Pass:  "password",
	}

	result, err := svc.AddUser(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != userId {
		t.Errorf("expected userId %v, got %v", userId, result)
	}
}

func TestAddUser_NameOrEmailExists(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockUserRepo{
		CheckUserNameFunc: func(ctx context.Context, name string) (bool, error) {
			return true, nil
		},
		CheckUserEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	user := models.UserAuth{
		Name:  "test",
		Email: "test@test.com",
		Pass:  "password",
	}

	_, err := svc.AddUser(ctx, user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckAndGetUser_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	pass := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	mockRepo := &MockUserRepo{
		UserByEmailFunc: func(ctx context.Context, userEmail string) (usersrepos.User, error) {
			return usersrepos.User{
				Id:    userId,
				Name:  "test",
				Email: "test@test.com",
				Pass:  string(hash),
			}, nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	user := models.UserAuth{
		Email: "test@test.com",
		Pass:  pass,
	}

	result, err := svc.CheckAndGetUser(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Id != userId {
		t.Errorf("expected userId %v, got %v", userId, result.Id)
	}
}

func TestCheckAndGetUser_WrongPassword(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)

	mockRepo := &MockUserRepo{
		UserByEmailFunc: func(ctx context.Context, userEmail string) (usersrepos.User, error) {
			return usersrepos.User{
				Id:    userId,
				Name:  "test",
				Email: "test@test.com",
				Pass:  string(hash),
			}, nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	user := models.UserAuth{
		Email: "test@test.com",
		Pass:  "wrong",
	}

	_, err := svc.CheckAndGetUser(ctx, user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockUserRepo{
		UserByIdFunc: func(ctx context.Context, id uuid.UUID) (usersrepos.User, error) {
			return usersrepos.User{
				Id:    userId,
				Name:  "test",
				Email: "test@test.com",
			}, nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	result, err := svc.GetUser(ctx, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Id != userId {
		t.Errorf("expected userId %v, got %v", userId, result.Id)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockUserRepo{
		UserByIdFunc: func(ctx context.Context, id uuid.UUID) (usersrepos.User, error) {
			return usersrepos.User{}, errors.New("not found")
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	_, err := svc.GetUser(ctx, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()

	mockRepo := &MockUserRepo{
		DeleteUserFunc: func(ctx context.Context, userId uuid.UUID) error {
			return nil
		},
	}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	err := svc.DeleteUser(ctx, userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateUser_EmptyName(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	emptyName := ""

	mockRepo := &MockUserRepo{}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	err := svc.UpdateUser(ctx, &emptyName, nil, nil, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateUser_EmptyPass(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	emptyPass := ""

	mockRepo := &MockUserRepo{}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	err := svc.UpdateUser(ctx, nil, &emptyPass, nil, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	badEmail := "notanemail"

	mockRepo := &MockUserRepo{}

	svc := NewUserService(mockRepo, noop.NewTracerProvider().Tracer("test"))

	err := svc.UpdateUser(ctx, nil, nil, &badEmail, userId)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
