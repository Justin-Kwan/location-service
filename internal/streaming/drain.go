package streaming

const (
	MaxBuffSize = 1024
)

// Drain is a bidirectional sink that reads and sends data.
type Drain struct {
	input  <-chan interface{}
	output chan interface{}
}

func NewDrain() *Drain {
	return &Drain{
		output: make(chan interface{}, MaxBuffSize),
	}
}

// SetInput sets a read-only channel as the input channel for the
// drain to read from.
func (d *Drain) SetInput(in <-chan interface{}) {
	d.input = in
}

// GetOutput returns the output channel that the drain sends data to,
// to set in a companion drain as an input channel.
func (d *Drain) GetOutput() <-chan interface{} {
	return d.output
}

// Send sends a message onto the output channel.
func (d *Drain) Send(data interface{}) {
	d.output <- data
}

// Read blocks and waits for a message to be received on the input
// channel. Then, the message is returned.
func (d *Drain) Read() interface{} {
	data := <-d.input
	return data
}

// Close closes the output channel for the drain.
func (d *Drain) Close() {
	close(d.output)
}
