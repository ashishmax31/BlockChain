package transactionvalidator

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"log"
	"sync"

	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
)

var mu sync.Mutex

func ValidateSignature(t datatypes.Transaction) (err error) {
	// Zero out the signature
	mu.Lock()
	defer mu.Unlock()
	signature, err := hex.DecodeString(t.Signature)
	t.Signature = ""
	message, _ := json.Marshal(t)
	hashed := sha256.Sum256(message)
	pubKey := getPublicKey(t.Sender)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		log.Println("Couldnt verify signature")
	}
	return
}

func getPublicKey(pKey string) *rsa.PublicKey {
	var s string
	for _, item := range pKey {
		s = s + string(item)
	}
	ln := splitSubN(s, 64)
	str := constructKey(ln)
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		log.Println("Damn!")
	}
	log.Printf("block.Type: %s\n", block.Type)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	key, _ := pub.(*rsa.PublicKey)
	return key
}

func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

func constructKey(s []string) (res string) {
	res = res + "-----BEGIN PUBLIC KEY-----\n"
	for _, substr := range s {
		res = res + substr + "\n"
	}
	return res + "-----END PUBLIC KEY-----"
}
