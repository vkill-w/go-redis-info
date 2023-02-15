# redis-info  
[![Build Status](https://app.travis-ci.com/itsshashank/redis-info.svg?branch=main)](https://app.travis-ci.com/itsshashank/redis-info)

Package to parse redis.Info to access different sections 

## Example

```go

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	redisinfo "github.com/vkill-wu/redis-info"
)


func main() {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:7777",
		DB:   0,
	})
	info, err := client.Info().Result()
	if err != nil {
		fmt.Println(err)
	}
	jsonStr := redisinfo.Parse(info)

	// Can use pre defined structure from package
	var rinfo redisinfo.Info

	// Use your own custom structure
	// var rinfo Info

	if err := json.Unmarshal(jsonStr, &rinfo); err != nil {
		panic(err)
	}

	// Access the required content
	// Note: All the values are string by default so parse it as needed.

	fmt.Println(rinfo.Client.ConnectedClients)
	fmt.Println(rinfo.Server.OS)
}
```