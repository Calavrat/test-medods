package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Calavrat/TestMedods/internal/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AutService struct {
	rry store.Authorization
}

func NewAuthService(rry store.Authorization) *AutService {
	return &AutService{rry: rry}
}

type tokenDetails struct {
	AccesToken   string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func (s *AutService) GenerateToken(userid string) (*tokenDetails, error) {
	//проверка есть ли user в бд по id
	user, err := s.rry.GetUser(userid)
	if err != nil {
		return nil, fmt.Errorf("error is not user by this id: %v", err)
	}

	td := &tokenDetails{}
	td.AtExpires = time.Now().Add(15 * time.Minute).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(24 * 7 * time.Hour).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading env variables: %v", err)
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.ID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["access_uuid"] = td.AccessUuid
	rtClaims["user_id"] = user.ID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)

	td.AccesToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("error access %v", err)
	}
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("Error refresh %v", err)
	}

	user.RefreshTk, err = HashToken(td.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error bcrypt refresh token")
	}

	if err := s.rry.AddRefresh(user); err != nil {
		return nil, fmt.Errorf("error when adding refresh token %v", err)
	}

	return td, err
}

func (s *AutService) ParseTokens(ts *map[string]string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("error loading env variables: %v", err)
	}
	//access claims
	AccessParse := make(map[string]string)

	token, _ := jwt.Parse((*ts)["access_token"], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Access token claims are not of type")
	}

	AccessParse["a_user_id"] = claims["user_id"].(string)
	AccessParse["a_uuid"] = claims["access_uuid"].(string)

	//refresh claims
	token, err := jwt.Parse((*ts)["refresh_token"], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return "", fmt.Errorf("error refresh token parse %v", err)
	}

	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Refresh token claims are not of type")
	}

	//проверка на взаимосвязь токенов
	if AccessParse["a_user_id"] != claims["user_id"].(string) || AccessParse["a_uuid"] != claims["access_uuid"].(string) {
		return "", errors.New("tokens are not linked")
	}

	return claims["user_id"].(string), nil
}

func HashToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		return "", fmt.Errorf("error create bcrypt hash")
	}

	return string(bytes), nil
}
