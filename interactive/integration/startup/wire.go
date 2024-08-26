//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"we_book/interactive/repository"
	"we_book/interactive/repository/cache"
	"we_book/interactive/repository/dao"
	"we_book/interactive/service"
)

var thirdProvider = wire.NewSet(InitDB, InitRedis, InitLog)

var interactiveProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCacheInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func NewInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveProvider)
	return service.NewInteractiveService(nil, nil)
}
