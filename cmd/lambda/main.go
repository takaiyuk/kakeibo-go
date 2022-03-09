package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/takaiyuk/kakeibo-go/pkg"
)

func main() {
	lambda.Start(pkg.Kakeibo)
}
