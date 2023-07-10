package auth

import (
	"context"
	"encoding/json"
	"time"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type SessionManager struct {
	Sql *sqlite3.Client
	Conn *sqlx.DB
}

type UserSession struct {
	Id int
	FirstName string
	LastName string
	Email string
}

type User struct {
	Id int
	FirstName string
	LastName string
	Email string
	Password string
}