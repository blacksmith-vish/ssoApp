package models

type App struct {
	Name   string
	Secret string
	ID     string
}

func (app App) GetID() string {
	return app.ID
}

func (app App) GetSecret() []byte {
	return []byte(app.Secret)
}
