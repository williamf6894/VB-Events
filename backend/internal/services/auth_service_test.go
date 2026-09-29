package services

import (
	"testing"
	"uuid"

	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/models"
)

func newTestAuthService(secret string) *AuthService {
	return NewAuthService(nil, config.Config{JWTSecret: secret})
}

func TestAuthService_TokenRoundTrip(t *testing.T) {
	svc := newTestAuthService("test-secret")

	participant := &models.Participant{Email: "token@example.com"}
	participant.ID = uuid.NewV7()

	token, err := svc.GenerateToken(participant)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %s", err)
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %s", err)
	}
	if userID != participant.ID {
		t.Fatalf("expected user ID %v, got %v", participant.ID, userID)
	}
}

func TestAuthService_ValidateTokenRejectsWrongSecret(t *testing.T) {
	issuer := newTestAuthService("secret-one")
	validator := newTestAuthService("secret-two")

	participant := &models.Participant{Email: "token@example.com"}
	participant.ID = uuid.NewV7()

	token, err := issuer.GenerateToken(participant)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %s", err)
	}

	if _, err := validator.ValidateToken(token); err == nil {
		t.Fatal("expected token signed with different secret to be rejected")
	}
}

func TestAuthService_ValidateTokenRejectsTamperedToken(t *testing.T) {
	svc := newTestAuthService("test-secret")

	if _, err := svc.ValidateToken("not-a-jwt"); err == nil {
		t.Fatal("expected malformed token to be rejected")
	}
}

func TestAuthService_GenerateTokenRequiresSecret(t *testing.T) {
	svc := newTestAuthService("")

	participant := &models.Participant{Email: "token@example.com"}
	participant.ID = uuid.NewV7()

	if _, err := svc.GenerateToken(participant); err == nil {
		t.Fatal("expected error when JWT_SECRET is empty")
	}
}
