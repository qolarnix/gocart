package auth

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"

)

type SessionManager struct {
	Conn *sqlx.DB
}

func (sm *SessionManager) RegisterUser(firstName, lastName, email, password string) error {
	userID := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = sm.Conn.Exec("INSERT INTO users (id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)",
		userID, firstName, lastName, email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

// Login

func (sm *SessionManager) Login(email, password string) (*UserSession, error) {
	var user User

	err := sm.Conn.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Error comparing passwords:", err)
		return nil, err
	}

	session := &UserSession{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return session, nil
}


type UserSession struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
}

type User struct {
	ID        string `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
}


