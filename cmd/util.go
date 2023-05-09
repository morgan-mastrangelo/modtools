package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v2"
)

const frozendepsFilename = "modtools_frozen.yml"

type Update struct {
	Path, Version string
}

type Module struct {
	Path     string
	Version  string
	Update   Update
	Indirect bool
}

type Exception struct {
	Path       string
	MinVersion string    `yaml:"minVersion"`
	ValidUntil time.Time `yaml:"validUntil"`
}

type Exceptions struct {
	filename   string
	list       []*Exception
	m          map[string]*Exception
	needSaving bool
	isNew      bool
}

func runCommand(name string, arg ...string) ([]byte, error) {
	var buf bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.Bytes(), err
}

func (e *Exceptions) Load() error {
	f, err := ioutil.ReadFile(e.filename)
	if err != nil {
		if os.IsNotExist(err) {
			e.isNew = true
			return nil
		}
		return err
	}
	err = yaml.Unmarshal(f, &e.list)
	if err == nil {
		e.filterExpired()
		e.m = make(map[string]*Exception)
		for _, item := range e.list {
			e.m[item.Path] = item
		}
	}

	return err
}

func (e *Exceptions) IsNew() bool {
	return e.isNew
}

func loadExceptions() (*Exceptions, error) {
	e := &Exceptions{
		filename: frozendepsFilename,
	}
	err := e.Load()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Exceptions) filterExpired() {
	now := time.Now()
	j := 0
	for i, item := range e.list {
		if item.ValidUntil.IsZero() || item.ValidUntil.Before(now) {
			continue
		}
		if i != j {
			e.list[j] = e.list[i]
		}
		j++
	}
	if j < len(e.list) {
		e.list = e.list[:j]
		e.needSaving = true
	} else {
		e.needSaving = false
	}
}

func (e *Exceptions) filterDuplicates() bool {
	j := 0
	for i, item := range e.list {
		if e.m[item.Path] != item {
			continue
		}
		if i != j {
			e.list[j] = e.list[i]
		}
		j++
	}
	if j < len(e.list) {
		e.list = e.list[:j]
		return true
	}
	return false
}

func (e *Exceptions) Save() error {
	if !e.filterDuplicates() && !e.needSaving {
		return nil
	}
	data, err := yaml.Marshal(e.list)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(e.filename, data, 0o644)
	if err == nil {
		e.needSaving = false
		e.isNew = false
	}
	return err
}

func readDeps(updates, directOnly bool) ([]Module, error) {
	format := "{{with .Module}}{{.Path}}{{end}}"
	if directOnly {
		format = "{{with .Module}}{{if not .Indirect}}{{.Path}}{{end}}{{end}}"
	}
	out, err := runCommand("go", "list", "-f", format, "all")
	if err != nil {
		return nil, err
	}
	args := []string{"list", "-m", "-json"}
	if updates {
		args = append(args, "-u", "-mod=readonly")
	}
	set := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		modpath := scanner.Text()
		if _, exists := set[modpath]; !exists {
			args = append(args, scanner.Text())
			set[modpath] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	out, err = runCommand("go", args...)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(out))
	var list []Module
	for {
		list = append(list, Module{})
		item := &list[len(list)-1]
		err := dec.Decode(item)
		if err != nil {
			if err == io.EOF {
				list = list[:len(list)-1]
				break
			}
			return nil, err
		}
	}
	return list, nil
}

func (e *Exceptions) Get(p string) *Exception {
	return e.m[p]
}

func (e *Exceptions) Add(ex *Exception) {
	e.list = append(e.list, ex)
	if e.m == nil {
		e.m = make(map[string]*Exception)
	}
	e.m[ex.Path] = ex
	e.needSaving = true
}

func (e *Exceptions) Remove(p string) {
	delete(e.m, p)
	e.needSaving = true
}
