build: main

main: codegen
	go build cmd/comment-like-server/main.go

codegen: cmd/comment-like-server/main.go

cmd/comment-like-server/main.go: api/comment_like/comment_like.yaml
	swagger generate server -f api/comment_like/comment_like.yaml

