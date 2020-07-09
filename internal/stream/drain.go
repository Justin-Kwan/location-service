package stream

const (
	MaxBuffSize = 1024
)

// Drain is a bidirectional sink that receives and sends data.
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
func (d *Drain) Send(data interface{}) bool {
	select {
	case d.output <- data:
		return true
	default:
		return false
	}
}

// Read chceks if a message is received on the input channel. Then,
// the message is returned, with a boolean value indication.
func (d *Drain) Read() (interface{}, bool) {
	select {
	case data := <-d.input:										// if received msg
		return data, true
	default:																	// otherwise no msg
		return nil, false
	}
}

// Close closes the output channel for the drain.
func (d *Drain) Close() {
	close(d.output)
}
