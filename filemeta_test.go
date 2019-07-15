package filemod

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"time"
)

func TestBlank(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Stat("")
	g.Expect(m.Exists()).NotTo(BeTrue())
}

func TestMissing(t *testing.T) {
	g := NewGomegaWithT(t)
	fs = osFacade{} // the real deal
	m := Stat("/etc/this-does-not-exist")
	g.Expect(m.Path()).To(Equal("/etc/this-does-not-exist"))
	g.Expect(m.Name()).To(Equal(""))
	g.Expect(m.Exists()).To(BeFalse())
	g.Expect(m.IsDir()).To(BeFalse())
	g.Expect(m.Mode()).To(BeEquivalentTo(0))
	g.Expect(m.Size()).To(BeEquivalentTo(0))
	g.Expect(m.ModTime().IsZero()).To(BeTrue())
	g.Expect(m.Sys()).To(BeNil())
	g.Expect(m.Err()).To(BeNil()) // nil error for files that do not exist
}

func TestHosts(t *testing.T) {
	g := NewGomegaWithT(t)
	fs = osFacade{} // the real deal
	m := Stat("/etc/hosts")
	g.Expect(m.Path()).To(Equal("/etc/hosts"))
	g.Expect(m.Name()).To(Equal("hosts"))
	g.Expect(m.Exists()).To(BeTrue())
	g.Expect(m.IsDir()).To(BeFalse())
	g.Expect(m.Mode()).NotTo(BeEquivalentTo(0))
	g.Expect(m.Size()).To(BeNumerically(">", 0))
	g.Expect(m.ModTime().IsZero()).To(BeFalse())
	g.Expect(m.Sys()).NotTo(BeNil())
	g.Expect(m.Err()).To(BeNil())
}

func TestEtc(t *testing.T) {
	g := NewGomegaWithT(t)
	fs = osFacade{} // the real deal
	m := Stat("/etc")
	g.Expect(m.Path()).To(Equal("/etc"))
	g.Expect(m.Name()).To(Equal("etc"))
	g.Expect(m.Exists()).To(BeTrue())
	g.Expect(m.IsDir()).To(BeTrue())
	g.Expect(m.Mode()).NotTo(BeEquivalentTo(0))
	g.Expect(m.Size()).To(BeNumerically(">", 0))
	g.Expect(m.ModTime().IsZero()).To(BeFalse())
	g.Expect(m.Sys()).NotTo(BeNil())
	g.Expect(m.Err()).To(BeNil())
}

func TestHappy(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now().UTC()
	fi := fileInfo{name: "foo", size: 123, modTime: now}
	fs = osStub{map[string]fileInfo{"/a/b/c/foo": fi}, nil} // global
	m := New("/a/b/c/foo")
	g.Expect(len(m)).To(Equal(1))
	g.Expect(m[0].Path()).To(Equal("/a/b/c/foo"))
	g.Expect(m[0].ModTime()).To(Equal(now))
	g.Expect(m[0].Exists()).To(BeTrue())
	g.Expect(m[0].Err()).To(BeNil())
}

func TestPartition(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = osStub{map[string]fileInfo{"/a/foo": fa, "/a/b/c/d": fd}, nil} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))
	files, dirs, absent := m.Partition()
	g.Expect(len(files)).To(Equal(1))
	g.Expect(len(dirs)).To(Equal(1))
	g.Expect(len(files)).To(Equal(1))
	g.Expect(len(absent)).To(Equal(1))
	g.Expect(files[0].IsDir()).To(BeFalse())
	g.Expect(files[0].Exists()).To(BeTrue())
	g.Expect(dirs[0].IsDir()).To(BeTrue())
	g.Expect(dirs[0].Exists()).To(BeTrue())
	g.Expect(absent[0].Exists()).To(BeFalse())
}

//func TestYoungest(t *testing.T) {
//	now := time.Now().UTC()
//	a := FileMetaInfo{"", "a", now.Add(-2 * time.Minute), ""}
//	b := FileMetaInfo{"", "b", now.Add(-1 * time.Minute), ""}
//	c := FileMetaInfo{"", "c", now, ""}
//
//	y1 := YoungestFile(c, a, b)
//	if y1 != c {
//		t.Errorf("Expected %v but got %v", c, y1)
//	}
//
//	y2 := YoungestFile(a, b, c)
//	if y2 != c {
//		t.Errorf("Expected %v but got %v", c, y2)
//	}
//}

//-------------------------------------------------------------------------------------------------

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	err     error
}

func (fi fileInfo) Name() string {
	return fi.name
}

func (fi fileInfo) Size() int64 {
	return fi.size
}

func (fi fileInfo) Mode() os.FileMode {
	return fi.mode
}

func (fi fileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi fileInfo) IsDir() bool {
	return fi.isDir
}

func (fi fileInfo) Sys() interface{} {
	return nil
}

var _ os.FileInfo = fileInfo{}

//-------------------------------------------------------------------------------------------------

type osStub struct {
	fi map[string]fileInfo
	e  error
}

func (o osStub) Stat(name string) (os.FileInfo, error) {
	v, exists := o.fi[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return v, o.e
}

var _ OS = osStub{}
