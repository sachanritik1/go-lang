package store

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainText
	p.hash = hashedBytes
	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"` // "-" to omit from JSON responses
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(id int) error
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, bio)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
		`
	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresUserStore) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE id = $1`
	user := &User{}
	err := store.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE username = $1`
	user := &User{}
	err := store.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *PostgresUserStore) UpdateUser(user *User) (*User, error) {
	query := `UPDATE users SET email = $1, password_hash = $2, bio = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at`
	err := store.db.QueryRow(query, user.Email, user.PasswordHash, user.Bio, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *PostgresUserStore) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := store.db.Exec(query, id)
	return err
}
