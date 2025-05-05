package auth

import (
	"errors"
	"strconv"
	"time"

	"app/internal/config"
	"app/internal/user"
	"app/pkg/logger"
	"app/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	userRepo *user.Repository
	emailSvc EmailService
	log      *logger.Logger
}

type EmailService interface {
	SendVerificationCode(email string, code string) error
}

func NewService(userRepo *user.Repository, emailSvc EmailService, log *logger.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		emailSvc: emailSvc,
		log:      log,
	}
}

func (s *AuthService) Register(input RegisterInput) error {
	if err := ValidatePassword(input.Password); err != nil {
		s.log.WithError(err).Error("Password validation failed")
		return err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		s.log.WithError(err).Error("Failed to hash password")
		return err
	}

	verifyToken, err := utils.GenerateVerifyToken()
	if err != nil {
		s.log.WithError(err).Error("Failed to generate verification token")
		return err
	}

	user := user.User{
		Email:       input.Email,
		Password:    hashedPassword,
		IsVerified:  false,
		VerifyToken: verifyToken,
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		s.log.WithError(err).Error("Failed to create user")
		return err
	}

	return s.emailSvc.SendVerificationCode(input.Email, verifyToken)
}

func (s *AuthService) Login(input LoginInput) (*user.User, error) {
	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user by email")
		return nil, err
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(input VerifyEmailInput) error {
	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
		s.log.WithError(err).Error("Failed to get user by email")
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.VerifyToken != input.Code {
		return errors.New("invalid verification code")
	}

	if err := s.userRepo.VerifyUser(input.Email); err != nil {
		s.log.WithError(err).Error("Failed to verify user")
		return err
	}

	return nil
}

func (s *AuthService) GenerateToken(user *user.User) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.New().JwtSecret))
}
