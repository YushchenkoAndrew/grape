package services_test

import (
	"api/config"
	"api/db"
	"api/interfaces"
	m "api/models"
	"api/service"
	"testing"

	"github.com/stretchr/testify/require"
)

const LINK_ID = 305

var link service.LinkService

func TestLinkCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Link
		err   bool
	}{
		{
			name:  "Create new link record",
			model: m.Link{ID: LINK_ID, Name: "main", Link: "test", ProjectID: 1},
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := link.Create(model)
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
}

func TestLinkRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.LinkQueryDto
		expected []m.Link
	}{
		{
			name:     "Read link record by ID",
			query:    m.LinkQueryDto{ID: LINK_ID},
			expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}},
		},
		{
			name:     "Read Link record by name & project_id",
			query:    m.LinkQueryDto{Name: "main", ProjectID: 1},
			expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}},
		},
		{
			name:     "Read link with pagination & query filter",
			query:    m.LinkQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Link{{Name: "main", Link: "test2", ProjectID: 1}, {Name: "main2", Link: "test", ProjectID: 1}},
		},
		{
			name:     "Read link with pagination (outside borders) & query filter",
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := link.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Link, models[i].Link)
			}
		})
	}
}

func TestLinkUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.LinkQueryDto
		model m.Link
		err   bool
	}{
		{
			name:  "Update link record by ID",
			query: m.LinkQueryDto{ID: LINK_ID},
			model: m.Link{ID: LINK_ID, Name: "main", Link: "test54", ProjectID: 1},
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := link.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.model.Link, models[0].Link)

			if tc.model.Name != "" {
				require.Equal(t, tc.model.Name, models[0].Name)
			}
		})
	}
}

func TestLinkCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.LinkQueryDto
		expected []m.Link
	}{
		{
			name:     "Check if cached Link record by ID was updated",
			query:    m.LinkQueryDto{ID: LINK_ID},
			expected: []m.Link{{Name: "main", Link: "test", ProjectID: 1}},
		},
		{
			name:     "Check if cached Link record by name & project_id was updated",
			query:    m.LinkQueryDto{Name: "main", ProjectID: 1},
			expected: []m.Link{{Name: "main", Link: "test", ProjectID: 1}},
		},
		{
			name:     "Check if cached Link record with pagination & project_id was updated",
			query:    m.LinkQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Link{{Name: "main", Link: "test", ProjectID: 1}, {Name: "main2", Link: "test", ProjectID: 1}},
		},
		{
			name:     "Check if cached Link record by project_id was updated",
			query:    m.LinkQueryDto{ProjectID: 1},
			expected: []m.Link{{Name: "main", Link: "test", ProjectID: 1}, {Name: "main2", Link: "test", ProjectID: 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := link.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Link, models[i].Link)
			}
		})
	}

}

func TestLinkDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.LinkQueryDto
		expected int
	}{
		{
			name:     "Delete Link record by ID",
			query:    m.LinkQueryDto{ID: LINK_ID},
			expected: 1,
		},
		{
			name:     "Delete Link record by name & project_id",
			query:    m.LinkQueryDto{Name: "main", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Delete Link record by project_id",
			query:    m.LinkQueryDto{ProjectID: 1},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := link.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestLinkFinalCacheState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.LinkQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by ID",
			query:    m.LinkQueryDto{ID: LINK_ID},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by name & project_id",
			query:    m.LinkQueryDto{Name: "main", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination & query filter",
			query:    m.LinkQueryDto{ProjectID: 1, Page: 0},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.LinkQueryDto{Name: "main", ProjectID: 1, Page: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by project_id",
			query:    m.LinkQueryDto{ProjectID: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := link.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, len(models))
		})
	}
}

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("../", ""),
	}).Init()

	db, client := db.Init([]interfaces.Table{
		m.NewLink(),
		m.NewProject(),
	})

	link = *service.NewLinkService(db, client)

	var project = *service.NewProjectService(db, client)
	project.Create(&m.Project{ID: 1, Name: "yes", Title: "js", Flag: "js", Desc: "js", Note: "js"})
}
