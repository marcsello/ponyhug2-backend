package views

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"time"
)

const (
	JWTIssuer = "ponyhug2"
)

var jwtSecretKey = []byte(env.StringOrPanic("JWT_SECRET"))

func validateJWT(logger *zap.Logger, tokenString string) (int32, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecretKey, nil
	}, jwt.WithIssuer(JWTIssuer), jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		logger.Error("failure while parsing the token", zap.Error(err))
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("claims has an invalid format")
		logger.Warn("failure reading claims from jwt", zap.Error(err))
		return 0, err
	}

	// TODO: verify expiration manually???

	sub, ok := claims["sub"].(int32)
	if !ok {
		err = fmt.Errorf("sub has an invalid format or missing")
		logger.Error("sub has an invalid type lol", zap.Error(err))
		return 0, err
	}

	return sub, nil

}

func generateToken(sub int32) (string, error) {
	nowUnix := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": JWTIssuer,
		"sub": sub,
		"nbf": nowUnix,
		"iat": nowUnix,
		"exp": time.Now().Add(time.Hour * 336).Unix(), // two weeks
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(jwtSecretKey)
}
