package models

type App struct {
	ID        int64  `db:"a_id"`
	Name      string `db:"name"`
	SecretKey string `db:"secret_key"`
}
