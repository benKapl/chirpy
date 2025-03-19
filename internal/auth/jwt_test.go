package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJwt(t *testing.T) {
	// Create data necessary for test
	userId, _ := uuid.Parse("65ab7f23-d66e-4543-b276-488ffaaad71c")
	validSecret := "this_is_the_secret_to_use"
	wrongSecret := "this_is_wrong_secret"
	validExpiration := 24 * time.Hour
	passedExpiration := -24 * time.Hour // yesterday

	// Create tests jwt
	jwt, _ := MakeJWT(userId, validSecret, validExpiration)
	expiredJWT, _ := MakeJWT(userId, validSecret, passedExpiration)

	tests := []struct {
		name    string
		jwt     string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid JWT",
			jwt:     jwt,
			secret:  validSecret,
			wantErr: false,
		},
		{
			name:    "Wrong secret JWT",
			jwt:     jwt,
			secret:  wrongSecret,
			wantErr: true,
		},
		{
			name:    "Expired JWT",
			jwt:     expiredJWT,
			secret:  validSecret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.jwt, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
