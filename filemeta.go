package filemod

import (
	"os"
	"sort"
	"strings"
	"time"
)

// FileMetaInfo holds information about one file (which may not exist).
type FileMetaInfo struct {
	Path    string
	Exists  bool
	ModTime time.Time
	Err     error
}

// Files holdsdata on a group of files.
type Files []FileMetaInfo

// New builds file information for one or more files.
func New(includeEmpties bool, paths ...string) Files {
	result := make(Files, len(paths))

	i := 0
	for _, p := range paths {
		fm := singleFileMeta(p)
		if fm.Exists {
			result[i] = fm
			i += 1
		} else {
			if includeEmpties {
				result[i] = fm
				i += 1
			}
		}
	}

	return result[:i]
}

func singleFileMeta(path string) FileMetaInfo {
	Debug("stat %q\n", path)
	if path == "" {
		return FileMetaInfo{Path: path, Exists: false}
	}

	info, err := fs.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			Debug("%q does not exist.\n", path)
			return FileMetaInfo{Path: path, Exists: false}
		} else {
			return FileMetaInfo{Path: path, Err: err}
		}
	}

	return FileMetaInfo{Path: path, Exists: true, ModTime: info.ModTime()}
}

func (file FileMetaInfo) Younger(other FileMetaInfo) FileMetaInfo {
	if other.ModTime.After(file.ModTime) {
		return other
	}
	return file
}

func (file FileMetaInfo) Older(other FileMetaInfo) FileMetaInfo {
	if other.ModTime.Before(file.ModTime) {
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

	if files[len(files)-1].ModTime.Before(other[0].ModTime) {
		return AllAreOlder
	}

	if other[len(other)-1].ModTime.Before(files[0].ModTime) {
		return AllAreYounger
	}

	return Overlapping
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
	return files[i].ModTime.Before(files[j].ModTime)
}

//-------------------------------------------------------------------------------------------------

// Errors gets any errors encountered when checking the files.
func (files Files) Errors() Errors {
	ee := make(Errors, 0, len(files))
	for _, f := range files {
		if f.Err != nil {
			ee = append(ee, f.Err)
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
