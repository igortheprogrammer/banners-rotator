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

type ViewEvent struct {
	SlotID   int64 `db:"slot_id" json:"slot_id"`
	BannerID int64 `db:"banner_id" json:"banner_id"`
	GroupID  int64 `db:"group_id" json:"group_id"`
	Date     int64 `db:"date" json:"date"`
}

type ClickEvent struct {
	SlotID   int64 `db:"slot_id" json:"slot_id"`
	BannerID int64 `db:"banner_id" json:"banner_id"`
	GroupID  int64 `db:"group_id" json:"group_id"`
	Date     int64 `db:"date" json:"date"`
}
