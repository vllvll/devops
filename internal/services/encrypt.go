package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type MetricEncrypt struct {
	publicKey rsa.PublicKey
}

type Encrypt interface {
	Encrypt(data []byte) ([]byte, error)
}

// NewMetricEncrypt Создание сервиса для шифрования по публичному ключу
func NewMetricEncrypt(path string) (Encrypt, error) {
	if path == "" {
		return nil, nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	publicKeyPem, _ := pem.Decode(bytes)
	if publicKeyPem == nil {
		return nil, fmt.Errorf("публичный ключ неправильного формата")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyPem.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("не смогли распарсить ключ")
	}

	return &MetricEncrypt{
		publicKey: *rsaPublicKey,
	}, nil
}

// Encrypt Получение хеша для метрики типа Gauge
func (c MetricEncrypt) Encrypt(data []byte) ([]byte, error) {

	return rsa.EncryptPKCS1v15(rand.Reader, &c.publicKey, data)
}
