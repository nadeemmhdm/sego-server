package storage

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Root struct{ ID, Name, Path string }
type Manager struct{ roots map[string]Root }

type Entry struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"is_dir"`
	Modified time.Time `json:"modified"`
}

func New(roots []Root) (*Manager, error) {
	m := &Manager{roots: map[string]Root{}}
	for _, r := range roots {
		abs, err := filepath.Abs(r.Path)
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll(abs, 0o750); err != nil {
			return nil, err
		}
		r.Path = filepath.Clean(abs)
		m.roots[r.ID] = r
	}
	return m, nil
}

func (m *Manager) Roots() []Root {
	out := make([]Root, 0, len(m.roots))
	for _, r := range m.roots {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (m *Manager) Resolve(rootID, rel string, mustExist bool) (string, error) {
	r, ok := m.roots[rootID]
	if !ok {
		return "", errors.New("unknown storage root")
	}
	rel = filepath.FromSlash(strings.TrimPrefix(rel, "/"))
	clean := filepath.Clean(rel)
	if clean == "." {
		clean = ""
	}
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", errors.New("path escapes storage root")
	}
	full := filepath.Join(r.Path, clean)
	check, err := filepath.Rel(r.Path, full)
	if err != nil || check == ".." || strings.HasPrefix(check, ".."+string(os.PathSeparator)) {
		return "", errors.New("path escapes storage root")
	}
	if err := m.rejectSymlinkChain(r.Path, full, mustExist); err != nil {
		return "", err
	}
	return full, nil
}

func (m *Manager) rejectSymlinkChain(root, full string, mustExist bool) error {
	rel, _ := filepath.Rel(root, full)
	cur := root
	for _, part := range strings.Split(rel, string(os.PathSeparator)) {
		if part == "." || part == "" {
			continue
		}
		cur = filepath.Join(cur, part)
		fi, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) && !mustExist {
				return nil
			}
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return errors.New("symbolic links are not allowed in protected paths")
		}
	}
	return nil
}

func (m *Manager) List(rootID, rel string) ([]Entry, error) {
	p, err := m.Resolve(rootID, rel, true)
	if err != nil {
		return nil, err
	}
	des, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}
	out := make([]Entry, 0, len(des))
	base := strings.Trim(strings.ReplaceAll(rel, "\\", "/"), "/")
	for _, de := range des {
		if de.Name() == ".sego-trash" {
			continue
		}
		info, err := de.Info()
		if err != nil {
			continue
		}
		rp := de.Name()
		if base != "" {
			rp = base + "/" + rp
		}
		out = append(out, Entry{Name: de.Name(), Path: rp, Size: info.Size(), IsDir: de.IsDir(), Modified: info.ModTime()})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Manager) Mkdir(rootID, rel string) error {
	p, e := m.Resolve(rootID, rel, false)
	if e != nil {
		return e
	}
	return os.MkdirAll(p, 0o750)
}

func (m *Manager) StreamWrite(rootID, rel string, src io.Reader, maxBytes int64) (int64, error) {
	p, e := m.Resolve(rootID, rel, false)
	if e != nil {
		return 0, e
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
		return 0, err
	}
	tmp := p + ".sego-upload"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o640)
	if err != nil {
		return 0, err
	}
	defer func() { f.Close() }()
	n, err := io.Copy(f, io.LimitReader(src, maxBytes+1))
	if err != nil {
		os.Remove(tmp)
		return n, err
	}
	if n > maxBytes {
		os.Remove(tmp)
		return n, errors.New("upload exceeds configured maximum")
	}
	if err := f.Sync(); err != nil {
		os.Remove(tmp)
		return n, err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return n, err
	}
	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return n, err
	}
	return n, nil
}

func (m *Manager) Open(rootID, rel string) (*os.File, fs.FileInfo, error) {
	p, e := m.Resolve(rootID, rel, true)
	if e != nil {
		return nil, nil, e
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, nil, e
	}
	info, e := f.Stat()
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	if info.IsDir() {
		f.Close()
		return nil, nil, errors.New("cannot download a directory")
	}
	return f, info, nil
}

func (m *Manager) Rename(rootID, oldRel, newRel string) error {
	a, e := m.Resolve(rootID, oldRel, true)
	if e != nil {
		return e
	}
	b, e := m.Resolve(rootID, newRel, false)
	if e != nil {
		return e
	}
	if err := os.MkdirAll(filepath.Dir(b), 0o750); err != nil {
		return err
	}
	return os.Rename(a, b)
}

func (m *Manager) Trash(rootID, rel string) error {
	p, e := m.Resolve(rootID, rel, true)
	if e != nil {
		return e
	}
	r := m.roots[rootID]
	trash := filepath.Join(r.Path, ".sego-trash")
	if err := os.MkdirAll(trash, 0o700); err != nil {
		return err
	}
	name := time.Now().UTC().Format("20060102T150405.000000000Z") + "-" + filepath.Base(p)
	return os.Rename(p, filepath.Join(trash, name))
}
