package geamodels

// Model 模型接口
type Model interface {
	LoadM2M()
}

type ModelBase struct{}

func (mb *ModelBase) LoadM2M() {}
