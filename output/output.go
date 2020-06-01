package output

import "github.com/bsladewski/lapis/stream"

// Output functions as a stream that passes data through while also compiling
// data to be output.
type Output interface {
	stream.Stream
	GetData() (interface{}, error)
}
