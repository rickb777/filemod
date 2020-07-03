package filemod

import (
	"errors"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestHappy(t *testing.T) {
	g := NewGomegaWithT(t)
	now := time.Now().UTC()
	fi := fileInfo{name: "foo", size: 123, modTime: now}
	fs = &osStub{[]fileInfo{fi}} // global

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
	fs = &osStub{[]fileInfo{fa, fd}} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	files, dirs, absent := m.Partition()

	// Then...
	g.Expect(len(files)).To(Equal(1))
	g.Expect(len(dirs)).To(Equal(1))
	g.Expect(len(absent)).To(Equal(1))
	g.Expect(files[0].IsDir()).To(BeFalse())
	g.Expect(files[0].Exists()).To(BeTrue())
	g.Expect(dirs[0].IsDir()).To(BeTrue())
	g.Expect(dirs[0].Exists()).To(BeTrue())
	g.Expect(absent[0].Exists()).To(BeFalse())
}

func TestDirectoriesOnly(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = &osStub{[]fileInfo{fa, fd}} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	dirs := m.DirectoriesOnly()

	// Then...
	g.Expect(len(dirs)).To(Equal(1))
	g.Expect(dirs[0].IsDir()).To(BeTrue())
	g.Expect(dirs[0].Exists()).To(BeTrue())
}

func TestFilesOnly(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = &osStub{[]fileInfo{fa, fd}} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	files := m.FilesOnly()

	// Then...
	g.Expect(len(files)).To(Equal(1))
	g.Expect(files[0].IsDir()).To(BeFalse())
	g.Expect(files[0].Exists()).To(BeTrue())
}

func TestAbsentOnly(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = &osStub{[]fileInfo{fa, fd}} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	absent := m.AbsentOnly()

	// Then...
	g.Expect(len(absent)).To(Equal(1))
	g.Expect(absent[0].Exists()).To(BeFalse())
}

func TestPresentOnly(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	fa := fileInfo{name: "foo", size: 123, modTime: now}
	fd := fileInfo{name: "d", size: 456, modTime: now, isDir: true}
	fs = &osStub{[]fileInfo{fa, fd}} // global
	m := New("/a/foo", "/a/b/c/d", "/a/x")
	g.Expect(len(m)).To(Equal(3))

	// When...
	present := m.PresentOnly()

	// Then...
	g.Expect(len(present)).To(Equal(2))
	g.Expect(present[0].IsDir()).To(BeFalse())
	g.Expect(present[0].Exists()).To(BeTrue())
	g.Expect(present[1].IsDir()).To(BeTrue())
	g.Expect(present[1].Exists()).To(BeTrue())
}

func TestSortedByModTime(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	ta1 := now.Add(-1 * time.Minute)
	ta2 := now.Add(-5 * time.Minute)
	tb1 := now.Add(-3 * time.Minute)
	tb2 := now.Add(-4 * time.Minute)

	a1 := fileInfo{name: "a1", size: 11, modTime: ta1}
	a2 := fileInfo{name: "a2", size: 22, modTime: ta2}
	b1 := fileInfo{name: "b1", size: 11, modTime: tb1}
	b2 := fileInfo{name: "b2", size: 22, modTime: tb2}
	x1 := fileInfo{name: "x1"}
	x2 := fileInfo{name: "x2"}
	fs = &osStub{[]fileInfo{a1, a2, b1, b2, x1, x2}}
	group := New("/a1", "/a2", "/b1", "/b2", "/x1", "/x2")

	// When...
	group.SortedByModTime()

	// Then...
	g.Expect(group[0].ModTime().IsZero()).To(BeTrue())
	g.Expect(group[1].ModTime().IsZero()).To(BeTrue())
	g.Expect(group[2].ModTime().Equal(ta2)).To(BeTrue())
	g.Expect(group[3].ModTime().Equal(tb2)).To(BeTrue())
	g.Expect(group[4].ModTime().Equal(tb1)).To(BeTrue())
	g.Expect(group[5].ModTime().Equal(ta1)).To(BeTrue())
}

func TestSortedByPath(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	a1 := fileInfo{name: "a1"}
	a2 := fileInfo{name: "a2"}
	b1 := fileInfo{name: "b1"}
	b2 := fileInfo{name: "b2"}
	x1 := fileInfo{name: "x1"}
	x2 := fileInfo{name: "x2"}
	fs = &osStub{[]fileInfo{a1, x1, b1, a2, b2, x2}}
	group := New("/a1", "/x1", "/b1", "/a2", "/b2", "/x2")

	// When...
	group.SortedByPath()

	// Then...
	g.Expect(group[0].Name()).To(Equal("a1"))
	g.Expect(group[1].Name()).To(Equal("a2"))
	g.Expect(group[2].Name()).To(Equal("b1"))
	g.Expect(group[3].Name()).To(Equal("b2"))
	g.Expect(group[4].Name()).To(Equal("x1"))
	g.Expect(group[5].Name()).To(Equal("x2"))
}

func TestSortedBySize(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	a1 := fileInfo{name: "a1", size: 11}
	a2 := fileInfo{name: "a2", size: 33}
	b1 := fileInfo{name: "b1", size: 44}
	b2 := fileInfo{name: "b2", size: 22}
	x1 := fileInfo{name: "x1"}
	fs = &osStub{[]fileInfo{a1, x1, b1, a2, b2}}
	group := New("/a1", "/x1", "/b1", "/a2", "/b2")

	// When...
	group.SortedBySize()

	// Then...
	g.Expect(group[0].Name()).To(Equal("x1"))
	g.Expect(group[1].Name()).To(Equal("a1"))
	g.Expect(group[2].Name()).To(Equal("b2"))
	g.Expect(group[3].Name()).To(Equal("a2"))
	g.Expect(group[4].Name()).To(Equal("b1"))
}

func TestFirstAndLast(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	a := fileInfo{name: "a", size: 11}
	b := fileInfo{name: "b", size: 44}
	c := fileInfo{name: "c", size: 33}
	d := fileInfo{name: "d", size: 22}
	fs = &osStub{[]fileInfo{a, b, c, d}}
	group := New("/a", "/b", "/c", "/d")

	// When...
	first := group.First()
	last := group.Last()

	// Then...
	g.Expect(first.Name()).To(Equal("a"))
	g.Expect(last.Name()).To(Equal("d"))
}

func TestEmptyFirstAndLast(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	group := Of()

	// When...
	first := group.First()
	last := group.Last()

	// Then...
	g.Expect(first).To(BeNil())
	g.Expect(last).To(BeNil())
}

func TestCompare(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	now := time.Now().UTC()
	empty := Of()
	a1 := fileInfo{name: "a1", size: 11, modTime: now.Add(-1 * time.Minute)}
	a2 := fileInfo{name: "a2", size: 22, modTime: now.Add(-2 * time.Minute)}
	b1 := fileInfo{name: "b1", size: 11, modTime: now.Add(-2)}
	b2 := fileInfo{name: "b2", size: 22, modTime: now.Add(-1)}
	//x1 := fileInfo{name: "x1"}
	//x2 := fileInfo{name: "x2"}
	fs = &osStub{[]fileInfo{a1, a2, b1, b2, a1, b2}}
	a1a2 := New("/a1", "/a2")
	b1b2 := New("/b1", "/b2")
	a1b2 := New("/a1", "/b2")

	// Then...
	g.Expect(a1a2.AllAreOlderThan(b1b2)).To(BeTrue())
	g.Expect(b1b2.AllAreOlderThan(a1a2)).To(BeFalse())

	g.Expect(b1b2.AllAreNewerThan(a1a2)).To(BeTrue())
	g.Expect(a1a2.AllAreNewerThan(b1b2)).To(BeFalse())

	g.Expect(b1b2.OverlapsWith(a1b2)).To(BeTrue())
	g.Expect(a1b2.OverlapsWith(b1b2)).To(BeTrue())

	g.Expect(a1a2.compare(empty)).To(Equal(undefined))
	g.Expect(empty.compare(b1b2)).To(Equal(undefined))
}

func TestErrors(t *testing.T) {
	g := NewGomegaWithT(t)
	// Given...
	a1 := fileInfo{err: errors.New("a1")}
	a2 := fileInfo{err: errors.New("a2")}
	fs = &osStub{[]fileInfo{a1, a2}} // global
	g1 := New("/a1", "/a2")

	// When...
	s := g1.Errors().Error()

	// Then...
	g.Expect(s).To(Equal("a1\na2"))
}
