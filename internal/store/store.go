package store

import (
	"strings"
	"sync"
	"time"

	"api-gin/internal/models"
)

// Store mantém o estado da aplicação em memória, protegido por um mutex
// para permitir acesso concorrente seguro entre as goroutines de requisição do Gin.
type Store struct {
	mu     sync.RWMutex
	salas  map[string]*models.Sala
	alunos map[string]*models.Aluno
	turmas map[string]*models.Turma
}

func New() *Store {
	return &Store{
		salas:  make(map[string]*models.Sala),
		alunos: make(map[string]*models.Aluno),
		turmas: make(map[string]*models.Turma),
	}
}

// ---------- Salas ----------

func (s *Store) CriarSala(sala models.Sala) (*models.Sala, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sala.ID == "" {
		return nil, ErrBadRequest("identificador da sala é obrigatório")
	}
	if sala.CapacidadeMaxima <= 0 {
		return nil, ErrBadRequest("capacidade máxima deve ser maior que zero")
	}
	if _, existe := s.salas[sala.ID]; existe {
		return nil, ErrConflict("já existe uma sala cadastrada com este identificador")
	}

	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}
	nova := sala
	s.salas[nova.ID] = &nova
	return &nova, nil
}

func (s *Store) ListarSalas() []models.Sala {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lista := make([]models.Sala, 0, len(s.salas))
	for _, sala := range s.salas {
		lista = append(lista, *sala)
	}
	return lista
}

func (s *Store) BuscarSala(id string) (*models.Sala, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sala, ok := s.salas[id]
	if !ok {
		return nil, ErrNotFound("sala não encontrada")
	}
	copia := *sala
	return &copia, nil
}

// AgendaDaSala retorna a grade de uso (turmas alocadas) de uma sala.
func (s *Store) AgendaDaSala(salaID string) ([]models.AgendaItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.salas[salaID]; !ok {
		return nil, ErrNotFound("sala não encontrada")
	}

	agenda := make([]models.AgendaItem, 0)
	for _, turma := range s.turmas {
		if turma.Alocacao != nil && turma.Alocacao.SalaID == salaID {
			agenda = append(agenda, models.AgendaItem{
				TurmaID:    turma.ID,
				TurmaNome:  turma.Nome,
				DiaSemana:  turma.Alocacao.DiaSemana,
				HoraInicio: turma.Alocacao.HoraInicio,
				HoraFim:    turma.Alocacao.HoraFim,
			})
		}
	}
	return agenda, nil
}

// ---------- Alunos ----------

func (s *Store) CriarAluno(aluno models.Aluno) (*models.Aluno, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if aluno.ID == "" {
		return nil, ErrBadRequest("identificador do aluno (matrícula) é obrigatório")
	}
	if aluno.Nome == "" {
		return nil, ErrBadRequest("nome do aluno é obrigatório")
	}
	if _, existe := s.alunos[aluno.ID]; existe {
		return nil, ErrConflict("já existe um aluno cadastrado com esta matrícula")
	}

	novo := aluno
	s.alunos[novo.ID] = &novo
	return &novo, nil
}

func (s *Store) ListarAlunos() []models.Aluno {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lista := make([]models.Aluno, 0, len(s.alunos))
	for _, aluno := range s.alunos {
		lista = append(lista, *aluno)
	}
	return lista
}

func (s *Store) BuscarAluno(id string) (*models.Aluno, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	aluno, ok := s.alunos[id]
	if !ok {
		return nil, ErrNotFound("aluno não encontrado")
	}
	copia := *aluno
	return &copia, nil
}

// ---------- Turmas ----------

func (s *Store) CriarTurma(turma models.Turma) (*models.TurmaResumo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if turma.ID == "" {
		return nil, ErrBadRequest("identificador da turma é obrigatório")
	}
	if turma.Nome == "" || turma.Disciplina == "" || turma.Professor == "" {
		return nil, ErrBadRequest("nome, disciplina e professor são obrigatórios")
	}
	if _, existe := s.turmas[turma.ID]; existe {
		return nil, ErrConflict("já existe uma turma cadastrada com este identificador")
	}

	nova := turma
	nova.AlunosIDs = []string{}
	nova.Alocacao = nil
	s.turmas[nova.ID] = &nova
	resumo := nova.Resumo()
	return &resumo, nil
}

func (s *Store) ListarTurmas() []models.TurmaResumo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lista := make([]models.TurmaResumo, 0, len(s.turmas))
	for _, turma := range s.turmas {
		lista = append(lista, turma.Resumo())
	}
	return lista
}

func (s *Store) BuscarTurma(id string) (*models.TurmaResumo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	turma, ok := s.turmas[id]
	if !ok {
		return nil, ErrNotFound("turma não encontrada")
	}
	resumo := turma.Resumo()
	return &resumo, nil
}

// MatricularAluno adiciona um aluno a uma turma, aplicando as regras de
// duplicidade, capacidade da sala alocada e conflito de agenda do aluno.
func (s *Store) MatricularAluno(turmaID, alunoID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	turma, ok := s.turmas[turmaID]
	if !ok {
		return ErrNotFound("turma não encontrada")
	}
	aluno, ok := s.alunos[alunoID]
	if !ok {
		return ErrNotFound("aluno não encontrado")
	}

	for _, id := range turma.AlunosIDs {
		if id == aluno.ID {
			return ErrConflict("aluno já está matriculado nesta turma")
		}
	}

	if turma.Alocacao != nil {
		sala := s.salas[turma.Alocacao.SalaID]
		if sala != nil && len(turma.AlunosIDs)+1 > sala.CapacidadeMaxima {
			return ErrUnprocessable("capacidade máxima da sala alocada foi atingida")
		}
	}

	if turma.Alocacao != nil {
		for _, outraID := range s.turmasDoAluno(aluno.ID) {
			outra := s.turmas[outraID]
			if outra == nil || outra.Alocacao == nil {
				continue
			}
			if horariosConflitam(turma.Alocacao, outra.Alocacao) {
				return ErrConflict("aluno possui conflito de agenda com outra turma no mesmo horário")
			}
		}
	}

	turma.AlunosIDs = append(turma.AlunosIDs, aluno.ID)
	return nil
}

