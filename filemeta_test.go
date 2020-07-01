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

func TestOf(t *testing.T) {
	g := NewGomegaWithT(t)
	m1 := Stat("a")
	m2 := Stat("2")
	ff := Of(m1, m2)
	g.Expect(ff).To(HaveLen(2))
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

func TestYoungerThan(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-11)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-12)}
	fs = osStub(map[string]fileInfo{"/a1": a1, "/a2": a2}) // global

	// When...
	m1 := Stat("/a1")
	m2 := Stat("/a2")

	y1 := m1.Newer(m2)
	y2 := m2.Newer(m1)

	// Then...
	g.Expect(y1.Name()).To(Equal("a1"))
	g.Expect(y2.Name()).To(Equal("a1"))
}

func TestOlderThan(t *testing.T) {
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
