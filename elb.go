package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

var loadBalancerName = flag.String("name", "", "the access point name of the load balancer")
var registerOperation = flag.Bool("register", false, "register this instance in the load balancer")
var deregisterOperation = flag.Bool("deregister", false, "deregister this instance in the load balancer")

func register(instanceID string, svc *elb.ELB) {
	params := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances: []*elb.Instance{ // Required
			{ // Required
				InstanceId: aws.String(instanceID),
			},
			// More values...
		},
		LoadBalancerName: aws.String(*loadBalancerName), // Required
	}
	_, err := svc.RegisterInstancesWithLoadBalancer(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Instance", instanceID, "registered.")
}

func deregister(instanceID string, svc *elb.ELB) {
	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances: []*elb.Instance{ // Required
			{ // Required
				InstanceId: aws.String(instanceID),
			},
			// More values...
		},
		LoadBalancerName: aws.String(*loadBalancerName), // Required
	}
	_, err := svc.DeregisterInstancesFromLoadBalancer(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Instance", instanceID, "deregistered.")
}

func main() {
	flag.Parse()

	session := session.New(&aws.Config{Region: aws.String("us-east-1")})
	svc := elb.New(session)
	metadata := ec2metadata.New(session)

	instanceID, _ := metadata.GetMetadata("instance-id")

	if instanceID == "unknown" {
		fmt.Printf("Unable to determine AWS instance ID.\n")
		os.Exit(1)
	}

	if *registerOperation == true && *deregisterOperation == false {
		register(instanceID, svc)
		os.Exit(0)
	}

	if *deregisterOperation == true && *registerOperation == false {
		deregister(instanceID, svc)
		os.Exit(0)
	}

	fmt.Printf("Unable to parse command flags.\n")
	os.Exit(1)
}