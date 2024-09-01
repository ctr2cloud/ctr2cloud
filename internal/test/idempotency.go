package test

import "github.com/stretchr/testify/require"

// RequireIdempotence runs the function f twice and asserts that it reports no change
// the second run
func RequireIdempotence(r *require.Assertions, f func() (bool, error)) {
	updated, err := f()
	r.NoError(err)
	r.True(updated)

	updated, err = f()
	r.NoError(err)
	r.False(updated)
}

// RequireNonIdempotence runs the function f twice and asserts that it reports a change
// the second run
func RequireNonIdempotence(r *require.Assertions, f func() (bool, error)) {
	updated, err := f()
	r.NoError(err)
	r.True(updated)

	updated, err = f()
	r.NoError(err)
	r.True(updated)
}
