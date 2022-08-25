package roles

import (
	"fmt"
)

type RoleVm struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RoleCreationModel struct {
	Name string
}

type RoleUpdatingModel struct {
	Id   string
	Name string
}

type RoleDeletingModel struct {
	Id string
}

type IRepo interface {
	getAllRoles() ([]RoleVm, error)
	insertNewRole(model RoleCreationModel) (string, error)
	deleteRole(id string) error
	updateRole(model RoleUpdatingModel) error
	isRoleExists(name string) (bool, error)
	isAnyUserBelongToRole(id string) (bool, error)
	getRoleById(id string) (*RoleVm, error)
}

type RoleService struct {
	roleRepo  IRepo
	validator RoleValidator
}

func NewRoleService(roleRepo IRepo, validator RoleValidator) RoleService {
	return RoleService{roleRepo, validator}
}

func (s RoleService) getAllRoles() ([]RoleVm, error) {
	return s.roleRepo.getAllRoles()
}

func (s RoleService) createRole(model RoleCreationModel) (string, error) {
	isNameExists, err := s.roleRepo.isRoleExists(model.Name)
	if err != nil {
		return "", err
	}

	err = s.validator.ValidateRoleCreationModel(model, isNameExists)
	if err != nil {
		return "", err
	}

	newId, err := s.roleRepo.insertNewRole(model)
	if err != nil {
		return "", fmt.Errorf("validation is Ok, but inserting error: %v", err)
	}

	return newId, nil
}

func (s RoleService) updateRole(model RoleUpdatingModel) error {
	isNameExists, err := s.roleRepo.isRoleExists(model.Name)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.getRoleById(model.Id)
	if err != nil {
		return err
	}

	err = s.validator.validateRoleUpdatingModel(model, isNameExists, role)
	if err != nil {
		return err
	}

	err = s.roleRepo.updateRole(model)
	if err != nil {
		return fmt.Errorf("validation is Ok, but updating error: %v", err)
	}

	return nil
}

func (s RoleService) deleteRole(model RoleDeletingModel) error {
	role, err := s.roleRepo.getRoleById(model.Id)
	if err != nil {
		return err
	}

	isAnyUserBelongToRole, err := s.roleRepo.isAnyUserBelongToRole(model.Id)
	if err != nil {
		return err
	}

	err = s.validator.validateRoleDeletingModel(model, role, isAnyUserBelongToRole)
	if err != nil {
		return err
	}

	err = s.roleRepo.deleteRole(model.Id)
	if err != nil {
		return fmt.Errorf("validation is Ok, but deleting error: %v", err)
	}

	return nil
}
