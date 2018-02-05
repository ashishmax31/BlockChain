package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ashishmax31/blockchain/core/nodes"
	"github.com/ashishmax31/blockchain/requesthandler"
)

var NodeAddress = "127.0.0.1"
var port = flag.String("p", "3000", "Port number")

func main() {
	flag.Parse()
	NodeAddress = NodeAddress + fmt.Sprintf(":%s", *port)
	nodes.CurrentNodeAddress = NodeAddress
	mux := http.NewServeMux()
	mux.HandleFunc("/chain", requesthandler.ShowChain)
	mux.HandleFunc("/new_transaction", requesthandler.NewTransaction)
	mux.HandleFunc("/mine", requesthandler.Mine)
	mux.HandleFunc("/register_node", requesthandler.RegisterNode)
	mux.HandleFunc("/show_nodes", requesthandler.ShowNodes)
	log.Println("Starting server...on :", NodeAddress)
	nodes.RegNode(NodeAddress)
	err := http.ListenAndServe(NodeAddress, mux)
	if err != nil {
		log.Println(err.Error())
	}
}
