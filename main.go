package main

import (
	"context"
	"fmt"
	"os"

	ory "github.com/ory/client-go"
)

// ใช้คอนเท็กซ์นี้เพื่อเข้าถึง Ory APIs ที่ต้องการ API Key
var oryAuthedContext = context.WithValue(context.Background(), ory.ContextAccessToken, os.Getenv("ORY_API_KEY"))
var namespace = "videos"
var object = "secret_post"
var relation = "view"
var subjectId = "Bob"

func main() {
	// กำหนด payload สำหรับสร้างความสัมพันธ์
	payload := ory.CreateRelationshipBody{
		Namespace: &namespace,
		Object:    &object,
		Relation:  &relation,
		SubjectId: &subjectId,
	}

	// การตั้งค่าเซิร์ฟเวอร์สำหรับ Write API (พอร์ต 4467)
	writeConfig := ory.NewConfiguration()
	writeConfig.Servers = []ory.ServerConfiguration{
		{
			URL: "http://localhost:4467", // Write API
		},
	}
	writeClient := ory.NewAPIClient(writeConfig)

	// การตั้งค่าเซิร์ฟเวอร์สำหรับ Read API (พอร์ต 4466)
	readConfig := ory.NewConfiguration()
	readConfig.Servers = []ory.ServerConfiguration{
		{
			URL: "http://localhost:4466", // Read API
		},
	}
	readClient := ory.NewAPIClient(readConfig)

	// ตรวจสอบการเชื่อมต่อกับ Write API
	_, r, err := writeClient.MetadataAPI.GetVersion(oryAuthedContext).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ข้อผิดพลาดในการเชื่อมต่อกับ Write API: %v\n", err)
		fmt.Fprintf(os.Stderr, "การตอบกลับ: %v\n", r)
		panic("ไม่สามารถเชื่อมต่อกับ Write API ได้")
	}
	fmt.Println("เชื่อมต่อกับ Write API สำเร็จ!")

	// สร้างความสัมพันธ์
	_, r, err = writeClient.RelationshipAPI.CreateRelationship(oryAuthedContext).CreateRelationshipBody(payload).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ข้อผิดพลาดในการสร้างความสัมพันธ์: %v\n", err)
		fmt.Fprintf(os.Stderr, "การตอบกลับ: %v\n", r)
		panic("ไม่สามารถสร้างความสัมพันธ์ได้")
	}
	fmt.Println("สร้างความสัมพันธ์สำเร็จ!")

	// ตรวจสอบการเชื่อมต่อกับ Read API
	_, r, err = readClient.MetadataAPI.GetVersion(oryAuthedContext).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ข้อผิดพลาดในการเชื่อมต่อกับ Read API: %v\n", err)
		fmt.Fprintf(os.Stderr, "การตอบกลับ: %v\n", r)
		panic("ไม่สามารถเชื่อมต่อกับ Read API ได้")
	}
	fmt.Println("เชื่อมต่อกับ Read API สำเร็จ!")

	// ตรวจสอบสิทธิ์
	check, r, err := readClient.PermissionAPI.CheckPermission(oryAuthedContext).
		Namespace(namespace).
		Object(object).
		Relation(relation).
		SubjectId(subjectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ข้อผิดพลาดในการตรวจสอบสิทธิ์: %v\n", err)
		fmt.Fprintf(os.Stderr, "การตอบกลับ: %v\n", r)
		panic("ไม่สามารถตรวจสอบสิทธิ์ได้")
	}
	if check.Allowed {
		fmt.Printf("%s สามารถ %s %s ได้\n", subjectId, relation, object)
	} else {
		fmt.Printf("%s ไม่สามารถ %s %s ได้\n", subjectId, relation, object)
	}
}