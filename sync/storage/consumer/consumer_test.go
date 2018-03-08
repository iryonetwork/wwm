package consumer

// Tests for consumer of storage sync messages coming from nats-streaming
// Nats-streaming server for test is started in TestMain that runs all the tests on Run() call and then shutdowns the server.

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats-streaming-server/server"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/mock"
	"github.com/iryonetwork/wwm/sync/storage/publisher"
)

var (
	// clusterID for test server.
	clusterID = "TestCluster"
	time1, _  = strfmt.ParseDateTime("2018-02-05T15:18:15.123Z")
	time2, _  = strfmt.ParseDateTime("2018-02-05T15:26:15.123Z")
	file1     = &storageSync.FileInfo{"bucket", "file1", "version", time1}
	file2     = &storageSync.FileInfo{"bucket", "file2", "version", time2}
)

func TestMain(m *testing.M) {
	// Start nats-streaming server
	s, err := server.RunServer(clusterID)
	if err != nil {
		os.Exit(1)
	}

	// Run all the tests
	c := m.Run()

	s.Shutdown()
	os.Exit(c)
}

func TestStartSuccess(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, context.Background(), "Consumer", h)
	defer cleanService()

	// start first consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	num := c.GetNumberOfSubsriptions()
	if num != 1 {
		t.Errorf("Expected number of subscriptions to be 1, got %d", num)
	}

	// start second consumer
	err = c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
	num = c.GetNumberOfSubsriptions()
	if num != 2 {
		t.Errorf("Expected number of subscriptions to be 2, got %d", num)
	}
}

func TestStartFailureInavlidType(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, context.Background(), "Consumer", h)
	defer cleanService()

	// start first consumer
	err := c.StartSubscription("invalid_type")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	num := c.GetNumberOfSubsriptions()
	if num != 0 {
		t.Errorf("Expected number of subscriptions to be 0, got %d", num)
	}
}

func TestStartFailureConnectionClosed(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, context.Background(), "Consumer", h)

	// cleanService closes connection
	cleanService()

	// Start first consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	num := c.GetNumberOfSubsriptions()
	if num != 0 {
		t.Errorf("Expected number of subscriptions to be 0, got %d", num)
	}
}

func TestMessageHandling(t *testing.T) {
	ctx := context.Background()
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, ctx, "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		SyncFile(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultSynced, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			called <- true
		}).
		Times(1)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(context.Background(), storageSync.FileNew, file1)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	select {
	case <-called:
		// all good
	case <-time.After(time.Duration(10 * time.Millisecond)):
		t.Error("Handler was not called during specified time")
	}
}

func TestMessageHandlingOnlyOnce(t *testing.T) {
	ctx := context.Background()
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, ctx, "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		SyncFile(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultSynced, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			called <- true
		}).
		Times(1)

	// start consumer 1
	err := c.StartSubscription(storageSync.FileUpdate)
	if err != nil {
		t.Fatal("Failed to start subscription 1")
	}
	err = c.StartSubscription(storageSync.FileUpdate)
	if err != nil {
		t.Fatal("Failed to start subscription 2")
	}

	err = p.Publish(context.Background(), storageSync.FileUpdate, file1)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	// Wait for called and some more.
	select {
	case <-called:
		// all good
	case <-time.After(time.Duration(10 * time.Millisecond)):
		t.Error("Handler was not called during specified time")
	}
	<-time.After(time.Duration(50 * time.Millisecond))
}

func TestMessageHandlingOnlyOnceSeparateConnections(t *testing.T) {
	ctx := context.Background()
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c1, cleanService1 := getTestService(t, ctx, "clientID1", h)
	defer cleanService1()
	c2, cleanService2 := getTestService(t, ctx, "clientID2", h)
	defer cleanService2()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		SyncFileDelete(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultSynced, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			called <- true
		}).
		Times(1)

	// start consumer 1
	err := c1.StartSubscription(storageSync.FileDelete)
	if err != nil {
		t.Fatal("Failed to start subscription 1")
	}

	err = c2.StartSubscription(storageSync.FileDelete)
	if err != nil {
		t.Fatal("Failed to start subscription 2")
	}

	err = p.Publish(context.Background(), storageSync.FileDelete, file1)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	// Wait for called and some more.
	select {
	case <-called:
		// all good
	case <-time.After(time.Duration(10 * time.Millisecond)):
		t.Error("Handler was not called during specified time")
	}
	<-time.After(time.Duration(50 * time.Millisecond))
}

