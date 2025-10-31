package models

type User struct {
    Id           uint   `json:"user_id" db:"user_id"`
    Name         string `json:"fullname" db:"name"`
    Email        string `json:"email" db:"email"`
    PasswordHash string `json:"-" db:"password_hash"`
}

type RegisterRequest struct {
    Name     string `json:"fullname" binding:"required"`
    UserName     string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    UserID uint   `json:"user_id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Token  string `json:"token"`
}