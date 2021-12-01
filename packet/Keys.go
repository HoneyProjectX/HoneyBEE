package packet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/google/uuid"
)

var (
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	publicKeySlice  []byte
	privateKeySlice []byte
	VerifyToken     = make([]byte, 4)
)

func GenerateKeys() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		Log.Error(err.Error())
	}
	privateKey.Precompute()
	publicKey = &privateKey.PublicKey
	publicKeySlice, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	rand.Read(VerifyToken)
	Log.Info("Key Generated!")
	// Get rid of vscode warnings
	_ = publicKey
	_ = privateKeySlice
}

func Auth(username string, sharedSecret []byte) (uuid.UUID, string) {
	PlayerUUID, response, autherr := Authenticate(username, "", sharedSecret, publicKeySlice)
	if autherr != nil {
		Log.Error("Auth Fail! ", autherr)
		return uuid.Nil, ""
	}
	return PlayerUUID, response.Name
}
