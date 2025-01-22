package auth

import (
	"context"
	"errors"
	"fmt"
	"gRPCAuthService/internal/domain/models"
	"gRPCAuthService/internal/lib/jwt"
	"gRPCAuthService/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	appProvider  AppProvider
	userProvider UserProvider
	userSaver    UserSaver
	token        time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, pinHash []byte) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
)

func NewAuth(log *slog.Logger, saver UserSaver, provider UserProvider, appProvider AppProvider, tokenTTl time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    saver,
		userProvider: provider,
		appProvider:  appProvider,
		token:        tokenTTl,
	}
}
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))
	log.Info("start")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %s", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid password", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %s", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)

	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		a.log.Error("failed", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %s", op, ErrInvalidAppID)
	}

	slog.Info("success", slog.String("email", email))

	token, err := jwt.NewToken(user, app, a.token)
	if err != nil {
		a.log.Error("failed", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %s", op, ErrInvalidCredentials)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("Register")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		return 0, fmt.Errorf("%s: %s", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, hash)
	if err != nil {
		slog.Error(err.Error())
		return 0, fmt.Errorf("%s: %s", op, err)
	}

	slog.Info("Successfully registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.String("userID", fmt.Sprint(userID)))

	log.Info("check")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return false, nil
		}
		slog.Error(err.Error())
		return false, fmt.Errorf("%s: %s", op, err)
	}

	slog.Info("isAdmin")

	return isAdmin, nil
}
