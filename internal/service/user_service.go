package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ユーザー用サービスの構造体
type UserService struct {
	repo *repository.UserRepository
}

// ユーザー用サービスのインスタンスを生成する関数
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ユーザー登録を実施する
func (s *UserService) Signup(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to hash password : Username="+user.Username, err)
	}

	id, err := s.repo.Create(ctx, user.Username, string(hashedPassword))
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// ログインを実施する
func (s *UserService) Login(ctx context.Context, user models.User) (string, int, error) {
	id, hashedPassword, err := s.repo.FindAuthByUsername(ctx, user.Username)
	if err != nil {
		return "", 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		return "", 0, apperror.NewAppError(apperror.TypeUnauthorized, "Invalid username or password : Username="+user.Username, err)
	}

	token, err := GenerateJWT(id)
	if err != nil {
		return "", 0, apperror.NewAppError(apperror.TypeInternalServer, "Failed to generate token : Username="+user.Username, err)
	}

	return token, id, nil
}

// JWTトークンを発行する
func GenerateJWT(userID int) (string, error) {
	// payloadの生成
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// JWTを生成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 署名付きトークン生成
	return token.SignedString(middleware.JwtKey)
}
