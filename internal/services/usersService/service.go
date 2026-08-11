package usersservice

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"todoapp/internal/api/models"
)

func (u *UserService) AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
	ctx, span := u.tracer.Start(ctx, "service.AddUser")
	defer span.End()

	exist, err := u.CheckNameAndEmail(ctx, user.Name, user.Email)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersservice.AddUser: %w", err)
	}
	if exist {
		err := fmt.Errorf("name or email already exists")
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersservice.AddUser: %w", err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersservice.AddUser: hash password: %w", err)
	}

	userId, err := u.userRepo.CreateUser(ctx, user.Name, string(hashPass), user.Email)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("usersservice.AddUser: create user: %w", err)
	}

	return userId, nil
}

func (u *UserService) CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
	ctx, span := u.tracer.Start(ctx, "service.CheckAndGetUser")
	defer span.End()

	foundUser, err := u.userRepo.UserByEmail(ctx, user.Email)
	if err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("usersservice.CheckAndGetUser: find by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Pass), []byte(user.Pass)); err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("usersservice.CheckAndGetUser: passwords not match")
	}

	return models.UserRequest{
		Id:    foundUser.Id,
		Name:  foundUser.Name,
		Email: foundUser.Email,
	}, nil
}

func (u *UserService) CheckNameAndEmail(ctx context.Context, userName string, userEmail string) (bool, error) {
	ctx, span := u.tracer.Start(ctx, "service.CheckNameAndEmail")
	defer span.End()

	exist, err := u.userRepo.CheckUserName(ctx, userName)
	if err != nil {
		span.RecordError(err)
		return true, fmt.Errorf("usersservice.CheckNameAndEmail: check name: %w", err)
	}
	if exist {
		return true, nil
	}

	exist, err = u.userRepo.CheckUserEmail(ctx, userEmail)
	if err != nil {
		span.RecordError(err)
		return true, fmt.Errorf("usersservice.CheckNameAndEmail: check email: %w", err)
	}

	return exist, nil
}

func (u *UserService) GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
	ctx, span := u.tracer.Start(ctx, "service.GetUser")
	defer span.End()

	user, err := u.userRepo.UserById(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("usersservice.GetUser: %w", err)
	}

	return models.UserRequest{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *UserService) UpdateUser(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error {
	ctx, span := u.tracer.Start(ctx, "service.UpdateUser")
	defer span.End()

	if name != nil && strings.TrimSpace(*name) == "" {
		err := fmt.Errorf("new name is empty")
		span.RecordError(err)
		return fmt.Errorf("usersservice.UpdateUser: %w", err)
	}

	if pass != nil && strings.TrimSpace(*pass) == "" {
		err := fmt.Errorf("new pass is empty")
		span.RecordError(err)
		return fmt.Errorf("usersservice.UpdateUser: %w", err)
	}

	if email != nil {
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *email)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: validate email: %w", err)
		}
		if !matched {
			err := fmt.Errorf("new email not valid")
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: %w", err)
		}
	}

	transaction, err := u.userRepo.DB().BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersservice.UpdateUser: begin tx: %w", err)
	}
	defer transaction.Rollback()

	if name != nil {
		exist, err := u.userRepo.CheckUserName(ctx, *name)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: check name: %w", err)
		}
		if exist {
			err := fmt.Errorf("new name already exists")
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: %w", err)
		}

		if err := u.userRepo.UpdateName(ctx, *name, userId, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: update name: %w", err)
		}
	}

	if pass != nil {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: hash password: %w", err)
		}

		if err := u.userRepo.UpdatePass(ctx, string(hashPass), userId, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: update pass: %w", err)
		}
	}

	if email != nil {
		exist, err := u.userRepo.CheckUserEmail(ctx, *email)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: check email: %w", err)
		}
		if exist {
			err := fmt.Errorf("new email already exists")
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: %w", err)
		}

		if err := u.userRepo.UpdateEmail(ctx, *email, userId, transaction); err != nil {
			span.RecordError(err)
			return fmt.Errorf("usersservice.UpdateUser: update email: %w", err)
		}
	}

	if err := transaction.Commit(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersservice.UpdateUser: commit: %w", err)
	}

	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	ctx, span := u.tracer.Start(ctx, "service.DeleteUser")
	defer span.End()

	if err := u.userRepo.DeleteUser(ctx, userId); err != nil {
		span.RecordError(err)
		return fmt.Errorf("usersservice.DeleteUser: %w", err)
	}

	return nil
}
