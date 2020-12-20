package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"pht/hndata"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//LambdaEvent event describing data
type LambdaEvent struct {
	storyType string
}

func getFirebaseData(ctx context.Context, event LambdaEvent) {
	fmt.Println(ctx)
	htmlData, storyType := hndata.FetchHackernewsData(event.storyType)

	if htmlData != "" {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(endpoints.ApSouth1RegionID),
		}))
		sqsQueue := sqs.New(sess)

		jsonBody, marshalError := json.Marshal(struct {
			HTML     string `json:"html"`
			TypeData string `json:"type"`
		}{htmlData, storyType})
		if marshalError != nil {
			panic(marshalError)
		}
		message := string(jsonBody)
		emailSenderURL := os.Getenv("EMAIL_SENDER_URL")
		sendMessageData := sqs.SendMessageInput{
			MessageBody: &message,
			QueueUrl:    &emailSenderURL,
		}
		sqsQueue.SendMessage(&sendMessageData)
	}

}

func main() {

	env := os.Getenv("ENV")
	if env == "PROD" {
		lambda.Start(getFirebaseData)
	} else {
		storyType := flag.String("storyType", "top", "Type of hn stories")
		flag.Parse()
		fmt.Println(*storyType)

		hndata.FetchHackernewsData(*storyType)
	}
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
