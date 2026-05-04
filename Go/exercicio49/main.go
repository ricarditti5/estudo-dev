package main

import "fmt"

type APP struct {
	AppName string
	Stack   []string
	WebServer
}
type WebServer struct {
	Port        int
	Host        string
	EnableHttps bool
}

func main() {
	stM := APP{
		AppName: "Student Manager",
		Stack:   []string{"React", "Typescript", "Shadcn", "Golang"},
		WebServer: WebServer{
			Port:        8080,
			Host:        "localhost:",
			EnableHttps: true,
		},
	}
	fmt.Println("The ", stM.AppName, "build with ", stM.Stack, "is running in", stM.Host, stM.Port)
}
