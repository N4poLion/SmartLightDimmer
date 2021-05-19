package main

import (
	"fmt"

	"github.com/N4poLion/SmartLightDimmer/iotserver"
)

func main() {
	fmt.Println("yolo")
	serverIot := iotserver.NewIotServer(4444)
	if err := serverIot.Start(); err != nil {
		fmt.Println("Error happend: ", err)
	}
}
