package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Vault struct {
	masterPassword string
	data           map[string]string
	filePath       string
	mu             sync.RWMutex
}

func NewVault(filePath string) *Vault {
	return &Vault{
		data:     make(map[string]string),
		filePath: filePath,
	}
}

func (v *Vault) IsMasterPasswordSet() bool {
	return v.masterPassword != ""
}

func (v *Vault) SetMasterPassword(password string) {
	v.masterPassword = password
}

func (v *Vault) CheckMasterPassword(password string) bool {
	return v.masterPassword == password
}

func (v *Vault) AddEntry(service, password string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[service] = password
}

func (v *Vault) GetEntry(service string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	password, found := v.data[service]
	if !found {
		return "", errors.New("entry not found")
	}
	return password, nil
}

func (v *Vault) Save() {
	// Serialize data to JSON or other formats here if you want
	content := ""
	for k, val := range v.data {
		encryptedValue, _ := encrypt(val, v.masterPassword)
		content += k + ":" + encryptedValue + "\n"
	}

	ioutil.WriteFile(v.filePath, []byte(content), 0644)
}

func (v *Vault) Load() {
	// Deserialize data from JSON or other formats here if you want
	if _, err := os.Stat(v.filePath); err == nil {
		content, _ := ioutil.ReadFile(v.filePath)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				decryptedValue, _ := decrypt(parts[1], v.masterPassword)
				v.data[parts[0]] = decryptedValue
			}
		}
	}
}

// Encrypt string to base64 crypto using AES
func encrypt(text string, key string) (string, error) {
	// Key
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Message
	message := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(message))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], message)

	// Convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt from base64 to decrypted string
func decrypt(cryptoText string, key string) (string, error) {
	// Key
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Decode base64 string
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	// Decrypt
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
