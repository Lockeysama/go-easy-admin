package geacontrollers

import (
	geamodels "github.com/lockeysama/go-easy-admin/geadmin/models"
)

type GEADataBase interface {
	GEADataBaseCount(model geamodels.Model, filters map[string]interface{}) (int64, error)
	GEADataBaseQueryList(
		model geamodels.Model,
		page int, limit int,
		filters map[string]interface{},
		order map[string]string,
		loadRel bool,
	) interface{}
	GEADataBaseQueryRow(model geamodels.Model, filters map[string]interface{}, loadRel bool) interface{}

	GEADataBaseInsert(model geamodels.Model) (int64, error)
	GEADataBaseUpdate(model geamodels.Model, filters map[string]interface{}, params map[string]interface{}) (int64, error)
	GEADataBaseDelete(model geamodels.Model, filters map[string]interface{}) (int64, error)
	GEADataM2MUpdate(model geamodels.Model, fieldName string, values []interface{}, action string) error
}
