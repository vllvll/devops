// Package file Функционал для работы с файлами
package file

import (
	"encoding/json"
	"os"

	"github.com/vllvll/devops/internal/types"
)

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

type ConsumerFile interface {
	ReadMetric() (*types.Metrics, error)
	Close() error
}

// NewFileConsumer Создание обработчика для чтения из файла
func NewFileConsumer(filename string) (ConsumerFile, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadMetric Чтение метрики из файла
func (c *Consumer) ReadMetric() (*types.Metrics, error) {
	metric := &types.Metrics{}
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, err
	}

	return metric, nil
}

// Close Закрытие файла
func (c *Consumer) Close() error {
	return c.file.Close()
}
