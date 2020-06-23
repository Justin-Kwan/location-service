package streaming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	client1 *Drain
	client2 *Drain
)

func setupDrainTests() {
	client1 = NewDrain()
	client2 = NewDrain()
	client1.SetInput(client2.GetOutput())
	client2.SetInput(client1.GetOutput())
}

var (
	SendNormalCases = []struct {
		data1 string
		data2 float64
    data3 int64
		data4 bool
	}{
		{"message1", 0.0001, 1, true},
    {"message2", 2.012222, -2, false},
    {"message3", -3.03, 2, true},
	}
)

func TestClient1SendClient2Read(t *testing.T) {
	setupDrainTests()

  for _, c := range SendNormalCases {
		// function under test
    client1.Send(c.data1)
  	client1.Send(c.data2)
  	client1.Send(c.data3)
    client1.Send(c.data4)
  }

  for _, c := range SendNormalCases {
    // function under test
    data1 := client2.Read().(string)
    data2 := client2.Read().(float64)
    data3 := client2.Read().(int64)
    data4 := client2.Read().(bool)

    assert.Equal(t, c.data1, data1)
    assert.Equal(t, c.data2, data2)
    assert.Equal(t, c.data3, data3)
    assert.Equal(t, c.data4, data4)
  }
}

func TestClient2SendClient1Read(t *testing.T) {
	setupDrainTests()

  // setup
  for _, c := range SendNormalCases {
		// function under test
    client2.Send(c.data1)
  	client2.Send(c.data2)
  	client2.Send(c.data3)
    client2.Send(c.data4)
  }

  for _, c := range SendNormalCases {
    // function under test
    data1 := client1.Read().(string)
    data2 := client1.Read().(float64)
    data3 := client1.Read().(int64)
    data4 := client1.Read().(bool)

    assert.Equal(t, c.data1, data1)
    assert.Equal(t, c.data2, data2)
    assert.Equal(t, c.data3, data3)
    assert.Equal(t, c.data4, data4)
  }
}
