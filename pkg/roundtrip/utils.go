package roundtrip

import "fmt"

type dump struct {
	data []byte
}

func newDump() *dump {
	return &dump{
		data: make([]byte, 0),
	}
}

func (d *dump) Write(p []byte) (n int, err error) {
	d.data = append(d.data, p...)
	return len(p), nil
}

type load struct {
	data []byte
}

func newLoad(data []byte) *load {
	return &load{
		data: data,
	}
}

func (l *load) Read(p []byte) (n int, err error) {
	if len(l.data) == 0 {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, l.data)
	l.data = l.data[n:]
	return n, nil
}
