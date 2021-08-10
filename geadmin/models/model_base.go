package geamodels

import "time"

// Model 模型接口
type Model interface {
	LoadM2M()
}

type ModelBase struct{}

func (mb *ModelBase) LoadM2M() {}

type NormalModel struct {
	ModelBase
	ID          int64     `orm:"column(id);auto;pk" description:"ID" display:"title=ID;pk=true"`
	Deleted     bool      `orm:"default(false);description(已删除)" description:"已删除"`
	CreatedTime time.Time `orm:"auto_now_add;type(datetime)" description:"创建时间" display:"title=创建时间;dbtype=Datetime"`
	UpdatedTime time.Time `orm:"auto_now;type(datetime)" description:"修改时间" display:"title=修改时间;dbtype=Datetime"`
}
