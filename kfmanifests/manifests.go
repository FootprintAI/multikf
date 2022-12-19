package kfmanifests

import (
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"io/fs"
)

//go:embed manifests/*
var manifestFs embed.FS

var (
	// looking for the following possible values: v1.4.1-lite, v1.4.1, v1.4
	versionRegexp = regexp.MustCompile(`(v[0-9]\.[0-9](\.[0-9])?(\-lite)?)`)
)

func versionRelFileName(version string) string {
	return fmt.Sprintf("manifests/kubeflow-manifest-%s-template.yaml", version)
}

func VersionBaseFileName(version string) string {
	return filepath.Base(versionRelFileName(version))
}

func ListVersions() []string {
	var versions []string
	fs.WalkDir(manifestFs, ".", func(path string, d fs.DirEntry, err error) error {
		matched := versionRegexp.FindAllString(path, -1)
		//fmt.Printf("path:%s, matched:%s\n", path, matched)
		if len(matched) != 1 {
			// ignore
			return nil
		}
		versions = append(versions, matched[0])
		return nil
	})
	return versions
}

func GetVersion(version string) (string, error) {
	// check the version is one of possible versions
	possbileVersions := ListVersions()
	versionExists := false
	for _, possibleVersion := range possbileVersions {
		if version == possibleVersion {
			versionExists = true
		}
	}
	if !versionExists {
		return "", errors.New("version not found")
	}
	manifestBytes, err := manifestFs.ReadFile(versionRelFileName(version))
	if err != nil {
		return "", err
	}
	return string(manifestBytes), nil
}
