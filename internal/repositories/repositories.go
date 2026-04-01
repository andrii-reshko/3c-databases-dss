package repositories

import "dss/internal/domain/entities"

type AlternativeRepository interface {
	FindAll() ([]*entities.Alternative, error)
	FindByID(id int64) (*entities.Alternative, error)
	Create(alt *entities.Alternative) (int64, error)
	Update(alt *entities.Alternative) error
	Delete(id int64) error
}

type CriterionRepository interface {
	FindAll() ([]*entities.Criterion, error)
	FindByID(id int64) (*entities.Criterion, error)
	Create(criterion *entities.Criterion) (int64, error)
	Update(criterion *entities.Criterion) error
	Delete(id int64) error
}

type EvaluationRepository interface {
	FindAll() ([]*entities.Evaluation, error)
	UpsertBatch(evaluations []*entities.Evaluation) error
}
