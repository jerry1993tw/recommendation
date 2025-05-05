package auth

import (
	"errors"

	"app/internal/user"
	"app/pkg/utils"
)

type AuthService struct {
	userRepo *user.Repository
	emailSvc EmailService
}

type EmailService interface {
	SendVerificationCode(email string, code string) error
}

func NewService(userRepo *user.Repository, emailSvc EmailService) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		emailSvc: emailSvc,
	}
}

func (s *AuthService) Register(input RegisterInput) error {
	if err := ValidatePassword(input.Password); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	verifyToken, err := utils.GenerateVerifyToken()
	if err != nil {
		return err
	}

	user := user.User{
		Email:       input.Email,
		Password:    hashedPassword,
		IsVerified:  false,
		VerifyToken: verifyToken,
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		return err
	}

	return s.emailSvc.SendVerificationCode(input.Email, verifyToken)
}

func (s *AuthService) Login(input LoginInput) (*user.User, error) {
	user, err := s.userRepo.GetUserByEmail(input.Email)
	if err != nil {
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
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.VerifyToken != input.Code {
		return errors.New("invalid verification code")
	}

	if err := s.userRepo.VerifyUser(input.Email); err != nil {
		return err
	}

	return nil
}
