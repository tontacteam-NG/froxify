package elastic

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/nxczje/froxy/proxify/pkg/types"
	"github.com/pkg/errors"
)

// Options contains necessary options required for elasticsearch communicaiton
type Options struct {
	// Address for elasticsearch instance
	Addr string `yaml:"addr"`
	// SSL enables ssl for elasticsearch connection
	SSL bool `yaml:"ssl"`
	// SSLVerification disables SSL verification for elasticsearch
	SSLVerification bool `yaml:"ssl-verification"`
	// Username for the elasticsearch instance
	Username string `yaml:"username"`
	// Password is the password for elasticsearch instance
	Password string `yaml:"password"`
	// IndexName is the name of the elasticsearch index
	IndexName string `yaml:"index-name"`
	// RedisAddr is the address of redis instance
	Redis *redis.Client `yaml:"redis"`
	// TLS passes tls configuration to elasticsearch
	Filter []string `yaml:"tls"`
}

// Client type for elasticsearch
type Client struct {
	index    string
	options  *Options
	esClient *elasticsearch.Client
	Redis    *redis.Client
	Filter   []string
}

// New creates and returns a new client for elasticsearch
func New(option *Options) (*Client, error) {
	scheme := "http://"
	if option.SSL {
		scheme = "https://"
	}

	elasticsearchClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{scheme + option.Addr},
		Username:  option.Username,
		Password:  option.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: option.SSLVerification,
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating elasticsearch client")
	}
	client := &Client{
		esClient: elasticsearchClient,
		index:    option.IndexName,
		options:  option,
		Redis:    option.Redis,
		Filter:   option.Filter,
	}
	return client,
		nil

}

// Store saves a passed log event in elasticsearch
func (c *Client) Save(data types.OutputData) error {
	if data.Userdata.HasResponse {
		exists, err := c.Redis.SIsMember(context.Background(), "hash_id", data.Name).Result()
		if err != nil {
			fmt.Println("error: ", err)
			return err
		} else {
			if exists {
				doc := map[string]interface{}{
					"response":  data.DataString,
					"timestamp": time.Now().Format(time.RFC3339),
				}

				body, err := json.Marshal(&map[string]interface{}{
					"doc":           doc,
					"doc_as_upsert": true,
				})
				if err != nil {
					return err
				}
				updateRequest := esapi.UpdateRequest{
					Index:      c.index,
					DocumentID: data.Name,
					Body:       bytes.NewReader(body),
				}
				res, err := updateRequest.Do(context.Background(), c.esClient)
				if err != nil || res == nil {
					return errors.New("error thrown by elasticsearch: " + err.Error())
				}
				if res.StatusCode >= 300 {
					return errors.New("elasticsearch responded with an error: " + string(res.String()))
				}
				// Drain response to reuse connection
				_, er := io.Copy(io.Discard, res.Body)
				res.Body.Close()
				if er != nil {
					return er
				}
			}
		}
	} else {
		method := strings.Split(data.DataString, " ")[0]
		if method == "CONNECT" {
			return nil
		}
		for _, line := range c.Filter {
			matched, err := regexp.MatchString(line, data.Userdata.Host)
			if err != nil {
				fmt.Println("regexp.MatchString ERROR:", err)

			}
			if matched {
				return nil
			}
		}
		hash := CaculatorHash(data)
		if hash == nil {
			return nil
		}
		exists, err := c.Redis.SIsMember(context.Background(), "hash", hash).Result()
		if err != nil {
			fmt.Println("error: ", err)
			return err
		} else {
			if exists {
				return nil
			} else {
				doc := map[string]interface{}{
					"request":   data.DataString,
					"timestamp": time.Now().Format(time.RFC3339),
				}

				body, err := json.Marshal(&map[string]interface{}{
					"doc":           doc,
					"doc_as_upsert": true,
				})
				if err != nil {
					return err
				}
				updateRequest := esapi.UpdateRequest{
					Index:      c.index,
					DocumentID: data.Name,
					Body:       bytes.NewReader(body),
				}
				res, err := updateRequest.Do(context.Background(), c.esClient)
				if err != nil || res == nil {
					return errors.New("error thrown by elasticsearch: " + err.Error())
				}
				if res.StatusCode >= 300 {
					return errors.New("elasticsearch responded with an error: " + string(res.String()))
				}
				// Drain response to reuse connection
				_, er := io.Copy(io.Discard, res.Body)
				res.Body.Close()
				if er != nil {
					return er
				}
				err = c.Redis.SAdd(context.Background(), "hash_id", data.Name).Err()
				if err != nil {
					return err
				}
				err = c.Redis.SAdd(context.Background(), "hash", hash).Err()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CaculatorHash(data types.OutputData) []byte {
	check := strings.Contains(data.DataString, "Repeater: True")
	if check {
		return nil
	} else {
		hasher := md5.New()
		hasher.Write(data.Data)
		md5Hash := hasher.Sum(nil)
		return md5Hash
	}
}
