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

const FILE_ID = 305

var file service.FileService

func TestFileCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.File
		err   bool
	}{
		{
			name:  "Create new file record",
			model: m.File{ID: FILE_ID, Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 1},
		},
		{
			name:  "Create new file record",
			model: m.File{Name: "test.txt", Type: "text/plain", Role: "assets", Path: "", ProjectID: 1},
		},
		{
			name:  "Check what will happend if create already existed file",
			model: m.File{Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 1},
		},
		{
			name:  "Check on project_id existens",
			model: m.File{Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 10000},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := file.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Type, model.Type)
			require.Equal(t, tc.model.Role, model.Role)
			require.Equal(t, tc.model.Path, model.Path)
		})
	}
}

func TestFileRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.FileQueryDto
		expected []m.File
	}{
		{
			name:     "Read file record by ID",
			query:    m.FileQueryDto{ID: FILE_ID},
			expected: []m.File{{Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 1}},
		},
		{
			name:     "Read File record by role & project_id",
			query:    m.FileQueryDto{Role: "assets", ProjectID: 1},
			expected: []m.File{{Name: "test.txt", Type: "text/plain", Role: "assets", Path: "", ProjectID: 1}},
		},
		{
			name:     "Read file with pagination & query filter",
			query:    m.FileQueryDto{ProjectID: 1, Page: 0},
			expected: []m.File{{Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 1}, {Name: "test.txt", Type: "text/plain", Role: "assets", Path: "", ProjectID: 1}},
		},
		{
			name:     "Read file with pagination (outside borders) & query filter",
			query:    m.FileQueryDto{Name: "test.txt", ProjectID: 1, Page: 1},
			expected: []m.File{},
		},
		{
			name:     "Read File record by project_id",
			query:    m.FileQueryDto{ProjectID: 1},
			expected: []m.File{{Name: "test.txt", Type: "text/plain", Role: "src", Path: "", ProjectID: 1}, {Name: "test.txt", Type: "text/plain", Role: "assets", Path: "", ProjectID: 1}},
		},
		{
			name:     "Read not existed values",
			query:    m.FileQueryDto{Name: "test.txt", ProjectID: 10000},
			expected: []m.File{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := file.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Type, models[i].Type)
				require.Equal(t, expected.Role, models[i].Role)
				require.Equal(t, expected.Path, models[i].Path)
			}
		})
	}
}

func TestFileUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.FileQueryDto
		model m.File
		err   bool
	}{
		{
			name:  "Update file record by ID",
			query: m.FileQueryDto{ID: FILE_ID},
			model: m.File{Name: "test.txt", Type: "text/plain", Role: "src", Path: "/test"},
		},
		{
			name:  "Update multiple files record by query",
			query: m.FileQueryDto{Role: "src", ProjectID: 1},
			model: m.File{Name: "index.json", Type: "application/js"},
		},
		{
			name:  "Update multiple files record by query",
			query: m.FileQueryDto{ProjectID: 1},
			model: m.File{Type: "application/json"},
		},
		{
			name:  "Update not existed record",
			query: m.FileQueryDto{Name: "test.txt", ProjectID: 10000},
			model: m.File{Name: "index.json", Type: "application/json"},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := file.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.model.Type, models[0].Type)

			if tc.model.Name != "" {
				require.Equal(t, tc.model.Name, models[0].Name)
			}

			if tc.model.Role != "" {
				require.Equal(t, tc.model.Role, models[0].Role)
			}

			if tc.model.Path != "" {
				require.Equal(t, tc.model.Path, models[0].Path)
			}
		})
	}
}

func TestFileCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.FileQueryDto
		expected []m.File
	}{
		{
			name:     "Check if cached File record by ID was updated",
			query:    m.FileQueryDto{ID: FILE_ID},
			expected: []m.File{{Name: "index.json", Type: "application/json", Role: "src", Path: "/test"}},
		},
		{
			name:     "Check if cached File record by role & project_id was updated",
			query:    m.FileQueryDto{Role: "assets", ProjectID: 1},
			expected: []m.File{{Name: "test.txt", Type: "application/json", Role: "assets", Path: ""}},
		},
		{
			name:     "Check if cached File record with pagination & project_id was updated",
			query:    m.FileQueryDto{ProjectID: 1, Page: 0},
			expected: []m.File{{Name: "index.json", Type: "application/json", Role: "src", Path: "/test"}, {Name: "test.txt", Type: "application/json", Role: "assets", Path: ""}},
		},
		{
			name:     "Check if cached File record by project_id was updated",
			query:    m.FileQueryDto{ProjectID: 1},
			expected: []m.File{{Name: "index.json", Type: "application/json", Role: "src", Path: "/test"}, {Name: "test.txt", Type: "application/json", Role: "assets", Path: ""}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := file.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Type, models[i].Type)
				require.Equal(t, expected.Role, models[i].Role)
				require.Equal(t, expected.Path, models[i].Path)
			}
		})
	}

}

func TestFileDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.FileQueryDto
		expected int
	}{
		{
			name:     "Delete File record by ID",
			query:    m.FileQueryDto{ID: FILE_ID},
			expected: 1,
		},
		{
			name:     "Delete File record by role & project_id",
			query:    m.FileQueryDto{Role: "src", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Delete File record by project_id",
			query:    m.FileQueryDto{ProjectID: 1},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := file.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestFinalCacheState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.FileQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by ID",
			query:    m.FileQueryDto{ID: FILE_ID},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by role & project_id",
			query:    m.FileQueryDto{Role: "src", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination & query filter",
			query:    m.FileQueryDto{ProjectID: 1, Page: 0},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.FileQueryDto{Role: "src", ProjectID: 1, Page: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by project_id",
			query:    m.FileQueryDto{ProjectID: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := file.Read(&tc.query)
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
		m.NewFile(),
	})

	file = *service.NewFileService(db, client)
}
