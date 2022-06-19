package template

import (
	"testing"

	_ "github.com/footprintai/multikf/kfmanifests"
	_ "github.com/stretchr/testify/assert"
)

type statisKfConfig struct{}

func (s statisKfConfig) GetDefaultPassword() string {
	return "12341234"
}

func TestKubeflowTemplate(t *testing.T) {
	//kt := NewKubeflowTemplate()
	//assert.NoError(t, kt.Populate(statisKfConfig{}))
	//buf := &bytes.Buffer{}
	//assert.NoError(t, kt.Execute(buf))
	//assert.EqualValues(t, kfmanifests.KF14, buf.String())
}
