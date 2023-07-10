package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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