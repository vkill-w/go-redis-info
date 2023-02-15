package main

import (
	"fmt"
	"github.com/vkill-w/go-redis-info"
)

func main() {

	parseInfo := redisinfo.LParseInfo("abc", "vm.vkill.cn:6379")
	fmt.Println(string(parseInfo))
	parseClients := redisinfo.LParseClients("abc", "vm.vkill.cn:6379")
	for _, bytes := range parseClients {
		fmt.Println(string(bytes))
	}

}
