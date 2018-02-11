package requestvalidators

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
)

// Validate ...function validates the incoming json from the network.
// When new types need to be validated they are to be added here.
func Validate(obj interface{}) (err error) {
	objType := fmt.Sprintf("%T", obj)
	switch objType {
	case "datatypes.Transaction":
		err = validateTransactionParams(obj)
	default:
		err = fmt.Errorf("Unknown type cannot validate.. Panicking")
		panic(err)
	}
	return
}

func validateTransactionParams(obj interface{}) (err error) {
	txnObj, _ := obj.(datatypes.Transaction)
	if txnObj.Amount <= 0 {
		err = fmt.Errorf("transaction amount must be greater than 0")
	} else if (txnObj.Receiver == "") || (txnObj.Sender == "") {
		err = fmt.Errorf("sender or reciever cannot be blank")
	}
	return
}

func GetPrivateKey() (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile("/Users/ashishanand/.ssh/id_rsa")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	der, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return der, err
}

func SignPayload(txnObj *datatypes.Transaction) error {

	privateKey, err := GetPrivateKey()
	if err != nil {
		panic(err)
	}
	message, err := json.Marshal(*txnObj)
	hashed := sha256.Sum256(message)
	bodyHash, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	txnObj.Signature = hex.EncodeToString(bodyHash)
	return err
}
