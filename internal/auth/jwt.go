package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type jwtClaim struct{
	jwt.RegisteredClaims
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	
	now := time.Now().UTC()

	regisClaim := jwtClaim{
		jwt.RegisteredClaims{
			Issuer: "chirpy-access",
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt:  jwt.NewNumericDate(now.Add(expiresIn)),
			Subject: userID.String(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &regisClaim)

	bytetoken := []byte(tokenSecret)
	signedtoken, err := token.SignedString(bytetoken)
	if err != nil{
		return "", err
	}

	return signedtoken, nil
}