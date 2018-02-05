package requesthandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashishmax31/blockchain/core/bchain"
	"github.com/ashishmax31/blockchain/core/nodes"
)

func ShowChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		payload, err := json.Marshal(bchain.BlockChain.ReadBlockChain())
		if err != nil {
			log.Fatalln(err.Error())
		}
		writeJSONResponse(payload, w)

	}
}

func NewTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)
		obj, err := decodeBody(decoder)
		if err != nil {
			writeErrorResponse(err.Error(), w)
			return
		}
		index := bchain.NewTransaction(obj)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(fmt.Sprintf("The transaction will be forged in block number %d \n", index)))
	}

}

func Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		bchain.Mine()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Mined"))
	}

}

func RegisterNode(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No route matches [%s] %q ", req.Method, req.URL.Path)
		return
	}
	var data map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		writeErrorResponse(err.Error(), w)
		return
	}

	payload, sync := decodeNodeAddress(data)
	// Post requests created by the nodes package will have sync set to true, whereas created by users/others
	// will/should have sync set to false.
	// Register the recieved node address only if the current node is the master node or if sync is set to true
	// ie The request is generated by the nodes package as part of the node address syncing process.

	// Testing this is going to be a PAIN IN THE ARSE!
	if nodes.CurrentNodeAddress == nodes.MasterNode || sync == true {
		status := nodes.RegNode(payload)
		if status {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(req.RemoteAddr))

		} else {
			log.Println("Already present!")
			writeErrorResponse("Already present", w)
		}

	} else {
		// Pass the recieved node address to the master node for registering
		log.Println("rerouting to master node!")
		nodes.SyncNode(nodes.MasterNode, payload, true)
	}

}

func ShowNodes(w http.ResponseWriter, req *http.Request) {
	nodeList := []string{}
	for addr := range nodes.NodeList() {
		nodeList = append(nodeList, addr)
	}
	data, err := json.Marshal(nodeList)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	w.Write(data)

}

func writeJSONResponse(payload []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func decodeBody(d *json.Decoder) (t bchain.Transaction, e error) {
	e = d.Decode(&t)
	return
}

func writeErrorResponse(e string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(e))
}

func decodeNodeAddress(data map[string]interface{}) (string, bool) {
	payload, _ := data["Address"].(string)
	sync, present := data["sync"].(bool)
	if !present {
		sync = false
	}
	return payload, sync
}
