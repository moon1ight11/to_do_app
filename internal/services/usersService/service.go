package usersservice

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"todoapp/internal/api/models"
	"todoapp/internal/storage/repos/users"
)

// добавление пользователя в БД +++
func (u *UserService) AddUser(user models.UserAuth) (uuid.UUID, error) {
	// проверка на уникальность имени и почты
	exist, err := u.CheckNameAndEmail(user.Name, user.Email)
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
	user_id, err := u.userRepo.CreateUser(user.Name, string(hashPass), user.Email)
	if err != nil {
		return uuid.Nil, err
	}

	return user_id, nil
}

// проверка и получение пользователя +++
func (u *UserService) CheckAndGetUser(user models.UserAuth) (models.UserRequest, error) {
	var foundUser users.User
	var err error

	// находим пользователя по почте
	foundUser, err = u.userRepo.UserByEmail(user.Email)
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

// проверка свободности имени и почты +++
func (u *UserService) CheckNameAndEmail(userName string, userEmail string) (bool, error) {
	// проверяем имя
	exist, err := u.userRepo.CheckUserName(userName)
	if err != nil {
		return true, err
	}
	if exist {
		return true, nil
	}

	// проверяем почту
	exist, err = u.userRepo.CheckUserEmail(userEmail)
	if err != nil {
		return true, err
	}

	return exist, nil
}



// получение данных пользователя
func (u *UserService) GetUser(UserId uuid.UUID) (users.User, error) {
	user, err := u.userRepo.UserById(UserId)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

// проверка существования пользователя по id
func (u *UserService) CheckUserByID(user_id uuid.UUID) (bool, error) {
	_, err := u.userRepo.UserById(user_id)
	if err != nil {
		return false, err
	}

	return true, nil
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
func (u *UserService) DeleteUser(user_id uuid.UUID) error {
	err := u.userRepo.DeleteUser(user_id)
	if err != nil {
		return err
	}

	return nil
}
