package helm3

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	m := NewTestMixin(t)

	err := m.Build()
	require.NoError(t, err)

	buildOutput := `RUN apt-get update && apt-get install -y curl
RUN curl https://get.helm.sh/helm-v3.1.2-linux-amd64.tar.gz --output helm3.tar.gz
RUN tar -xvf helm3.tar.gz
RUN mv linux-amd64/helm /usr/local/bin/helm3`

	t.Run("build with a valid config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-valid-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")

		wantOutput := buildOutput + "\nRUN helm3 repo add stable kubernetes-charts --username username --password password"

		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with invalid config", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-invalid-config.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)

		err = m.Build()
		require.NoError(t, err, "build failed")
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, buildOutput, gotOutput)
	})
}
