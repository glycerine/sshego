package sshego

import (
	"bytes"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
)

func TestShovelStops(t *testing.T) {

	cv.Convey("a Shovel should stop when requested", t, func() {

		s := newShovel(false)

		a := newMockRwc([]byte("hello_from_a"))
		b := newMockRwc([]byte("hello_from_b"))

		s.Start(b, a, "b<-a")
		<-s.Ready
		time.Sleep(100 * time.Millisecond)
		s.Stop()
		cv.So(b.sink.String(), cv.ShouldResemble, "hello_from_a")
		cv.So(a.sink.String(), cv.ShouldResemble, "")
	})

	cv.Convey("a ShovelPair should stop when requested", t, func() {

		s := newShovelPair(false)

		a := newMockRwc([]byte("hello_from_a"))
		b := newMockRwc([]byte("hello_from_b"))

		s.Start(a, b, "a<-b", "b->a")
		<-s.Ready
		time.Sleep(1 * time.Millisecond)
		s.Stop()
		cv.So(b.sink.String(), cv.ShouldResemble, "hello_from_a")
		cv.So(a.sink.String(), cv.ShouldResemble, "hello_from_b")
	})

}

type mockRwc struct {
	src  *bytes.Buffer
	sink *bytes.Buffer
}

func newMockRwc(src []byte) *mockRwc {
	return &mockRwc{
		src:  bytes.NewBuffer(src),
		sink: bytes.NewBuffer(nil),
	}
}

func (m *mockRwc) Read(p []byte) (n int, err error) {
	return m.src.Read(p)
}

func (m *mockRwc) Write(p []byte) (n int, err error) {
	return m.sink.Write(p)
}

func (m *mockRwc) Close() error {
	return nil
}
