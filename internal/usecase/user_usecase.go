package usecase

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"one-repo/internal/domain"
	"one-repo/internal/domain/entity"
	"one-repo/internal/domain/repository"
	"strconv"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, name, email, password string) (*entity.User, string, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, name, email, password string) (*entity.User, string, error) {
	// Emailの重複チェック
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", domain.ErrUserAlreadyExists
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// ユーザーの作成
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	// JWTトークンの生成
	token, err := GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	// Emailによるユーザーの検索
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidPassword
	}

	// JWTトークンの生成
	return GenerateToken(strconv.Itoa(int(user.ID)))
}
