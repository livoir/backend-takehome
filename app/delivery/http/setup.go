package http

import (
	"app/repository"
	"app/usecase"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, jwtPrivateKey, jwtPublicKey string) (*gin.Engine, error) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))
	transactor := repository.NewSQLTransactor(db)
	tokenRepository := repository.NewTokenRepositoryJWT(jwtPrivateKey, jwtPublicKey)
	userRepository := repository.NewUserRepositoryMySQL(db)
	postRepository := repository.NewPostRepositoryMySQL(db)

	authUseCase := usecase.NewAuthUseCaseImpl(userRepository, tokenRepository, transactor)
	postUseCase := usecase.NewPostUsecaseImpl(postRepository, transactor)

	middleware := NewMiddlewareHandler(authUseCase)
	NewAuthHandler(r.Group(""), authUseCase)
	NewPostHandler(r.Group("/posts"), middleware, postUseCase)
	return r, nil
}
