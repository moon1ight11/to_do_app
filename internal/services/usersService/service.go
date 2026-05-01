package usersservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"todoapp/internal/api/models"
)

// добавление пользователя в БД
func (u *UserService) AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
	ctx, span := u.tracer.Start(ctx, "service.AddUser")
	defer span.End()

	// проверка на уникальность имени и почты
	exist, err := u.CheckNameAndEmail(ctx, user.Name, user.Email)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("error in AddUser: %w", err)
	}
	if exist {
		span.RecordError(fmt.Errorf("error in AddUser: name or email already exists"))
		return uuid.Nil, fmt.Errorf("error in AddUser: name or email already exists")
	}

	// хэширование пароля
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("error in AddUser hash password: %w", err)
	}

	// добавление в базу данных
	userId, err := u.userRepo.CreateUser(ctx, user.Name, string(hashPass), user.Email)
	if err != nil {
		span.RecordError(err)
		return uuid.Nil, fmt.Errorf("error in AddUser: %w", err)
	}

	return userId, nil
}

// проверка и получение пользователя
func (u *UserService) CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
	ctx, span := u.tracer.Start(ctx, "service.CheckAndGetUser")
	defer span.End()

	// находим пользователя по почте
	foundUser, err := u.userRepo.UserByEmail(ctx, user.Email)
	if err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("error in CheckAndGetUser: %w", err)
	}

	// сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Pass), []byte(user.Pass))
	if err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("error in CheckAndGetUser: passwords not match")
	}

	// приводим тип для экспорта
	var userApi models.UserRequest
	userApi.Id = foundUser.Id
	userApi.Name = foundUser.Name
	userApi.Email = foundUser.Email

	// если все ок - возвращаем найденного пользователя
	return userApi, nil
}

// проверка свободности имени и почты
func (u *UserService) CheckNameAndEmail(ctx context.Context, userName string, userEmail string) (bool, error) {
	ctx, span := u.tracer.Start(ctx, "service.CheckNameAndEmail")
	defer span.End()

	// проверяем имя
	exist, err := u.userRepo.CheckUserName(ctx, userName)
	if err != nil {
		span.RecordError(err)
		return true, fmt.Errorf("error in CheckNameAndEmail: %w", err)
	}
	if exist {
		span.RecordError(fmt.Errorf("name already exist"))
		return true, nil
	}

	// проверяем почту
	exist, err = u.userRepo.CheckUserEmail(ctx, userEmail)
	if err != nil {
		span.RecordError(err)
		return true, fmt.Errorf("error in CheckNameAndEmail: %w", err)
	}

	return exist, nil
}

// получение данных пользователя
func (u *UserService) GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
	ctx, span := u.tracer.Start(ctx, "service.GetUser")
	defer span.End()

	// запрашиваем пользователя в репозитории
	user, err := u.userRepo.UserById(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return models.UserRequest{}, fmt.Errorf("error in GetUser: %w", err)
	}

	// приводим тип
	var userApi models.UserRequest
	userApi.Id = user.Id
	userApi.Name = user.Name
	userApi.Email = user.Email

	return userApi, nil
}

// обновление полей пользователя
func (u *UserService) UpdateUser(ctx context.Context, name *string, pass *string, email *string, userId uuid.UUID) error {
	ctx, span := u.tracer.Start(ctx, "service.UpdateUser")
	defer span.End()

	// если обновляется имя - чтобы было не пустое
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			span.RecordError(fmt.Errorf("error in UpdateUser: new name is empty"))
			return fmt.Errorf("error in UpdateUser: new name is empty")
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if pass != nil {
		if strings.TrimSpace(*pass) == "" {
			span.RecordError(fmt.Errorf("error in UpdateUser: new pass is empty"))
			return fmt.Errorf("error in UpdateUser: new pass is empty")
		}
	}

	// если обновляется почта
	if email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *email)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in check email in UpdateUser: %w", err)
		}

		// если нет - отклоняем
		if !matched {
			span.RecordError(fmt.Errorf("error in UpdateUser: new email not looks like email"))
			return fmt.Errorf("error in UpdateUser: new email not looks like email")
		}
	}

	// открываем транзакцию
	transaction, err := u.userRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in UpdateUser BeginTx: %w", err)
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем имя
	if name != nil {
		// проверяем, не занято ли новое имя
		NameExist, err := u.userRepo.CheckUserName(ctx, *name)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
		if NameExist {
			span.RecordError(fmt.Errorf("error in UpdateUser: new name already exist"))
			return fmt.Errorf("error in UpdateUser: new name already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateName(ctx, *name, userId, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
	}

	// если меняем пароль
	if pass != nil {
		// хэширование пароля
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in hash password in UpdateUser: %w", err)
		}

		// изменяем пароль
		err = u.userRepo.UpdatePass(ctx, string(hashPass), userId, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
	}

	// если меняем почту
	if email != nil {
		// проверяем, не занята ли новая почта
		EmailExist, err := u.userRepo.CheckUserEmail(ctx, *email)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
		if EmailExist {
			span.RecordError(fmt.Errorf("error in UpdateUser: new email already exist"))
			return fmt.Errorf("error in UpdateUser: new email already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateEmail(ctx, *email, userId, transaction)
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("error in UpdateUser: %w", err)
		}
	}

	// если все ок - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in UpdateUser commit: %w", err)
	}

	return nil
}

// удаление пользователя
func (u *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	ctx, span := u.tracer.Start(ctx, "service.DeleteUser")
	defer span.End()

	err := u.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("error in DeleteUser: %w", err)
	}

	return nil
}
