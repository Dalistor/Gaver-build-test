package patterns

import (
	"test/internal/database"
)

type Repository interface {
	BeforeCreate(model *Model) error
	Create(model *Model) error
	AfterCreate(model *Model) error

	BeforeGetByID(id string) error
	GetByID(id string, model *Model) error
	AfterGetByID(model *Model) error

	BeforeGetAll() error
	GetAll() (error, []*Model)
	AfterGetAll(models []*Model) error

	BeforeUpdate(model *Model) error
	Update(model *Model) error
	AfterUpdate(model *Model) error

	BeforeDelete(id string) error
	Delete(id string) error
	AfterDelete(id string) error
}

type DefaultRepository[M Model] struct{}

func NewRepository[M Model]() *DefaultRepository[M] {
	return &DefaultRepository[M]{}
}

func (dr *DefaultRepository[M]) BeforeCreate(model *M) error {
	return nil
}

func (dr *DefaultRepository[M]) AfterCreate(model *M) error {
	return nil
}

func (dr *DefaultRepository[M]) Create(model *M) error {
	var m M = *model

	if err := dr.BeforeCreate(model); err != nil {
		return err
	}

	if err := m.Validate(); err != nil {
		return err
	}

	if err := database.DB.Create(&m).Error; err != nil {
		return err
	}

	if err := dr.AfterCreate(model); err != nil {
		return err
	}

	return nil
}

func (dr *DefaultRepository[M]) BeforeGetByID(id string) error {
	return nil
}

func (dr *DefaultRepository[M]) AfterGetByID(model *M) error {
	return nil
}

func (dr *DefaultRepository[M]) GetByID(id string, model *M) error {
	var m M

	if err := dr.BeforeGetByID(id); err != nil {
		return err
	}

	if err := database.DB.Where("id = ?", id).First(&m).Error; err != nil {
		return err
	}

	if err := dr.AfterGetByID(&m); err != nil {
		return err
	}

	*model = m
	return nil
}

func (dr *DefaultRepository[M]) BeforeGetAll() error {
	return nil
}

func (dr *DefaultRepository[M]) AfterGetAll(models []*M) error {
	return nil
}

func (dr *DefaultRepository[M]) GetAll() (error, []*Model) {
	var m []*M

	if err := dr.BeforeGetAll(); err != nil {
		return err, nil
	}

	if err := database.DB.Find(&m).Error; err != nil {
		return err, nil
	}

	var dtoList []*Model
	for _, _m := range m {
		var dto Model = *_m
		dtoList = append(dtoList, &dto)
	}

	if err := dr.AfterGetAll(m); err != nil {
		return err, nil
	}

	return nil, dtoList
}

func (dr *DefaultRepository[M]) BeforeUpdate(model *M) error {
	return nil
}

func (dr *DefaultRepository[M]) AfterUpdate(model *M) error {
	return nil
}

func (dr *DefaultRepository[M]) Update(model *M) error {
	var m M = *model

	if err := dr.BeforeUpdate(model); err != nil {
		return err
	}

	if err := m.Validate(); err != nil {
		return err
	}

	if err := dr.AfterUpdate(model); err != nil {
		return err
	}

	if err := database.DB.Save(&m).Error; err != nil {
		return err
	}

	return nil
}

func (dr *DefaultRepository[M]) BeforeDelete(id string) error {
	return nil
}

func (dr *DefaultRepository[M]) AfterDelete(id string) error {
	return nil
}

func (dr *DefaultRepository[M]) Delete(id string) error {
	var m M

	if err := dr.BeforeDelete(id); err != nil {
		return err
	}

	if err := database.DB.Where("id = ?", id).First(&m).Error; err != nil {
		return err
	}

	if err := database.DB.Delete(&m).Error; err != nil {
		return err
	}

	if err := dr.AfterDelete(id); err != nil {
		return err
	}

	return nil
}
