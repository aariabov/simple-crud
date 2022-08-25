package roles

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type RoleRepo struct {
	db *sql.DB
}

func NewRoleRepo(db *sql.DB) RoleRepo {
	return RoleRepo{db}
}

func (r RoleRepo) getAllRoles() ([]RoleVm, error) {
	rows, err := r.db.Query("select id, name from asp_net_roles")
	if err != nil {
		return nil, fmt.Errorf("getAllRoles error: %v", err)
	}
	defer rows.Close()

	roles := []RoleVm{}
	for rows.Next() {
		p := RoleVm{}
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, p)
	}

	return roles, nil
}

func (r RoleRepo) insertNewRole(model RoleCreationModel) (string, error) {
	newId := uuid.New().String()
	_, err := r.db.Exec("insert into asp_net_roles (id, name) values ($1, $2) returning id",
		newId, model.Name)
	if err != nil {
		return "", fmt.Errorf("insertNewRole with id '%s' name '%s' error: %v", newId, model.Name, err)
	}

	return newId, nil
}

func (r RoleRepo) deleteRole(id string) error {
	_, err := r.db.Exec("delete from asp_net_roles where id = $1", id)
	if err != nil {
		return fmt.Errorf("deleteRole with id '%s' error: %v", id, err)
	}

	return nil
}

func (r RoleRepo) updateRole(model RoleUpdatingModel) error {
	_, err := r.db.Exec("update asp_net_roles set name = $1 where id = $2",
		model.Name, model.Id)
	if err != nil {
		return fmt.Errorf("updateRole with id '%s' name '%s' error: %v", model.Id, model.Name, err)
	}

	return nil
}

func (r RoleRepo) isRoleExists(name string) (bool, error) {
	var id string
	err := r.db.QueryRow("select id from asp_net_roles where name = $1", name).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, fmt.Errorf("isRoleExists with name '%s' error: %v", name, err)
		}

		return false, nil
	}

	return true, nil
}

func (r RoleRepo) isAnyUserBelongToRole(id string) (bool, error) {
	var isExists bool
	err := r.db.QueryRow("select exists(select user_id from asp_net_user_roles where role_id = $1)", id).
		Scan(&isExists)
	if err != nil {
		return false, fmt.Errorf("isAnyUserBelongToRole with roleId '%s' error: %v", id, err)
	}

	return isExists, nil
}

func (r RoleRepo) getRoleById(id string) (*RoleVm, error) {
	var role RoleVm
	err := r.db.QueryRow("select id, name from asp_net_roles where id = $1", id).
		Scan(&role.Id, &role.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("getRoleById '%s' error: %v", id, err)
		}

		return nil, nil
	}

	return &role, nil
}
