package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

func RunStreamTmux() {
	list_tmux := make([]string, 0)
	pts := exec.Command("tmux", "list-panes", "-F", "'#{pane_tty}'")

	out, err := pts.StdoutPipe()
	scanner := bufio.NewScanner(out)
	if err != nil {
		color.Red("[ERROR] %s", err.Error())
		return
	}
	go func() {
		for scanner.Scan() {
			list_tmux = append(list_tmux, scanner.Text())
		}
	}()
	err = pts.Run()
	if err != nil {
		color.Red("[ERROR] %s", err.Error())
		return
	}
	for _, pty := range list_tmux {

		go RunStream(strings.ReplaceAll(pty, "'", ""))
	}
}

func RunStream(pty string) {
	fmt.Println(pty)
	for i := range fu_channel {
		cmd := exec.Command("osmedeusdev", "scan", "-f", "froxy", "-t", i, "-p=Std="+pty+"", "--debug")
		color.Cyan("[DEBUG] %s", cmd.String())
		err := cmd.Start()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err.Error())
		}

	}
}

func AddRun(target string) {
	fu_channel <- target
}
