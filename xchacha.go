package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
	"golang.org/x/crypto/chacha20poly1305"
)

type Xchacha struct {
	Key     []byte `json:"key"`
	Nonce   []byte `json:"nonce"`
	Padding []byte `json:"padding"`
	Aead    cipher.AEAD
}

var (
	xchacha     Xchacha
	paddingData = []byte("padding data")
)

func setupXchacha() error {
	var err error

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

			xchacha, err = Deserialize(keyData)
			if err != nil {
				return fmt.Errorf("failed to deserialize: %v", err)
			}

			return nil
		})

		if err != nil {
			return err
		}

	default:
		xchacha.Key = make([]byte, chacha20poly1305.KeySize)
		if _, err := rand.Read(xchacha.Key); err != nil {
			return fmt.Errorf("failed to generate key: %v", err)
		}

		xchacha.Nonce = make([]byte, chacha20poly1305.NonceSizeX)
		if _, err := rand.Read(xchacha.Nonce); err != nil {
			return fmt.Errorf("failed to generate nonce: %v", err)
		}

		xchachaBIN, err := Serialize(xchacha)
		if err != nil {
			return fmt.Errorf("failed to serialize: %v", err)
		}

		err = db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte("xchacha"))
			if bucket == nil {
				return nil
			}

			return bucket.Put([]byte("data"), xchachaBIN)
		})

		if err != nil {
			return err
		}
	}

	xchacha.Aead, err = chacha20poly1305.NewX(xchacha.Key)
	if err != nil {
		return fmt.Errorf("failed to generate aead: %v", err)
	}

	xchacha.Padding = paddingData

	return nil
}

func encrypt(data string) (string, error) {
	var err error

	plaintext := []byte(data)

	cipherBytes := xchacha.Aead.Seal(nil, xchacha.Nonce, plaintext, xchacha.Padding)
	cipherText := hex.EncodeToString(cipherBytes)

	return cipherText, err
}

func decrypt(cipherText []byte) (string, error) {
	var err error

	if len(cipherText) < xchacha.Aead.NonceSize() {
		return "", errors.New("too short data")
	}

	cipherBytes, err := hex.DecodeString(string(cipherText))
	if err != nil {
		return "", fmt.Errorf("failed string to bytes: %v", err)
	}
	decrypted, err := xchacha.Aead.Open(nil, xchacha.Nonce, cipherBytes, xchacha.Padding)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(decrypted), nil
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
