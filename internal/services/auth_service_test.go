package services

import (
	"testing"
	"time"
)

func TestGenerateSessionID(t *testing.T) {
	testCases := map[string]struct {
		expectError bool
	}{
		"successfulIDGeneration": {expectError: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sessionID, err := generateSessionID()

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tc.expectError {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if sessionID <= 1000 {
					t.Errorf("Generated session ID is invalid: %d", sessionID)
				}
			}
		})
	}
}

func TestValidateCredentials(t *testing.T) {
	testCases := map[string]struct {
		username    string
		password    string
		expectError bool
	}{
		"emptyUsername":    {username: "", password: "password", expectError: true},
		"emptyPassword":    {username: "username", password: "", expectError: true},
		"validCredentials": {username: "username", password: "password", expectError: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateCredentials(tc.username, tc.password)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

func TestTransformSessionDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		sessionToken int64
		userID       int64
		expiresAt    int64
	}{
		"validData":  {sessionToken: 1234567890, userID: 1, expiresAt: time.Now().Unix() + 3600},
		"zeroValues": {sessionToken: 0, userID: 0, expiresAt: 0},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			session := transformSessionDTOToModel(tc.sessionToken, tc.userID, tc.expiresAt)

			if session.SessionToken != tc.sessionToken {
				t.Errorf("Expected sessionToken %d, got %d", tc.sessionToken, session.SessionToken)
			}
			if session.UserID != tc.userID {
				t.Errorf("Expected userID %d, got %d", tc.userID, session.UserID)
			}
			if session.ExpiresAt != tc.expiresAt {
				t.Errorf("Expected expiresAt %d, got %d", tc.expiresAt, session.ExpiresAt)
			}
		})
	}
}

func TestValidateSessionToken(t *testing.T) {
	testCases := map[string]struct {
		sessionToken string
		expectError  bool
	}{
		"validSessionToken":   {sessionToken: "16960400123456", expectError: false},
		"invalidSessionToken": {sessionToken: "invalid_token", expectError: true},
		"emptySessionToken":   {sessionToken: "", expectError: true},
		"validSmallToken":     {sessionToken: "123", expectError: false},
		"validNegativeToken":  {sessionToken: "-16960400123456", expectError: false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := validateSessionToken(tc.sessionToken)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
