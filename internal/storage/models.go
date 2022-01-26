package storage

type Slot struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
}

type Banner struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
}

type Group struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
}
