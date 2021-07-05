package confighelper

type Type int

const (
	ENV Type = iota
	YAML
	// CONF
)

type Config interface {
}

func LoadConfig(from Type, file string, to Config) {
	switch Type {
	case ENV:
		configFromEnv(to)
	case YAML:
		configFromYAML(file, to)
	default:
		panic("config type unsupported now")
	}
}

func configFromEnv(to Config) {
	v := reflect.ValueOf(to).Elem()
	if v.Kind() == reflect.Invalid {
		v = reflect.New(reflect.TypeOf(to).Elem()).Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		defaultValue := t.Field(i).Tag.Get("default")
		if defaultValue != "" {
			continue
		}
	}
}

func configFromYAML(file string, to Config) {

}
