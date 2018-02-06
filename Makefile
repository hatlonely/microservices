codegen: cmd/comment_like/main.go

cmd/comment_like/main.go: api/comment_like/comment_like.yaml
	swagger generate server -f api/comment_like/comment_like.yaml

