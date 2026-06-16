package service

import (
	"regexp"

	"hexagonalarchitecture/internal/core/port"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

const errIDRequired = "%w: id is required"

type appService struct { // userService คือ struct ที่ implement interface port.UserService  (เป็น Implementation)
	repo      port.AppRepository
	publisher port.UserEventPublisher
	logger    port.Logger
	ids       port.IDGenerator
	clock     port.Clock
}

type AppServiceDeps struct { // UserServiceDeps คือ struct ที่ใช้สำหรับเก็บ dependencies ของ userService (Dependency Injection)
	Repo      port.AppRepository
	Publisher port.UserEventPublisher
	Logger    port.Logger
	IDs       port.IDGenerator
	Clock     port.Clock
}

func NewAppService(deps AppServiceDeps) port.AppService {
	return &appService{
		repo:      deps.Repo,
		publisher: deps.Publisher,
		logger:    deps.Logger,
		ids:       deps.IDs,
		clock:     deps.Clock,
	}
}
