package roles

import (
	"fmt"
	"unicode/utf8"
)

type ValidationErrors struct {
	ModelErrors map[string]string `json:"modelErrors"`
}

func (e *ValidationErrors) Error() string {
	return e.Error()
}

type RoleValidator struct{}

func NewRoleValidator() RoleValidator {
	return RoleValidator{}
}

func (v RoleValidator) ValidateRoleCreationModel(model RoleCreationModel, isNameExists bool) error {
	errors := make(map[string]string)
	if model.Name == "" {
		errors["name"] = "Название роли не может быть пустым"
		return &ValidationErrors{errors}
	}

	const nameMinLen, nameMaxLen int = 3, 100
	nameLen := utf8.RuneCountInString(model.Name)
	if nameLen < nameMinLen || nameLen > nameMaxLen {
		msg := fmt.Sprintf("Название роли должно быть от %d до %d символов", nameMinLen, nameMaxLen)
		errors["name"] = msg
		return &ValidationErrors{errors}
	}

	if isNameExists {
		errors["name"] = "Название роли уже существует"
		return &ValidationErrors{errors}
	}

	return nil
}

func (v RoleValidator) validateRoleUpdatingModel(model RoleUpdatingModel, isNameExists bool, role *RoleVm) error {
	errors := make(map[string]string)
	if model.Name == "" {
		errors["name"] = "Название роли не может быть пустым"
		return &ValidationErrors{errors}
	}

	const nameMinLen, nameMaxLen int = 3, 100
	nameLen := utf8.RuneCountInString(model.Name)
	if nameLen < nameMinLen || nameLen > nameMaxLen {
		msg := fmt.Sprintf("Название роли должно быть от %d до %d символов", nameMinLen, nameMaxLen)
		errors["name"] = msg
		return &ValidationErrors{errors}
	}

	if model.Id == "" {
		errors["name"] = "Идентификатор роли не может быть пустым"
		return &ValidationErrors{errors}
	}

	if isNameExists {
		errors["name"] = "Название роли уже существует"
		return &ValidationErrors{errors}
	}

	if role == nil {
		msg := fmt.Sprintf("Роль с идентификатором %s не найдена", model.Id)
		errors["name"] = msg
		return &ValidationErrors{errors}
	}

	if role.Name == "Admin" {
		errors["name"] = "Роль 'Admin' нельзя редактировать"
		return &ValidationErrors{errors}
	}

	return nil
}

func (v RoleValidator) validateRoleDeletingModel(model RoleDeletingModel, role *RoleVm, isAnyUserBelongToRole bool) error {
	errors := make(map[string]string)

	if model.Id == "" {
		errors["name"] = "Идентификатор роли не может быть пустым"
		return &ValidationErrors{errors}
	}

	if role == nil {
		msg := fmt.Sprintf("Роль с идентификатором %s не найдена", model.Id)
		errors["name"] = msg
		return &ValidationErrors{errors}
	}

	if isAnyUserBelongToRole {
		errors["id"] = fmt.Sprintf("Пользователи используют роль '%s'", role.Name)
		return &ValidationErrors{errors}
	}

	return nil
}
