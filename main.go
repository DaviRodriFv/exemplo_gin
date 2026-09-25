package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"api-gin/internal/handlers"
	"api-gin/internal/store"
)

func main() {

	r := gin.New()

	// Uso dos Middlewares globais nativos e personalizados
	r.Use(gin.Recovery())

	dataStore := store.New()
	h := handlers.New(dataStore)

	// 4. Mapeamento de Rotas sob Grupo Versionado
	v1 := r.Group("/api/v1")
	{
		// Monitoramento da API
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now(),
				"version":   "1.0.0",
			})
		})

		// Domínio de Salas (inventário físico)
		v1.POST("/salas", h.CriarSala)
		v1.GET("/salas", h.ListarSalas)
		v1.GET("/salas/:id/agenda", h.AgendaDaSala)

		// Domínio de Alunos
		v1.POST("/alunos", h.CriarAluno)
		v1.GET("/alunos", h.ListarAlunos)
		v1.GET("/alunos/:id", h.BuscarAluno)

		// Domínio de Turmas (Classes)
		v1.POST("/turmas", h.CriarTurma)
		v1.GET("/turmas", h.ListarTurmas)
		v1.GET("/turmas/:id", h.BuscarTurma)
		v1.POST("/turmas/:id/alunos", h.MatricularAluno)
		v1.GET("/turmas/:id/alunos", h.ListarAlunosDaTurma)
		v1.POST("/turmas/:id/alocar", h.AlocarSala)
	}

	r.Run(":8080")
}
