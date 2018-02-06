codegen: cmd/hatlonely-me-server/main.go

cmd/hatlonely-me-server/main.go: api/hatlonely.me/hatlonely.me.yaml
	swagger generate server -f api/hatlonely.me/hatlonely.me.yaml	

