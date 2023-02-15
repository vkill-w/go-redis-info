package redisinfo

import (
	"encoding/json"
	"fmt"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/go-redis/redis"
	"strings"
)

type Info struct {
	Server Server `json:"server"`
	Client Client `json:"clients"`
}

type Server struct {
	RedisVersion    string `json:"redis_version"`
	RedisGitSha1    string `json:"redis_git_sha1"`
	RedisGitDirty   string `json:"redis_git_dirty"`
	RedisBuildID    string `json:"redis_build_id"`
	RedisMode       string `json:"redis_mode"`
	OS              string `json:"os"`
	ArchBits        string `json:"arch_bits"`
	MultiplexingAPI string `json:"multiplexing_api"`
	GCCVersion      string `json:"gcc_version"`
	ProcessID       string `json:"process_id"`
	RunID           string `json:"run_id"`
	TCPPort         string `json:"tcp_port"`
	UptimeInSeconds string `json:"uptime_in_seconds"`
	UptimeInDays    string `json:"uptime_in_days"`
	HZ              string `json:"hz"`
	LRUClock        string `json:"lru_clock"`
	Executable      string `json:"executable"`
	ConfigFile      string `json:"config_file"`
	IOThreadsActive string `json:"io_threads_active"`
}

type Client struct {
	ConnectedClients            string `json:"connected_clients"`
	ClientRecentMaxOutputBuffer string `json:"client_recent_max_output_buffer"`
	ClientRecentMaxInputBuffer  string `json:"client_recent_max_input_buffer"`
	BlockedClients              string `json:"blocked_clients"`
	ClientsInTimeoutTable       string `json:"clients_in_timeout_table"`
	TrackingClients             string `json:"tracking_clients"`
}

type Memory struct {
	UsedMemory             string `json:"used_memory"`
	UsedMemoryHuman        string `json:"used_memory_human"`
	UsedMemoryRss          string `json:"used_memory_rss"`
	UsedMemoryRssHuman     string `json:"used_memory_rss_human"`
	UsedMemoryPeak         string `json:"used_memory_peak"`
	UsedMemoryPeakHuman    string `json:"used_memory_peak_human"`
	TotalSystemMemory      string `json:"total_system_memory"`
	TotalSystemMemoryHuman string `json:"total_system_memory_human"`
	UsedMemoryLua          string `json:"used_memory_lua"`
	UsedMemoryLuaHuman     string `json:"used_memory_lua_human"`
	Maxmemory              string `json:"maxmemory"`
	MaxmemoryHuman         string `json:"maxmemory_human"`
	MaxmemoryPolicy        string `json:"maxmemory_policy"`
	MemFragmentationRatio  string `json:"mem_fragmentation_ratio"`
	MemAllocator           string `json:"mem_allocator"`
}

type Persistence struct {
	Loading                  string `json:"loading"`
	RdbChangesSinceLastSave  string `json:"rdb_changes_since_last_save"`
	RdbBgsaveInProgress      string `json:"rdb_bgsave_in_progress"`
	RdbLastSaveTime          string `json:"rdb_last_save_time"`
	RdbLastBgsaveStatus      string `json:"rdb_last_bgsave_status"`
	RdbLastBgsaveTimeSec     string `json:"rdb_last_bgsave_time_sec"`
	RdbCurrentBgsaveTimeSec  string `json:"rdb_current_bgsave_time_sec"`
	AofEnabled               string `json:"aof_enabled"`
	AofRewriteInProgress     string `json:"aof_rewrite_in_progress"`
	AofRewriteScheduled      string `json:"aof_rewrite_scheduled"`
	AofLastRewriteTimeSec    string `json:"aof_last_rewrite_time_sec"`
	AofCurrentRewriteTimeSec string `json:"aof_current_rewrite_time_sec"`
	AofLastBgrewriteStatus   string `json:"aof_last_bgrewrite_status"`
	AofLastWriteStatus       string `json:"aof_last_write_status"`
}

type Stats struct {
	TotalConnectionsReceived string `json:"total_connections_received"`
	TotalCommandsProcessed   string `json:"total_commands_processed"`
	InstantaneousOpsPerSec   string `json:"instantaneous_ops_per_sec"`
	TotalNetInputBytes       string `json:"total_net_input_bytes"`
	TotalNetOutputBytes      string `json:"total_net_output_bytes"`
	InstantaneousInputKbps   string `json:"instantaneous_input_kbps"`
	InstantaneousOutputKbps  string `json:"instantaneous_output_kbps"`
	RejectedConnections      string `json:"rejected_connections"`
	SyncFull                 string `json:"sync_full"`
	SyncPartialOk            string `json:"sync_partial_ok"`
	SyncPartialErr           string `json:"sync_partial_err"`
	ExpiredKeys              string `json:"expired_keys"`
	EvictedKeys              string `json:"evicted_keys"`
	KeyspaceHits             string `json:"keyspace_hits"`
	KeyspaceMisses           string `json:"keyspace_misses"`
	PubsubChannels           string `json:"pubsub_channels"`
	PubsubPatterns           string `json:"pubsub_patterns"`
	LatestForkUsec           string `json:"latest_fork_usec"`
	MigrateCachedSockets     string `json:"migrate_cached_sockets"`
}

