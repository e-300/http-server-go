package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	
	now := time.Now().UTC()
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt:  jwt.NewNumericDate(now.Add(expiresIn)),
		Subject: userID.String(),

	})

	bytetoken := []byte(tokenSecret)
	signedtoken, err := token.SignedString(bytetoken)
	if err != nil{
		return "", err
	}

	return signedtoken, nil
}