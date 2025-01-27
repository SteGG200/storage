package db

type UploadSession struct {
	Token     string `db:"token"`
	Path      string `db:"path"`
	Filename  string `db:"filename"`
	CreatedAt int64  `db:"created_at"`
}
