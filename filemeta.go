package filemod

import (
	"os"
	"sort"
	"strings"
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

// Files holdsdata on a group of files.
type Files []FileMetaInfo

// New builds file information for one or more files. If filesystem errors
// arise, these are held in the files returned and can be inspected later.
func New(paths ...string) Files {
	result := make(Files, len(paths))

	for i, p := range paths {
		fm := Stat(p)
		result[i] = fm
	}

	return result
}

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

func (file FileMetaInfo) Younger(other FileMetaInfo) FileMetaInfo {
	if other.ModTime().After(file.ModTime()) {
		return other
	}
	return file
}

func (file FileMetaInfo) Older(other FileMetaInfo) FileMetaInfo {
	if other.ModTime().Before(file.ModTime()) {
		return other
	}
	return file
}

//-------------------------------------------------------------------------------------------------

// Comparison enumerates how two sets of files compare.
type Comparison int

const (
	Undefined Comparison = iota
	AllAreYounger
	Overlapping
	AllAreOlder
)

// Compare sorts the files and another set of files, then compares the youngest
// and oldest in each set.
func (files Files) Compare(other Files) Comparison {
	if len(files) == 0 || len(other) == 0 {
		return Undefined
	}

	files.Sorted()
	other.Sorted()

	if files[len(files)-1].ModTime().Before(other[0].ModTime()) {
		return AllAreOlder
	}

	if other[len(other)-1].ModTime().Before(files[0].ModTime()) {
		return AllAreYounger
	}

	return Overlapping
}

//-------------------------------------------------------------------------------------------------

// Partition separates files that exist from those that don't.
func (files Files) Partition() (allFiles, allDirs, absent Files) {
	nf, nd, na := 0, 0, 0
	for _, f := range files {
		if f.Exists() && f.IsDir() {
			nd++
		} else if f.Exists() {
			nf++
		} else {
			na++
		}
	}

	allFiles = make(Files, 0, nf)
	allDirs = make(Files, 0, nd)
	absent = make(Files, 0, na)

	for _, f := range files {
		if f.Exists() && f.IsDir() {
			allDirs = append(allDirs, f)
		} else if f.Exists() {
			allFiles = append(allFiles, f)
		} else {
			absent = append(absent, f)
		}
	}

	return allFiles, allDirs, absent
}

//-------------------------------------------------------------------------------------------------

// Sorted rearranges the files into modification-time order with the oldest first.
func (files Files) Sorted() {
	sort.Stable(byModTime(files))
}

type byModTime Files

func (files byModTime) Len() int {
	return len(files)
}

func (files byModTime) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files byModTime) Less(i, j int) bool {
	return files[i].ModTime().Before(files[j].ModTime())
}

//-------------------------------------------------------------------------------------------------

// Errors gets any errors encountered when checking the files.
func (files Files) Errors() Errors {
	ee := make(Errors, 0, len(files))
	for _, f := range files {
		if f.err != nil {
			ee = append(ee, f.err)
		}
	}
	return ee
}

//-------------------------------------------------------------------------------------------------

// Errors holds a sequence of errors.
type Errors []error

// Error gets the error string, built from each error conjoined with a newline.
func (ee Errors) Error() string {
	buf := &strings.Builder{}
	for i, e := range ee {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(e.Error())
	}
	return buf.String()
}

//-------------------------------------------------------------------------------------------------

// Debug is a function that prints trace information. By default it does nothing;
// set it to (e.g.) 'fmt.Printf' to enable messages.
var Debug = func(message string, args ...interface{}) {}
