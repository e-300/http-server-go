package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)
// hashing password 
func HashPassword(password string) (string, error){
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil{
		log.Fatal(err)
	}
	return hash, nil
}
// Compare the pw in HTTP request with the password that is stored in the database
func CheckPasswordHash(password, hash string) (bool, error){
	b , err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil{
		log.Fatal(err)
	}
	return b, nil
}