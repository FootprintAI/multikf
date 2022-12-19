package kfmanifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifests(t *testing.T) {
	versions := ListVersions()
	assert.EqualValues(t, []string{"v1.4.1", "v1.5.1", "v1.6.0", "v1.6.1-lite"}, versions)

	for _, version := range versions {
		foundVersionManifest, err := GetVersion(version)
		assert.NoError(t, err)
		assert.NotNil(t, foundVersionManifest)
	}
}
