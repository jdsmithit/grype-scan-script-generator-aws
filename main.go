package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func main() {
	awsAccountNumber := "1"
	argLength := len(os.Args[1:])
	fmt.Printf("Arg length is %d\n", argLength)

	for i, region := range os.Args[1:] {
		fmt.Printf("Region %d is %s\n", i+1, region)

		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String(region),
		},
		)

		svc := ecr.New(sess)

		input := &ecr.DescribeRepositoriesInput{}

		result, err := svc.DescribeRepositories(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ecr.ErrCodeServerException:
					fmt.Println(ecr.ErrCodeServerException, aerr.Error())
				case ecr.ErrCodeInvalidParameterException:
					fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
				case ecr.ErrCodeRepositoryNotFoundException:
					fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}

		for _, repos := range result.Repositories {
			ECR_BatchGetImage(*repos.RepositoryName, region, *sess)
		}
	}
}

func ECR_BatchGetImage(repositoryName string, region string, sess session.Session) {
	svc := ecr.New(&sess)

	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
		MaxResults:     aws.Int64(100),
		RegistryId:     aws.String(awsAccountNumber),
	}

	result, err := svc.ListImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	for _, image := range result.ImageIds {
		image := awsAccountNumber + ".dkr.ecr.eu-central-1.amazonaws.com/" + repositoryName + "@" + *image.ImageDigest
		logName := strings.ReplaceAll(strings.ReplaceAll(image, "/", "-"), ":", "-") + ".log"
		scanCmd := "/usr/local/bin/grype -s AllLayers"
		fmt.Println(scanCmd, image, "--file", logName)
	}
}
