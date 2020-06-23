package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"location-service/internal/testutil"
	"location-service/internal/types"
)

var (
	keyDB *KeyDB
)

type TestKey struct {
	Key        string
	Val        string
	UpdatedVal string
}

func getTestKeys() []TestKey {
	return []TestKey{
		TestKey{
			Key: "test_id1",
			Val: `test_field1: test_val1,
					  test_field2: test_val2,
						test_field3: test_val3,`,
			UpdatedVal: "updated_val1",
		},
		TestKey{
			Key: "test_id2",
			Val: `test_field1: ")(*)$%^654",
						test_field2: hello,
						test_field3: {
          	address: '{{integer(100, 999)}} {{street()}}, {{city()}}, {{state()}}, {{integer(100, 10000)}}',
          	about: '{{lorem(1, "paragraphs")}}',
          	registered: '{{date(new Date(2014, 0, 1), new Date(), "YYYY-MM-ddThh:mm:ss Z")}}',
          	latitude: '{{floating(-90.000001, 90)}}',
          	longitude: '{{floating(-180.000001, 180)}}',
          	tags: [
            	'{{repeat(7)}}',
            	'{{lorem(1, "words")}}'
          	],
          	friends: [
            	{
              	'{{repeat(3)}}',
              	Key: '{{index()}}',
              	name: '{{firstName()}} {{surname()}}'
            	}
          	]
        	}`,
			UpdatedVal: "updated_val2",
		},
		TestKey{
			Key:        "test_id3",
			Val:        "test_val3",
			UpdatedVal: "",
		},
	}
}

func setupKeyDBTests() func() {
	cfg := testutil.GetConfig()
	KeyDBPool := NewPool(&(*cfg).RedisKeyDB)

	keyDB = NewKeyDB(KeyDBPool)
	keyDB.Clear()

	return func() {
		keyDB.Clear()
	}
}

func populateKeyDB(t *testing.T) {
	teardownTests := setupKeyDBTests()

	testKeys := getTestKeys()

	for _, tk := range testKeys {
		if err := keyDB.Set(tk.Key, tk.Val); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		val, err := keyDB.Get(tk.Key)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, tk.Val, val)
	}
}

func TestKeyDBSetNormalCases(t *testing.T) {
	teardownTests := setupKeyDBTests()
	defer teardownTests()

	testKeys := getTestKeys()

	for _, tk := range testKeys {
		// function under test
		if err := keyDB.Set(tk.Key, tk.Val); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert key and values in map match the inserted ones
		val, err := keyDB.Get(tk.Key)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, tk.Val, val)
	}
}

var (
	SetIfExistsKeyNotFoundCases = []struct {
		inputKey    string
		inputVal    string
		expectedErr error
	}{
		{"non_existent_member", "random_value1", types.ErrKeyNotFound},
		{" ", "random_value2", types.ErrKeyNotFound},
		{"", "random_value3", types.ErrKeyNotFound},
		{"*", "random_value4", types.ErrKeyNotFound},
	}
)

func TestKeyDBSetIfExistsNormalCases(t *testing.T) {
	teardownTests := setupKeyDBTests()
	defer teardownTests()

	populateKeyDB(t)
	testKeys := getTestKeys()

	for _, tk := range testKeys {
		// function under test
		if err := keyDB.SetIfExists(tk.Key, tk.UpdatedVal); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		val, err := keyDB.Get(tk.Key)
		if err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		assert.Equal(t, tk.UpdatedVal, val)
	}
}

func TestKeyDBSetIfExistsKeyNotFoundCases(t *testing.T) {
	setupKeyDBTests()

	for _, c := range SetIfExistsKeyNotFoundCases {
		// function under test
		err := keyDB.SetIfExists(c.inputKey, c.inputVal)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

var (
	getKeyDBKeyNotFoundCases = []struct {
		inputKey    string
		expectedErr error
	}{
		{"non_existent_key", types.ErrKeyNotFound},
		{" ", types.ErrKeyNotFound},
		{"", types.ErrKeyNotFound},
		{"*", types.ErrKeyNotFound},
	}
)

func TestKeyDBGet(t *testing.T) {
	teardownTests := setupKeyDBTests()
	defer teardownTests()

	for _, c := range getKeyDBKeyNotFoundCases {
		// function under test
		_, err := keyDB.Get(c.inputKey)
		assert.EqualError(t, err, c.expectedErr.Error())
	}
}

func TestKeyDBDelete(t *testing.T) {
	teardownTests := setupKeyDBTests()
	defer teardownTests()

	populateKeyDB(t)
	testKeys := getTestKeys()

	// test
	for _, tk := range testKeys { // THIS DOESN'T WORK
		// function under test
		if err := keyDB.Delete(tk.Key); err != nil {
			teardownTests()
			t.Fatalf(err.Error())
		}

		// assert point of interest is deleted
		_, err := keyDB.Get(tk.Key)
		assert.EqualError(t, err, types.ErrKeyNotFound.Error())
	}
}
