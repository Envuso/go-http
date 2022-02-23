package cache

import (
	"github.com/rclone/rclone/lib/cache"

	"gohttp/lib/storage/fs"
)

type (
	Store interface {
		Clear()
		Entries() int
		Put(path string, file fs.File) error
		Rename(oldPath, newPath string) (file fs.File, found bool)
		Get(path string, create CreateFunc) (file fs.File, err error)
	}

	CreateFunc cache.CreateFunc
)
