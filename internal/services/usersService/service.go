package usersservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"strings"
	"todoapp/internal/api/models"
)

// добавление пользователя в БД
func (u *UserService) AddUser(ctx context.Context, user models.UserAuth) (uuid.UUID, error) {
	// проверка на уникальность имени и почты
	exist, err := u.CheckNameAndEmail(ctx, user.Name, user.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exist {
		return uuid.Nil, fmt.Errorf("name or email already exists")
	}

	// хэширование пароля
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in hash password: %w", err)
	}

	// добавление в базу данных
	userId, err := u.userRepo.CreateUser(ctx, user.Name, string(hashPass), user.Email)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

// проверка и получение пользователя
func (u *UserService) CheckAndGetUser(ctx context.Context, user models.UserAuth) (models.UserRequest, error) {
	// находим пользователя по почте
	foundUser, err := u.userRepo.UserByEmail(ctx, user.Email)
	if err != nil {
		return models.UserRequest{}, err
	}

	// сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Pass), []byte(user.Pass))
	if err != nil {
		return models.UserRequest{}, fmt.Errorf("Passwords not match")
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
	// проверяем имя
	exist, err := u.userRepo.CheckUserName(ctx, userName)
	if err != nil {
		return true, err
	}
	if exist {
		return true, nil
	}

	// проверяем почту
	exist, err = u.userRepo.CheckUserEmail(ctx, userEmail)
	if err != nil {
		return true, err
	}

	return exist, nil
}

// получение данных пользователя
func (u *UserService) GetUser(ctx context.Context, userId uuid.UUID) (models.UserRequest, error) {
	// запрашиваем пользователя в репозитории
	user, err := u.userRepo.UserById(ctx, userId)
	if err != nil {
		return models.UserRequest{}, err
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
	// если обновляется имя - чтобы было не пустое
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			log.Println("New name is empty")
			return fmt.Errorf("new name is empty")
		}
	}

	// если обновляется пароль - чтобы не был пустым
	if pass != nil {
		if strings.TrimSpace(*pass) == "" {
			log.Println("New pass is empty")
			return fmt.Errorf("new pass is empty")
		}
	}

	// если обновляется почта
	if email != nil {
		// проверяем, похожа ли новая почта на почту
		pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(pattern, *email)
		if err != nil {
			log.Println("Error in MatchString", err)
			return fmt.Errorf("error in check email: %w", err)
		}

		// если нет - отклоняем
		if !matched {
			log.Println("New email not looks like email")
			return fmt.Errorf("new email not looks like email")
		}
	}

	// открываем транзакцию
	transaction, err := u.userRepo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем имя
	if name != nil {
		// проверяем, не занято ли новое имя
		NameExist, err := u.userRepo.CheckUserName(ctx, *name)
		if err != nil {
			return err
		}
		if NameExist {
			return fmt.Errorf("New name already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateName(ctx, *name, userId, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем пароль
	if pass != nil {
		// хэширование пароля
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error in hash password: %w", err)
		}

		// изменяем пароль
		err = u.userRepo.UpdatePass(ctx, string(hashPass), userId, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем почту
	if email != nil {
		// проверяем, не занята ли новая почта
		EmailExist, err := u.userRepo.CheckUserEmail(ctx, *email)
		if err != nil {
			return err
		}
		if EmailExist {
			return fmt.Errorf("New email already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateEmail(ctx, *email, userId, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("error in update user commit: %w", err)
	}

	return nil
}

// удаление пользователя
func (u *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	err := u.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
