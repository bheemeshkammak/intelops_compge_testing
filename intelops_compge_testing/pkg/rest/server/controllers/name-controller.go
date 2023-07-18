package controllers

import (
	"github.com/bheemeshkammak/intelops_compge_testing/intelops_compge_testing/pkg/rest/server/models"
	"github.com/bheemeshkammak/intelops_compge_testing/intelops_compge_testing/pkg/rest/server/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type NameController struct {
	nameService *services.NameService
}

func NewNameController() (*NameController, error) {
	nameService, err := services.NewNameService()
	if err != nil {
		return nil, err
	}
	return &NameController{
		nameService: nameService,
	}, nil
}

func (nameController *NameController) CreateName(context *gin.Context) {
	// validate input
	var input models.Name
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger name creation
	if _, err := nameController.nameService.CreateName(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Name created successfully"})
}

func (nameController *NameController) UpdateName(context *gin.Context) {
	// validate input
	var input models.Name
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger name update
	if _, err := nameController.nameService.UpdateName(id, &input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Name updated successfully"})
}

func (nameController *NameController) FetchName(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger name fetching
	name, err := nameController.nameService.GetName(id)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, name)
}

func (nameController *NameController) DeleteName(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger name deletion
	if err := nameController.nameService.DeleteName(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Name deleted successfully",
	})
}

func (nameController *NameController) ListNames(context *gin.Context) {
	// trigger all names fetching
	names, err := nameController.nameService.ListNames()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, names)
}

func (*NameController) PatchName(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*NameController) OptionsName(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*NameController) HeadName(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}