type Replication struct {
	Role                       string `json:"role"`
	ConnectedSlaves            string `json:"connected_slaves"`
	Slaves                     string `json:"slave"`
	MasterReplOffset           string `json:"master_repl_offset"`
	ReplBacklogActive          string `json:"repl_backlog_active"`
	ReplBacklogSize            string `json:"repl_backlog_size"`
	ReplBacklogFirstByteOffset string `json:"repl_backlog_first_byte_offset"`
	ReplBacklogHistLen         string `json:"repl_backlog_histlen"`
}

type CPU struct {
	UsedCPUSys          string `json:"used_cpu_sys"`
	UsedCPUUser         string `json:"used_cpu_user"`
	UsedCPUSysChildren  string `json:"used_cpu_sys_children"`
	UsedCPUUserChildren string `json:"used_cpu_user_children"`
}

type Modules struct{}

type Cluster struct {
	ClusterEnabled string `json:"cluster_enabled"`
}

type Keyspace struct {
	DB      string `json:"db,omitempty" index:"true"`
	Keys    string `json:"keys"`
	Expires string `json:"expires"`
	AvgTTL  string `json:"avg_ttl"`
}

func Parse(info string) []byte {

	contentPerCategory := make(map[string]map[string]string)
	content := strings.ReplaceAll(info, "\r", "")
	lines := strings.Split(content, "\n")
	category := ""
	for _, line := range lines {
		if line == "" {
			continue
		} else if strings.HasPrefix(line, "#") {
			category = strings.ToLower(line[2:])
			contentPerCategory[category] = make(map[string]string)
		} else {
			contents := strings.Split(line, ":")
			contentPerCategory[category][firstUpper(contents[0])] = contents[1]
		}
	}
	redisInfo, err := json.Marshal(contentPerCategory)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return redisInfo
}

func LParseInfo(redisInstanceName, redisEndpoint string) []byte {
	client := redisClient(redisEndpoint)
	result, err := client.Info().Result()
	if err != nil {
		fmt.Printf("LParseInfo error : %s", err.Error())
	}
	parseBytes := Parse(result)
	json := jsonvalue.MustUnmarshal(parseBytes)
	json.Set(redisInstanceName).At("RedisInstaceName")
	redisInfo, _ := json.Marshal()
	return redisInfo
}

func LParseClients(redisInstanceName, redisEndpoint string) [][]byte {
	client := redisClient(redisEndpoint)

	result, err := client.ClientList().Result()
	if err != nil {
		fmt.Printf("LParseInfo error : %s", err.Error())
	}
	parseBytes := ParseClients(result)
	redisInfo := addInstanceName(parseBytes, redisInstanceName)
	return redisInfo
}

func addInstanceName(parseBytes []byte, redisInstance string) [][]byte {
	var redisInfo [][]byte
	mustUnmarshal := jsonvalue.MustUnmarshal(parseBytes)
	for _, obj := range mustUnmarshal.ForRangeArr() {
		obj.Set(redisInstance).At("RedisInstaceName")
		marshal, _ := obj.Marshal()
		redisInfo = append(redisInfo, marshal)
	}
	return redisInfo
}

func ParseClients(info string) []byte {

	contentPerCategory := make([]map[string]string, 0)
	content := strings.ReplaceAll(info, "\r", "")
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		} else {
			mTemp := map[string]string{}
			contentsElem := strings.Split(line, " ")
			for _, elem := range contentsElem {
				contents := strings.Split(elem, "=")
				if strings.Contains(contents[1], ":") {
					mTemp[firstUpper(contents[0])] = strings.Split(contents[1], ":")[0]
				} else {
					mTemp[firstUpper(contents[0])] = contents[1]
				}
			}
			contentPerCategory = append(contentPerCategory, mTemp)
		}
	}
	clientsInfo, err := json.Marshal(contentPerCategory)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return clientsInfo
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func redisClient(endpoint string) *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr: endpoint,
		DB:   0,
	})

	return redisClient
}
