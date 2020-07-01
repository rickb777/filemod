package filemod

import (
	"sort"
	"strings"
)

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

// Of simply constructs a list of files. This is a convenience function.
func Of(files ...FileMetaInfo) Files {
	return files
}

//-------------------------------------------------------------------------------------------------

// Comparison enumerates how two sets of files compare.
type comparison int

const (
	undefined comparison = iota
	allAreNewer
	overlapping
	allAreOlder
)

// compare sorts the files and another set of files, then compares the youngest
// and oldest in each set.
func (files Files) compare(other Files) comparison {
	if len(files) == 0 || len(other) == 0 {
		return undefined
	}

	files.SortedByModTime()
	other.SortedByModTime()

	// if the last of 'files' is before the first of 'other'...
	if files[len(files)-1].ModTime().Before(other[0].ModTime()) {
		return allAreOlder
	}

	// if the last of 'other' is before the first of 'files'...
	if other[len(other)-1].ModTime().Before(files[0].ModTime()) {
		return allAreNewer
	}

	return overlapping
}

func (files Files) AllAreOlderThan(other Files) bool {
	return files.compare(other) == allAreOlder
}

func (files Files) OverlapsWith(other Files) bool {
	return files.compare(other) == overlapping
}

func (files Files) AllAreNewerThan(other Files) bool {
	return files.compare(other) == allAreNewer
}

//-------------------------------------------------------------------------------------------------

// Partition separates files and riectories that exist from those that don't.
func (files Files) Partition() (allFiles, allDirs, absent Files) {
	// to avoid unnecessary memory allocation, the first pass counts the items
	nf, nd, na := 0, 0, 0
	for _, f := range files {
		if f.Exists() {
			if f.IsDir() {
				nd++
			} else {
				nf++
			}
		} else {
			na++
		}
	}

	allFiles = make(Files, 0, nf)
	allDirs = make(Files, 0, nd)
	absent = make(Files, 0, na)

	// the second pass builds the results
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

// SortedByModTime rearranges the files into modification-time order with the oldest first.
func (files Files) SortedByModTime() {
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

// SortedByPath rearranges the files into path order
func (files Files) SortedByPath() {
	sort.Stable(byPath(files))
}

type byPath Files

func (files byPath) Len() int {
	return len(files)
}

func (files byPath) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files byPath) Less(i, j int) bool {
	return files[i].path < files[j].path
}

//-------------------------------------------------------------------------------------------------

// SortedBySize rearranges the files into size order
func (files Files) SortedBySize() {
	sort.Stable(bySize(files))
}

type bySize Files

func (files bySize) Len() int {
	return len(files)
}

func (files bySize) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (files bySize) Less(i, j int) bool {
	return files[i].Size() < files[j].Size()
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
