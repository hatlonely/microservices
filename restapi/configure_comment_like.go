package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/tylerb/graceful"

	"microservices/restapi/operations"
	"microservices/restapi/operations/comment"
	"microservices/restapi/operations/like"
	"microservices/internal/comment_like"
	"microservices/models"
	"strings"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../api/comment_like/comment_like.yaml

func configureFlags(api *operations.CommentLikeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CommentLikeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.LikeCountLikeHandler = like.CountLikeHandlerFunc(func(params like.CountLikeParams) middleware.Responder {
		count, err := comment_like.CountLike(params.Title)
		if err != nil {
			return like.NewCountLikeBadRequest()
		}
		return like.NewCountLikeOK().WithPayload(&models.CountLikeModel{
			Count: count,
			Title: params.Title,
		})
	})

	api.CommentDoCommentHandler = comment.DoCommentHandlerFunc(func(params comment.DoCommentParams) middleware.Responder {
		ip := strings.Split(params.HTTPRequest.RemoteAddr, ":")[0]
		ua := params.HTTPRequest.UserAgent()
		var nickname, mail string
		if params.Nickname != nil {
			nickname = *params.Nickname
		}
		if params.Mail != nil {
			mail = *params.Mail
		}

		if err := comment_like.DoComment(ip, ua, params.Title, params.Comment, nickname, mail); err != nil {
			return comment.NewDoCommentBadRequest().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		return comment.NewDoCommentBadRequest().WithPayload(&models.ErrorModel{
			Code:    http.StatusOK,
			Message: "ok",
		})
	})

	api.LikeDoLikeHandler = like.DoLikeHandlerFunc(func(params like.DoLikeParams) middleware.Responder {
		ip := strings.Split(params.HTTPRequest.RemoteAddr, ":")[0]
		ua := params.HTTPRequest.UserAgent()

		if err := comment_like.DoLike(ip, ua, params.Title); err != nil {
			return like.NewDoLikeBadRequest().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		return like.NewDoLikeBadRequest().WithPayload(&models.ErrorModel{
			Code:    http.StatusOK,
			Message: "ok",
		})
	})

	api.LikeDoUnlikeHandler = like.DoUnlikeHandlerFunc(func(params like.DoUnlikeParams) middleware.Responder {
		ip := strings.Split(params.HTTPRequest.RemoteAddr, ":")[0]
		ua := params.HTTPRequest.UserAgent()

		if err := comment_like.DoUnlike(ip, ua, params.Title); err != nil {
			return like.NewDoUnlikeBadRequest().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		return like.NewDoUnlikeBadRequest().WithPayload(&models.ErrorModel{
			Code:    http.StatusOK,
			Message: "ok",
		})
	})

	api.CommentShowCommentHandler = comment.ShowCommentHandlerFunc(func(params comment.ShowCommentParams) middleware.Responder {
		//comments, err := comment_like.ShowComment(params.Title)
		//code := int64(200)
		//message := "ok"
		//if err != nil {
		//	code = int64(400)
		//	message = err.Error()
		//	return comment.NewShowCommentBadRequest().WithPayload(&models.ErrorModel{
		//		Code:    &code,
		//		Message: &message,
		//	})
		//}

		return middleware.NotImplemented("operation comment.ShowComment has not yet been implemented")
	})

	api.LikeShowLikeHandler = like.ShowLikeHandlerFunc(func(params like.ShowLikeParams) middleware.Responder {
		ip := strings.Split(params.HTTPRequest.RemoteAddr, ":")[0]
		ua := params.HTTPRequest.UserAgent()

		isLike, err := comment_like.ShowLike(ip, ua, params.Title)
		if err != nil {
			return like.NewShowLikeBadRequest().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}

		return like.NewShowLikeOK().WithPayload(&models.ShowLikeModel{
			IP: ip,
			Ua: ua,
			Title: params.Title,
			Islike: isLike,
		})
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
