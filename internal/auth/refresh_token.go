package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error){
	// byte slice of len 32 
	key := make([]byte, 32)
	// read random bytes into key in place
	_, err := rand.Read(key)
	if err != nil{
		return "", err
	}

	hexKey := hex.EncodeToString(key)

	return hexKey, nil
}



