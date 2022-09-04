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

func (p *Producer) WriteMetric(metrics *types.Metrics) error {
	return p.encoder.Encode(&metrics)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
