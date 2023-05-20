package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoGrammar(t *testing.T) {
	t.Run("Successfully parse application version,name,url from source string", func(t *testing.T) {
		app, err := EvalInfo(`@aloe version v1
@aloe name cli
@aloe url https://tfadeyi.github.io`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
	})
	t.Run("Successfully parse application semver version v1.0.0", func(t *testing.T) {
		app, err := EvalInfo(`@aloe version v1.0.0
@aloe name cli
@aloe url https://tfadeyi.github.io`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
	})
	t.Run("Successfully parse application semver version v1.0.0-alpha1", func(t *testing.T) {
		app, err := EvalInfo(`@aloe version v1.0.0-alpha1
@aloe name cli
@aloe url https://tfadeyi.github.io`)
		require.NoError(t, err)
		assert.EqualValues(t, "v1.0.0-alpha1", app.Version)
		assert.EqualValues(t, "cli", app.Name)
		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
	})
	t.Run("Fails to parse application info if the version is missing", func(t *testing.T) {
		_, err := EvalInfo(`@aloe name cli
@aloe url https://tfadeyi.github.io`)
		require.ErrorIs(t, err, ErrMissingRequiredField)
	})
	t.Run("Fails to parse application info if the name is missing", func(t *testing.T) {
		_, err := EvalInfo(`@aloe version v1.0.0-alpha1
@aloe url https://tfadeyi.github.io`)
		require.ErrorIs(t, err, ErrMissingRequiredField)
	})
	t.Run("Fails to parse application info if the url is missing", func(t *testing.T) {
		_, err := EvalInfo(`@aloe version v1.0.0-alpha1
@aloe name cli`)
		require.ErrorIs(t, err, ErrMissingRequiredField)
	})
	t.Run("Fails to parse invalid source string", func(t *testing.T) {
		_, err := EvalInfo(`please stop writing bad code`)
		require.ErrorIs(t, err, ErrParseSource)
	})

}
