package store

import (
	"database/sql"
	"time"

	"github.com/sachanritik1/go-lang/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int64, ttlMinutes int64, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int64, scope string) error
}

func (pts *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`
	_, err := pts.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (pts *PostgresTokenStore) CreateNewToken(userID int64, ttlMinutes int64, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, time.Duration(ttlMinutes)*time.Minute, scope)
	if err != nil {
		return nil, err
	}

	err = pts.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (pts *PostgresTokenStore) DeleteAllTokensForUser(userID int64, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2
	`
	_, err := pts.db.Exec(query, userID, scope)
	return err
}
