package devices

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aculclasure/aoc2022/ds"
)

// fileInfoRgx defines what a file information line looks like.
var fileInfoRgx = regexp.MustCompile(`^(\d+) .*$`)

// Directory represents a directory in the filesystem of the elf's device.
type Directory struct {
	Name     string
	Files    []File
	Children map[string]*Directory
	Parent   *Directory
}

// File represents a file in the filesystem of the elf's device.
type File struct {
	Name string
	Size int
}

// TotalSize returns the sum of the size of this directory and the sizes of
// all its subdirectories.
func (d *Directory) TotalSize() int {
	var stk ds.Stack[*Directory]
	stk.Push(d)
	sum := 0
	for stk.Size() > 0 {
		next, _ := stk.Pop()
		for _, c := range next.Children {
			stk.Push(c)
		}
		for _, f := range next.Files {
			sum += f.Size
		}
	}
	return sum
}

// AddSubdir accepts a Directory and adds it as a subdirectory to the current
// directory if it is not already a subdirectory of the current directory.
func (d *Directory) AddSubdir(subdir *Directory) {
	if d.Children == nil {
		d.Children = make(map[string]*Directory)
	}

	_, ok := d.Children[subdir.Name]
	if !ok {
		subdir.Parent = d
		d.Children[subdir.Name] = subdir
	}
}

// AllDescendans returns a slice of all descendants of the current directory. A
// nil slice is returned if the current directory has no descendants.
func (d *Directory) AllDescendants() []*Directory {
	if d == nil {
		return nil
	}

	var (
		stk         ds.Stack[*Directory]
		descendants []*Directory
	)
	for _, v := range d.Children {
		stk.Push(v)
	}
	for stk.Size() > 0 {
		nextChild, _ := stk.Pop()
		descendants = append(descendants, nextChild)
		for _, v := range nextChild.Children {
			stk.Push(v)
		}
	}

	return descendants
}

// BestDirectoryToCleanup returns the smallest directory that could be removed
// in order to provide the minimum needed disk space for running a system update.
// Nil is returned if the current directory is nil or no descendant directory
// meeting the cleanup criteria can be found.
func (d *Directory) BestDirectoryToCleanup(minSystemFreeSpace int) *Directory {
	if d == nil {
		return nil
	}
	const totalAvailableSystemSpace = 70000000
	currentUsedSpace := d.TotalSize()
	allDirs := []*Directory{d}
	allDirs = append(allDirs, d.AllDescendants()...)
	var potential []*Directory
	for _, dir := range allDirs {
		freedSpace := (totalAvailableSystemSpace - currentUsedSpace) + dir.TotalSize()
		if freedSpace >= minSystemFreeSpace {
			potential = append(potential, dir)
		}
	}
	sort.Slice(potential, func(i, j int) bool {
		return potential[i].TotalSize() < potential[j].TotalSize()
	})
	if len(potential) == 0 {
		return nil
	}
	return potential[0]
}

