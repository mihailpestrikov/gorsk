package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// Другие потенциальные импорты (например, для логирования, конфигурации)
)

// Server содержит экземпляр gin-движка и другие утилиты, связанные с сервером.
type Server struct {
	// Добавьте поля, если они существуют, например, логгер, конфигурация
	// Log *zlog.Logger
	// Cfg *config.Config
}

// New создает новый экземпляр сервера. (Предполагается, что существует конструктор)
// func New(...) *Server {
//     return &Server{...}
// }

// Health обрабатывает запросы проверки работоспособности, возвращая 200 OK с JSON {"status": "ok"}.
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Здесь могут быть другие методы, связанные с сервером, такие как Run, Shutdown и т.д.
