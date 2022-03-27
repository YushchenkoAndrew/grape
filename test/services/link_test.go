package services_test

import (
	"api/config"
	"api/db"
	"api/interfaces"
	m "api/models"
	"api/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	ID    = 305
	DELAY = time.Second
)

var (
	s service.LinkService
)

func TestDefault(t *testing.T) {
	var tests = struct {
		create []struct {
			name  string
			model m.Link
			err   bool
		}

		read []struct {
			name     string
			query    m.LinkQueryDto
			expected []m.Link
		}

		update []struct {
			name  string
			query m.LinkQueryDto
			model m.Link
			err   bool
		}
	}{
		create: []struct {
			name  string
			model m.Link
			err   bool
		}{
			{
				name:  "Create new link record",
				model: m.Link{ID: ID, Name: "main", Link: "test", ProjectID: 1},
			},
			{
				name:  "Create new link record",
				model: m.Link{Name: "main2", Link: "test", ProjectID: 1},
			},
			{
				name:  "Update already existed link record with 'Create' func",
				model: m.Link{Name: "main", Link: "test2", ProjectID: 1},
			},
			{
				name:  "Check on project_id existens",
				model: m.Link{Name: "main", Link: "test", ProjectID: 10000},
				err:   true,
			},
		},

		read: []struct {
			name     string
			query    m.LinkQueryDto
			expected []m.Link
		}{
			{
				name:     "Read link record by ID",
				query:    m.LinkQueryDto{ID: ID},
				expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}},
			},
			{
				name:     "Read Link record by name & project_id",
				query:    m.LinkQueryDto{Name: "main", ProjectID: 1},
				expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}},
			},
			{
				name:     "Read link with pagination",
				query:    m.LinkQueryDto{Name: "main", ProjectID: 1, Page: 1},
				expected: []m.Link{},
			},
			{
				name:     "Read Link record by project_id",
				query:    m.LinkQueryDto{ProjectID: 1},
				expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}, {Name: "main2", Link: "test", ProjectID: 1}},
			},
			{
				name:     "Read not existed values",
				query:    m.LinkQueryDto{Name: "main", ProjectID: 10000},
				expected: []m.Link{},
			},
		},

		update: []struct {
			name  string
			query m.LinkQueryDto
			model m.Link
			err   bool
		}{
			{
				name:  "Update link record by ID",
				query: m.LinkQueryDto{ID: ID},
				model: m.Link{ID: ID, Name: "main", Link: "test54", ProjectID: 1},
			},
			{
				name:  "Update link record by query",
				query: m.LinkQueryDto{Name: "main2", ProjectID: 1},
				model: m.Link{Name: "main2", Link: "test3", ProjectID: 1},
			},
			{
				name:  "Update not existed record",
				query: m.LinkQueryDto{Name: "main", ProjectID: 10000},
				model: m.Link{Name: "main", Link: "test2", ProjectID: 1},
				err:   true,
			},
			{
				name:  "Update multiple records at once",
				query: m.LinkQueryDto{ProjectID: 1},
				model: m.Link{Link: "test"},
			},
		},
	}

	for _, tc := range tests.create {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := s.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Link, model.Link)
		})
	}

	for _, tc := range tests.read {
		t.Run(tc.name, func(t *testing.T) {
			models, err := s.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))
			if len(models) != 1 {
				return
			}

			for _, model := range tc.expected {
				require.Equal(t, model.Name, model.Name)
				require.Equal(t, model.Link, model.Link)
			}
		})
	}

	for _, tc := range tests.update {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := s.Update(&tc.query, model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Link, model.Link)
		})
	}
}

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./"),
	}).Init()

	db, client := db.Init([]interfaces.Table{
		m.NewLink(),
	})

	s = *service.NewLinkService(db, client)
}
