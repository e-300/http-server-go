package auth

import (
	"testing"
	"time"
	"net/http"
	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}


func TestValidateJwt(t *testing.T){
	userID := uuid.New()
	validTok, _ := MakeJWT(userID, "secret", time.Hour)

	userID2 := uuid.New()
	expiredTok, _ := MakeJWT(userID2, "secret", -time.Hour)

	tests := []struct{
		name        string
		tokenString string 
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name: 		"Valid token",
			tokenString: validTok,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name: 		"Invalid token",
			tokenString: "Invalid Token String",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,

		},
		{
			name: 		"Wrong Secret",
			tokenString: validTok,
			tokenSecret: "Wrong_Secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name: 		"Expired Token",
			tokenString: expiredTok,
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != (tt.wantErr){
				t.Errorf("ValidateJWT() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID{
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestBearer(t *testing.T){

	tests := []struct {
		name 		string
		headers 	http.Header
		wantToken	string
		wantErr		bool
	}{
		{
			name: "Correct Bearer",
			headers: http.Header{"Authorization" : []string{"Bearer sometoken"}},
			wantToken: "sometoken",
			wantErr: false,

		},
		{
			name: "No Bearer Token",
			headers: http.Header{"Authorization": []string{"Bearer"}},
			wantToken: "",
			wantErr: true,
		},
		{
			name: "No Bearer",
			headers: http.Header{},
			wantToken: "",
			wantErr: true,

		},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}