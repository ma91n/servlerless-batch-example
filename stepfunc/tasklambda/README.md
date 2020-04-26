# Deploy

```bash
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -buildid=" -o ./bin/lambda main.go
zip -j bin/lambda.zip bin/lambda

# 初回
aws lambda create-function --profile my_servlerless_batch --region ap-northeast-1 \
  --function-name tasklambda1 --zip-file fileb://bin/lambda.zip --handler lambda --runtime go1.x \
  --role arn:aws:iam::123456789012:role/serverless-batch-lambdarole

# 2回目以降
aws lambda update-function-code --profile my_servlerless_batch --region ap-northeast-1 --function-name tasklambda --zip-file fileb://bin/lambda.zip
```
