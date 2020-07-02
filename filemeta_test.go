package filemod

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"time"
)

func TestBlank(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Stat("")
	g.Expect(m1.Exists()).NotTo(BeTrue())

	m2 := Lstat("")
	g.Expect(m2.Exists()).NotTo(BeTrue())
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
	m1 := Stat("/etc/hosts")

	// Then...
	g.Expect(m1.Path()).To(Equal("/etc/hosts"))
	g.Expect(m1.Name()).To(Equal("hosts"))
	g.Expect(m1.Exists()).To(BeTrue())
	g.Expect(m1.IsDir()).To(BeFalse())
	g.Expect(m1.Mode()).NotTo(BeEquivalentTo(0))
	g.Expect(m1.Size()).To(BeNumerically(">", 0))
	g.Expect(m1.ModTime().IsZero()).To(BeFalse())
	g.Expect(m1.Sys()).NotTo(BeNil())
	g.Expect(m1.Err()).To(BeNil())

	// When...
	m2 := m1.Refresh()

	// Then...
	g.Expect(m2.Path()).To(Equal("/etc/hosts"))
	g.Expect(m2.Name()).To(Equal("hosts"))
	g.Expect(m2.Exists()).To(BeTrue())
	g.Expect(m2.IsDir()).To(BeFalse())
	g.Expect(m2.Mode()).NotTo(BeEquivalentTo(0))
	g.Expect(m2.Size()).To(BeNumerically(">", 0))
	g.Expect(m2.ModTime().IsZero()).To(BeFalse())
	g.Expect(m2.Sys()).NotTo(BeNil())
	g.Expect(m2.Err()).To(BeNil())
}

func TestLstatHosts(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	fs = osFacade{} // the real deal

	// When...
	m1 := Lstat("/etc/hosts")

	// Then...
	g.Expect(m1.Path()).To(Equal("/etc/hosts"))
	g.Expect(m1.Name()).To(Equal("hosts"))
	g.Expect(m1.Exists()).To(BeTrue())
	g.Expect(m1.IsDir()).To(BeFalse())
	g.Expect(m1.Mode()).NotTo(BeEquivalentTo(0))
	g.Expect(m1.Size()).To(BeNumerically(">", 0))
	g.Expect(m1.ModTime().IsZero()).To(BeFalse())
	g.Expect(m1.Sys()).NotTo(BeNil())
	g.Expect(m1.Err()).To(BeNil())
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

func TestRefresh(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	t1 := fileInfo{name: "t1", size: 11, modTime: now.Add(-2)}
	t2 := fileInfo{name: "t2", size: 22, modTime: now.Add(-1)}
	fs = &osStub{[]fileInfo{t1, t2}} // global

	// When...
	m1 := Stat("/t")
	m2 := m1.Refresh()

	// Then...
	g.Expect(m1.Path()).To(Equal("/t"))
	g.Expect(m2.Path()).To(Equal("/t"))
	g.Expect(m1.ModTime().Before(m2.ModTime())).To(BeTrue())
}

func TestYoungerThan(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-11)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-12)}
	fs = &osStub{[]fileInfo{a1, a2}} // global

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
	fs = &osStub{[]fileInfo{a1, a2}} // global

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

// stub works as a stack of file infos from which items are popped on each Stat/Lstat call
type osStub struct {
	fi []fileInfo
}

func (stub *osStub) Stat(_ string) (os.FileInfo, error) {
	return stub.fakeStat()
}

func (stub *osStub) Lstat(_ string) (os.FileInfo, error) {
	return stub.fakeStat()
}

func (stub *osStub) fakeStat() (os.FileInfo, error) {
	if len(stub.fi) == 0 {
		return nil, os.ErrNotExist
	}
	v := stub.fi[0]
	stub.fi = stub.fi[1:]
	return v, v.err
}

var _ statter = &osStub{}
