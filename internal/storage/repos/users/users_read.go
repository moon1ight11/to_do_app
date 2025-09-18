package users

import "log"

// получение списка всех пользователей
func (db *Base) AllUsers() ([]User, error) {
	query := `
				SELECT id, name, email
				FROM users
			`
	var users []User

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Println("Error in AllUsers query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			log.Println("Error in AllUsers scan", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// получение количества пользователей
func (db *Base) AmountUsers() (int, error) {
	query := `
				SELECT COUNT(*)
				FROM users
			`
	var amount int
	err := db.DB.QueryRow(query).Scan(&amount)
	if err != nil {
		log.Println("Error in AmountUsers query", err)
		return 0, err
	}

	return amount, nil
}

// получение пользователя по id
func (db *Base) UserById(id int) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM users
				WHERE id = $1
			`
	var user User
	err := db.DB.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		log.Println("Error in UserById query", err)
		return User{}, err
	}

	return user, nil
}

// получение пользователя по имени
func (db *Base) UserByName(name string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM users
				WHERE name = $1
			`
	var user User
	err := db.DB.QueryRow(query, name).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		log.Println("Error in UserByName query", err)
		return User{}, err
	}

	return user, nil
}

// получение пользователя по почте
func (db *Base) UserByEmail(email string) (User, error) {
	query := `
				SELECT id, name, email, pass
				FROM users
				WHERE email = $1
			`
	var user User
	err := db.DB.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		log.Println("Error in UserByEmail query", err)
		return User{}, err
	}

	return user, nil
}