func TestMessageHandlingNack(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, context.Background(), "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	nack := make(chan bool)
	ok := make(chan bool)
	nackCall := h.EXPECT().
		SyncFile(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultError, fmt.Errorf("error")).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			nack <- true
		}).
		Times(1)
	h.EXPECT().
		SyncFile(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultSynced, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			ok <- true
		}).
		Times(1).
		After(nackCall)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(context.Background(), storageSync.FileNew, file1)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	select {
	case <-nack:
		// Wait 1 second (minimum AckWait time) for redelivery.
		select {
		case <-ok:
			// all good
		case <-time.After(time.Duration(1020 * time.Millisecond)):
			t.Fatal("Handler (mock 'ok' on redelivery) was not called during specified time")
		}
	case <-time.After(time.Duration(10 * time.Millisecond)):
		t.Fatal("Handler (mock 'nack') was not called during specified time")
	}
}

func TestDurability(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, context.Background(), "Consumer_1", h)
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	ok := make(chan bool)
	firstHandlerCall := h.EXPECT().
		SyncFile(gomock.Any(), file1.BucketID, file1.FileID, file1.Version, time1).
		Return(storageSync.ResultSyncNotNeeded, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			ok <- true
		}).
		Times(1)
	h.EXPECT().
		SyncFile(gomock.Any(), file2.BucketID, file2.FileID, file2.Version, time2).
		Return(storageSync.ResultSynced, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			ok <- true
		}).
		Times(1).
		After(firstHandlerCall)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(context.Background(), storageSync.FileNew, file1)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	select {
	case <-ok:
		// all good
	case <-time.After(time.Duration(10 * time.Millisecond)):
		t.Fatal("Handler was not called during specified time")
	}

	// Stop consumer
	cleanService()
	// wait for disconnection
	<-time.After(time.Duration(50 * time.Millisecond))

	// Publish another one
	err = p.Publish(context.Background(), storageSync.FileNew, file2)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	// Start new consumer
	c, cleanService = getTestService(t, context.Background(), "Consumer_2", h)
	defer cleanService()

	// start consumer
	err = c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	select {
	case <-ok:
		// all good
	case <-time.After(time.Duration(50 * time.Millisecond)):
		t.Fatal("Handler was not called during specified time")
	}
}

func TestContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, ctx, "Consumer", h)
	defer cleanService()

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}
	// start consumer
	err = c.StartSubscription(storageSync.FileUpdate)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	cancel()
	time.Sleep(time.Duration(50 * time.Millisecond))

	if c.GetNumberOfSubsriptions() != 0 {
		t.Fatal("Close was not called on cancel context")
	}
}

func getMockHandlers(t *testing.T) (*mock.MockHandlers, func()) {
	mockHandlersCtrl := gomock.NewController(t)
	mockHandlers := mock.NewMockHandlers(mockHandlersCtrl)

	cleanup := func() {
		mockHandlersCtrl.Finish()
	}

	return mockHandlers, cleanup
}

func getTestService(t *testing.T, ctx context.Context, clientID string, h storageSync.Handlers) (*stanConsumer, func()) {
	conn, err := stan.Connect(clusterID, clientID)
	if err != nil {
		t.Fatal("Connection to test stan-straming server failed")
	}

	// initalize consumer
	cfg := Cfg{
		Connection: conn,
		AckWait:    time.Duration(time.Second),
		Handlers:   h,
	}

	c := New(ctx, cfg, zerolog.New(os.Stdout), GetPrometheusMetricsCollection())

	cleanup := func() {
		c.Close()
	}

	return c.(*stanConsumer), cleanup
}

func getTestPublisher(t *testing.T) (storageSync.Publisher, func()) {
	conn, err := stan.Connect(clusterID, "Publisher")
	if err != nil {
		t.Fatal("Connection to test stan-straming server failed")
	}

	cfg := publisher.Cfg{
		Connection:      conn,
		Retries:         5,
		StartRetryWait:  time.Duration(time.Millisecond),
		RetryWaitFactor: 1.0,
	}

	p := publisher.New(context.Background(), cfg, zerolog.New(os.Stdout), publisher.GetPrometheusMetricsCollection())

	cleanup := func() {
		p.Close()
	}

	return p, cleanup
}
