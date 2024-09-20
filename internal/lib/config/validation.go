package config

type Validator interface {
	validate() error
}

func mustValidate(val ...Validator) {
	for i := range val {
		if err := val[i].validate(); err != nil {
			panic("failed to validate config: " + err.Error())
		}
	}
}
