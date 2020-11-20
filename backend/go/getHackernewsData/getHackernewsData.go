package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"pht/hndata"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endPoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//LambdaEvent event describing data
type LambdaEvent struct {
	storyType string
}

func getFirebaseData(ctx context.Context, event LambdaEvent) {
	fmt.Prinln(ctx)
	htmlData, storyType := hndata.FetchHackernewsData(event.storyType)

	if htmlData != "" {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(endPoints.ApSouth1RegionID),
		}))
		sqsQueue := sqs.New(sess)

		jsonBody, marshalError := json.Marshal(struct {
			html     string
			typeData string
		}{htmlData, storyType})
		if marshalError != nil {
			panic(marshalError)
		}
		message := string(jsonBody)
		emailSenderURL := os.Getenv("EMAIL_SENDER_URL")
		sendMessageData := sqs.SendMessageInput{
			MessageBody: message,
			QueueUrl:    emailSenderURL,
		}
		sqsQueue.SendMessage(&sendMessageData)
	}

}

func main() {

	lambda.Start()
	// htmlFile.Write([]byte(htmlTemplate))

}

func prepareArgs(args []string) map[string]string {
	parsedArgs := make(map[string]string)
	for _, argument := range args {

		keyVal := strings.Split(argument, "=")
		if len(keyVal) > 1 {

			parsedArgs[keyVal[0]] = strings.Join(keyVal[1:], ",")
		} else {
			parsedArgs[keyVal[0]] = ""
		}
	}
	return parsedArgs
}
