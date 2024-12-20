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
	commentRepository := repository.NewCommentRepositoryMySQL(db)

	authUseCase := usecase.NewAuthUseCaseImpl(userRepository, tokenRepository, transactor)
	postUseCase := usecase.NewPostUsecaseImpl(postRepository, userRepository, transactor)
	commentUseCase := usecase.NewCommentUseCaseImpl(commentRepository, userRepository, postRepository, transactor)

	middleware := NewMiddlewareHandler(authUseCase)
	authGroup := r.Group("")
	postGroup := r.Group("/posts")
	commentGroup := r.Group("/posts/:postID/comments")
	NewAuthHandler(authGroup, authUseCase)
	NewPostHandler(postGroup, middleware, postUseCase)
	NewCommentHandler(commentGroup, middleware, commentUseCase)
	return r, nil
}
