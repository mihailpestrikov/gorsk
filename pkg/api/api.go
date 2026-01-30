package api

import (
	"github.com/gin-gonic/gin"

	// Существующие импорты на основе структуры проекта
	"pkg/utl/middleware" // Заглушка, предполагается, что в проекте есть пакет middleware
	"pkg/utl/rbac"       // Заглушка, предполагается, что в проекте есть пакет rbac
	"pkg/utl/server"     // Новый импорт для пакета server
	// Примеры других импортов API
	// "pkg/api/auth"
	// "pkg/api/user"
)

// Handler содержит сервисы API.
type Handler struct {
	// Существующие поля сервисов (например, authService, userService)
	// Auth    auth.Service
	// User    user.Service
}

// NewHTTPHandler создает новый обработчик HTTP API.
func NewHTTPHandler(
	// Существующие параметры для сервисов
	// authSvc auth.Service,
	// userSvc user.Service,
) *Handler {
	return &Handler{
		// Инициализируйте существующие поля сервисов
		// Auth: authSvc,
		// User: userSvc,
	}
}

// RegisterRoutes регистрирует все маршруты API в gin-движке.
// Он принимает экземпляр server.Server для получения обработчика Health.
// Предполагается, что `*middleware.Middleware` и `*rbac.RBAC` также передаются, что типично для таких настроек.
func (h *Handler) RegisterRoutes(app *gin.Engine, mw *middleware.Middleware, rbac *rbac.RBAC, srv *server.Server) {
	// Добавление публичной конечной точки проверки работоспособности
	// Эта конечная точка не требует аутентификации, поэтому ее можно добавить непосредственно в приложение.
	app.GET("/health", srv.Health)

	// --- Существующая регистрация маршрутов (пример, фактическое содержимое будет отличаться) ---
	// v1 := app.Group("/api/v1")
	// v1.Use(mw.AuthJWT()) // Пример middleware
	// {
	//    user.RegisterRoutes(v1, h.User, mw, rbac)
	//    // Другие аутентифицированные маршруты
	// }

	// public := app.Group("/api")
	// {
	//    auth.RegisterRoutes(public, h.Auth, mw, rbac)
	//    // Другие публичные маршруты API, если таковые имеются
	// }
	// --- Конец примера существующей регистрации маршрутов ---
}
