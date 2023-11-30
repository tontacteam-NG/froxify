package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/hibiken/asynq"
)

// var fu_channel = make(chan string, 2)
var list_tmux chan string
var client *asynq.Client

// var round *RoundRobin

func init() {
	pts := exec.Command("tmux", "list-panes", "-F", "'#{pane_tty}'")
	var list []string
	out, err := pts.StdoutPipe()
	scanner := bufio.NewScanner(out)
	if err != nil {
		color.Red("[ERROR] %s", err.Error())
		return
	}
	go func() {
		for scanner.Scan() {
			a := scanner.Text()
			list = append(list, a)
			fmt.Println(a)
		}
	}()
	pts.Run()
	redisConnection := asynq.RedisClientOpt{
		Addr: "localhost:6379", // Redis server address
	}
	client = asynq.NewClient(redisConnection)

	pts.Wait()
	time.Sleep(time.Second * 3)
	list_tmux = make(chan string, len(list))
	for _, v := range list {
		list_tmux <- v

	}
	go SetUp()

	// round = NewRoundRobin(list_tmux)
}

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
	// payload := OmsPayload{
	// 	RawRequest: file.Name(),
	// 	Std:        strings.ReplaceAll(round.Next(), "'", ""),
	// }
	// p, err := json.Marshal(payload)
	// if err != nil {
	// 	return err
	// }
	task := asynq.NewTask(RawRequest, []byte(file.Name()))

	_, err = client.Enqueue(task, asynq.MaxRetry(0))

	if err != nil {
		log.Println(err)
	}
	// go func(filename string) {
	// 	fu_channel <- filename
	// 	color.Yellow("[PUSH] %s", filename)
	// }(file.Name())
	return file.Close()
}

// func Run() {
// 	var max = make(chan string, 5)
// 	for i := range fu_channel {
// 		max <- "Job"
// 		cmd := exec.Command("osmedeus", "scan", "-f", "froxy", "-t", i)
// 		go func(cmd exec.Cmd) {
// 			color.Cyan("[DEBUG] %s", cmd.String())
// 			err := cmd.Start()
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			err = cmd.Wait()
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			<-max
// 		}(*cmd)
// 	}
// }

// func RunStreamTmux() {
// 	list_tmux := make([]string, 0)
// 	pts := exec.Command("tmux", "list-panes", "-F", "'#{pane_tty}'")

// 	out, err := pts.StdoutPipe()
// 	scanner := bufio.NewScanner(out)
// 	if err != nil {
// 		color.Red("[ERROR] %s", err.Error())
// 		return
// 	}
// 	go func() {
// 		for scanner.Scan() {
// 			list_tmux = append(list_tmux, scanner.Text())
// 		}
// 	}()
// 	err = pts.Run()
// 	if err != nil {
// 		color.Red("[ERROR] %s", err.Error())
// 		return
// 	}
// for _, pty := range list_tmux {

// 	// go RunStream(strings.ReplaceAll(pty, "'", ""))
// }
// }

// func RunStream(pty string) {
// 	fmt.Println(pty)
// 	for i := range fu_channel {
// 		cmd := exec.Command("osmedeusdev", "scan", "-f", "froxy", "-t", i, "-p=Std="+pty+"", "--debug")
// 		color.Cyan("[DEBUG] %s", cmd.String())
// 		err := cmd.Start()
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		err = cmd.Wait()
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}

// 	}
// }

// func AddRun(target string) {
// 	fu_channel <- target
// }
