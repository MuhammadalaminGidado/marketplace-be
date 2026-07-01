package utils

import "gorm.io/gorm"

type Permission struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func GetPermissionNamesByEntityID(db *gorm.DB, entityID uint) ([]string, error) {
	var permissionNames []string

	err := db.Model(&Permission{}).
		Select("permissions.name").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN entity_roles er ON er.role_id = rp.role_id").
		Where("er.entity_id = ?", entityID).
		Pluck("name", &permissionNames).
		Error

	if err != nil {
		return nil, err
	}

	return permissionNames, nil
}
func HasPermission(db *gorm.DB, entityID uint, perm string) (bool, error) {
	permissionNames, err := GetPermissionNamesByEntityID(db, entityID)
	if err != nil {
		return false, err
	}

	for _, name := range permissionNames {
		if name == perm {
			return true, nil
		}
	}

	return false, nil
}
