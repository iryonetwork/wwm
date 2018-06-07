package waitlist

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/service/waitlist/mock"
	"github.com/iryonetwork/wwm/utils"
)

func TestCreateItem(t *testing.T) {
	s, storage, cleanup := getTestService(t)
	defer cleanup()

	// waitlistID to test
	waitlistID, _ := uuid.NewV4()
	// create item to test
	patient1ID, _ := uuid.NewV4()
	item1 := &models.Item{
		Priority:  swag.Int64(1),
		PatientID: swag.String(patient1ID.String()),
	}
	item2 := &models.Item{
		Priority:  swag.Int64(4),
		PatientID: swag.String(patient1ID.String()),
	}
	patient2ID, _ := uuid.NewV4()
	item3 := &models.Item{
		Priority:  swag.Int64(1),
		PatientID: swag.String(patient2ID.String()),
	}

	// #1 Test adding item succesfully
	storage.EXPECT().ListItems(waitlistID.Bytes()).Return([]*models.Item{item1}, nil)
	storage.EXPECT().AddItem(waitlistID.Bytes(), item3).Return(item3, nil).Times(1)
	_, err := s.CreateItem(waitlistID.Bytes(), item3)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// #2 Test failure on conflicting patientID
	storage.EXPECT().ListItems(waitlistID.Bytes()).Return([]*models.Item{item1}, nil)
	_, err = s.CreateItem(waitlistID.Bytes(), item2)
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrConflict {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrConflict, uErr.Code())
	}
}

func getTestService(t *testing.T) (Service, *mock.MockStorage, func()) {
	// setup storage mock
	storageCtrl := gomock.NewController(t)
	storage := mock.NewMockStorage(storageCtrl)

	svc := &service{
		storage: storage,
		logger:  zerolog.New(os.Stdout),
	}

	cleanup := func() {
		storageCtrl.Finish()
	}

	return svc, storage, cleanup
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
