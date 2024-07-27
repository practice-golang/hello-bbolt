package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
)

type Xchacha struct {
	Key            []byte `json:"key"`
	Nonce          []byte `json:"nonce"`
	AdditionalData []byte `json:"additional-data"`
	Aead           cipher.AEAD
}

type Encryptor interface {
	Encrypt(data string) (string, error)
	Decrypt(cipherText []byte) (string, error)
}

func generateKeyFromString(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}

func Serialize(data Xchacha) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	result := []byte(hex.EncodeToString(buffer.Bytes()))
	return result, nil
}

func Deserialize(data []byte) (Xchacha, error) {
	var result Xchacha

	decodedData, err := hex.DecodeString(string(data))
	if err != nil {
		return Xchacha{}, err
	}

	var buffer bytes.Buffer
	buffer.Write(decodedData)

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&result)
	if err != nil {
		return Xchacha{}, err
	}
	return result, nil
}

func (x *Xchacha) Encrypt(data string) (string, error) {
	var err error

	plaintext := []byte(data)

	cipherBytes := x.Aead.Seal(nil, x.Nonce, plaintext, x.AdditionalData)
	cipherText := hex.EncodeToString(cipherBytes)

	return cipherText, err
}

func (x *Xchacha) Decrypt(cipherText []byte) (string, error) {
	var err error

	if len(cipherText) < x.Aead.NonceSize() {
		return "", errors.New("too short data")
	}

	cipherBytes, err := hex.DecodeString(string(cipherText))
	if err != nil {
		return "", fmt.Errorf("failed string to bytes: %v", err)
	}
	decrypted, err := x.Aead.Open(nil, x.Nonce, cipherBytes, x.AdditionalData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(decrypted), nil
}
