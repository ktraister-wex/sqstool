package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)


//this is working to pull down queue names
func listQueues(ENV string) []string {
    	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files

	sess, err := session.NewSession(&aws.Config{
	    Region: aws.String("us-east-1")},
	)

        // Create a SQS service client.
        svc := sqs.New(sess)

	//have to create a session object first
	output, err := svc.ListQueues(&sqs.ListQueuesInput{
	    QueueNamePrefix: aws.String(ENV),
            })
	if err != nil { panic(err) }

	queues := output.QueueUrls
	final_queues := []string{}

	for _, i := range queues {
	    fmt.Println(string(*i))
	    final_queues = append(final_queues, *i)
        }
	return final_queues
}

func writeMsg(QUEUE, MSG string) int {
        //just return 0 for now so we don't actually send anything
        return 0    

        sess, err := session.NewSession(&aws.Config{
	    Region: aws.String("us-east-1")},
	)

        // Create a SQS service client.
        svc := sqs.New(sess)
	_, err = svc.SendMessage(&sqs.SendMessageInput{
	    DelaySeconds: aws.Int64(10),
	    /*
	    MessageAttributes: map[string]*sqs.MessageAttributeValue{
		"Title": &sqs.MessageAttributeValue{
		    DataType:    aws.String("String"),
		    StringValue: aws.String("The Whistler"),
		},
		"Author": &sqs.MessageAttributeValue{
		    DataType:    aws.String("String"),
		    StringValue: aws.String("John Grisham"),
		},
		"WeeksOn": &sqs.MessageAttributeValue{
		    DataType:    aws.String("Number"),
		    StringValue: aws.String("6"),
		},
	    },
	    */
	    MessageBody: aws.String(MSG),
	    QueueUrl:    &QUEUE,
	})

	if err != nil {
	    fmt.Println("Error", err)
	    return 1
	} else {
	    return 0
	}
}

func main () {
    env := os.Getenv("ENVIRONMENT")

    fmt.Printf("Available Queues in %s:\n", env)
    fmt.Println("------------------------")
    queues := listQueues(env)
    fmt.Println("                        ")
    fmt.Println("                        ")
    fmt.Println("Do you want to send a message to any of these queues?[Y/Yes]")

    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    // convert CRLF to LF
    text = strings.Replace(text, "\n", "", -1)

    if strings.Compare("Yes", text) == 0  || strings.Compare("Y", text) == 0{
	fmt.Println("Select Queue:")
        fmt.Println("------------------------")
        text, _ = reader.ReadString('\n')
        text = strings.Replace(text, "\n", "", -1)
	target_queue := ""
	for _, v := range queues {
	    if strings.Compare(text, v) == 0 {
		target_queue = v
		break
	    }
	}
        if strings.Compare(target_queue, "") == 0 {
	    fmt.Println("Requeusted queue not found! Exiting!")
	    os.Exit(1)
	} else {
	    fmt.Println("Sending message to queue -> ", target_queue)
	    fmt.Println("")
	    fmt.Println("Input message body")
	    fmt.Println("---------------------")
            text, _ = reader.ReadString('\n')
            text = strings.Replace(text, "\n", "", -1)
	    rc := writeMsg(target_queue, text)
	    if rc != 0 {
		os.Exit(2)
            } else {
		fmt.Println("Message sent to queue successfully! Have a nice day!")
	    }
	}
    }
}
