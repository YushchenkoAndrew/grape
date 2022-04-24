package services_test

import (
	"api/client"
	"api/config"
	"api/interfaces"
	i "api/interfaces/service"
	m "api/models"
	"api/service"
	"testing"

	"github.com/stretchr/testify/require"
)

const PATTERN_ID = 305

var pattern i.Default[m.Pattern, m.PatternQueryDto]

func TestPatternCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Pattern
		err   bool
	}{
		{
			name:  "Create new pattern record",
			model: m.Pattern{ID: PATTERN_ID, Mode: "test", Colors: 3, Path: "<path />"},
		},
		{
			name:  "Create new pattern record",
			model: m.Pattern{Mode: "test", Colors: 4, Path: "<path t='5' />"},
		},
		{
			name:  "Check what will happend if create already existed pattern",
			model: m.Pattern{Mode: "test", Colors: 4, Path: "<path />"},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := pattern.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Mode, model.Mode)
			require.Equal(t, tc.model.Colors, model.Colors)
			require.Equal(t, tc.model.Path, model.Path)
		})
	}
}

func TestPatternRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.PatternQueryDto
		expected []m.Pattern
	}{
		{
			name:     "Read Pattern record by ID",
			query:    m.PatternQueryDto{ID: PATTERN_ID},
			expected: []m.Pattern{{ID: PATTERN_ID, Mode: "test", Colors: 3, Path: "<path />"}},
		},
		{
			name:     "Read Pattern record by mode",
			query:    m.PatternQueryDto{Mode: "test"},
			expected: []m.Pattern{{Mode: "test", Colors: 4, Path: "<path t='5' />"}, {ID: PATTERN_ID, Mode: "test", Colors: 3, Path: "<path />"}},
		},
		{
			name:     "Read pattern with pagination & query filter",
			query:    m.PatternQueryDto{Page: 0, Mode: "test"},
			expected: []m.Pattern{{Mode: "test", Colors: 4, Path: "<path t='5' />"}, {ID: PATTERN_ID, Mode: "test", Colors: 3, Path: "<path />"}},
		},
		{
			name:     "Read file with pagination (outside borders) & query filter",
			query:    m.PatternQueryDto{Mode: "testxt", Page: 1},
			expected: []m.Pattern{},
		},
		{
			name:     "Read Pattern record by colors",
			query:    m.PatternQueryDto{Colors: 3},
			expected: []m.Pattern{{ID: PATTERN_ID, Mode: "test", Colors: 3, Path: "<path />"}},
		},
		{
			name:     "Read not existed values",
			query:    m.PatternQueryDto{Colors: 3, Mode: "test2"},
			expected: []m.Pattern{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := pattern.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Mode, models[i].Mode)
				require.Equal(t, expected.Colors, models[i].Colors)
				require.Equal(t, expected.Path, models[i].Path)
			}
		})
	}
}

func TestPatternUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.PatternQueryDto
		model m.Pattern
		err   bool
	}{
		{
			name:  "Update Pattern record by ID",
			query: m.PatternQueryDto{ID: PATTERN_ID},
			model: m.Pattern{Colors: 5},
		},
		{
			name:  "Update multiple Patterns record by query",
			query: m.PatternQueryDto{Mode: "test"},
			model: m.Pattern{Colors: 6},
		},
		{
			name:  "Update multiple Patterns record by query",
			query: m.PatternQueryDto{Mode: "test"},
			model: m.Pattern{Path: "<path />"},
			err:   true,
		},
		{
			name:  "Update not existed record",
			query: m.PatternQueryDto{Colors: 100},
			model: m.Pattern{},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := pattern.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tc.model.Mode != "" {
				require.Equal(t, tc.model.Mode, models[0].Mode)
			}

			if tc.model.Colors != 0 {
				require.Equal(t, tc.model.Colors, models[0].Colors)
			}

			if tc.model.Path != "" {
				require.Equal(t, tc.model.Path, models[0].Path)
			}
		})
	}
}

func TestPatternCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.PatternQueryDto
		expected []m.Pattern
	}{
		{
			name:     "Check if cached Pattern record by ID was updated",
			query:    m.PatternQueryDto{ID: PATTERN_ID},
			expected: []m.Pattern{{ID: PATTERN_ID, Mode: "test", Colors: 6, Path: "<path />"}},
		},
		{
			name:     "Check if cached Pattern record by mode & project_id was updated",
			query:    m.PatternQueryDto{Mode: "test"},
			expected: []m.Pattern{{Mode: "test", Colors: 6, Path: "<path t='5' />"}, {ID: PATTERN_ID, Mode: "test", Colors: 6, Path: "<path />"}},
		},
		{
			name:     "Check if cached Pattern record with pagination & project_id was updated",
			query:    m.PatternQueryDto{Mode: "test", Page: 0},
			expected: []m.Pattern{{Mode: "test", Colors: 6, Path: "<path t='5' />"}, {ID: PATTERN_ID, Mode: "test", Colors: 6, Path: "<path />"}},
		},
		{
			name:     "Check if cached Pattern record by colors was updated",
			query:    m.PatternQueryDto{Colors: 3},
			expected: []m.Pattern{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := pattern.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Mode, models[i].Mode)
				require.Equal(t, expected.Colors, models[i].Colors)
				require.Equal(t, expected.Path, models[i].Path)
			}
		})
	}

}

func TestPatternDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.PatternQueryDto
		expected int
	}{
		{
			name:     "Delete Pattern record by ID",
			query:    m.PatternQueryDto{ID: PATTERN_ID},
			expected: 1,
		},
		{
			name:     "Delete Pattern record by role & project_id",
			query:    m.PatternQueryDto{Mode: "test"},
			expected: 1,
		},
		{
			name:     "Delete Pattern record by colors",
			query:    m.PatternQueryDto{Colors: 3},
			expected: 0,
		},
		{
			name:     "Delete Pattern record by Color 2",
			query:    m.PatternQueryDto{Colors: 6},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := pattern.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestPatternFinalCacheState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.PatternQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by ID",
			query:    m.PatternQueryDto{ID: PATTERN_ID},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by role & project_id",
			query:    m.PatternQueryDto{Mode: "test"},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination & query filter",
			query:    m.PatternQueryDto{Mode: "test", Page: 0},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.PatternQueryDto{Mode: "test", Page: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by colors",
			query:    m.PatternQueryDto{Colors: 3},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by colors",
			query:    m.PatternQueryDto{Colors: 6},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := pattern.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, len(models))
		})
	}
}

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./", ""),
	}).Init()

	redis := client.ConnRedis()
	db := client.ConnDB([]interfaces.Table{
		m.NewPattern(),
	})

	pattern = service.NewPatternService(db, redis)
}
