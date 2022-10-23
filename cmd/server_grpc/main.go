package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/storage"
	"github.com/vllvll/devops/internal/storage/file"
	"github.com/vllvll/devops/internal/types"
	"github.com/vllvll/devops/pkg/postgres"

	// Импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/vllvll/devops/proto"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	pb.UnimplementedMetricsServer

	repository repositories.StatsRepository // Сервис для чтения и записи данных метрик
	signer     services.Signer              // Сервис для создания подписи
	db         *sql.DB                      // База данных
	decrypt    services.Decrypt             // Сервис для расшифрования данных
}

func (s *MetricsServer) BulkSaveMetrics(ctx context.Context, in *pb.AddBulkMetricsRequest) (*pb.AddBulkMetricsResponse, error) {
	var response pb.AddBulkMetricsResponse

	var metrics []types.Metrics
	var counters = types.Counters{}
	var gauges = types.Gauges{}

	inBulkMetrics := in.GetMetrics()
	inMetrics := inBulkMetrics.GetMetrics()

	for _, inMetric := range inMetrics {
		var metricType string
		switch inMetric.MType {
		case pb.Metric_COUNTER:
			metricType = dictionaries.CounterType
		case pb.Metric_GAUGE:
			metricType = dictionaries.GaugeType
		}

		metrics = append(metrics, types.Metrics{
			ID:    inMetric.Id,
			MType: metricType,
			Delta: inMetric.Delta,
			Value: inMetric.Value,
			Hash:  *inMetric.Hash,
		})
	}

	for _, metric := range metrics {
		switch metric.MType {
		case dictionaries.GaugeType:
			if !s.signer.IsEqualHashGauge(metric.ID, *metric.Value, metric.Hash) {
				return nil, status.Error(codes.InvalidArgument, "Hash not equal for gauge type")
			}

			gauges[metric.ID] = types.Gauge(*metric.Value)
		case dictionaries.CounterType:
			if !s.signer.IsEqualHashCounter(metric.ID, *metric.Delta, metric.Hash) {
				return nil, status.Error(codes.InvalidArgument, "Hash not equal for counter type")
			}

			counters[metric.ID] += types.Counter(*metric.Delta)
		}
	}

	err := s.repository.UpdateAll(gauges, counters)
	if err != nil {
		return nil, status.Error(codes.Internal, "Can't save metrics")
	}

	return &response, nil
}

func trustSubnetInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	config, err := conf.CreateServerConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	var ip string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("ip")
		if len(values) > 0 {
			ip = values[0]
		}
	}

	if config.TrustedSubnet != "" {
		if ip == "" {
			return nil, status.Error(codes.InvalidArgument, "Missing ip")
		}

		ipNet := net.ParseIP(ip)

		_, cidrNet, err := net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			return nil, status.Error(codes.Internal, "Can't parse CIDR")
		}

		if !cidrNet.Contains(ipNet) {
			return nil, status.Error(codes.PermissionDenied, "Can't parse CIDR")
		}
	}

	return handler(ctx, req)
}

func main() {
	config, err := conf.CreateServerConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	db, err := postgres.ConnectDatabase(config.DatabaseDsn)
	if err != nil {
		log.Fatalf("Error with database: %v", err)
	}
	defer db.Close()

	var storeTick = time.Tick(config.StoreInterval)

	statsRepository := repositories.NewStatsDatabaseRepository(db)
	if config.DatabaseDsn == "" {
		statsRepository = repositories.NewStatsMemoryRepository()
	}

	decrypt, err := services.NewMetricDecrypt(config.CryptoKey)
	if err != nil {
		log.Fatalf("Ошибка с инициализацией сервиса шифрования: %v", err)
	}

	signer := services.NewMetricSigner(config.Key)

	consumer, err := file.NewFileConsumer(config.StoreFile)
	if err != nil {
		log.Fatalf("Error with file consumer: %v", err)
	}

	producer, err := file.NewFileProducer(config.StoreFile)
	if err != nil {
		log.Fatalf("Error with file producer: %v", err)
	}
	defer producer.Close()

	fileStorage := storage.NewStatsStorage(config, consumer, producer)

	defer fileStorage.Save(statsRepository)

	statsRepository, err = fileStorage.Start(statsRepository)
	if err != nil {
		log.Fatalf("Error with file file storage: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(trustSubnetInterceptor))

	go func() {
		listen, err := net.Listen("tcp", config.Address)
		if err != nil {
			log.Fatal(err)
		}

		pb.RegisterMetricsServer(s, &MetricsServer{
			repository: statsRepository,
			signer:     signer,
			db:         db,
			decrypt:    decrypt,
		})

		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			s.GracefulStop()
			fileStorage.Save(statsRepository)

			return
		case <-storeTick:
			fileStorage.Save(statsRepository)
		}
	}
}
