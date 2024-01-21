package producer

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mocks.go -package=mock_producer "eventsproxy/internal/service/producer" NatsProducer
