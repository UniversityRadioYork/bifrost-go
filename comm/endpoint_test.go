package comm

import (
	"context"
	"testing"

	"github.com/UniversityRadioYork/bifrost-go/message"
)

// File bifrost/endpoint_test.go contains tests for the Endpoint struct.

// Tests that a pair of endpoints produced by NewEndpointPair connect to each other correctly.
func TestNewEndpointPair_TxRx(t *testing.T) {
	l, r := NewEndpointPair()

	testEndpointTxRx(t, l.Tx, r.Rx)
	testEndpointTxRx(t, r.Tx, l.Rx)
}

// Tests one side of an endpoint pair Tx/Rx connection.
func testEndpointTxRx(t *testing.T, tx chan<- message.Message, rx <-chan message.Message) {
	t.Helper()

	want := message.New("foo", "bar").AddArgs("baz")
	go func() { tx <- *want }()
	got := <-rx

	message.AssertMessagesEqual(t, "tx/rx", &got, want)
}

func TestEndpoint_SendReceive(t *testing.T) {
	l, r := NewEndpointPair()
	ctx, cancel := context.WithCancel(context.Background())

	want := message.New("!", "jam").AddArgs("on", "toast")

	go func() {
		if !l.Send(ctx, *want) {
			t.Error("send failed unexpectedly")
		}
	}()

	got, err := r.Recv(ctx)
	if err != nil {
		t.Fatalf("recv errored: %v", err)
	}
	message.AssertMessagesEqual(t, "send/rx", got, want)

	// After cancelling, sends and receives should fail.
	cancel()

	if l.Send(ctx, *want) {
		t.Error("send succeeded unexpectedly")
	}

	if _, err := r.Recv(ctx); err == nil {
		t.Error("recv succeeded unexpectedly")
	}
}
