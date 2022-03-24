package basemodels

import "time"

// Model 模型接口
type Model interface {
	LoadM2M()
}

type ModelBase struct{}

func (mb *ModelBase) LoadM2M() {}

type NormalModel struct {
	ModelBase
	ID          int64     `gorm:"column(id);auto;pk" description:"ID" gea:"title=ID;pk=true"`
	Deleted     bool      `gorm:"default(false);description(已删除)" description:"已删除"`
	CreatedTime time.Time `gorm:"autoCreateTime" description:"创建时间" gea:"title=创建时间;dbtype=Datetime"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" description:"修改时间" gea:"title=修改时间;dbtype=Datetime"`
}
