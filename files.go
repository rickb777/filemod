package filemod

import (
	"sort"
	"strings"
)

// Files holds data on a group of files.
// Use the builtin 'append' if required; also use slicing as required.
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

// AllAreOlderThan compares the file modification timestamps of two lists of files,
// returning true if all of 'files' are older than all of 'other'.
func (files Files) AllAreOlderThan(other Files) bool {
	return files.compare(other) == allAreOlder
}

// OverlapsWith compares the file modification timestamps
// of two lists of files, returning true if the range of modification times of
// 'files' overlaps with the range of modification times of the 'other' list.
func (files Files) OverlapsWith(other Files) bool {
	return files.compare(other) == overlapping
}

// AllAreNewerThan compares the file modification timestamps of two lists of files,
// returning true if all of 'files' are newer than all of 'other'.
func (files Files) AllAreNewerThan(other Files) bool {
	return files.compare(other) == allAreNewer
}

//-------------------------------------------------------------------------------------------------

// Partition separates files and directories that exist from those that don't.
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

// Filter returns only those items for which a predicate p returns true.
func (files Files) Filter(p func(info FileMetaInfo) bool) Files {
	// to avoid unnecessary memory allocation, the first pass counts the items
	n := 0
	for _, f := range files {
		if p(f) {
			n++
		}
	}

	result := make(Files, 0, n)

	// the second pass builds the results
	for _, f := range files {
		if p(f) {
			result = append(result, f)
		}
	}

	return result
}

// FilesOnly returns only files that exist (i.e. not directories).
func (files Files) FilesOnly() Files {
	return files.Filter(func(f FileMetaInfo) bool {
		return f.Exists() && !f.IsDir()
	})
}

// DirectoriesOnly returns only directories that exist.
func (files Files) DirectoriesOnly() Files {
	return files.Filter(func(f FileMetaInfo) bool {
		return f.Exists() && f.IsDir()
	})
}

// PresentOnly returns only those files/directories that exist.
func (files Files) PresentOnly() Files {
	return files.Filter(func(f FileMetaInfo) bool {
		return f.Exists()
	})
}

// AbsentOnly returns only those files/directories that don't exist.
func (files Files) AbsentOnly() Files {
	return files.Filter(func(f FileMetaInfo) bool {
		return !f.Exists()
	})
}

// First returns the first file, equivalent to files[0] with a check
// for an empty slice.
func (files Files) First() *FileMetaInfo {
	if len(files) == 0 {
		return nil
	}
	return &files[0]
}

// Last returns the last file, equivalent to files[n-1] with a check
// for an empty slice.
func (files Files) Last() *FileMetaInfo {
	if len(files) == 0 {
		return nil
	}
	return &files[len(files)-1]
}

//-------------------------------------------------------------------------------------------------

// SortedByModTime rearranges the files into modification-time order with the oldest first.
// It returns the modified list.
func (files Files) SortedByModTime() Files {
	sort.Stable(byModTime(files))
	return files
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

// SortedByPath rearranges the files into path order, similar to comparing pairs of strings.
// It returns the modified list.
func (files Files) SortedByPath() Files {
	sort.Stable(byPath(files))
	return files
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

// SortedBySize rearranges the files into size order. It returns the modified list.
func (files Files) SortedBySize() Files {
	sort.Stable(bySize(files))
	return files
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

// Errors holds a slice of errors and is itself an error.
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
