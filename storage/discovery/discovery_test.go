package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/iryonetwork/wwm/gen/discovery/models"
	"github.com/iryonetwork/wwm/storage/discovery/db"
	"github.com/iryonetwork/wwm/storage/discovery/db/mock"
	"github.com/iryonetwork/wwm/utils"
)

var (
	uuid1      = strfmt.UUID("7F1572F5-503A-464C-BF67-A77DB49BB090")
	uuid2      = strfmt.UUID("E2772B4B-EC6B-4509-9960-49F61DDDD08D")
	noErrors   = false
	withErrors = true
	time1, _   = strfmt.ParseDateTime("2018-01-18T15:22:46.123Z")
	card       = &models.Card{
		PatientID: uuid1,
		Connections: models.Connections{
			&models.ConnectionsItems{
				Key:   "K1",
				Value: "V1",
			},
			&models.ConnectionsItems{
				Key:   "K2",
				Value: "V2",
			},
		},
		Locations: models.Locations{uuid1},
	}
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		title         string
		calls         func(*mock.MockDB)
		expected      *models.Card
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull create",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(&patient{
						PatientID: uuid1.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(&connection{
						PatientID: uuid1.String(),
						Key:       "K1",
						Value:     "V1",
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(&connection{
						PatientID: uuid1.String(),
						Key:       "K2",
						Value:     "V2",
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(&location{
						PatientID:  uuid1.String(),
						LocationID: uuid1.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			&models.Card{
				PatientID: uuid1,
				Connections: models.Connections{
					&models.ConnectionsItems{Key: "K1", Value: "V1"},
					&models.ConnectionsItems{Key: "K2", Value: "V2"},
				},
				Locations: models.Locations{uuid1},
			},
			noErrors,
			nil,
		},
		{
			"Failed to create transaction",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")))
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Failed insert",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(&patient{
						PatientID: uuid1.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")),
					db.EXPECT().Rollback().Return(db))
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Failed commit",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Create(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")))
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Create(models.Connections{
				&models.ConnectionsItems{Key: "K1", Value: "V1"},
				&models.ConnectionsItems{Key: "K2", Value: "V2"},
			}, models.Locations{uuid1})

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		title         string
		calls         func(*mock.MockDB)
		expected      *models.Card
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull update",
			func(db *mock.MockDB) {
				p := &patient{
					PatientID: uuid1.String(),
					Connections: []connection{
						connection{PatientID: uuid1.String(), Key: "K1", Value: "V1"},
						connection{PatientID: uuid1.String(), Key: "K2", Value: "V2"},
					},
					Locations: []location{
						location{PatientID: uuid1.String(), LocationID: uuid1.String()},
					},
				}
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// read existing data
					db.EXPECT().Preload("Connections").Return(db),
					db.EXPECT().Preload("Locations").Return(db),
					db.EXPECT().First(&patient{PatientID: uuid1.String()}).SetArg(0, *p).Return(db),
					db.EXPECT().GetError().Return(nil),

					// remove a connection
					db.EXPECT().Delete(&connection{
						PatientID: uuid1.String(),
						Key:       "K1",
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// replace existing connection
					db.EXPECT().Save(&connection{
						PatientID: uuid1.String(),
						Key:       "K2",
						Value:     "V2.2",
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// add new connection
					db.EXPECT().Create(&connection{
						PatientID: uuid1.String(),
						Key:       "K3",
						Value:     "V3",
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// remove a location
					db.EXPECT().Delete(&location{
						PatientID:  uuid1.String(),
						LocationID: uuid1.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// add a new location
					db.EXPECT().Create(&location{
						PatientID:  uuid1.String(),
						LocationID: uuid2.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			&models.Card{
				PatientID: uuid1,
				Connections: models.Connections{
					&models.ConnectionsItems{Key: "K2", Value: "V2.2"},
					&models.ConnectionsItems{Key: "K3", Value: "V3"},
				},
				Locations: models.Locations{uuid2},
			},
			noErrors,
			nil,
		},
		{
			"Patient not found",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// read existing data
					db.EXPECT().Preload("Connections").Return(db),
					db.EXPECT().Preload("Locations").Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound),
					db.EXPECT().Rollback())
			},
			nil,
			withErrors,
			ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Update(uuid1, models.Connections{
				&models.ConnectionsItems{Key: "K2", Value: "V2.2"},
				&models.ConnectionsItems{Key: "K3", Value: "V3"},
			}, models.Locations{uuid2})

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		title         string
		calls         func(*mock.MockDB)
		expected      *models.Card
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull read",
			func(db *mock.MockDB) {
				p := &patient{
					PatientID: uuid1.String(),
					Connections: []connection{
						connection{PatientID: uuid1.String(), Key: "K1", Value: "V1"},
						connection{PatientID: uuid1.String(), Key: "K2", Value: "V2"},
					},
					Locations: []location{
						location{PatientID: uuid1.String(), LocationID: uuid1.String()},
					},
				}
				gomock.InOrder(
					db.EXPECT().Preload("Connections").Return(db),
					db.EXPECT().Preload("Locations").Return(db),
					db.EXPECT().First(&patient{PatientID: uuid1.String()}).SetArg(0, *p).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&models.Card{
				PatientID: uuid1,
				Connections: models.Connections{
					&models.ConnectionsItems{Key: "K1", Value: "V1"},
					&models.ConnectionsItems{Key: "K2", Value: "V2"},
				},
				Locations: models.Locations{uuid1},
			},
			noErrors,
			nil,
		},
		{
			"Patient not found",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// read existing data
					db.EXPECT().Preload("Connections").Return(db),
					db.EXPECT().Preload("Locations").Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound))
			},
			nil,
			withErrors,
			ErrNotFound,
		},
		{
			"First read fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// read existing data
					db.EXPECT().Preload("Connections").Return(db),
					db.EXPECT().Preload("Locations").Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")))
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Get(uuid1)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		title         string
		calls         func(*mock.MockDB)
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull delete",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(&patient{PatientID: uuid1.String()}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// delete rows
					db.EXPECT().Delete(connection{}, "patient_id = ?", uuid1.String()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(location{}, "patient_id = ?", uuid1.String()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(patient{}, "patient_id = ?", uuid1.String()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			noErrors,
			nil,
		},
		{
			"Begin transaction fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")))
			},
			withErrors,
			nil,
		},
		{
			"Patient not found",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound))
			},
			withErrors,
			ErrNotFound,
		},
		{
			"Connection delete fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// delete rows
					db.EXPECT().Delete(connection{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")),
					db.EXPECT().Rollback())
			},
			withErrors,
			nil,
		},
		{
			"location delete fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// delete rows
					db.EXPECT().Delete(connection{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(location{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")),
					db.EXPECT().Rollback())
			},
			withErrors,
			nil,
		},
		{
			"Patient delete fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// delete rows
					db.EXPECT().Delete(connection{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(location{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(patient{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")),
					db.EXPECT().Rollback())
			},
			withErrors,
			nil,
		},
		{
			"Commit fails",
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// load patient
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// delete rows
					db.EXPECT().Delete(connection{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(location{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),
					db.EXPECT().Delete(patient{}, gomock.Any(), gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("Error")))
			},
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			err := s.Delete(uuid1)

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		title         string
		query         string
		calls         func(sqlmock.Sqlmock)
		expected      models.Cards
		errorExpected bool
		exactError    error
	}{
		{
			"Singe token",
			"single",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT patient_id FROM \"connections\" WHERE \\(value ILIKE \\$1\\)").
					WithArgs("%single%").
					WillReturnRows(sqlmock.NewRows([]string{"patient_id"}).AddRow(uuid1.String()))
				mock.ExpectQuery("SELECT \\* FROM \"patients\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id"}).AddRow(uuid1.String()))
				mock.ExpectQuery("SELECT \\* FROM \"connections\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id", "key", "value"}).
						AddRow(uuid1.String(), "K1", "V1").
						AddRow(uuid1.String(), "K2", "V2"))
				mock.ExpectQuery("SELECT \\* FROM \"locations\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id", "location_id"}).
						AddRow(uuid1.String(), uuid1.String()))
			},
			models.Cards{
				&models.Card{
					PatientID: uuid1,
					Connections: models.Connections{
						&models.ConnectionsItems{Key: "K1", Value: "V1"},
						&models.ConnectionsItems{Key: "K2", Value: "V2"},
					},
					Locations: models.Locations{uuid1},
				},
			},
			noErrors,
			nil,
		},
		{
			"Composed token",
			"composed token",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT patient_id FROM \"connections\" WHERE \\(value ILIKE \\$1 OR value ILIKE \\$2\\)").
					WithArgs("%composed%", "%token%").
					WillReturnRows(sqlmock.NewRows([]string{"patient_id"}).AddRow(uuid1.String()))
				mock.ExpectQuery("SELECT \\* FROM \"patients\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id"}).AddRow(uuid1.String()))
				mock.ExpectQuery("SELECT \\* FROM \"connections\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id", "key", "value"}).
						AddRow(uuid1.String(), "K1", "V1").
						AddRow(uuid1.String(), "K2", "V2"))
				mock.ExpectQuery("SELECT \\* FROM \"locations\" .+").
					WithArgs(uuid1.String()).
					WillReturnRows(sqlmock.NewRows([]string{"patient_id", "location_id"}).
						AddRow(uuid1.String(), uuid1.String()))
			},
			models.Cards{
				&models.Card{
					PatientID: uuid1,
					Connections: models.Connections{
						&models.ConnectionsItems{Key: "K1", Value: "V1"},
						&models.ConnectionsItems{Key: "K2", Value: "V2"},
					},
					Locations: models.Locations{uuid1},
				},
			},
			noErrors,
			nil,
		},
		{
			"First select no results",
			"noResults",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT patient_id .*").
					WillReturnRows(sqlmock.NewRows([]string{"patient_id"}))
			},
			models.Cards{},
			noErrors,
			nil,
		},
		{
			"First select fails",
			"fails",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT patient_id .+").
					WillReturnError(fmt.Errorf("Failed"))
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestDB(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Find(tc.query)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n\t%s\nto equal\n\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestLink(t *testing.T) {
	testCases := []struct {
		title         string
		locationID    strfmt.UUID
		calls         func(*mock.MockDB)
		getCardRes    *models.Card
		getCardError  error
		expected      models.Locations
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull link",
			uuid2,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// create a link
					db.EXPECT().Create(&location{
						PatientID:  uuid1.String(),
						LocationID: uuid2.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			card,
			nil,
			models.Locations{uuid1, uuid2},
			noErrors,
			nil,
		},
		{
			"Patient not found",
			uuid2,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			nil,
			ErrNotFound,
			nil,
			withErrors,
			ErrNotFound,
		},
		{
			"Fetching patient fails",
			uuid2,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			nil,
			fmt.Errorf("error"),
			nil,
			withErrors,
			nil,
		},
		{
			"Link already exists",
			uuid1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			card,
			nil,
			models.Locations{uuid1},
			noErrors,
			nil,
		},
		{
			"Commit fails",
			uuid2,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// create a link
					db.EXPECT().Create(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")))
			},
			card,
			nil,
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, mdb, c := getTestStorage(t)
			defer c()

			origGetCard := getCard
			getCard = func(_ db.DB, _ strfmt.UUID) (*models.Card, error) {
				return tc.getCardRes, tc.getCardError
			}
			defer func() {
				getCard = origGetCard
			}()

			// collect mocked calls
			tc.calls(mdb)

			// call the method
			out, err := s.Link(uuid1, tc.locationID)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestUnlink(t *testing.T) {
	testCases := []struct {
		title         string
		locationID    strfmt.UUID
		calls         func(*mock.MockDB)
		getCardRes    *models.Card
		getCardError  error
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull unlink",
			uuid1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// create a link
					db.EXPECT().Delete(&location{
						PatientID:  uuid1.String(),
						LocationID: uuid1.String(),
					}).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			card,
			nil,
			noErrors,
			nil,
		},
		{
			"Patient not found",
			uuid1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			nil,
			ErrNotFound,
			withErrors,
			ErrNotFound,
		},
		{
			"Fetching patient fails",
			uuid1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			nil,
			fmt.Errorf("error"),
			withErrors,
			nil,
		},
		{
			"Link does not exist",
			uuid2,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// rollback changes
					db.EXPECT().Rollback())
			},
			card,
			nil,
			withErrors,
			ErrNotFound,
		},
		{
			"Commit fails",
			uuid1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					// start a transaction
					db.EXPECT().Begin().Return(db),
					db.EXPECT().GetError().Return(nil),

					// create a link
					db.EXPECT().Delete(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(nil),

					// commit changes
					db.EXPECT().Commit().Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")))
			},
			card,
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, mdb, c := getTestStorage(t)
			defer c()

			origGetCard := getCard
			getCard = func(_ db.DB, _ strfmt.UUID) (*models.Card, error) {
				return tc.getCardRes, tc.getCardError
			}
			defer func() {
				getCard = origGetCard
			}()

			// collect mocked calls
			tc.calls(mdb)

			// call the method
			err := s.Unlink(uuid1, tc.locationID)

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestGetCodes(t *testing.T) {
	cat := "CAT"
	id := "ID"
	title := "TITLE"
	code := &models.Code{
		Category: &cat,
		ID:       &id,
		Locale:   "LOC",
		Title:    &title,
	}
	codeWithParent := &models.Code{
		Category: &cat,
		ID:       &id,
		Locale:   "LOC",
		Title:    &title,
		ParentID: "PARENT",
	}

	testCases := []struct {
		title         string
		category      string
		query         string
		locale        string
		parentID      string
		calls         func(sqlmock.Sqlmock)
		expected      models.Codes
		errorExpected bool
		exactError    error
	}{
		{
			"Basic call",
			"CAT",
			"",
			"",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.locale = \\$2 ORDER BY").
					WithArgs("CAT", "en").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", nil))
			},
			models.Codes{code},
			noErrors,
			nil,
		},
		{
			"Custom locale",
			"CAT",
			"",
			"LOC",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.locale = \\$2 ORDER BY").
					WithArgs("CAT", "LOC").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", nil))
			},
			models.Codes{code},
			noErrors,
			nil,
		},
		{
			"With query",
			"CAT",
			"QS",
			"",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.locale = \\$2 AND ct\\.title ILIKE \\$3 ORDER BY").
					WithArgs("CAT", "en", "%QS%").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", nil))
			},
			models.Codes{code},
			noErrors,
			nil,
		},
		{
			"With parent ID",
			"CAT",
			"",
			"",
			"PARENT",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.locale = \\$2 AND c\\.parent_id = \\$3 ORDER BY").
					WithArgs("CAT", "en", "PARENT").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", "PARENT"))
			},
			models.Codes{codeWithParent},
			noErrors,
			nil,
		},
		{
			"Failed query",
			"CAT",
			"",
			"",
			"PARENT",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+").
					WillReturnError(fmt.Errorf("Error"))
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, mock, c := getTestDB(t)
			defer c()

			// collect mocked calls
			tc.calls(mock)

			// call the method
			out, err := s.CodesGet(tc.category, tc.query, tc.parentID, tc.locale)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestGetCode(t *testing.T) {
	cat := "CAT"
	id := "ID"
	title := "TITLE"
	code := &models.Code{
		Category: &cat,
		ID:       &id,
		Locale:   "LOC",
		Title:    &title,
	}
	codeWithParent := &models.Code{
		Category: &cat,
		ID:       &id,
		Locale:   "LOC",
		Title:    &title,
		ParentID: "PARENT",
	}

	testCases := []struct {
		title         string
		category      string
		ID            string
		locale        string
		calls         func(sqlmock.Sqlmock)
		expected      *models.Code
		errorExpected bool
		exactError    error
	}{
		{
			"Success with basic call",
			"CAT",
			"ID",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.code_id = \\$2 AND ct\\.locale = \\$3").
					WithArgs("CAT", "ID", "en").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", nil))
			},
			code,
			noErrors,
			nil,
		},
		{
			"Success with custom locale",
			"CAT",
			"ID",
			"LOC",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.code_id = \\$2 AND ct\\.locale = \\$3").
					WithArgs("CAT", "ID", "LOC").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}).
						AddRow("CAT", "ID", "TITLE", "LOC", "PARENT"))
			},
			codeWithParent,
			noErrors,
			nil,
		},
		{
			"Failed query",
			"CAT",
			"ID",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+").
					WillReturnError(fmt.Errorf("Error"))
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Code not found",
			"CAT",
			"ID",
			"",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM codes AS c INNER JOIN code_titles AS ct .+ "+
					"c\\.category_id = \\$1 AND ct\\.code_id = \\$2 AND ct\\.locale = \\$3").
					WithArgs("CAT", "ID", "en").
					WillReturnRows(sqlmock.NewRows([]string{"category_id", "code_id", "title", "locale", "parent_id"}))
			},
			nil,
			withErrors,
			utils.NewError(utils.ErrNotFound, "code_does_not_exist"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, mock, c := getTestDB(t)
			defer c()

			// collect mocked calls
			tc.calls(mock)

			// call the method
			out, err := s.CodeGet(tc.category, tc.ID, tc.locale)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func getTestStorage(t *testing.T) (*storage, *mock.MockDB, func()) {
	// mock getNewUUID
	origGetNewUUID := getNewUUID
	getNewUUID = func() strfmt.UUID {
		return uuid1
	}

	// mock getCurrentTime
	origGetCurrentTime := getCurrentTime
	getCurrentTime = func() time.Time {
		return time.Time(time1)
	}

	// setup minio mock
	dbCtrl := gomock.NewController(t)
	db := mock.NewMockDB(dbCtrl)

	s := &storage{
		db:         db,
		logger:     zerolog.New(os.Stdout),
		locationID: uuid1.String(),
	}

	cleanup := func() {
		getNewUUID = origGetNewUUID
		getCurrentTime = origGetCurrentTime
		dbCtrl.Finish()
	}

	return s, db, cleanup
}

func getTestDB(t *testing.T) (*storage, sqlmock.Sqlmock, func()) {
	// mock getNewUUID
	origGetNewUUID := getNewUUID
	getNewUUID = func() strfmt.UUID {
		return uuid1
	}

	// mock getCurrentTime
	origGetCurrentTime := getCurrentTime
	getCurrentTime = func() time.Time {
		return time.Time(time1)
	}

	// setup minio mock

	mdb, mock, _ := sqlmock.New()
	gdb, _ := gorm.Open("postgres", mdb)

	s := &storage{
		gdb:        gdb,
		logger:     zerolog.New(os.Stdout),
		locationID: uuid1.String(),
		db:         db.New(gdb),
	}

	cleanup := func() {
		getNewUUID = origGetNewUUID
		getCurrentTime = origGetCurrentTime
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}

	return s, mock, cleanup
}

func toJSON(in interface{}) string {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(in)
	return buf.String()
}
