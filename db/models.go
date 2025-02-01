package db

type UploadSession struct {
	Token     string `db:"token"`
	Path      string `db:"path"`
	Filename  string `db:"filename"`
	CreatedAt int64  `db:"created_at"`
}

type ItemPassword struct {
	ID       int64  `db:"id"`
	Path     string `db:"path"`
	Password string `db:"password"`
}
