package redisinfo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-redis/redis"
)

var client *redis.Client

func connect_redis() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:7777",
		DB:   0,
	})
}

func TestParse(t *testing.T) {
	connect_redis()

	tests := []struct {
		name string
		args string
		want Client
	}{
		{
			"test_Clients",
			"clients",
			Client{
				ConnectedClients:            "1",
				ClientRecentMaxInputBuffer:  "8",
				ClientRecentMaxOutputBuffer: "0",
				BlockedClients:              "0",
				TrackingClients:             "0",
				ClientsInTimeoutTable:       "0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := client.Info(tt.args).Result()
			if err != nil {
				fmt.Println(err)
			}
			var rinfo Info
			got := Parse(info)
			if err := json.Unmarshal(got, &rinfo); err != nil {
				panic(err)
			}
			if !reflect.DeepEqual(rinfo.Client, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
