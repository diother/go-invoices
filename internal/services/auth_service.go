package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/diother/go-invoices/internal/custom_errors"
	"github.com/diother/go-invoices/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetUser(username string) (*models.User, error)
	InsertSession(session *models.Session) error
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Authenticate(username, password string) (user *models.User, err error) {
	if err := validateCredentials(username, password); err != nil {
		return nil, custom_errors.NewCredentialsError("validate credentials failed: %v", err)
	}
	user, err = s.repo.GetUser(username)
	if err != nil {
		var credentialsError *custom_errors.CredentialsError
		if errors.As(err, &credentialsError) {
			return nil, err
		}
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, custom_errors.NewCredentialsError("username or password is invalid")
	}
	return
}

func (s *AuthService) GenerateSession(user *models.User) (session *models.Session, err error) {
	sessionToken, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("session token generation failed: %w", err)
	}
	expiresAt := time.Now().Unix() + (30 * 24 * 60 * 60)

	session = transformSessionDTOToModel(sessionToken, user.ID, expiresAt)
	if err = s.repo.InsertSession(session); err != nil {
		return nil, fmt.Errorf("session insertion failed: %w", err)
	}
	return
}

func generateSessionID() (int64, error) {
	timestamp := time.Now().Unix()
	randomNum, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}
	sessionID := timestamp*10000 + randomNum.Int64() + 1000
	return sessionID, nil
}

func validateCredentials(username, password string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if password == "" {
		return fmt.Errorf("password is empty")
	}
	return nil
}

func transformSessionDTOToModel(sessionToken, userID, expiresAt int64) *models.Session {
	return models.NewSession(sessionToken, userID, expiresAt)
}
