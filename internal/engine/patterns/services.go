package patterns

import "fmt"

type Service interface {
	BeforeCreate(model *Model) error
	Create(model *Model) error
	AfterCreate(model *Model) error

	BeforeGetByID(id string) error
	GetByID(id string, model *Model) error
	AfterGetByID(model *Model) error

	BeforeGetAll() error
	GetAll(models *[]Model) error
	AfterGetAll(models []*Model) error

	BeforeUpdate(model *Model) error
	Update(model *Model, dto *DTO) error
	AfterUpdate(model *Model) error

	BeforeDelete(id string) error
	Delete(id string) error
	AfterDelete(id string) error
}

type DefaultService[M Model] struct {
	repo Repository
}

func NewService[M Model](repo Repository) *DefaultService[M] {
	return &DefaultService[M]{repo: repo}
}

func (ds *DefaultService[M]) BeforeCreate(model *M) error {
	return nil
}

func (ds *DefaultService[M]) AfterCreate(model *M) error {
	return nil
}

func (ds *DefaultService[M]) Create(model *M) error {
	var m M = *model

	if err := m.Validate(); err != nil {
		return err
	}

	if err := ds.BeforeCreate(model); err != nil {
		return err
	}

	var parsedModel Model = m
	if err := ds.repo.Create(&parsedModel); err != nil {
		return err
	}

	if err := ds.AfterCreate(model); err != nil {
		return err
	}

	*model = parsedModel.(M)
	return nil
}

func (ds *DefaultService[M]) BeforeGetByID(id string) error {
	return nil
}

func (ds *DefaultService[M]) AfterGetByID(model *M) error {
	return nil
}

func (ds *DefaultService[M]) GetByID(id string, dto *DTO) error {
	if err := ds.BeforeGetByID(id); err != nil {
		return err
	}

	var model Model
	if err := ds.repo.GetByID(id, &model); err != nil {
		return err
	}

	m, ok := model.(M)
	if !ok {
		return fmt.Errorf("cannot assert type to M")
	}

	if err := ds.AfterGetByID(&m); err != nil {
		return err
	}

	*dto = *(m.ToDTO())
	return nil
}

func (ds *DefaultService[M]) BeforeGetAll() error {
	return nil
}

func (ds *DefaultService[M]) AfterGetAll(models []*M) error {
	return nil
}

func (ds *DefaultService[M]) GetAll(dtoList *[]DTO) error {
	var err error
	var models []*Model

	if err := ds.BeforeGetAll(); err != nil {
		return err
	}

	if err, models = ds.repo.GetAll(); err != nil {
		return err
	}

	var mList []*M
	for _, model := range models {
		m, ok := (*model).(M)
		if !ok {
			return fmt.Errorf("cannot assert type to M")
		}
		mList = append(mList, &m)
	}

	if err := ds.AfterGetAll(mList); err != nil {
		return err
	}

	for _, model := range models {
		m, ok := (*model).(M)
		if !ok {
			return fmt.Errorf("cannot assert type to M")
		}
		*dtoList = append(*dtoList, *m.ToDTO())
	}

	return nil
}

func (ds *DefaultService[M]) BeforeUpdate(model *M) error {
	return nil
}

func (ds *DefaultService[M]) AfterUpdate(model *M) error {
	return nil
}

func (ds *DefaultService[M]) Update(model *M, dto *DTO) error {
	if err := (*model).Validate(); err != nil {
		return err
	}

	if err := ds.BeforeUpdate(model); err != nil {
		return err
	}

	var parsedModel Model = *model
	if err := ds.repo.Update(&parsedModel); err != nil {
		return err
	}

	*dto = *parsedModel.ToDTO()
	return nil
}

func (ds *DefaultService[M]) BeforeDelete(id string) error {
	return nil
}

func (ds *DefaultService[M]) AfterDelete(id string) error {
	return nil
}

func (ds *DefaultService[M]) Delete(id string) error {
	if err := ds.BeforeDelete(id); err != nil {
		return err
	}

	if err := ds.repo.Delete(id); err != nil {
		return err
	}

	if err := ds.AfterDelete(id); err != nil {
		return err
	}

	return nil
}
