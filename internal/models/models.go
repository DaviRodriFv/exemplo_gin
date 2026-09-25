package models

// Sala representa um espaço físico (sala de aula ou laboratório) do campus.
type Sala struct {
	ID               string   `json:"id"`
	Nome             string   `json:"nome"`
	CapacidadeMaxima int      `json:"capacidade_maxima"`
	Recursos         []string `json:"recursos"`
}

// Aluno representa um estudante matriculado na instituição.
type Aluno struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Alocacao representa a alocação física de uma turma em uma sala,
// definindo o dia da semana e o intervalo de horário da aula.
type Alocacao struct {
	SalaID     string `json:"sala_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

// Turma representa uma turma acadêmica, com seus alunos matriculados
// e, opcionalmente, uma alocação física.
type Turma struct {
	ID         string    `json:"id"`
	Nome       string    `json:"nome"`
	Disciplina string    `json:"disciplina"`
	Professor  string    `json:"professor"`
	AlunosIDs  []string  `json:"-"`
	Alocacao   *Alocacao `json:"alocacao,omitempty"`
}

// TurmaResumo é a representação de uma turma exposta pela API de listagem,
// incluindo a quantidade de alunos matriculados e o status de alocação.
type TurmaResumo struct {
	ID                string    `json:"id"`
	Nome              string    `json:"nome"`
	Disciplina        string    `json:"disciplina"`
	Professor         string    `json:"professor"`
	QuantidadeAlunos  int       `json:"quantidade_alunos"`
	Alocada           bool      `json:"alocada"`
	Alocacao          *Alocacao `json:"alocacao,omitempty"`
}

func (t *Turma) Resumo() TurmaResumo {
	return TurmaResumo{
		ID:               t.ID,
		Nome:             t.Nome,
		Disciplina:       t.Disciplina,
		Professor:        t.Professor,
		QuantidadeAlunos: len(t.AlunosIDs),
		Alocada:          t.Alocacao != nil,
		Alocacao:         t.Alocacao,
	}
}

// AgendaItem representa um item da grade de uso de uma sala.
type AgendaItem struct {
	TurmaID    string `json:"turma_id"`
	TurmaNome  string `json:"turma_nome"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}
