package file

import (
	"encoding/json"
	"os"

	"github.com/vllvll/devops/internal/types"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

type ProducerFile interface {
	WriteMetric(metrics *types.Metrics) error
	Close() error
}

// NewFileProducer Создание обработчика для записи в файл
func NewFileProducer(filename string) (ProducerFile, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteMetric Запись метрики в файл
func (p *Producer) WriteMetric(metrics *types.Metrics) error {
	return p.encoder.Encode(&metrics)
}

// Close Закрытие файла
func (p *Producer) Close() error {
	return p.file.Close()
}
