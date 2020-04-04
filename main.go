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
	License     string
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

	longestName := 0
	byLicense := map[string][]*Dependency{}

	for _, dep := range pkg.Dependencies {
		if len(dep.Name) > longestName {
			longestName = len(dep.Name)
		}

		if _, ok := byLicense[dep.License]; !ok {
			byLicense[dep.License] = make([]*Dependency, 0)
		}

		byLicense[dep.License] = append(byLicense[dep.License], dep)
	}

	for license, deps := range byLicense {
		fmt.Printf("  %s\n", license)

		for _, dep := range deps {
			fmt.Printf("    * %-*s - %s\n", longestName, dep.Name, dep.Description)
		}

		fmt.Print("\n")
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

	license := "N/A"
	if r.License != nil {
		license = r.License.GetName()
	}

	return &Dependency{
		Name:        *r.Name,
		Description: *r.Description,
		URL:         *r.HTMLURL,
		License:     license,
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
