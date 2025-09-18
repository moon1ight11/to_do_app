package users

import "log"

// обновление имени 
func (db *Base) UpdateName (newName string, id int) error {
	query := `
				UPDATE users
				SET name = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newName, id)
	if err != nil {
		log.Println("Error in UpdateName query", err)
		return err
	}

	return nil
}

// обновление пароля
func (db *Base) UpdatePass (newPass string, id int) error {
	query := `
				UPDATE users
				SET pass = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newPass, id)
	if err != nil {
		log.Println("Error in UpdatePass query", err)
		return err
	}
	return nil
}

// обновление почты
func (db *Base) UpdateEmail (newEmail string, id int) error {
	query := `
				UPDATE users
				SET email = $1, updated_at = NOW()
				WHERE id = $2
			`
	_, err := db.DB.Exec(query, newEmail, id)
	if err != nil {
		log.Println("Error in UpdateEmail query", err)
		return err
	}
	return nil
}