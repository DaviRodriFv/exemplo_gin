package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api-gin/internal/models"
	"api-gin/internal/store"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

// responderErro traduz um erro de negócio (*store.AppError) para a resposta
// HTTP correspondente; qualquer outro erro cai como 500.
func responderErro(c *gin.Context, err error) {
	if appErr, ok := err.(*store.AppError); ok {
		c.JSON(appErr.Status, gin.H{"erro": appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro interno inesperado"})
}

// ---------- Salas ----------

type criarSalaRequest struct {
	ID               string   `json:"id" binding:"required"`
	Nome             string   `json:"nome" binding:"required"`
	CapacidadeMaxima int      `json:"capacidade_maxima" binding:"required"`
	Recursos         []string `json:"recursos"`
}

func (h *Handler) CriarSala(c *gin.Context) {
	var req criarSalaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	sala, err := h.store.CriarSala(models.Sala{
		ID:               req.ID,
		Nome:             req.Nome,
		CapacidadeMaxima: req.CapacidadeMaxima,
		Recursos:         req.Recursos,
	})
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, sala)
}

func (h *Handler) ListarSalas(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListarSalas())
}

func (h *Handler) AgendaDaSala(c *gin.Context) {
	agenda, err := h.store.AgendaDaSala(c.Param("id"))
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, agenda)
}

// ---------- Alunos ----------

type criarAlunoRequest struct {
	ID    string `json:"id" binding:"required"`
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) CriarAluno(c *gin.Context) {
	var req criarAlunoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	aluno, err := h.store.CriarAluno(models.Aluno{
		ID:    req.ID,
		Nome:  req.Nome,
		Email: req.Email,
	})
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, aluno)
}

func (h *Handler) ListarAlunos(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListarAlunos())
}

func (h *Handler) BuscarAluno(c *gin.Context) {
	aluno, err := h.store.BuscarAluno(c.Param("id"))
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, aluno)
}

// ---------- Turmas ----------

type criarTurmaRequest struct {
	ID         string `json:"id" binding:"required"`
	Nome       string `json:"nome" binding:"required"`
	Disciplina string `json:"disciplina" binding:"required"`
	Professor  string `json:"professor" binding:"required"`
}

func (h *Handler) CriarTurma(c *gin.Context) {
	var req criarTurmaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	turma, err := h.store.CriarTurma(models.Turma{
		ID:         req.ID,
		Nome:       req.Nome,
		Disciplina: req.Disciplina,
		Professor:  req.Professor,
	})
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, turma)
}

func (h *Handler) ListarTurmas(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.ListarTurmas())
}

func (h *Handler) BuscarTurma(c *gin.Context) {
	turma, err := h.store.BuscarTurma(c.Param("id"))
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, turma)
}

type matricularAlunoRequest struct {
	AlunoID string `json:"aluno_id" binding:"required"`
}

func (h *Handler) MatricularAluno(c *gin.Context) {
	var req matricularAlunoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := h.store.MatricularAluno(c.Param("id"), req.AlunoID); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "aluno matriculado com sucesso"})
}

func (h *Handler) ListarAlunosDaTurma(c *gin.Context) {
	alunos, err := h.store.ListarAlunosDaTurma(c.Param("id"))
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, alunos)
}

type alocarSalaRequest struct {
	SalaID     string `json:"sala_id" binding:"required"`
	DiaSemana  string `json:"dia_semana" binding:"required"`
	HoraInicio string `json:"hora_inicio" binding:"required"`
	HoraFim    string `json:"hora_fim" binding:"required"`
}

func (h *Handler) AlocarSala(c *gin.Context) {
	var req alocarSalaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	turma, err := h.store.AlocarSala(c.Param("id"), req.SalaID, req.DiaSemana, req.HoraInicio, req.HoraFim)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, turma)
}
