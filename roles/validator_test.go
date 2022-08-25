package roles

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestValidateRoleCreationModel_Invalid(t *testing.T) {
	const isNameExists = true
	const nameMinLen, nameMaxLen int = 3, 100
	lenRangeError := fmt.Sprintf("Название роли должно быть от %d до %d символов", nameMinLen, nameMaxLen)
	tooShortName := strings.Repeat("c", nameMinLen-1)
	tooLongName := strings.Repeat("c", nameMaxLen+1)

	var tests = []struct {
		name string
		err  string
	}{
		{"", "Название роли не может быть пустым"},
		{tooShortName, lenRangeError},
		{tooLongName, lenRangeError},
		{"alreadyExistingRole", "Название роли уже существует"},
	}

	validator := NewRoleValidator()
	for _, test := range tests {
		model := RoleCreationModel{test.name}
		validationErrors := validator.ValidateRoleCreationModel(model, isNameExists)
		expectedValidationErrors := &ValidationErrors{map[string]string{"name": test.err}}
		if !reflect.DeepEqual(validationErrors, expectedValidationErrors) {
			t.Errorf("got %v, expected %v", validationErrors, expectedValidationErrors)
		}
	}
}

func TestValidateRoleCreationModel_Valid(t *testing.T) {
	const isNameExists = false
	var names = []string{
		"validRole",
		"valid_role",
		"valid-role",
		"valid role",
		"valid role 42",
	}

	validator := NewRoleValidator()
	for _, name := range names {
		model := RoleCreationModel{name}
		validationErrors := validator.ValidateRoleCreationModel(model, isNameExists)
		if validationErrors != nil {
			t.Errorf("got %v, expected %v", validationErrors, nil)
		}
	}
}
