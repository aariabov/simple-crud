package roles

//go:generate mockgen -source controller.go -destination tests/mock_controller.go -package tests

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService RoleService
}

func NewRoleController(roleService RoleService) RoleController {
	return RoleController{roleService}
}

func (rc RoleController) GetAllRoles(c *gin.Context) {
	allRoles, err := rc.roleService.getAllRoles()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, allRoles)
}

func (rc RoleController) CreateRole(c *gin.Context) {
	var model RoleCreationModel
	if err := c.BindJSON(&model); err != nil {
		c.Error(fmt.Errorf("binding model error: %v", err))
		return
	}

	newId, err := rc.roleService.createRole(model)
	if err != nil {
		if validationErr, ok := err.(*ValidationErrors); ok {
			c.Error(validationErr)
			return
		}

		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, newId)
}

func (rc RoleController) UpdateRole(c *gin.Context) {
	var model RoleUpdatingModel
	if err := c.BindJSON(&model); err != nil {
		c.Error(fmt.Errorf("binding model error: %v", err))
	}

	err := rc.roleService.updateRole(model)
	if err != nil {
		if validationErr, ok := err.(*ValidationErrors); ok {
			c.Error(validationErr)
			return
		}

		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (rc RoleController) DeleteRole(c *gin.Context) {
	var model RoleDeletingModel
	if err := c.BindJSON(&model); err != nil {
		c.Error(fmt.Errorf("binding model error: %v", err))
	}

	err := rc.roleService.deleteRole(model)
	if err != nil {
		if validationErr, ok := err.(*ValidationErrors); ok {
			c.Error(validationErr)
			return
		}

		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
