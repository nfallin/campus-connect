package Model

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const tokenLifetime = 24 * time.Hour

func (repository *Repository) RegisterUser(request Types.RegisterUserRequest) (Types.UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return Types.UserResponse{}, err
	}

	result, err := repository.db.Exec(`
		INSERT INTO users (user_name, user_email, password_hash)
		VALUES (?, ?, ?)
	`, request.Username, request.Email, string(passwordHash))
	if err != nil {
		return Types.UserResponse{}, translateUserWriteError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Types.UserResponse{}, err
	}

	return Types.UserResponse{
		ID:       int(id),
		Username: request.Username,
		Name:     request.Username,
		Email:    request.Email,
	}, nil
}

func (repository *Repository) LoginUser(request Types.LoginUserRequest) (Types.AuthResponse, error) {
	var user Types.UserResponse
	var passwordHash string
	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email, password_hash
		FROM users
		WHERE user_name = ? OR user_email = ?
	`, request.Identifier, request.Identifier).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Types.AuthResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return Types.AuthResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(request.Password)) != nil {
		return Types.AuthResponse{}, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(tokenLifetime).UTC()
	token := newToken(user.ID, expiresAt)

	user.Name = user.Username
	return Types.AuthResponse{User: user, Token: token}, nil
}

func (repository *Repository) AuthenticateToken(token string) (Types.UserResponse, error) {
	if token == "" {
		return Types.UserResponse{}, ErrUnauthorized
	}

	userID, expiresAt, err := parseToken(token)
	if err != nil || !expiresAt.After(time.Now().UTC()) {
		return Types.UserResponse{}, ErrUnauthorized
	}

	var user Types.UserResponse
	err = repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, ErrUnauthorized
	}
	if err != nil {
		return Types.UserResponse{}, err
	}

	user.Name = user.Username
	return user, nil
}

func (repository *Repository) GetUser(id int) (Types.UserResponse, error) {
	if err := requirePositiveID(id); err != nil {
		return Types.UserResponse{}, err
	}
	var user Types.UserResponse
	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, ErrNotFound
	}
	if err != nil {
		return Types.UserResponse{}, err
	}
	user.Name = user.Username
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
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		user.Name = user.Username
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repository *Repository) UpdateUser(id int, request Types.UpdateUserRequest) (Types.UserResponse, error) {
	if err := requirePositiveID(id); err != nil {
		return Types.UserResponse{}, err
	}
	current, err := repository.GetUser(id)
	if err != nil {
		return Types.UserResponse{}, err
	}
	if request.Username != nil {
		current.Username = *request.Username
		current.Name = *request.Username
	}
	if request.Email != nil {
		current.Email = *request.Email
	}
	_, err = repository.db.Exec(`
		UPDATE users SET user_name = ?, user_email = ? WHERE user_id = ?
	`, current.Username, current.Email, id)
	if err != nil {
		return Types.UserResponse{}, translateUserWriteError(err)
	}
	return current, nil
}

func (repository *Repository) DeleteUser(id int) error {
	if err := requirePositiveID(id); err != nil {
		return err
	}
	result, err := repository.db.Exec(`DELETE FROM users WHERE user_id = ?`, id)
	if err != nil {
		return err
	}
	return requireAffected(result)
}
