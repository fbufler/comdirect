package cache

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/fbufler/comdirect/pkg/comdirect"
)

type Cache struct {
	encryptionKey string
	storagePath   string
}

func NewCache(storagePath string, encryptionKey string) *Cache {
	return &Cache{
		encryptionKey: encryptionKey,
		storagePath:   storagePath,
	}
}

func (c *Cache) Load() (*comdirect.AuthToken, error) {
	slog.Debug("Loading token")
	data, err := os.ReadFile(c.storagePath)
	if err != nil {
		return nil, err
	}

	token, err := c.decrypt(string(data))
	if err != nil {
		slog.Warn("Failed to decrypt token")
		return nil, err
	}

	return token, nil
}

func (c *Cache) Save(token *comdirect.AuthToken) error {
	slog.Debug("Storing token")
	encryptedToken, err := c.encrypt(token)
	if err != nil {
		return err
	}

	return os.WriteFile(c.storagePath, []byte(encryptedToken), 0600)
}

func (c *Cache) encrypt(token *comdirect.AuthToken) (string, error) {
	slog.Debug("Encrypting token")
	serializedToken, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	encryptedToken := xorEncrypt(serializedToken, c.encryptionKey)
	return base64.StdEncoding.EncodeToString(encryptedToken), nil
}

func (c *Cache) decrypt(data string) (*comdirect.AuthToken, error) {
	slog.Debug("Decrypting token")
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	decryptedData := xorEncrypt(decodedData, c.encryptionKey)

	token := &comdirect.AuthToken{}
	err = json.Unmarshal(decryptedData, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func xorEncrypt(data []byte, key string) []byte {
	keyBytes := []byte(key)
	encrypted := make([]byte, len(data))
	for i := range data {
		encrypted[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}
	return encrypted
}
