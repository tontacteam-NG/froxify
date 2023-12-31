package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/nxczje/froxy/proxify"
	"github.com/nxczje/froxy/proxify/pkg/certs"
	"github.com/nxczje/froxy/proxify/pkg/logger/elastic"
	"github.com/nxczje/froxy/proxify/pkg/logger/kafka"
	"github.com/nxczje/froxy/proxify/pkg/types"

	"github.com/go-redis/redis/v8"
	"github.com/kr/pretty"
	"github.com/projectdiscovery/gologger"
)

func main() {
	certs.LoadCerts("./")

	// Connect redis client
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Println(err.Error())
		return
	}

	Filter := filterReq()
	addKafka := "127.0.0.1:9092"
	TopicKafka := "nothing"
	pr, _ := proxify.NewProxy(&proxify.Options{
		Verbosity: types.VerbosityDefault,
		Elastic: &elastic.Options{
			Addr:            "127.0.0.1:9200",
			IndexName:       "nothing",
			SSL:             false,
			SSLVerification: false,
			Username:        "",
			Password:        "",
			Redis:           redisClient,
			Filter:          Filter,
		},
		Kafka: &kafka.Options{
			Addr:   addKafka,
			Topic:  TopicKafka,
			Redis:  redisClient,
			Filter: Filter,
		},
		CertCacheSize:    254,
		ListenAddrHTTP:   "0.0.0.0:8888",
		ListenAddrSocks5: "0.0.0.0:10080",
	})

	pretty.Println("ListenAddrHTTP on port 8888")
	pretty.Println("ListenAddrSocks5 on port 10080")
	// Create redis client

	// run Osmedeus server
	// cmd := exec.Command("osmedeus", "server")
	// pretty.Println("Server Osmedeus running on port 8000")
	// errr := cmd.Run()
	// if errr != nil {
	// 	pretty.Println("Osmedeus running server error")
	// }

	err = pr.Run()
	if err != nil {
		gologger.Fatal().Msgf("Could not run proxify: %s\n", err)
	}

}

func filterReq() []string {
	f, err := os.OpenFile("../config/TLS-Pass-Through-List.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		pretty.Println("Open filter.txt error:", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	lines := []string{}
	for sc.Scan() {
		line := sc.Text()
		lines = append(lines, line)
	}
	return lines
}
