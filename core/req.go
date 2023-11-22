package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
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
		color.Yellow("[PUSH] %s", filename)
	}(file.Name())
	return file.Close()
}

func Run() {
	var max = make(chan string, 5)
	for i := range fu_channel {
		max <- "Job"
		cmd := exec.Command("osmedeus", "scan", "-f", "froxy", "-t", i)
		go func(cmd exec.Cmd) {
			color.Cyan("[DEBUG] %s", cmd.String())
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
