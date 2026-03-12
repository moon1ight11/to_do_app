package usersservice

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"todoapp/internal/api/models"
)

// добавление пользователя в БД
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

// проверка и получение пользователя
func (u *UserService) CheckAndGetUser(user models.UserAuth) (models.UserRequest, error) {
	// находим пользователя по почте
	foundUser, err := u.userRepo.UserByEmail(user.Email)
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
func (u *UserService) GetUser(userId uuid.UUID) (models.UserRequest, error) {
	// запрашиваем пользователя в репозитории
	user, err := u.userRepo.UserById(userId)
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
func (u *UserService) UpdateUser(name *string, pass *string, email *string, userId uuid.UUID) error {
	// открываем транзакцию
	transaction, err := u.userRepo.DB.Begin()
	if err != nil {
		return err
	}

	// отложенно откатываем транзакцию
	defer transaction.Rollback()

	// если меняем имя
	if name != nil {
		// проверяем, не занято ли новое имя
		NameExist, err := u.userRepo.CheckUserName(*name)
		if err != nil {
			return err
		}
		if NameExist {
			return fmt.Errorf("New name already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateName(*name, userId, transaction)
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
		err = u.userRepo.UpdatePass(string(hashPass), userId, transaction)
		if err != nil {
			return err
		}
	}

	// если меняем почту
	if email != nil {
		// проверяем, не занята ли новая почта
		EmailExist, err := u.userRepo.CheckUserEmail(*email)
		if err != nil {
			return err
		}
		if EmailExist {
			return fmt.Errorf("New email already exist")
		}

		// если ок - меняем
		err = u.userRepo.UpdateEmail(*email, userId, transaction)
		if err != nil {
			return err
		}
	}

	// если все ок - подтверждаем транзакцию
	transaction.Commit()
	return nil
}

// проверка существования пользователя по id
func (u *UserService) CheckUserByID(userId uuid.UUID) (bool, error) {
	// получаем пользователя по id
	_, err := u.userRepo.UserById(userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

// удаление пользователя
func (u *UserService) DeleteUser(userId uuid.UUID) error {
	err := u.userRepo.DeleteUser(userId)
	if err != nil {
		return err
	}

	return nil
}
