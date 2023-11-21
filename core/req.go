package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var fu_channel = make(chan string, 2)

func SaveRequest(req []byte) error {
	file, err := os.CreateTemp("/tmp/sqlmap", "sqlmap")
	if err != nil {
		file.Close()
		return err
	}
	if err != nil {
		log.Println(err.Error())
	}
	_, err = file.Write(req)
	if err != nil {
		log.Println(err.Error())
	}
	go func(filename string) {
		fu_channel <- filename
		fmt.Println("pushed")
	}(file.Name())
	return file.Close()
}

// func PushFuzz(filename string) {
// 	fu_channel <- filename
// 	fmt.Println("pushed")
// }

func Run() {
	fmt.Println("Started Fuzzing")
	var max = make(chan string, 5)
	for i := range fu_channel {
		max <- "Job"
		cmd := exec.Command("osmedeus", "scan", "-f", "froxy", "-t", i)
		go func(cmd exec.Cmd) {
			fmt.Println(cmd.String())

			err := cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
			}
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err.Error())
			}
			<-max
		}(*cmd)

	}
}
