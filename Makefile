.PHONY: mock
mock:
	@mockgen  -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
	@mockgen  -source=internal/service/code.go -package=svcmocks -destination=internal/service/mocks/code.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=svcmocks -destination=internal/repository/cache/mocks/user.mock.go
	@mockgen -source=internal/repository/cache/code.go -package=svcmocks -destination=internal/repository/cache/mocks/code.mock.go
	@go mod tidy