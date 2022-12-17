package devices_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/devices"
	"github.com/google/go-cmp/cmp"
)

// Provides custom comparing logic when comparing 2 Directory structs to each
// other. For our test cases, it's sufficient to say 2 directories are equal
// when they have the same name. When using this comparer, other fields in the
// Directory struct are ignored.
var compareByDirName = cmp.Comparer(func(d1, d2 *devices.Directory) bool {
	if d1 == nil && d2 == nil {
		return true
	}
	if d1 == nil || d2 == nil {
		return false
	}
	return d1.Name == d2.Name
})

func TestDirectory_TotalSizeGivenNestedDirectoriesReturnsExpectedSum(t *testing.T) {
	t.Parallel()
	b := &devices.Directory{
		Name: "b",
		Files: []devices.File{
			{Name: "file4", Size: 4},
			{Name: "file5", Size: 5},
		},
	}
	c := &devices.Directory{
		Name: "c",
		Files: []devices.File{
			{Name: "file6", Size: 6},
			{Name: "file7", Size: 7},
		},
	}
	a := &devices.Directory{
		Name: "a",
		Files: []devices.File{
			{Name: "file2", Size: 2},
			{Name: "file3", Size: 3},
		},
		Children: map[string]*devices.Directory{
			"b": b,
			"c": c,
		},
	}
	want := 27
	got := a.TotalSize()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestDirectory_AddSubdirGivenNonExistentSubdirAddsSubdir(t *testing.T) {
	t.Parallel()
	parentDir := &devices.Directory{Name: "a"}
	subDir := &devices.Directory{Name: "b"}
	parentDir.AddSubdir(subDir)
	_, ok := parentDir.Children["b"]
	if !ok {
		t.Error(`want parent to have child "b" but it did not`)
	}
}

func TestDirectory_AddSubdirGivenExistingSubdirDoesNotOverrideExistingSubdir(t *testing.T) {
	t.Parallel()
	parentDir := &devices.Directory{Name: "a"}
	subDir := &devices.Directory{Name: "b", Files: []devices.File{{Name: "file0.txt", Size: 10}}}
	parentDir.AddSubdir(subDir)
	newSubDir := &devices.Directory{Name: "b", Files: []devices.File{{Name: "file1.txt", Size: 20}}}
	parentDir.AddSubdir(newSubDir)
	got := parentDir.Children["b"]
	if subDir != got {
		t.Errorf("got unexpected subdirectory %+v", subDir)
	}
}

func TestDirectory_AllDescendantLeaves(t *testing.T) {
	t.Parallel()
	rootDir := &devices.Directory{Name: "/"}
	rootDir.AddSubdir(&devices.Directory{Name: "a"})
	b := &devices.Directory{Name: "b"}
	b.AddSubdir(&devices.Directory{Name: "c"})
	rootDir.AddSubdir(b)
	want := []*devices.Directory{
		{Name: "a"},
		{Name: "c"},
	}
	got := rootDir.AllDescendantLeaves()
	if got == nil {
		t.Fatalf("want %+v, got nil", want)
	}
	sort.Slice(got, func(i, j int) bool {
		return got[i].Name < got[j].Name
	})
	if !cmp.Equal(want, got, compareByDirName) {
		t.Error(cmp.Diff(want, got, compareByDirName))
	}
}

func TestDirectory_BestDirectoryToCleanup(t *testing.T) {
	t.Parallel()
	rootDir := buildTreeFromExample()
	want := &devices.Directory{Name: "d"}
	const minSystemFreeSpace = 30000000
	got := rootDir.BestDirectoryToCleanup(minSystemFreeSpace)
	if !cmp.Equal(want, got, compareByDirName) {
		t.Error(cmp.Diff(want, got, compareByDirName))
	}
}

func TestAllDescendants(t *testing.T) {
	t.Parallel()
	rootDir := &devices.Directory{Name: "/"}
	child1 := &devices.Directory{Name: "a"}
	child2 := &devices.Directory{Name: "b"}
	rootDir.AddSubdir(child1)
	rootDir.AddSubdir(child2)
	grandChild := &devices.Directory{Name: "c"}
	child1.AddSubdir(grandChild)
	grandChild.AddSubdir(&devices.Directory{Name: "d"})

	testCases := map[string]struct {
		input *devices.Directory
		want  []*devices.Directory
	}{
		"Nil root directory returns nil": {
			input: nil,
			want:  nil,
		},
		"Root directory with no children returns nil": {
			input: child2,
			want:  nil,
		},
		"Root directory with one descendant returns descendant": {
			input: grandChild,
			want:  []*devices.Directory{{Name: "d"}},
		},
		"Root directory with multiple generations of descendants returns expected descendants": {
			input: rootDir,
			want: []*devices.Directory{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
				{Name: "d"},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := devices.AllDescendants(tc.input)
			sort.Slice(got, func(i, j int) bool {
				return got[i].Name < got[j].Name
			})
			if !cmp.Equal(tc.want, got, compareByDirName) {
				t.Error(cmp.Diff(tc.want, got, compareByDirName))
			}
		})
	}
}

func TestDirectoriesSmallerThan(t *testing.T) {
	t.Parallel()
	rootDir := buildTreeFromExample()
	want := []*devices.Directory{
		{Name: "a"},
		{Name: "e"},
	}
	got := devices.DirectoriesSmallerThan(rootDir, 100000)
	sort.Slice(got, func(i, j int) bool {
		return got[i].Name < got[j].Name
	})
	if !cmp.Equal(want, got, compareByDirName) {
		t.Error(cmp.Diff(want, got, compareByDirName))
	}
}

func TestTreeFromTerminalOutput(t *testing.T) {
	t.Parallel()
	output := strings.NewReader(`$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
`)
	got, err := devices.TreeFromTerminalOutput(output)
	if err != nil {
		t.Fatal(err)
	}

	wantName := "/"
	if wantName != got.Name {
		t.Errorf("want directory name %s, got %s", wantName, got.Name)
	}

	wantChildDirNames := []string{"a", "d"}
	var gotChildDirNames []string
	for k := range got.Children {
		gotChildDirNames = append(gotChildDirNames, k)
	}
	sort.Slice(gotChildDirNames, func(i, j int) bool {
		return gotChildDirNames[i] < gotChildDirNames[j]
	})
	if !cmp.Equal(wantChildDirNames, gotChildDirNames) {
		t.Error(cmp.Diff(wantChildDirNames, gotChildDirNames))
	}

	wantFiles := []devices.File{
		{Name: "b.txt", Size: 14848514},
		{Name: "c.dat", Size: 8504156},
	}
	gotFiles := got.Files
	if !cmp.Equal(wantFiles, gotFiles) {
		t.Error(cmp.Diff(wantFiles, gotFiles))
	}
}

func TestDirFromLine(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input string
		want  *devices.Directory
	}{
		"Valid change directory line returns expected Directory": {
			input: "$ cd a",
			want:  &devices.Directory{Name: "a"},
		},
		"Valid directory listing line returns expected Directory": {
			input: "dir a",
			want:  &devices.Directory{Name: "a"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := devices.DirFromLine(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestDirFromLineErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Line that has nothing to do with a directory returns error":        "1234 a.txt",
		"Change directory line with wrong number of fields returns error":   "$ cd ",
		"Change directory line with directory backout target returns error": "$ cd ..",
		"Change directory line with root directory target returns error":    "$ cd /",
		"Directory listing line with wrong number of fields returns error":  "dir",
		"Directory listing line with invalid dir field returns error":       "dirrrr a",
	}
	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := devices.DirFromLine(input)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

func TestFileFromLineGivenValidLineReturnsFile(t *testing.T) {
	t.Parallel()
	line := "1234 a.txt"
	want := devices.File{Name: "a.txt", Size: 1234}
	got, err := devices.FileFromLine(line)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFileFromLineErrorCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]string{
		"Line with only a single field returns error":      "1234",
		"Line with invalid size field value returns error": "4kb a.txt",
	}
	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := devices.FileFromLine(input)
			if err == nil {
				t.Error("expected an error but did not get one")
			}
		})
	}
}

// buildTreeFromExample is a test helper that builds the directory tree provided
// in the day 7 example. It returns the root directory of the tree.
func buildTreeFromExample() *devices.Directory {
	rootDir := &devices.Directory{
		Name: "/",
		Files: []devices.File{
			{Name: "b.txt", Size: 14848514},
			{Name: "c.dat", Size: 8504156},
		}}
	a := &devices.Directory{
		Name: "a",
		Files: []devices.File{
			{Name: "f", Size: 29116},
			{Name: "g", Size: 2557},
		},
	}
	e := &devices.Directory{
		Name: "e",
		Files: []devices.File{
			{Name: "i", Size: 584},
		},
	}
	a.AddSubdir(e)
	d := &devices.Directory{
		Name: "d",
		Files: []devices.File{
			{Name: "j", Size: 4060174},
			{Name: "d.log", Size: 8033020},
			{Name: "d.ext", Size: 5626152},
			{Name: "k", Size: 7214296},
		},
	}
	rootDir.AddSubdir(a)
	rootDir.AddSubdir(d)
	return rootDir
}
