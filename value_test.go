package argsx

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestString(t *testing.T) {
	os.Args = append(os.Args,
		"--string.value", "string value",
		"--string.must",
		"--string.slice", "A,B,C,D",
		"--string.slice.delimiter", "E-F-G-H",
		"--string.slice.empty",
	)

	value, err := Fetch("string.value").String()
	require.NoError(t, err)
	require.Equal(t, "string value", value)

	value = Fetch("string.must").MustString()
	require.Equal(t, "", value)

	slice, err := Fetch("string.slice").StringSlice()
	require.NoError(t, err)
	require.Equal(t, []string{"A", "B", "C", "D"}, slice)

	slice, err = Fetch("string.slice.delimiter").StringSlice(WithDelimiter[string]("-"))
	require.NoError(t, err)
	require.Equal(t, []string{"E", "F", "G", "H"}, slice)

	slice, err = Fetch("string.slice.empty").StringSlice()
	require.NotNil(t, err)
	require.Equal(t, 0, len(slice))

	slice = Fetch("string.slice.empty").MustStringSlice()
	require.Equal(t, 0, len(slice))

	slice, err = Fetch("string.slice.default").StringSlice(WithDefault[string]("Z", "Y"))
	require.NoError(t, err)
	require.Equal(t, []string{"Z", "Y"}, slice)
}

func TestValue(t *testing.T) {
	os.Args = append(os.Args,
		"--int.value", "123",
		"--int.empty",
		"--int.default",
		"--int.equals=987",
		"--int.must.value", "12345",
		"--int.must.empty",
		"--int.must.default",
	)

	// int of value
	value, err := Fetch("int.value").Int()
	require.NoError(t, err)
	require.Equal(t, 123, value)

	// int no value returns error
	value, err = Fetch("int.empty").Int()
	require.NotNil(t, err)
	require.Equal(t, 0, value)

	// int of default value
	value, err = Fetch("int.default").Int(1234)
	require.NoError(t, err)
	require.Equal(t, 1234, value)

	value, err = Fetch("int.equals").Int()
	require.NoError(t, err)
	require.Equal(t, 987, value)

	// must int of value
	value = Fetch("int.must.value").MustInt()
	require.Equal(t, 12345, value)

	// mustInt no value 0
	value = Fetch("int.must.empty").MustInt()
	require.Equal(t, 0, value)

	// mustInt of default value
	value = Fetch("int.must.default").MustInt(123456)
	require.Equal(t, 123456, value)
}
