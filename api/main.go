package api

import (
	"context"
	"fmt"
	"net/http"
	"new-refferal/models"
	"reflect"
	"strconv"
	"time"

	"firebase.google.com/go/v4/auth"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"new-refferal/conf"
	"new-refferal/helpers/null"
	"new-refferal/log"
	"new-refferal/services"
)

const (
	firebaseFilePath = "./firebase-auth.json"
	userContextKey   = "user"
)

type (
	API struct {
		router       *gin.Engine
		server       *http.Server
		cfg          conf.Config
		services     services.Service
		queryDecoder *schema.Decoder
		auth         *auth.Client
	}

	// Route stores an API route data
	Route struct {
		Path       string
		Method     string
		Func       func(http.ResponseWriter, *http.Request)
		Middleware []negroni.HandlerFunc
	}
)

func NewAPI(cfg conf.Config, s services.Service) (*API, error) {
	queryDecoder := schema.NewDecoder()
	queryDecoder.IgnoreUnknownKeys(true)
	queryDecoder.RegisterConverter(null.Time{}, func(s string) reflect.Value {
		timestamp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}
		}
		t := null.NewTime(time.Unix(timestamp, 0))
		return reflect.ValueOf(t)
	})
	api := &API{
		cfg:          cfg,
		services:     s,
		queryDecoder: queryDecoder,
	}

	api.initialize()
	return api, nil
}

// Run starts the http server and binds the handlers.
func (api *API) Run() error {
	return api.startServe()
}

func (api *API) Stop() error {
	return api.server.Shutdown(context.Background())
}

func (api *API) Title() string {
	return "API"
}

func (api *API) initialize() {
	api.router = gin.Default()

	api.router.Use(gin.Logger())

	api.router.Use(gin.Recovery())

	api.router.Use(cors.New(cors.Config{
		AllowOrigins:     api.cfg.API.CORSAllowedOrigins,
		AllowCredentials: true,
		AllowMethods: []string{
			http.MethodPost, http.MethodHead, http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete,
		},
		AllowHeaders: []string{
			"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "User-Env", "Access-Control-Request-Headers", "Access-Control-Request-Method",
		},
	}))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",
		Key:              []byte("secret key"),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24 * 7,
		IdentityKey:      "wallet_address",
		SigningAlgorithm: "HS512",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"wallet_address": v.WalletAddress,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				WalletAddress: claims["wallet_address"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var user models.User
			if err := c.ShouldBind(&user); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			usr, err := api.services.LogInOrRegister(&user)
			if err != nil {
				return nil, err
			}

			return usr, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			usr, err := api.services.LogInOrRegister(data.(*models.User))
			if err != nil || usr.ID == 0 {
				return false
			}

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// public routes
	api.router.GET("/", api.Index)
	api.router.GET("/health", api.Health)

	api.router.POST("/register", authMiddleware.LoginHandler)
	api.router.POST("/delegate", api.Delegate)

	api.router.GET("/rewards_total_stats", api.GetTotalRewardStats)
	api.router.GET("/all_rewards", api.GetAllRewards)
	api.router.GET("/invitations_stats", api.GetInvitationsStats)

	api.router.GET("/my_link", api.GetMyLink)
	api.router.GET("/invited", api.GetInvitedFriends)

	api.router.GET("/open_box", api.OpenBox)

	api.router.POST("/gets", api.Gets)

	authGroup := api.router.Group("/auth")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/refresh", authMiddleware.RefreshHandler)
	}

	adminGroup := authGroup.Group("/admin")
	adminGroup.Use(api.AdminMiddleware())
	{
		adminGroup.GET("/total_stats", api.GetTotalStats)
		adminGroup.GET("/total_stake_stats", api.GetTotalStakeStats)
		adminGroup.GET("/friends_stake_stats", api.GetFriendsStakeStats)
		adminGroup.GET("/reward_payment_stats", api.GetRewardPaymentStats)
	}

	//mGroup := api.router.Group("/m")
	//mGroup.Use(authMiddleware.MiddlewareFunc() /*api.SomeMiddleware()*/)
	//{
	//	mGroup.POST("/name", api.Name)
	//	mGroup.GET("/read/:id", api.Read)
	//}
	api.server = &http.Server{Addr: fmt.Sprintf(":%d", api.cfg.API.ListenOnPort), Handler: api.router}
}

func (api *API) startServe() error {
	log.Info("Start listening server on port", zap.Uint64("port", api.cfg.API.ListenOnPort))
	err := api.server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Warn("API server was closed")
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot run API service: %s", err.Error())
	}
	return nil
}
