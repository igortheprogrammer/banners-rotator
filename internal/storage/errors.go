package storage

import "errors"

var (
	ErrSlotNotCreated       = errors.New("slot not created")
	ErrBannerNotCreated     = errors.New("banner not created")
	ErrGroupNotCreated      = errors.New("group not created")
	ErrRotationNotCreated   = errors.New("rotation not created")
	ErrRotationNotDeleted   = errors.New("rotation not deleted")
	ErrViewEventNotCreated  = errors.New("view event not created")
	ErrClickEventNotCreated = errors.New("click event not created")
)
