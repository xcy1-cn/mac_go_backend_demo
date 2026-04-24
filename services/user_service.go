package services

import (
	"database/sql"
	"demo/day6-9/db"
	"demo/day6-9/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(req models.RegisterRequest) error {
	// 1. 参数基础校验
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}

	// 2. 检查用户名是否已存在
	var count int
	checkSQL := "SELECT COUNT(*) FROM users WHERE username = ?"
	err := db.DB.QueryRow(checkSQL, req.Username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("username already exists")
	}

	// 是否保存
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 3. 插入用户
	// insertSQL := "INSERT INTO users (username, password, nickname) VALUES (?, ?, ?)"
	// _, err = db.DB.Exec(insertSQL, req.Username, req.Password, req.Nickname)
	insertSQL := "INSERT INTO users (username, password, nickname) VALUES (?, ?, ?)"
	_, err = db.DB.Exec(insertSQL, req.Username, string(hashedPassword), req.Nickname)
	if err != nil {
		return err
	}

	return nil
}

func LoginUser(req models.LoginRequest) (*models.User, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	var user models.User
	querySQL := "SELECT id, username, password, nickname, created_at, updated_at FROM users WHERE username = ?"

	err := db.DB.QueryRow(querySQL, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Nickname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func GetUserByID(userID int64) (*models.User, error) {
	var user models.User

	querySQL := "SELECT id, username, nickname, created_at, updated_at FROM users WHERE id = ?"
	err := db.DB.QueryRow(querySQL, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
