# Build and zip s3 image store function
build_s3_store:
	@ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/s3_store internal/s3_store/handler/handler.go
	@ zip -jrm build/s3_store/s3_store build/s3_store

# Build and zip compress image function