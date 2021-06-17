/*
Copyright 2021 The UnDistro authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type config struct {
	Projects []project `json:"projects,omitempty"`
}

type cmd struct {
	Name string   `json:"name,omitempty"`
	Args []string `json:"args,omitempty"`
}

type project struct {
	Name                   string `json:"name,omitempty"`
	Repo                   string `json:"repo,omitempty"`
	Version                string `json:"version,omitempty"`
	BeforeReleaseCommand   cmd    `json:"beforeReleaseCommand,omitempty"`
	ReleaseCommand         cmd    `json:"releaseCommand,omitempty"`
	AfterReleaseCommand    cmd    `json:"afterReleaseCommand,omitempty"`
	PackageImagesCommand   cmd    `json:"packageImagesCommand,omitempty"`
	PackageBinariesCommand cmd    `json:"packageBinariesCommand,omitempty"`
}

func parseVersion(v string) (*semver.Version, error) {
	vLessV := strings.TrimPrefix(v, "v")
	if _, err := semver.StrictNewVersion(vLessV); err != nil {
		return nil, err
	}
	return semver.NewVersion(v)
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		_, err = parseVersion(info.Name())
		if err != nil {
			logrus.Warnf("ignoring %s is not semver", info.Name())
			return nil
		}
		filesystem := os.DirFS(path)
		f, err := filesystem.Open(filepath.Join(path, "config.yaml"))
		if err != nil {
			return fmt.Errorf("failed to open config: %v", err)
		}
		cfg := config{}
		err = yaml.NewYAMLToJSONDecoder(f).Decode(&cfg)
		if err != nil {
			return fmt.Errorf("failed to read config: %v", err)
		}
		for _, p := range cfg.Projects {
			logrus.Infof("running into project %v\n", p.Name)
			cloneCmd := exec.Command("git", "clone", p.Repo, p.Name)
			cloneCmd.Stdin = os.Stdin
			cloneCmd.Stderr = os.Stderr
			cloneCmd.Stdout = os.Stdout
			err = cloneCmd.Run()
			if err != nil {
				return fmt.Errorf("failed to run git clone: %v", err)
			}
			err = os.Chdir(filepath.Join(path, p.Name))
			if err != nil {
				return fmt.Errorf("failed to cd into project: %v", err)
			}
			tagcmd := exec.Command("git", "checkout", fmt.Sprintf("tags/%v", p.Version))
			tagcmd.Stdin = os.Stdin
			tagcmd.Stderr = os.Stderr
			tagcmd.Stdout = os.Stdout
			err = tagcmd.Run()
			if err != nil {
				return fmt.Errorf("failed to run checkout: %v", err)
			}
			if p.BeforeReleaseCommand.Name != "" {
				bcmd := exec.Command(p.BeforeReleaseCommand.Name, p.BeforeReleaseCommand.Args...)
				bcmd.Stdin = os.Stdin
				bcmd.Stderr = os.Stderr
				bcmd.Stdout = os.Stdout
				err = bcmd.Run()
				if err != nil {
					return fmt.Errorf("failed to run beforerelease: %v", err)
				}
			}
			if p.ReleaseCommand.Name != "" {
				bcmd := exec.Command(p.ReleaseCommand.Name, p.ReleaseCommand.Args...)
				bcmd.Stdin = os.Stdin
				bcmd.Stderr = os.Stderr
				bcmd.Stdout = os.Stdout
				err = bcmd.Run()
				if err != nil {
					return fmt.Errorf("failed to run release: %v", err)
				}
			}
			if p.AfterReleaseCommand.Name != "" {
				bcmd := exec.Command(p.AfterReleaseCommand.Name, p.AfterReleaseCommand.Args...)
				bcmd.Stdin = os.Stdin
				bcmd.Stderr = os.Stderr
				bcmd.Stdout = os.Stdout
				err = bcmd.Run()
				if err != nil {
					return fmt.Errorf("failed to run afterrelease: %v", err)
				}
			}
			if p.PackageBinariesCommand.Name != "" {
				bcmd := exec.Command(p.PackageBinariesCommand.Name, p.PackageBinariesCommand.Args...)
				bcmd.Stdin = os.Stdin
				bcmd.Stderr = os.Stderr
				bcmd.Stdout = os.Stdout
				err = bcmd.Run()
				if err != nil {
					return fmt.Errorf("failed to run packagebinaries: %v", err)
				}
			}
			if p.PackageImagesCommand.Name != "" {
				bcmd := exec.Command(p.PackageImagesCommand.Name, p.PackageImagesCommand.Args...)
				bcmd.Stdin = os.Stdin
				bcmd.Stderr = os.Stderr
				bcmd.Stdout = os.Stdout
				err = bcmd.Run()
				if err != nil {
					return fmt.Errorf("failed to run packageimages: %v", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		logrus.Fatal(err)
	}
}
