package authn

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/config"
	"github.com/muhrifqii/tuskar/internal/repository"
	"github.com/muhrifqii/tuskar/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	userRepository repository.UserRepository
	log            *zap.Logger
	jwtConf        config.JwtConfig
}

func NewService(zap *zap.Logger, jwtConf config.JwtConfig, userRepository repository.UserRepository) *Service {
	return &Service{
		log:            zap,
		jwtConf:        jwtConf,
		userRepository: userRepository,
	}
}

func generateAccessToken(tokenExpiration int, userID, secret string) (string, time.Time, error) {
	accessTokenExpiry := time.Now().Add(time.Minute * time.Duration(tokenExpiration))
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": accessTokenExpiry.Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return accessTokenString, accessTokenExpiry, nil
}

func generateRefreshToken(refreshTokenExpiration int, secret string) (string, time.Time, error) {
	refreshTokenExpiry := time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenExpiration))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": refreshTokenExpiry.Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return refreshTokenString, refreshTokenExpiry, nil
}

func (s *Service) Login(ctx context.Context, req domain.AuthnRequest) (domain.AuthnResponse, error) {
	var (
		user *domain.User
		err  error
	)

	user, err = s.userRepository.GetByUsername(ctx, req.Username)

	if err != nil {
		return domain.AuthnResponse{}, domain.ErrInvalidCredentials
	}
	if err = utils.CheckPassword(user.Password, req.Password); err != nil {
		return domain.AuthnResponse{}, domain.ErrInvalidCredentials
	}

	accessTokenString, accessTokenExpiry, err := generateAccessToken(s.jwtConf.Expiration, user.Username, s.jwtConf.Secret)
	if err != nil {
		return domain.AuthnResponse{}, err
	}

	refreshTokenString, refreshTokenExpiry, err := generateRefreshToken(s.jwtConf.RefreshExpirationInDays, s.jwtConf.RefreshSecret)
	if err != nil {
		return domain.AuthnResponse{}, err
	}

	return domain.AuthnResponse{
		AccessToken:           accessTokenString,
		AccessTokenExpiresAt:  accessTokenExpiry,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: refreshTokenExpiry,
	}, nil
}

func (s *Service) Logout(ctx context.Context) error {
	panic("not implemented")
}
