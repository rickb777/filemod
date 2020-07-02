package filemod

import (
	"os"
	"time"
)

// FileMetaInfo holds information about one file (which may not exist).
type FileMetaInfo struct {
	path string
	err  error
	fi   os.FileInfo // absent if file does not exist
}

// Tests whether the file exists.
func (file FileMetaInfo) Exists() bool {
	return file.fi != nil && file.err == nil
}

// Gets the file path.
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

// Err gets the error, if any, from Stat.
// If there is an error, it will be of type *os.PathError.
func (file FileMetaInfo) Err() error {
	return file.err
}

//-------------------------------------------------------------------------------------------------

// Stat tests a file path using the operating system.
func Stat(path string) FileMetaInfo {
	Debug("stat %q\n", path)
	if path == "" {
		return FileMetaInfo{path: path}
	}

	info, err := fs.Stat(path)

	return newFileMetaInfo(path, err, info)
}

// Lstat tests a file path using the operating system.
// If the file is a symbolic link, the returned FileInfo
// describes the symbolic link. Lstat makes no attempt to follow the link.
func Lstat(path string) FileMetaInfo {
	Debug("lstat %q\n", path)
	if path == "" {
		return FileMetaInfo{path: path}
	}

	info, err := fs.Lstat(path)

	return newFileMetaInfo(path, err, info)
}

func newFileMetaInfo(path string, err error, info os.FileInfo) FileMetaInfo {
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

// Refresh queries the operating system for the status of the file again.
// A new FileMetaInfo is returned that contains the current status of the file.
func (file FileMetaInfo) Refresh() FileMetaInfo {
	return Stat(file.path)
}

//-------------------------------------------------------------------------------------------------

// Newer compares the modification timestamps and returns the
// file that is newer.
func (file FileMetaInfo) Newer(other FileMetaInfo) FileMetaInfo {
	if file.NewerThan(other) {
		return file
	}
	return other
}

// Older compares the modification timestamps and returns the
// file that is older.
func (file FileMetaInfo) Older(other FileMetaInfo) FileMetaInfo {
	if file.OlderThan(other) {
		return file
	}
	return other
}

// NewerThan compares the modification timestamps and returns true
// if this file is newer than the other.
func (file FileMetaInfo) NewerThan(other FileMetaInfo) bool {
	return file.ModTime().After(other.ModTime())
}

// OlderThan compares the modification timestamps and returns true
// if this file is older than the other.
func (file FileMetaInfo) OlderThan(other FileMetaInfo) bool {
	return file.ModTime().Before(other.ModTime())
}

//-------------------------------------------------------------------------------------------------

// Debug is a function that prints trace information. By default it does nothing;
// set it to (e.g.) 'fmt.Printf' to enable messages.
var Debug = func(message string, args ...interface{}) {}
