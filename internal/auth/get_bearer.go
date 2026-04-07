package auth

import (
	"net/http"
	"strings"
	"errors"
)

func GetBearerToken(headers http.Header) (string, error){
	token_string := headers.Get("Authorization")
	if (token_string == "") || (!strings.HasPrefix(token_string, "Bearer")){
		return "", errors.New("No Bearer")
	}
	ts_noPrefix := strings.TrimPrefix(token_string, "Bearer")
	if ts_noPrefix == ""{
		return "", errors.New("Empty Token")
	}
	ts_noWhiteSpace := strings.TrimSpace(ts_noPrefix)
	return ts_noWhiteSpace, nil
}