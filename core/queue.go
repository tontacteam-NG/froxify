package core

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/hibiken/asynq"
)

const (
	RawRequest = "request:raw"
)

type OmsPayload struct {
	RawRequest string
	Std        string
}

// Asynq worker server

func SetUp() {
	redisConnection := asynq.RedisClientOpt{
		Addr: "localhost:6379", // Redis server address
	}
	worker := asynq.NewServer(redisConnection, asynq.Config{
		// Specify how many concurrent workers to use.
		Concurrency: len(list_tmux),
		// Specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(RawRequest, HandleRunStream)
	if err := worker.Run(mux); err != nil {
		log.Println(err)
	}
}

func HandleRunStream(c context.Context, t *asynq.Task) error {
	// var p OmsPayload
	// if err := json.Unmarshal(t.Payload(), &p); err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }
	std := <-list_tmux

	defer func() {
		color.Green("[Std] %s", std)
		list_tmux <- std

	}()
	cmd := exec.Command("osmedeusdev", "scan", "-f", "froxy", "-t", string(t.Payload()), "-p=Std="+strings.ReplaceAll(std, "'", "")+"", "--debug")
	color.Cyan("[DEBUG] %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }

	return nil
}

// func Difference(a, b []string) (diff []string) {
// 	m := make(map[string]bool)

// 	for _, item := range b {
// 		m[item] = true
// 	}

// 	for _, item := range a {
// 		if _, ok := m[item]; !ok {
// 			diff = append(diff, item)
// 		}
// 	}
// 	return
// }