// TreeFromTerminalOutput accepts an io.Reader that points to output captured
// from the terminal on the elf's device, builds a directory tree from the
// terminal data and returns the root directory of the tree. An error is returned
// if the terminal argument is nil or if a problem occurs when analyzing the
// terminal output.
func TreeFromTerminalOutput(terminal io.Reader) (*Directory, error) {
	if terminal == nil {
		return nil, errors.New("terminal must be non-nil")
	}

	var (
		stk  ds.Stack[*Directory]
		line string
	)
	rootDir := &Directory{Name: "/"}
	stk.Push(rootDir)
	scn := bufio.NewScanner(terminal)
	for scn.Scan() {
		line = scn.Text()
		switch {
		case line == "$ cd /":
			for stk.Size() > 1 {
				stk.Pop()
			}
		case line == "$ cd ..":
			if stk.Size() <= 1 {
				continue
			}
			stk.Pop()
		case strings.HasPrefix(line, "$ cd "):
			cwd, ok := stk.Peek()
			if !ok {
				continue
			}
			dir, err := DirFromLine(line)
			if err != nil {
				log.Print(err)
				continue
			}
			if cwd.Children == nil {
				dir.Parent = cwd
				cwd.Children = map[string]*Directory{dir.Name: dir}
				stk.Push(dir)
				continue
			}
			d, ok := cwd.Children[dir.Name]
			if !ok {
				dir.Parent = cwd
				cwd.Children[dir.Name] = dir
				stk.Push(dir)
				continue
			}
			stk.Push(d)
		case fileInfoRgx.MatchString(line):
			cwd, ok := stk.Peek()
			if !ok {
				continue
			}
			f, err := FileFromLine(line)
			if err != nil {
				log.Print(err)
				continue
			}
			cwd.Files = append(cwd.Files, f)
		case strings.HasPrefix(line, "dir"):
			cwd, ok := stk.Peek()
			if !ok {
				continue
			}
			dir, err := DirFromLine(line)
			if err != nil {
				log.Print(err)
				continue
			}
			if cwd.Children == nil {
				dir.Parent = cwd
				cwd.Children = map[string]*Directory{dir.Name: dir}
				continue
			}
			_, ok = cwd.Children[dir.Name]
			if !ok {
				dir.Parent = cwd
				cwd.Children[dir.Name] = dir
			}
		}

	}
	err := scn.Err()
	if err != nil {
		return nil, err
	}

	return rootDir, nil
}

// DirectoriesSmallerThan accepts a root directory and an int representing
// the maximum total size and returns a slice of all directories starting from
// the root directory whose total size is no larger than the maximum total
// size. A nil slice is returned if the root directory is nil or if all directories
// found are larger than the maximum total size.
func DirectoriesSmallerThan(root *Directory, maxTotalSize int) []*Directory {
	if root == nil {
		return nil
	}
	var matches []*Directory
	q := ds.NewQueue[*Directory]()
	q.Enqueue(root)
	for q.Size() > 0 {
		next, _ := q.Dequeue()
		if next.TotalSize() > maxTotalSize {
			for _, c := range next.Children {
				q.Enqueue(c)
			}
			continue
		}
		matches = append(matches, next)
		desc := next.AllDescendants()
		matches = append(matches, desc...)
	}
	return matches
}

// DirFromLine accepts a line from the device's terminal output that is expected
// to be in the form:
// "$ cd <dirname>" or "dir <dirname>"
// and returns a Directory struct. An error is returned if the line cannot be
// parsed into a Directory struct.
func DirFromLine(line string) (*Directory, error) {
	if !strings.HasPrefix(line, "dir") && !strings.HasPrefix(line, "$ cd ") {
		return nil, fmt.Errorf(`line must be in the form "$ cd <dirname>" or "dir <dirname>" (got line %s)`, line)
	}

	var fields []string
	if strings.HasPrefix(line, "dir") {
		fields = strings.Fields(line)
		if len(fields) < 2 || fields[0] != "dir" {
			return nil, fmt.Errorf(`line must be in the form "dir <dirname>" (got line %s)`, line)
		}
		return &Directory{Name: fields[1]}, nil
	}

	fields = strings.Fields(line)
	if len(fields) < 3 || fields[2] == "/" || fields[2] == ".." {
		return nil, fmt.Errorf(`line must be in the form "$ cd <dirname>" (got line %s)`, line)
	}
	return &Directory{Name: fields[2]}, nil
}

// FileFromLine accepts a line from the device's terminal output that is expected
// to be in the form "<filesize> <filename>" where filesize is an integer representing
// the size of the file and filename is a string. It parses this data and returns
// a File struct. An error is returned if the line cannot be parsed into a File
// struct.
func FileFromLine(line string) (File, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return File{}, fmt.Errorf(`line must be in the form <filesize> <filename> (got line %s)`, line)
	}
	sz, err := strconv.Atoi(fields[0])
	if err != nil {
		return File{}, err
	}

	return File{Name: fields[1], Size: sz}, nil
}
