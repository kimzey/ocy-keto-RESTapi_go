package main

import (
	"context"
	"fmt"
	"os"

	ory "github.com/ory/client-go"
)

var oryAuthedContext = context.WithValue(context.Background(), ory.ContextAccessToken, os.Getenv("ORY_API_KEY"))

func main() {
	// ตั้งค่า configuration สำหรับ Read API
	readConfig := ory.NewConfiguration()
	readConfig.Servers = []ory.ServerConfiguration{
		{
			URL: "http://localhost:4466",
		},
	}
	readClient := ory.NewAPIClient(readConfig)

	// ตรวจสอบการเชื่อมต่อ
	_, r, err := readClient.MetadataAPI.GetVersion(oryAuthedContext).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to Read API: %v\n", err)
		fmt.Fprintf(os.Stderr, "Response: %v\n", r)
		panic("Could not connect to Read API")
	}
	fmt.Println("Successfully connected to Read API!")

	// ตรวจสอบสิทธิ์โดยใช้ CheckPermission API
	checkPermissionEfficient(readClient, "ccc2", "write")
	checkPermissionEfficient(readClient, "ccc2", "read")
	checkPermissionEfficient(readClient, "asd123", "read")
	checkPermissionEfficient(readClient, "asd123", "create")
}

func checkPermissionEfficient(client *ory.APIClient, subjectId string, permission string) {
	// ใช้ CheckPermission API ซึ่งจัดการกับการตรวจสอบแบบ recursive โดยอัตโนมัติ
	check, r, err := client.PermissionAPI.CheckPermission(oryAuthedContext).
		Namespace("permission").
		Object(permission).
		Relation("role").
		SubjectId(subjectId).
		Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking permission: %v\n", err)
		fmt.Fprintf(os.Stderr, "Response: %v\n", r)
		return
	}

	if check.Allowed {
		fmt.Printf("%s สามารถ %s ได้\n", subjectId, permission)
	} else {
		fmt.Printf("%s ไม่สามารถ %s ได้\n", subjectId, permission)
	}
}