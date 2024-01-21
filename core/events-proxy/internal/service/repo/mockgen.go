package repo

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/mocks.go -package=mock_repo "eventsproxy/internal/service/repo" UserRepo
