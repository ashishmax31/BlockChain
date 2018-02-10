package consensus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
)

func GetChain(address string) (chain []datatypes.Block) {
	body := hitServer(address)
	err := json.Unmarshal(body, &chain)
	if err != nil {
		log.Println("Unmarshaling JSON from another node failed!")
		panic(err.Error())
	}
	return chain
}

func hitServer(address string) []byte {
	address = fmt.Sprintf("http://%s/chain", address)
	log.Printf("Hitting server: %s \n", address)
	client := &http.Client{}
	resp, err := client.Get(address)
	if err != nil {
		log.Println("Hitting another server failed!")
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body

}
