package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/yourusername/sso-grpc/api"
	"github.com/yourusername/sso-grpc/internal/model"
)

var jwtSecret = []byte("your-secret-key")

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthServiceServer {
	return &AuthServiceServer{db: db}
}

// Register – регистрация пользователя
func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Базовая валидация
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("все поля обязательны")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var id int
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = s.db.QueryRowContext(ctx, query, req.Username, req.Email, hashedPassword).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения пользователя: %v", err)
	}

	return &pb.RegisterResponse{
		Id:      int32(id),
		Message: "Пользователь успешно зарегистрирован",
	}, nil
}

// Login – аутентификация пользователя, возвращает JWT
func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user model.User
	query := "SELECT id, password, username, email FROM users WHERE username = $1"
	err := s.db.QueryRowContext(ctx, query, req.Username).Scan(&user.ID, &user.Password, &user.Username, &user.Email)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return nil, errors.New("неверный пароль")
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

// Profile – получение профиля пользователя по JWT
func (s *AuthServiceServer) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	claims, err := validateJWT(req.Token)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("ошибка парсинга токена")
	}

	var user model.User
	query := "SELECT id, username, email FROM users WHERE id = $1"
	err = s.db.QueryRowContext(ctx, query, int(userID)).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return &pb.ProfileResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// --- Вспомогательные функции ---

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("недействительный токен")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("ошибка извлечения данных из токена")
	}
	return claims, nil
}
