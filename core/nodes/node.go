package nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type Node struct {
	Identifier string
	Address    string
}

type nodes struct {
	sync.Mutex
	Nodes nodeSet
}
type nodeSet map[string]Node

var nodeList = nodes{
	Nodes: make(nodeSet),
}

//MasterNode ... Hard code the masternode address for now. In future this will be elected by the nodes.
var MasterNode = "127.0.0.1:3000"

// CurrentNodeAddress ... Contains the address of the current node.
var CurrentNodeAddress string

// SyncList ... A channel to send the address of the newly registered node.
var SyncList = make(chan string, 1000)

func RegNode(address string) bool {
	_, present := nodeList.Nodes[address]
	if present {
		return false
	}
	node := Node{
		Identifier: GenerateUUID(),
		Address:    address,
	}
	nodeList.Lock()
	defer nodeList.Unlock()
	nodeList.Nodes[address] = node
	if address != CurrentNodeAddress {
		SyncNodes(address)
	}
	return true
}

func SyncNodes(address string) {
	// Syncing node lists are done exclusively by the master node.
	if CurrentNodeAddress == MasterNode {
		// Put the newly registered node's address to a channel so that it is synced with all the nodes
		// presently in the node list.
		SyncList <- address

		// This function syncs the current node list to the newly registered node.
		syncNodeListToNewNode(address)
	}
}
func syncNodeListToNewNode(address string) {
	node_list := []string{}
	node_list = getNodeAddresses(node_list, address)
	for _, val := range node_list {
		NodeToSync := val
		if NodeToSync != address {
			go SyncNode(address, NodeToSync, true)
		}
	}

}

func getNodeAddresses(node_list []string, address string) []string {
	for nodeAddress := range nodeList.Nodes {
		if nodeAddress != address {
			node_list = append(node_list, nodeAddress)
		}
	}
	return node_list
}

func GenerateUUID() string {
	uuid := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", uuid)
}

// A go routine which syncs the newly added node's address with the nodes currently present in the node list.
func syncNodeList() {
	for {
		addr := <-SyncList
		if len(nodeList.Nodes) != 1 {
			for address := range nodeList.Nodes {
				syncing_addr := address
				if (syncing_addr != CurrentNodeAddress) && (syncing_addr != addr) {
					// log.Println("Syncing the newly registered node to all the nodes in the node_list(From channel)")
					go SyncNode(syncing_addr, addr, true)
				}
			}
		}
	}

}

func SyncNode(address, addr string, syncFlag bool) {
	if address == addr {
		return
	}
	log.Printf("[%s] Hitting %s with address %s \n", CurrentNodeAddress, address, addr)
	body := map[string]interface{}{"Address": addr, "sync": syncFlag}
	url := fmt.Sprintf("http://%s/register_node", address)
	status := performPostRequest(body, url)
	log.Printf("[%s] %d\n", CurrentNodeAddress, status)
	if status == http.StatusUnprocessableEntity {
		// Maybe the node already exists in the master node but not in the intended node,so send
		// the intended nodes address to the masternode and let the masternode sync
		// its nodes with the intended node.
		log.Printf("[%s]Hitting masternode %s with the current node address: %s \n", CurrentNodeAddress, MasterNode, CurrentNodeAddress)
		url = fmt.Sprintf("http://%s/register_node", MasterNode)
		body := map[string]interface{}{"Address": CurrentNodeAddress, "sync": syncFlag}
		performPostRequest(body, url)
	}
}

func NodeList() nodeSet {
	return nodeList.Nodes
}

func performPostRequest(body map[string]interface{}, url string) int {
	jsonVal, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonVal))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func init() {
	go syncNodeList()
}
