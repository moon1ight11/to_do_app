package users

import "log"

// Добавление пользователя в DB и возврат его id
func (db *Base) AddUser(NewUser User) (int, error) {
	query := `
				INSERT INTO users (name, pass, email)
				VALUES ($1, $2, $3)
				RETURNING id
			`
	var user_id int
	
	err := db.DB.QueryRow(query, NewUser.Name, NewUser.Pass, NewUser.Email).Scan(&user_id)
	if err != nil {
		log.Println("Error in AddUser query", err)
		return 0, err
	}

	return user_id, nil
}