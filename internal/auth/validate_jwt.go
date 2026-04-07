package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(*jwt.Token)(interface{}, error){
		byteToken := []byte(tokenSecret)
		return byteToken, nil
	})

	if err != nil{
		return uuid.Nil, err
	}

	subject , err := token.Claims.GetSubject()
	if err != nil{
		return uuid.Nil, err
	}
	user, err := uuid.Parse(subject)
	if err != nil{
		return user, err
	}

	return user, nil
	
}
