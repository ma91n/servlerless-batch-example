#############
# Build
#############
.PHONY: build

build:
	go mod download
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -buildid=" -trimpath -o bin/lambda main.go

deploy: build
	zip -j bin/lambda.zip bin/lambda
	aws lambda update-function-code --profile mano --function-name scatter --zip-file fileb://bin/lambda.zip
