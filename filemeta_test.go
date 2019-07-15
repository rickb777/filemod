package filemod

import (
	"errors"
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

func TestStatMissing(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	fs = osFacade{} // the real deal

	// When...
	m := Stat("/etc/this-does-not-exist")

	// Then...
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

func TestStatHosts(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	fs = osFacade{} // the real deal

	// When...
	m := Stat("/etc/hosts")

	// Then...
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

func TestStatEtc(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	fs = osFacade{} // the real deal

	// When...
	m := Stat("/etc")

	// Then...
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
	fs = osStub(map[string]fileInfo{"/a/b/c/foo": fi}) // global

	// When...
	m := New("/a/b/c/foo")

	// Then...
	g.Expect(len(m)).To(Equal(1))
	g.Expect(m[0].Path()).To(Equal("/a/b/c/foo"))
	g.Expect(m[0].ModTime()).To(Equal(now))
	g.Expect(m[0].Exists()).To(BeTrue())
	g.Expect(m[0].Err()).To(BeNil())
}

func TestPartition(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = osStub(map[string]fileInfo{"/a/foo": fa, "/a/b/c/d": fd}) // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	files, dirs, absent := m.Partition()

	// Then...
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

func TestCompare(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-11)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-12)}
	b1 := fileInfo{name: "b1", size: 11, modTime: now.Add(-2)}
	b2 := fileInfo{name: "b2", size: 22, modTime: now.Add(-1)}
	fs = osStub(map[string]fileInfo{"/a1": a1, "/a2": a2, "/b1": b1, "/b2": b2}) // global
	g1 := New("/a1", "/a2")
	g2 := New("/b1", "/b2")
	ab := New("/a1", "/b2")
	x := New()

	// When...
	v1 := g1.Compare(g2)
	v2 := g2.Compare(g1)
	v3 := g2.Compare(ab)
	v4a := g2.Compare(x)
	v4b := x.Compare(g2)

	// Then...
	g.Expect(v1).To(Equal(AllAreOlder))
	g.Expect(v2).To(Equal(AllAreYounger))
	g.Expect(v3).To(Equal(Overlapping))
	g.Expect(v4a).To(Equal(Undefined))
	g.Expect(v4b).To(Equal(Undefined))
}

func TestYounger(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-11)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-12)}
	fs = osStub(map[string]fileInfo{"/a1": a1, "/a2": a2}) // global

	// When...
	m1 := Stat("/a1")
	m2 := Stat("/a2")

	y1 := m1.Younger(m2)
	y2 := m2.Younger(m1)

	// Then...
	g.Expect(y1.Name()).To(Equal("a1"))
	g.Expect(y2.Name()).To(Equal("a1"))
}

func TestOlder(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-11)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-12)}
	fs = osStub(map[string]fileInfo{"/a1": a1, "/a2": a2}) // global

	// When...
	m1 := Stat("/a1")
	m2 := Stat("/a2")

	y1 := m1.Older(m2)
	y2 := m2.Older(m1)

	// Then...
	g.Expect(y1.Name()).To(Equal("a2"))
	g.Expect(y2.Name()).To(Equal("a2"))
}

func TestErrors(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	a1 := fileInfo{err: errors.New("a1")}
	a2 := fileInfo{err: errors.New("a2")}
	fs = osStub(map[string]fileInfo{"/a1": a1, "/a2": a2}) // global
	g1 := New("/a1", "/a2")

	// When...
	s := g1.Errors().Error()

	// Then...
	g.Expect(s).To(Equal("a1\na2"))
}

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

type osStub map[string]fileInfo

func (stub osStub) Stat(name string) (os.FileInfo, error) {
	v, exists := stub[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return v, v.err
}

var _ OS = osStub{}
