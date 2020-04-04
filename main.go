package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/github"
	"github.com/sirkon/goproxy/gomod"
)

type Dependency struct {
	Name        string
	Description string
	URL         string
}

type Package struct {
	Name         string
	Dependencies []*Dependency
}

var dataSource = map[string]func(u, r string) (*Dependency, error){
	"github.com": getFromGithub,
}

var ghClient = github.NewClient(nil)

func main() {
	m, err := parseGoMod("go.mod")
	if err != nil {
		panic(err)
	}

	pkg := Package{
		Name: m.Name,
	}

	for dep := range m.Require {
		repoParts := strings.Split(dep, "/")
		source, user, path := repoParts[0], repoParts[1], strings.Join(repoParts[2:], "/")

		if f, ok := dataSource[source]; ok {
			d, err := f(user, path)
			if err != nil {
				continue
			}

			pkg.Dependencies = append(pkg.Dependencies, d)
		}
	}

	fmt.Printf("Dependencies for %s\n\n", pkg.Name)

	for _, dep := range pkg.Dependencies {
		fmt.Printf("  %s - %s\n", dep.Name, dep.Description)
	}
}

func getFromGithub(userOrOrg, repo string) (*Dependency, error) {
	r, _, err := ghClient.Repositories.Get(context.Background(), userOrOrg, repo)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, errors.New("no error - no response")
	}

	return &Dependency{
		Name:        *r.Name,
		Description: *r.Description,
		URL:         *r.HTMLURL,
	}, nil
}

func parseGoMod(f string) (*gomod.Module, error) {
	fileData, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	modFile, err := gomod.Parse(f, fileData)
	if err != nil {
		return nil, err
	}

	return modFile, nil
}
