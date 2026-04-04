package auth

import (
	"github.com/google/uuid"
	"time"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	return "",nil
}