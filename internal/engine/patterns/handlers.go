package patterns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	BeforeCreate(c *gin.Context)
	Create(c *gin.Context)
	AfterCreate(c *gin.Context)

	BeforeGetByID(c *gin.Context)
	GetByID(c *gin.Context)
	AfterGetByID(c *gin.Context)

	BeforeGetAll(c *gin.Context)
	GetAll(c *gin.Context)
	AfterGetAll(c *gin.Context)

	BeforeUpdate(c *gin.Context)
	Update(c *gin.Context)
	AfterUpdate(c *gin.Context)

	BeforeDelete(c *gin.Context)
	Delete(c *gin.Context)
	AfterDelete(c *gin.Context)
}

type DefaultHandler[M Model] struct {
	service Service
}

func NewHandler[M Model](service Service) *DefaultHandler[M] {
	return &DefaultHandler[M]{service: service}
}

func (dh *DefaultHandler[M]) BeforeCreate(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) Create(c *gin.Context) {
	var dto DTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var model Model = *dto.ToModel()
	if err := dh.service.Create(&model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ToDTO())
}

func (dh *DefaultHandler[M]) AfterCreate(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) BeforeGetByID(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) GetByID(c *gin.Context) {
	id := c.Param("id")

	var model Model
	if err := dh.service.GetByID(id, &model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ToDTO())
}

func (dh *DefaultHandler[M]) AfterGetByID(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) BeforeGetAll(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) GetAll(c *gin.Context) {
	var models []Model

	if err := dh.service.GetAll(&models); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dtoList []DTO
	for _, model := range models {
		dtoList = append(dtoList, *model.ToDTO())
	}

	c.JSON(http.StatusOK, dtoList)
}

func (dh *DefaultHandler[M]) AfterGetAll(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) BeforeUpdate(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) Update(c *gin.Context) {
	id := c.Param("id")
	var model Model
	if err := dh.service.GetByID(id, &model); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dto DTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ToDTO())
}

func (dh *DefaultHandler[M]) AfterUpdate(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) BeforeDelete(c *gin.Context) error {
	return nil
}

func (dh *DefaultHandler[M]) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := dh.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func (dh *DefaultHandler[M]) AfterDelete(c *gin.Context) error {
	return nil
}
