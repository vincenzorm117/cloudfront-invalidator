package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func FormatError(err error) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse {
        Body: string(err.Error()),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        StatusCode: 500,
    }, nil
}

func Handler (ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    cloudfrontDistId := e.QueryStringParameters["cfid"]

    fmt.Println("cloudfrontDistId: " + string(cloudfrontDistId))

    // Initialize cloudfront instance with session
    sess := session.Must(session.NewSession())
    svc := cloudfront.New(sess)

    // Perform Cloudfront Invalidation
    now := time.Now()
    resp, err := svc.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(cloudfrontDistId),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(now.UTC().Format(time.RFC3339)),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items: []*string{
					aws.String("/*"),
				},
			},
		},
	})

    if err != nil {
        return FormatError(err)
    }

    json, err := json.Marshal(resp)

    if err != nil {
        return FormatError(err)
    }

    return events.APIGatewayProxyResponse {
		Body: string(json),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200,
	}, nil
}

func main() {
    lambda.Start(Handler)
}