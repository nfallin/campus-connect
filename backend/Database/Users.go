package Database

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
)

// Users
func (repository *Repository) CreateUser(request Types.CreateUserRequest) (Types.UserResponse, bool, error) {
	var user Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_name = ? AND user_email = ?
	`, request.Name, request.Email).Scan(&user.ID, &user.Name, &user.Email)

	if err == nil {
		return user, false, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, false, err
	}

	result, err := repository.db.Exec(`
		INSERT INTO users (user_name, user_email)
		VALUES (?, ?)
	`, request.Name, request.Email)

	if err != nil {
		return Types.UserResponse{}, false, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Types.UserResponse{}, false, err
	}

	return Types.UserResponse{
		ID:    int(id),
		Name:  request.Name,
		Email: request.Email,
	}, true, nil
}

func (repository *Repository) GetUser(id int) (Types.UserResponse, error) {
	var user Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, ErrNotFound
	}

	if err != nil {
		return Types.UserResponse{}, err
	}

	return user, nil
}

func (repository *Repository) GetAllUsers() ([]Types.UserResponse, error) {
	rows, err := repository.db.Query(`
		SELECT user_id, user_name, user_email
		FROM users
		ORDER BY user_id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]Types.UserResponse, 0)

	for rows.Next() {
		var user Types.UserResponse

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *Repository) UpdateUser(id int, request Types.UpdateUserRequest) (Types.UserResponse, error) {
	current, err := repository.GetUser(id)

	if err != nil {
		return Types.UserResponse{}, err
	}

	if request.Name != nil {
		current.Name = *request.Name
	}

	if request.Email != nil {
		current.Email = *request.Email
	}

	_, err = repository.db.Exec(`
		UPDATE users
		SET user_name = ?, user_email = ?
		WHERE user_id = ?
	`, current.Name, current.Email, id)

	if err != nil {
		return Types.UserResponse{}, err
	}

	return current, nil
}

func (repository *Repository) DeleteUser(id int) error {
	result, err := repository.db.Exec(`DELETE FROM users WHERE user_id = ?`, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}
