package usersvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/mertani/config"
	"github.com/arfan21/mertani/internal/entity"
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/internal/user"
	"github.com/arfan21/mertani/pkg/constant"
	"github.com/arfan21/mertani/pkg/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo user.Repository
}

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Register(ctx context.Context, req model.UserRegisterRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to validate request: %w", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to hash password: %w", err)
		return
	}

	data := entity.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to register user: %w", err)
		return
	}

	return
}

func (s Service) Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to validate request: %w", err)
		return
	}

	data, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to get user by email: %w", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = constant.ErrEmailOrPasswordInvalid
		}
		err = fmt.Errorf("user.service.Login: failed to compare password: %w", err)
		return
	}

	accessTokenExpire := time.Duration(config.Get().JWT.AccessTokenExpireIn) * time.Second

	accessToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Email,
		config.Get().JWT.AccessTokenSecret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to create access token: %w", err)
		return
	}

	res = model.UserLoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(accessTokenExpire.Seconds()),
		TokenType:   "Bearer",
	}

	return
}

func (s Service) CreateJWTWithExpiry(id, email, secret string, expiry time.Duration) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Get().Service.Name,
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("usecase: failed to create jwt token: %w", err)
		return
	}

	return
}
