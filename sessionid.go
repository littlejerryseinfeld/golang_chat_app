package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/google/uuid"
)

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func generateSignature(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Create hash of data
	h := sha256.New()
	h.Write(data)
	digest := h.Sum(nil)

	// Sign the hash with private key
	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		digest,
	)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func verifySignature(data []byte, signature []byte, privateKey *rsa.PrivateKey) error {
	return rsa.VerifyPKCS1v15(
		&privateKey.PublicKey,
		crypto.SHA256,
		data,
		signature,
	)
}

func decryptSignature(signature []byte, privateKey *rsa.PrivateKey) ([]byte, error) {

	return rsa.DecryptPKCS1v15(nil, privateKey, signature)
}

func getNewSessionId() string {

	id := uuid.New().String()
	return id
}

func getHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)
	return string(digest)
}
