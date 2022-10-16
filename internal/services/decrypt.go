package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type MetricDecrypt struct {
	privateKey rsa.PrivateKey
}

type Decrypt interface {
	Decrypt(data []byte) ([]byte, error)
}

// NewMetricDecrypt Создание сервиса для расшифрования по приватному ключу
func NewMetricDecrypt(path string) (Decrypt, error) {
	if path == "" {
		return nil, nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyPem, _ := pem.Decode(bytes)
	if privateKeyPem == nil {
		return nil, fmt.Errorf("праватный ключ неправильного формата")
	}

	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyPem.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := parsedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("не смогли распарсить ключ")
	}

	return &MetricDecrypt{
		privateKey: *rsaPrivateKey,
	}, nil
}

// Decrypt Получение хеша для метрики типа Gauge
func (d MetricDecrypt) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, &d.privateKey, data)
}
