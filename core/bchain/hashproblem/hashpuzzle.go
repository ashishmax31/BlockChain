package hashproblem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func ValidProof(prevProof int64, proof int64) bool {
	str := fmt.Sprintf("%v%v", proof, prevProof)
	hasher := sha256.New()
	hasher.Write([]byte(str))
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))
	// println(sha1Hash)
	if sha1Hash[:4] == "0000" {
		return true
	}
	return false
}
