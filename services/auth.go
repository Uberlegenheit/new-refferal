package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"new-refferal/models"
	"os"
	"strings"
	"time"
)

func (s *ServiceFacade) CreateToken(walletAddr string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["wallet_address"] = walletAddr
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["wallet_address"] = walletAddr
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *ServiceFacade) CreateAuth(walletAddr string, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := s.dao.AddAuthToken(td.AccessUuid, walletAddr, at.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(td.RefreshUuid, walletAddr, rt.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(fmt.Sprintf("%s_access", td.RefreshUuid), td.AccessUuid, rt.Sub(now))
	if err != nil {
		return err
	}
	err = s.dao.AddAuthToken(fmt.Sprintf("%s_refresh", td.AccessUuid), td.RefreshUuid, rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceFacade) ExtractTokenMetadata(c *gin.Context) (*models.AccessDetails, error) {
	token, err := s.VerifyToken(c.Request)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		walletAddr := fmt.Sprintf("%s", claims["wallet_address"])
		if err != nil {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			WalletAddr: walletAddr,
		}, nil
	}
	return nil, err
}

func (s *ServiceFacade) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	if tokenString == "" {
		return nil, fmt.Errorf("no token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *ServiceFacade) Refresh(r *http.Request) (*models.TokenDetails, error) {
	refreshToken := r.Header.Get("Authorization")
	parts := strings.Split(refreshToken, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("error: %s", "cannot get the refresh token")
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("error: %s", "invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("error: %s", "invalid token")
		}
		walletAddr := fmt.Sprintf("%s", claims["wallet_address"])
		accessUuid, ok, err := s.dao.GetAuthToken(fmt.Sprintf("%s_access", refreshUuid))
		if err != nil || !ok {
			return nil, fmt.Errorf("error: %s", "cannot get access token: invalid refresh_access token")
		}
		err = s.DeleteAuth(refreshUuid,
			accessUuid.(string),
			fmt.Sprintf("%s_access", refreshUuid),
			fmt.Sprintf("%s_refresh", accessUuid.(string)))
		if err != nil {
			return nil, fmt.Errorf("error: %s", "invalid token provided")
		}

		ts, err := s.CreateToken(walletAddr)
		if err != nil {
			return nil, fmt.Errorf("error: %s", "cannot create token")
		}

		saveErr := s.CreateAuth(walletAddr, ts)
		if saveErr != nil {
			return nil, fmt.Errorf("error: %s", "cannot create auth")
		}

		return ts, nil
	}

	return nil, fmt.Errorf("error: %s", "cannot refresh tokens")
}

func extractToken(r *http.Request) string {
	accessToken := r.Header.Get("Authorization")
	token := strings.Split(accessToken, " ")
	if len(token) != 2 {
		return ""
	} else if token[0] != "Bearer" {
		return ""
	}
	return token[1]
}
