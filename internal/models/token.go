package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	ScopeAuthentication = "authentication"
)

// Token is type for authentication token
type Token struct {
	PlanText string    `json:"token"`
	UserID   int64     `json:"-"`
	Hash     []byte    `json:"-"`
	Expiry   time.Time `json:"expiry"`
	Scope    string    `json:"-"`
}

// generateToken generates a token lasts for ttl and return it
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlanText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256(([]byte(token.PlanText)))
	token.Hash = hash[:]
	return token, nil
}

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, canle := context.WithTimeout(context.Background(), 3*time.Second)
	defer canle()

	//delete existing token
	_, err := m.DB.ExecContext(ctx, "DELETE FROM tokens WHERE user_id=?", u.ID)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO tokens (user_id, name, email, token_hash, expiry, created_at, updated_at) 
			VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err = m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		t.Expiry,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	tokenHash := sha256.Sum256([]byte(token))
	var user User

	query := `SELECT u.id, u.first_name, u.last_name, u.email
			 FROM users u INNER JOIN tokens t ON (u.id=t.user_id)
			 WHERE t.token_hash=?
			 AND t.expiry > ?`
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &user, nil

}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancle()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email=?", email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("Incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil

}
