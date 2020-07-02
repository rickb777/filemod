package filemod

import "os"

// statter is a seam for testing
type statter interface {
	Stat(name string) (os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
}

type osFacade struct{}

func (o osFacade) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (o osFacade) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

var fs statter = osFacade{}
