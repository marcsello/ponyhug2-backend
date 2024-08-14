package views

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	JWTIssuer = "ponyhug2"
)

var jwtSecretKey = []byte(env.StringOrPanic("JWT_SECRET"))

func validateJWT(logger *zap.Logger, tokenString string) (int16, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecretKey, nil
	}, jwt.WithIssuer(JWTIssuer), jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	if err != nil {
		logger.Error("failure while parsing the token", zap.Error(err))
		return 0, err
	}

	if !token.Valid {
		err = fmt.Errorf("token is invalid")
		logger.Warn("the token is not valid", zap.Error(err))
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("claims has an invalid format")
		logger.Warn("failure reading claims from jwt", zap.Error(err))
		return 0, err
	}

	subStr, ok := claims["sub"].(string)
	if !ok {
		err = fmt.Errorf("sub has an invalid format or missing")
		logger.Error("sub has an invalid type lol", zap.Error(err))
		return 0, err
	}

	var sub int64
	sub, err = strconv.ParseInt(subStr, 10, 16)
	if err != nil {
		logger.Error("sub is not a valid int16", zap.Error(err))
		return 0, err
	}

	return int16(sub), nil

}

func generateToken(sub int16) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    JWTIssuer,
		Subject:   fmt.Sprintf("%d", sub),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 336)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(jwtSecretKey)
}
