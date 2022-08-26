package HandlerFunc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/olivere/elastic/v7"
)

const esHost = "http://p-gofood-butler-elasticsearch-client-c-09:9200"

const thresholdCpu = 20

var excludeNodeList = []string{
	"p-gofood-butler-elasticsearch-client-a-09",
	"p-gofood-butler-elasticsearch-client-b-09",
	"p-gofood-butler-elasticsearch-client-c-09",
	"p-gofood-butler-elasticsearch-master-a-02",
	"p-gofood-butler-elasticsearch-master-b-02",
	"p-gofood-butler-elasticsearch-master-c-02",
}

var excludeIndicesList = []string{
	".tasks",
	".ltrstore",
	".kibana_task_manager_1",
	".kibana_1",
	".apm-agent-configuration",
}

type NoReAllocationResponse struct {
	Val string `json:"status"`
}

type shards struct {
	Index      string
	Shard      int
	Store      float64
	QueryTotal int64
}

type reallocate struct {
	ProblematicNode  string `json:"problematic_node"`
	SwapShardNode    string `json:"swap_shard_node"`
	SwapShard        shards `json:"swap_shard"`
	ProblematicShard shards `json:"problematic_shard"`
}

type shardInfo struct {
	Index      string `json:"index"`
	Shard      string `json:"shard"`
	Node       string `json:"node"`
	Store      string `json:"store"`
	QueryTotal int64  `json:"query_total"`
}

type NodeInfo struct {
	Name     string
	CpuUsage int
}

//All sorting Function
func sortNodes(Nodes []NodeInfo) []NodeInfo {
	sort.Slice(Nodes, func(p, q int) bool {
		return Nodes[p].CpuUsage < Nodes[q].CpuUsage
	})
	return Nodes
}

func sortShardAscending(val []shards) []shards {
	sort.Slice(val, func(p, q int) bool {
		if val[p].QueryTotal == val[q].QueryTotal {
			return val[p].Store < val[q].Store
		}
		return val[p].QueryTotal < val[q].QueryTotal
	})
	return val
}

func sortShardDescending(val []shards) []shards {
	sort.Slice(val, func(p, q int) bool {
		if val[p].QueryTotal == val[q].QueryTotal {
			return val[p].Store > val[q].Store
		}
		return val[p].QueryTotal > val[q].QueryTotal
	})
	return val
}

// FindProblematicNode is used to Find the problematic Node
func FindProblematicNode(Nodes []NodeInfo) string {
	var problematicNode string
	var maxCpu = 0
	for _, item := range Nodes {
		if item.CpuUsage > thresholdCpu {
			if maxCpu == 0 {
				maxCpu = item.CpuUsage
				problematicNode = item.Name
			} else if item.CpuUsage > maxCpu {
				maxCpu = item.CpuUsage
				problematicNode = item.Name
			}
		}
	}
	return problematicNode
}

func possibleIndex(index string) bool {
	for _, indices := range excludeIndicesList {
		if indices == index {
			return false
		}
	}
	return true
}

