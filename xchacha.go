package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
	"golang.org/x/crypto/chacha20poly1305"
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

var (
	passphrase     = string("my passsphrase")
	additionalData = []byte("padding data")
)

func SetupXchacha() (*Xchacha, error) {
	var err error
	var result Xchacha

	switch {
	case KeyExists("xchacha", []byte("data")):
		err = db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("xchacha"))
			if bucket == nil {
				return nil
			}

			keyData := bucket.Get([]byte("data"))
			if keyData == nil {
				return nil
			}

			result, err = deserialize(keyData)
			if err != nil {
				return fmt.Errorf("failed to deserialize: %v", err)
			}

			return nil
		})

		enteredKey := generateKeyFromString(passphrase)
		if !bytes.Equal(enteredKey, result.Key) {
			result = Xchacha{}
			return nil, errors.New("wrong passphrase")
		}

		if err != nil {
			return nil, err
		}

	default:
		// result.Key = make([]byte, chacha20poly1305.KeySize)
		// if _, err := rand.Read(result.Key); err != nil {
		// 	return nil, fmt.Errorf("failed to generate key: %v", err)
		// }
		result.Key = generateKeyFromString(passphrase)
		result.Nonce = make([]byte, chacha20poly1305.NonceSizeX)
		if _, err := rand.Read(result.Nonce); err != nil {
			return nil, fmt.Errorf("failed to generate nonce: %v", err)
		}

		xchachaBIN, err := serialize(result)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize: %v", err)
		}

		err = db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("xchacha"))
			if bucket == nil {
				return nil
			}

			return bucket.Put([]byte("data"), xchachaBIN)
		})

		if err != nil {
			return nil, err
		}
	}

	result.Aead, err = chacha20poly1305.NewX(result.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate aead: %v", err)
	}

	result.AdditionalData = additionalData

	return &result, nil
}

func generateKeyFromString(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
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

func serialize(data Xchacha) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	result := []byte(hex.EncodeToString(buffer.Bytes()))
	return result, nil
}

func deserialize(data []byte) (Xchacha, error) {
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