// ListarAlunosDaTurma retorna os alunos matriculados em uma turma.
func (s *Store) ListarAlunosDaTurma(turmaID string) ([]models.Aluno, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	turma, ok := s.turmas[turmaID]
	if !ok {
		return nil, ErrNotFound("turma não encontrada")
	}

	alunos := make([]models.Aluno, 0, len(turma.AlunosIDs))
	for _, id := range turma.AlunosIDs {
		if aluno, ok := s.alunos[id]; ok {
			alunos = append(alunos, *aluno)
		}
	}
	return alunos, nil
}

// AlocarSala associa uma sala a uma turma em um dia/horário específico,
// validando existência, capacidade e conflitos de agenda na sala.
func (s *Store) AlocarSala(turmaID, salaID, diaSemana, horaInicio, horaFim string) (*models.TurmaResumo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	turma, ok := s.turmas[turmaID]
	if !ok {
		return nil, ErrNotFound("turma não encontrada")
	}
	sala, ok := s.salas[salaID]
	if !ok {
		return nil, ErrNotFound("sala não encontrada")
	}

	inicioMin, fimMin, err := validarIntervalo(horaInicio, horaFim)
	if err != nil {
		return nil, err
	}
	dia := normalizarDia(diaSemana)
	if dia == "" {
		return nil, ErrBadRequest("dia da semana inválido")
	}

	if sala.CapacidadeMaxima < len(turma.AlunosIDs) {
		return nil, ErrUnprocessable("capacidade da sala é menor que a quantidade de alunos já matriculados na turma")
	}

	nova := &models.Alocacao{
		SalaID:     salaID,
		DiaSemana:  dia,
		HoraInicio: horaInicio,
		HoraFim:    horaFim,
	}

	for _, outra := range s.turmas {
		if outra.ID == turma.ID || outra.Alocacao == nil {
			continue
		}
		if outra.Alocacao.SalaID != salaID {
			continue
		}
		if outra.Alocacao.DiaSemana != dia {
			continue
		}
		outroInicio, outroFim, _ := validarIntervalo(outra.Alocacao.HoraInicio, outra.Alocacao.HoraFim)
		if inicioMin < outroFim && fimMin > outroInicio {
			return nil, ErrConflict("sala já está alocada para outra turma em horário sobreposto neste dia")
		}
	}

	turma.Alocacao = nova
	resumo := turma.Resumo()
	return &resumo, nil
}

// ---------- Auxiliares internos ----------

func (s *Store) turmasDoAluno(alunoID string) []string {
	ids := make([]string, 0)
	for _, turma := range s.turmas {
		for _, id := range turma.AlunosIDs {
			if id == alunoID {
				ids = append(ids, turma.ID)
				break
			}
		}
	}
	return ids
}

func horariosConflitam(a, b *models.Alocacao) bool {
	if a == nil || b == nil {
		return false
	}
	if !strings.EqualFold(a.DiaSemana, b.DiaSemana) {
		return false
	}
	aInicio, aFim, errA := validarIntervalo(a.HoraInicio, a.HoraFim)
	bInicio, bFim, errB := validarIntervalo(b.HoraInicio, b.HoraFim)
	if errA != nil || errB != nil {
		return false
	}
	return aInicio < bFim && bInicio < aFim
}

// validarIntervalo faz o parse de horários no formato HH:MM e retorna os
// minutos desde a meia-noite, garantindo que o início seja anterior ao fim.
func validarIntervalo(inicio, fim string) (int, int, error) {
	inicioMin, err := parseHoraParaMinutos(inicio)
	if err != nil {
		return 0, 0, ErrBadRequest("horário de início inválido, use o formato HH:MM")
	}
	fimMin, err := parseHoraParaMinutos(fim)
	if err != nil {
		return 0, 0, ErrBadRequest("horário de término inválido, use o formato HH:MM")
	}
	if inicioMin >= fimMin {
		return 0, 0, ErrBadRequest("horário de início deve ser anterior ao horário de término")
	}
	return inicioMin, fimMin, nil
}

func parseHoraParaMinutos(hora string) (int, error) {
	t, err := time.Parse("15:04", hora)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}

var diasValidos = map[string]string{
	"SEGUNDA": "SEGUNDA", "SEGUNDA-FEIRA": "SEGUNDA",
	"TERCA": "TERCA", "TERÇA": "TERCA", "TERCA-FEIRA": "TERCA", "TERÇA-FEIRA": "TERCA",
	"QUARTA": "QUARTA", "QUARTA-FEIRA": "QUARTA",
	"QUINTA": "QUINTA", "QUINTA-FEIRA": "QUINTA",
	"SEXTA": "SEXTA", "SEXTA-FEIRA": "SEXTA",
	"SABADO": "SABADO", "SÁBADO": "SABADO",
	"DOMINGO": "DOMINGO",
}

// normalizarDia padroniza o dia da semana informado, aceitando variações
// comuns de escrita (com/sem acento, com/sem sufixo "-feira").
func normalizarDia(dia string) string {
	chave := strings.ToUpper(strings.TrimSpace(dia))
	if normalizado, ok := diasValidos[chave]; ok {
		return normalizado
	}
	return ""
}