// GetAllNodes is used to Find all Information about Nodes and shards allocated to nodes
func GetAllNodes(client *elastic.Client) ([]NodeInfo, map[string][]shards) {
	var Nodes []NodeInfo
	rest, err := client.NodesStats().Do(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	//Traversing all the nodes in the cluster and storing the name and CPU usage of Node
	for _, item := range rest.Nodes {
		var temp = NodeInfo{}
		temp.Name = item.Name
		temp.CpuUsage = item.OS.CPU.Percent
		Nodes = append(Nodes, temp)
	}

	information, _ := client.IndexStats().Do(context.Background())
	index := make(map[string]int64)
	for name, item := range information.Indices {
		index[name] = item.Total.Search.QueryTotal
	}

	//Getting all the information about shards present in the cluster
	response, err := http.Get(fmt.Sprintf("%s/_cat/shards?h=index,shard,state,store,node,ip&format=json", esHost))

	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var info []shardInfo
	err = json.Unmarshal(responseData, &info)
	if err != nil {
		log.Fatal(err)
	}

	//Storing all shards that are allocated to a Node
	memo := make(map[string][]shards)
	for _, item := range info {
		var temp shards
		if possibleIndex(item.Index) {
			temp.Index = item.Index
			number, err1 := strconv.Atoi(item.Shard)
			if err1 != nil {
				log.Fatal(err1)
			}
			temp.Shard = number
			s := item.Store
			length := len(s)
			var mem float64
			var substring2 string
			var substring1 string
			if length <= 1 {
				mem = 0
			} else if length == 2 {
				substring1 = s[0:1]
				tem, err := strconv.ParseFloat(substring1, 32)
				if err != nil {
					log.Fatal(err)
				}
				mem = tem
			} else {
				substring1 = s[0 : length-2]
				substring2 = s[length-2 : length]
				tem, err := strconv.ParseFloat(substring1, 64)
				if err != nil {
					log.Fatal(err)
				}
				if substring2 == "gb" {
					mem = tem * 1024 * 1024 * 1024
				} else if substring2 == "kb" {
					mem = tem * 1024
				} else if substring2 == "mb" {
					mem = tem * 1024 * 1024
				} else {
					substring := s[0 : length-1]
					val, err := strconv.ParseFloat(substring, 64)
					if err != nil {
						log.Fatal(err)
					}
					mem = val
				}
			}

			temp.Store = mem
			temp.QueryTotal = index[temp.Index]
			memo[item.Node] = append(memo[item.Node], temp)
		}
	}

	return Nodes, memo
}

func possibleSwap(problematicNode string, destination string, memo map[string][]shards) (bool, shards) {

	var swap shards
	// Traversing the shards of destination Node
	for _, item := range memo[destination] {
		var possible = true

		//Traversing shards of problematicNode and checking if current Shard is present in ProblematicNode or not
		for _, idx := range memo[problematicNode] {
			if idx.Index == item.Index && idx.Shard == item.Shard {
				possible = false
				break
			}
		}
		//If current Shard is possible to swap then we will return true and the Shard to be swapped
		if possible {
			return true, item
		}
	}
	// We didn't Found any Shard to be swapped from the current Node
	return false, swap
}

// Finding the Shard which can be swapped
func findShardToBeSwapped(problematicNode string, memo map[string][]shards, Nodes []NodeInfo) (bool, string, shards, shards) {

	var name string
	var temp shards
	var temp2 shards

	//Traversing all the shards of Problematic Node
	for _, item := range memo[problematicNode] {
		//Traversing all the Nodes
		for _, val := range Nodes {
			var NodePossible = true
			var currentNode = val.Name
			//Traversing all shards of current Node
			for _, idx := range memo[currentNode] {
				if idx.Index == item.Index && idx.Shard == item.Shard {
					NodePossible = false
					break
				}
			}

			//Current Shard of problematicNode is not present in current Node
			if NodePossible {
				//Checking if there is any Shard in the current Node that can be swapped
				condition, shardTobeSwapped := possibleSwap(problematicNode, currentNode, memo)
				if condition {
					return true, currentNode, item, shardTobeSwapped
				}
			}
		}
	}
	return false, name, temp, temp2
}

func reallocateShard(problematicNode string, destinationNode string, problematicShard shards, SwapShard shards, client *elastic.Client) {

	_, err := client.ClusterReroute().Add(elastic.NewMoveAllocationCommand(problematicShard.Index, problematicShard.Shard, problematicNode, destinationNode)).Add(elastic.NewMoveAllocationCommand(SwapShard.Index, SwapShard.Shard, destinationNode, problematicNode)).Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

}

func filterNodes(Nodes []NodeInfo) []NodeInfo {

	var newList []NodeInfo
	for _, item := range Nodes {
		possible := false
		for _, idx := range excludeNodeList {
			if idx == item.Name {
				possible = true
				break
			}
		}
		if !possible {
			newList = append(newList, item)
		}
	}
	return newList
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	//setting up a client
	w.Header().Set("Content-Type", "application/json")

	client, err := elastic.NewClient(elastic.SetURL(esHost), elastic.SetSniff(false))
	if err != nil {
		fmt.Println(err.Error())
	}

	//Getting All the information about all Nodes in the cluster
	nodes, allShards := GetAllNodes(client)

	//Filtering out the nodes on which No operations to be made
	nodes = filterNodes(nodes)

	//Finding the problematic Node
	problematicNode := FindProblematicNode(nodes)

	//Sorting the nodes on the basis of CPU Usage
	nodes = sortNodes(nodes)

	//Sorting the all shards allocated to a Node in ascending order considering the throughput of indices
	for idx, item := range allShards {
		fmt.Println(idx)
		temp := sortShardAscending(item)
		allShards[idx] = temp
	}

	// Sorting All shards of problematic Node in descending order considering the throughput of indices
	allShards[problematicNode] = sortShardDescending(allShards[problematicNode])

	// This function will return the destinationNode and the shards to be swapped , If Swapping is not possible then it will return false
	possible, destinationNode, problematicShard, swapShard := findShardToBeSwapped(problematicNode, allShards, nodes)
	if possible {
		var information reallocate
		information.ProblematicNode = problematicNode
		information.SwapShardNode = destinationNode
		information.ProblematicShard = problematicShard
		information.SwapShard = swapShard
		err := json.NewEncoder(w).Encode(information)
		if err != nil {
			log.Fatal(err.Error())
		}
		// reallocateShard(problematicNode, destinationNode, problematicShard, swapShard, client)
	} else {
		var information NoReAllocationResponse
		value := "No reallocation required"
		information.Val = value
		err := json.NewEncoder(w).Encode(information)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
