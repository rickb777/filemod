package filemod

import (
	"os"
	"time"
)

// FileMetaInfo holds information about one file (which may not exist).
type FileMetaInfo struct {
	path string
	err  error
	fi   os.FileInfo
}

func (file FileMetaInfo) Exists() bool {
	return file.fi != nil && file.err == nil
}

func (file FileMetaInfo) Path() string {
	return file.path
}

func (file FileMetaInfo) Name() string {
	if file.fi == nil {
		return ""
	}
	return file.fi.Name()
}

func (file FileMetaInfo) Size() int64 {
	if file.fi == nil {
		return 0
	}
	return file.fi.Size()
}

func (file FileMetaInfo) Mode() os.FileMode {
	if file.fi == nil {
		return 0
	}
	return file.fi.Mode()
}

func (file FileMetaInfo) ModTime() time.Time {
	if file.fi == nil {
		return time.Time{}
	}
	return file.fi.ModTime()
}

func (file FileMetaInfo) IsDir() bool {
	if file.fi == nil {
		return false
	}
	return file.fi.IsDir()
}

func (file FileMetaInfo) Sys() interface{} {
	if file.fi == nil {
		return nil
	}
	return file.fi.Sys()
}

func (file FileMetaInfo) Err() error {
	return file.err
}

//-------------------------------------------------------------------------------------------------

func Stat(path string) FileMetaInfo {
	Debug("stat %q\n", path)
	if path == "" {
		return FileMetaInfo{path: path}
	}

	info, err := fs.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			Debug("%q does not exist.\n", path)
			return FileMetaInfo{path: path}
		} else {
			Debug("%q stat error %v.\n", path, err)
			return FileMetaInfo{path: path, err: err}
		}
	}

	return FileMetaInfo{
		path: path,
		fi:   info,
	}
}

func (file FileMetaInfo) Newer(other FileMetaInfo) FileMetaInfo {
	if file.NewerThan(other) {
		return file
	}
	return other
}

func (file FileMetaInfo) Older(other FileMetaInfo) FileMetaInfo {
	if file.OlderThan(other) {
		return file
	}
	return other
}

func (file FileMetaInfo) NewerThan(other FileMetaInfo) bool {
	return file.ModTime().After(other.ModTime())
}

func (file FileMetaInfo) OlderThan(other FileMetaInfo) bool {
	return file.ModTime().Before(other.ModTime())
}

//-------------------------------------------------------------------------------------------------

// Debug is a function that prints trace information. By default it does nothing;
// set it to (e.g.) 'fmt.Printf' to enable messages.
var Debug = func(message string, args ...interface{}) {}
