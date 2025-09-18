package users

import "log"

// удаление пользователя по id
func (db *Base) DeleteUser(id int) error {
	query := `
				DELETE FROM users
				WHERE id = $1
			`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Error in DeleteUser query", err)
		return err
	}
	return nil
}