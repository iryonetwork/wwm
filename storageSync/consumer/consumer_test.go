package consumer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats-streaming-server/server"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storageSync"
	"github.com/iryonetwork/wwm/storageSync/consumer/mock"
	"github.com/iryonetwork/wwm/storageSync/publisher"
)

var (
	// clusterID for test server.
	clusterID = "TestCluster"
	time1, _  = strfmt.ParseDateTime("2018-02-05T15:16:15.123Z")
	file1     = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time1,
		Name:        "file11",
		Path:        "BUCKET/file11/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
	}
	file2 = &models.FileDescriptor{
		Archetype:   "",
		Checksum:    "CHS",
		ContentType: "image/jpeg",
		Created:     time1,
		Name:        "Image",
		Path:        "BUCKET/Image/V1",
		Version:     "V1",
		Size:        15698,
		Operation:   "w",
	}
)

func TestMain(m *testing.M) {
	// Start nats-streaming server
	s, err := server.RunServer(clusterID)
	if err != nil {
		os.Exit(1)
	}

	c := m.Run()

	s.Shutdown()
	os.Exit(c)
}

func TestStartSuccess(t *testing.T) {
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, "Consumer", h)
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
	c, cleanService := getTestService(t, "Consumer", h)
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
	c, cleanService := getTestService(t, "Consumer", h)

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
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		FileNew(gomock.Eq(file1)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
			called <- true
		}).
		Times(1)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(storageSync.FileNew, file1)
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
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c, cleanService := getTestService(t, "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		FileUpdate(gomock.Eq(file1)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
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

	err = p.Publish(storageSync.FileUpdate, file1)
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
	h, cleanHandlers := getMockHandlers(t)
	defer cleanHandlers()
	c1, cleanService1 := getTestService(t, "clientID1", h)
	defer cleanService1()
	c2, cleanService2 := getTestService(t, "clientID2", h)
	defer cleanService2()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	called := make(chan bool)
	h.EXPECT().
		FileDelete(gomock.Eq(file1)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
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

	err = p.Publish(storageSync.FileDelete, file1)
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
	c, cleanService := getTestService(t, "Consumer", h)
	defer cleanService()
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	nack := make(chan bool)
	ok := make(chan bool)
	nackCall := h.EXPECT().
		FileNew(gomock.Eq(file1)).
		Return(fmt.Errorf("error")).
		Do(func(fd *models.FileDescriptor) {
			nack <- true
		}).
		Times(1)
	h.EXPECT().
		FileNew(gomock.Eq(file1)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
			ok <- true
		}).
		Times(1).
		After(nackCall)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(storageSync.FileNew, file1)
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
	c, cleanService := getTestService(t, "Consumer_1", h)
	p, cleanPublisher := getTestPublisher(t)
	defer cleanPublisher()

	// Expect handler call
	ok := make(chan bool)
	firstHandlerCall := h.EXPECT().
		FileNew(gomock.Eq(file1)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
			ok <- true
		}).
		Times(1)
	h.EXPECT().
		FileNew(gomock.Eq(file2)).
		Return(nil).
		Do(func(fd *models.FileDescriptor) {
			ok <- true
		}).
		Times(1).
		After(firstHandlerCall)

	// start consumer
	err := c.StartSubscription(storageSync.FileNew)
	if err != nil {
		t.Fatal("Failed to start subscription")
	}

	err = p.Publish(storageSync.FileNew, file1)
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

	// Publish another one
	err = p.Publish(storageSync.FileNew, file2)
	if err != nil {
		t.Fatal("Failed to publish to test nats-streaming server")
	}

	// Start new consumer
	c, cleanService = getTestService(t, "Consumer_2", h)
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

func getMockHandlers(t *testing.T) (*mock.MockHandlers, func()) {
	mockHandlersCtrl := gomock.NewController(t)
	mockHandlers := mock.NewMockHandlers(mockHandlersCtrl)

	cleanup := func() {
		mockHandlersCtrl.Finish()
	}

	return mockHandlers, cleanup
}

func getTestService(t *testing.T, clientID string, h Handlers) (*stanConsumer, func()) {
	conn, err := stan.Connect(clusterID, clientID)
	if err != nil {
		t.Fatal("Connection to test stan-straming server failed")
	}

	c := &stanConsumer{
		conn:     conn,
		handlers: h,
		ackWait:  time.Duration(time.Second),
		logger:   zerolog.New(os.Stdout),
	}

	cleanup := func() {
		c.Close()
	}

	return c, cleanup
}

func getTestPublisher(t *testing.T) (storageSync.Publisher, func()) {
	conn, err := stan.Connect(clusterID, "Publisher")
	if err != nil {
		t.Fatal("Connection to test stan-straming server failed")
	}

	p, _ := publisher.New(conn, 5, time.Duration(time.Millisecond), 1.0, zerolog.New(os.Stdout))

	cleanup := func() {
		p.Close()
	}

	return p, cleanup
}
