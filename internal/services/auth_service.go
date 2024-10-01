package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/diother/go-invoices/internal/custom_errors"
	"github.com/diother/go-invoices/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID int64) (*models.User, error)
	GetSession(sessionToken int64) (*models.Session, error)
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
	user, err = s.repo.GetUserByUsername(username)
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
	expiresAt := time.Now().Unix() + (12 * 30 * 24 * 60 * 60)

	session = transformSessionDTOToModel(sessionToken, user.ID, expiresAt)
	if err = s.repo.InsertSession(session); err != nil {
		return nil, fmt.Errorf("session insertion failed: %w", err)
	}
	return
}

func (s *AuthService) ValidateSession(sessionTokenString string) (user *models.User, err error) {
	sessionToken, err := validateSessionToken(sessionTokenString)
	if err != nil {
		return nil, fmt.Errorf("session token invalid: %w", err)
	}
	session, err := s.repo.GetSession(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("get session failed: %w", err)
	}
	user, err = s.repo.GetUserByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	return
}

// needs unit test
func validateSessionToken(sessionToken string) (token int64, err error) {
	token, err = strconv.ParseInt(sessionToken, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid session token: %w", err)
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
