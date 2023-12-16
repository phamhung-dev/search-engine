package tokenpvd

import (
	"os"
	"time"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/golang-jwt/jwt"
)

type jwtProvider struct {
	accessTokenKey  string
	refreshTokenKey string
}

func NewJWTProvider() *jwtProvider {
	accessTokenKey := os.Getenv("JWT_ACCESS_TOKEN_KEY")
	refreshTokenKey := os.Getenv("JWT_REFRESH_TOKEN_KEY")

	if accessTokenKey == "" || refreshTokenKey == "" {
		beeLogger.Log.Fatal(ErrProviderIsNotConfigured.Error())
	}

	return &jwtProvider{
		accessTokenKey:  accessTokenKey,
		refreshTokenKey: refreshTokenKey,
	}
}

type claims struct {
	Payload TokenPayload `json:"payload"`
	jwt.StandardClaims
}

func (provider *jwtProvider) Generate(payload TokenPayload, expiredIn int) (*Token, error) {
	accessToken, err := generateToken(payload, expiredIn, provider.accessTokenKey)

	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(payload, 2*expiredIn, provider.refreshTokenKey)

	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  accessToken,
		CreatedAt:    time.Now(),
		ExpiredIn:    expiredIn,
		RefreshToken: refreshToken,
	}, nil
}

func (provider *jwtProvider) ValidateAccessToken(token string) (*TokenPayload, error) {
	return validateToken(token, provider.accessTokenKey)
}

func (provider *jwtProvider) ValidateRefreshToken(token string) (*TokenPayload, error) {
	return validateToken(token, provider.refreshTokenKey)
}

func generateToken(payload TokenPayload, expiredIn int, secretKey string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(expiredIn)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	})

	token, err := t.SignedString([]byte(secretKey))

	if err != nil {
		return "", ErrEncodingToken
	}

	return token, nil
}

func validateToken(token string, secretKey string) (*TokenPayload, error) {
	res, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !res.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := res.Claims.(*claims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return &claims.Payload, nil
}
