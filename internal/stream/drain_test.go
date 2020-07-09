package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	setupClient *Drain
	testClient  *Drain
)

func setupDrainTests() {
	setupClient = NewDrain()
	testClient = NewDrain()
	setupClient.SetInput(testClient.GetOutput())
	testClient.SetInput(setupClient.GetOutput())
}

var (
	SendRead_NormalCases = []struct {
		data interface{}
	}{
		{"data1"},
		{-3.03},
		{1},
		{true},
		{false},
		{[]byte("data_bytes")},
		{struct{ val string }{val: "string"}},
	}
)

func TestSendRead_NormalCases(t *testing.T) {
	setupDrainTests()

	// setup
	for _, c := range SendRead_NormalCases {
		setupClient.Send(c.data)
	}

	for _, c := range SendRead_NormalCases {
		// function under test
		data, msgReceived := testClient.Read()

		// assert data matches and msgReceived bool is true
		assert.Equal(t, c.data, data)
		assert.True(t, msgReceived)
	}
}

func TestRead_NoDataCases(t *testing.T) {
	// function under test
	data, msgReceived := testClient.Read()

	// assert data is nil and msgReceived bool is false
	assert.Nil(t, data)
	assert.False(t, msgReceived)
}
