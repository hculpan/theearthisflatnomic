package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hculpan/theearthisflatnomic/entity"
)

var jwtSecret []byte = []byte{}

// Claims defines the keys we want to put
// in the JWT token
type Claims struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	FullName    string `json:"fullname"`
	DisplayName string `json:"displayname"`
	jwt.StandardClaims
}

// Authenticate authenticates the user
func Authenticate(username, password string) (*entity.User, error) {
	user, found := entity.FindUserByUsername(username)
	if !found {
		return nil, fmt.Errorf("Invalid username/password")
	}
	result := user.VerifyPassword(password)
	if result == nil {
		user.LastLogin = time.Now()
		if err := user.Save(); err != nil {
			result = err
		}
	}
	return user, result
}

func getSecretKey() {
	if len(jwtSecret) == 0 {
		if os.Getenv("TEIFN_SECRET_KEY") == "" {
			panic("TEIFN_SECRET_KEY is not setup correctly")
		}
		jwtSecret = []byte(os.Getenv("TEIFN_SECRET_KEY"))
	}
}

// CreateToken create a jwt token
func CreateToken(u entity.User) (string, error) {
	getSecretKey()

	expireTime := time.Now().Add((90 * 24) * time.Hour)
	claims := Claims{
		Username:    u.Username,
		Password:    string(u.Password),
		FullName:    u.FullName,
		DisplayName: u.DisplayName,
		StandardClaims: jwt.StandardClaims{
			//Expiration time
			ExpiresAt: expireTime.Unix(),
			//Designated token publisher
			Issuer: "the_earth_is_flat_nomic",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(jwtSecret))
}

// DecodeToken decodes a JWT token
func DecodeToken(t string) (*Claims, error) {
	getSecretKey()

	result := &Claims{}
	tkn, err := jwt.ParseWithClaims(t, result, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return result, err
	}

	if !tkn.Valid {
		return result, fmt.Errorf("Invalid token")
	}

	return result, nil
}
