package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type MetricSigner struct {
	key []byte // Ключ для подписи
}

type Signer interface {
	Hash(content string, key []byte) string
	IsEqual(content string, key []byte, compareSum string) bool

	GetHashGauge(name string, value float64) string
	GetHashCounter(name string, value int64) string
	IsEqualHashGauge(name string, value float64, compareSum string) bool
	IsEqualHashCounter(name string, value int64, compareSum string) bool
}

// NewMetricSigner Создание сервиса для работы с подписью
func NewMetricSigner(key string) Signer {
	return &MetricSigner{
		key: []byte(key),
	}
}

// GetHashGauge Получение хеша для метрики типа Gauge
func (s MetricSigner) GetHashGauge(name string, value float64) string {
	return s.Hash(fmt.Sprintf("%s:gauge:%f", name, value), s.key)
}

// GetHashCounter Получение хеша для метрики типа Counter
func (s MetricSigner) GetHashCounter(name string, value int64) string {
	return s.Hash(fmt.Sprintf("%s:counter:%d", name, value), s.key)
}

// IsEqualHashGauge Проверка хеша для метрики типа Gauge
func (s MetricSigner) IsEqualHashGauge(name string, value float64, compareSum string) bool {
	if string(s.key) == "" {
		return true
	}

	return s.IsEqual(fmt.Sprintf("%s:gauge:%f", name, value), s.key, compareSum)
}

// IsEqualHashCounter Проверка хеша для метрики типа Counter
func (s MetricSigner) IsEqualHashCounter(name string, value int64, compareSum string) bool {
	if string(s.key) == "" {
		return true
	}

	return s.IsEqual(fmt.Sprintf("%s:counter:%d", name, value), s.key, compareSum)
}

// Hash Создание хеша
func (s MetricSigner) Hash(content string, key []byte) string {
	sign := hmac.New(sha256.New, key)
	sign.Write([]byte(content))
	sum := sign.Sum(nil)

	return hex.EncodeToString(sum)
}

// IsEqual Проверка хеша
func (s MetricSigner) IsEqual(content string, key []byte, compareSum string) bool {
	sign := hmac.New(sha256.New, key)
	sign.Write([]byte(content))
	sum := sign.Sum(nil)

	compareByte, _ := hex.DecodeString(compareSum)

	return hmac.Equal(sum, compareByte)
}
