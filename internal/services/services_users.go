package services

import (
	"fmt"
	"github.com/google/uuid"
	"todoapp/internal/storage/repos/users"
)

type UserService struct {
	userRepo *users.Base
}

func NewUserService(userRepo *users.Base) *UserService {
	return &UserService{userRepo: userRepo}
}

// проверка свободности имени пользователя
func (u *UserService) CheckName(user_name string) (bool, error) {
	exist, err := u.userRepo.CheckUserName(user_name)
	if err != nil {
		return true, err
	}

	return exist, nil
}

// проверка свободности почты
func (u *UserService) CheckEmail(user_email string) (bool, error) {
	exist, err := u.userRepo.CheckUserEmail(user_email)
	if err != nil {
		return true, err
	}

	return exist, nil
}

// проверка существования пользователя по id
func (u *UserService) CheckUserByID(id uuid.UUID) (bool, error) {
	_, err := u.userRepo.UserById(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// добавление пользователя в БД
func (u *UserService) CreateUser(User users.User) (uuid.UUID, error) {
	user_id, err := u.userRepo.AddUser(User)
	if err != nil {
		return uuid.Nil, err
	}

	return user_id, nil
}

// проверка пользователя
func (u *UserService) CheckUser(User users.User) (bool, users.User, error) {
	var foundUser users.User
	var err error

	foundUser, err = u.userRepo.UserByEmail(User.Email)
	if err != nil {
		return false, users.User{}, err
	}

	if foundUser.Pass != User.Pass {
		return false, users.User{}, fmt.Errorf("Passwords not match")
	}

	return true, foundUser, nil
}

// обновление полей пользователя
func (u *UserService) UpdateUser(user_name *string, user_pass *string, user_email *string, user_id uuid.UUID) error {
	// открываем транзакцию
	transaction, err := u.userRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем имя
	if user_name != nil {
		// проверяем, не занято ли новое имя
		NameExist, err := u.userRepo.CheckUserName(*user_name)
		if err != nil {
			return err
		}
		if NameExist {
			return fmt.Errorf("New name already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateName(*user_name, user_id, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем пароль
	if user_pass != nil {
		err := u.userRepo.UpdatePass(*user_pass, user_id, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем почту
	if user_email != nil {
		// проверяем, не занята ли новая почта
		EmailExist, err := u.CheckEmail(*user_email)
		if err != nil {
			return err
		}
		if EmailExist {
			return fmt.Errorf("New email already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateEmail(*user_email, user_id, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// удаление пользователя
func (u *UserService) DeleteUser(id uuid.UUID) error {
	err := u.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
