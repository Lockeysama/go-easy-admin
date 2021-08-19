package geamodels

// Model 模型接口
type Model interface{}

type M2MModel interface {
	LoadM2M()
}

type ModelBase struct{}

func (mb *ModelBase) LoadM2M() {}
