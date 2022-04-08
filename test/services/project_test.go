package services_test

import (
	"api/client"
	"api/config"
	"api/interfaces"
	m "api/models"
	"api/service"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	project service.ProjectService

	createdTo   = time.Now().Add(time.Hour)
	createdFrom = time.Now().Add(-1 * time.Hour)
)

func TestProjectCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Project
		err   bool
	}{
		{
			name:  "Create new project record",
			model: m.Project{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"},
		},
		{
			name:  "Create new file record",
			model: m.Project{Name: "_test2", Title: "yes", Flag: "js", Desc: "yes", Note: "yes"},
		},
		{
			name:  "Check what will happend if create already existed project name",
			model: m.Project{Name: "_test"},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := project.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Title, model.Title)
			require.Equal(t, tc.model.Flag, model.Flag)
			require.Equal(t, tc.model.Desc, model.Desc)
			require.Equal(t, tc.model.Note, model.Note)
		})
	}
}

func TestProjectRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected []m.Project
	}{
		{
			name:     "Read project record by Name",
			query:    m.ProjectQueryDto{Name: "_test"},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project record by Name 2",
			query:    m.ProjectQueryDto{Name: "_test2"},
			expected: []m.Project{{Name: "_test2", Title: "yes", Flag: "js", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project record by role & project_id",
			query:    m.ProjectQueryDto{Flag: "aa", CreatedTo: createdTo, CreatedFrom: createdFrom},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project with pagination & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 0},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read file with pagination (outside borders) & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 1},
			expected: []m.Project{},
		},
		{
			name:     "Read not existed values",
			query:    m.ProjectQueryDto{Name: "___test"},
			expected: []m.Project{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Title, models[i].Title)
				require.Equal(t, expected.Flag, models[i].Flag)
				require.Equal(t, expected.Desc, models[i].Desc)
				require.Equal(t, expected.Note, models[i].Note)
			}
		})
	}
}

func TestProjectUpdateStat(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Project
		err   bool
	}{
		{
			name:  "Create new project record",
			model: m.Project{Name: "_test5", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := project.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Title, model.Title)
			require.Equal(t, tc.model.Flag, model.Flag)
			require.Equal(t, tc.model.Desc, model.Desc)
			require.Equal(t, tc.model.Note, model.Note)
		})
	}
}

func TestProjectCheckStat(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected []m.Project
	}{
		{
			name:     "Read project record by Name",
			query:    m.ProjectQueryDto{Name: "_test"},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project record by Name 2",
			query:    m.ProjectQueryDto{Name: "_test2"},
			expected: []m.Project{{Name: "_test2", Title: "yes", Flag: "js", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project record by role & project_id",
			query:    m.ProjectQueryDto{Flag: "aa", CreatedTo: createdTo, CreatedFrom: createdFrom},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}, {Name: "_test5", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project with pagination & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 0},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}, {Name: "_test5", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read file with pagination (outside borders) & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 1},
			expected: []m.Project{},
		},
		{
			name:     "Read project with pagination & limit",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 0, Limit: 1},
			expected: []m.Project{{Name: "_test5", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project with pagination & limit",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 1, Limit: 1},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Title, models[i].Title)
				require.Equal(t, expected.Flag, models[i].Flag)
				require.Equal(t, expected.Desc, models[i].Desc)
				require.Equal(t, expected.Note, models[i].Note)
			}
		})
	}
}

func TestProjectDeleteStat(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected int
	}{
		{
			name:     "Delete Project record by Name",
			query:    m.ProjectQueryDto{Name: "_test5"},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := project.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestProjectCheckDeleteStat(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected []m.Project
	}{
		{
			name:     "Read project with pagination & limit",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 0, Limit: 1},
			expected: []m.Project{{Name: "_test", Title: "yes", Flag: "aa", Desc: "yes", Note: "yes"}},
		},
		{
			name:     "Read project with pagination & limit",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 1, Limit: 1},
			expected: []m.Project{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Title, models[i].Title)
				require.Equal(t, expected.Flag, models[i].Flag)
				require.Equal(t, expected.Desc, models[i].Desc)
				require.Equal(t, expected.Note, models[i].Note)
			}
		})
	}
}

func TestProjectUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.ProjectQueryDto
		model m.Project
		err   bool
	}{
		{
			name:  "Update project record by Name",
			query: m.ProjectQueryDto{Name: "_test"},
			model: m.Project{Name: "_test", Title: "yes2", Flag: "__", Desc: "yes2", Note: "yes2"},
		},
		{
			name:  "Update project record by Name 2",
			query: m.ProjectQueryDto{Name: "_test2"},
			model: m.Project{Name: "_test2", Title: "yes2", Flag: "__", Desc: "yes2", Note: "yes2"},
		},
		{
			name:  "Update multiple projects record by query",
			query: m.ProjectQueryDto{Flag: "__"},
			model: m.Project{Title: "yes54"},
		},
		{
			name:  "Update project name (with the name that already exists)",
			query: m.ProjectQueryDto{Name: "_test2"},
			model: m.Project{Name: "_test", Title: "yes2", Flag: "aa", Desc: "yes2", Note: "yes2"},
			err:   true,
		},
		{
			name:  "Update project name",
			query: m.ProjectQueryDto{Name: "_test2"},
			model: m.Project{Name: "_test3"},
		},
		{
			name:  "Update not existed record",
			query: m.ProjectQueryDto{Name: "___test"},
			model: m.Project{Name: "_test", Title: "yes2", Flag: "__", Desc: "yes2", Note: "yes2"},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tc.model.Name != "" {
				require.Equal(t, tc.model.Name, models[0].Name)
			}

			if tc.model.Title != "" {
				require.Equal(t, tc.model.Title, models[0].Title)
			}

			if tc.model.Flag != "" {
				require.Equal(t, tc.model.Flag, models[0].Flag)
			}

			if tc.model.Desc != "" {
				require.Equal(t, tc.model.Desc, models[0].Desc)
			}

			if tc.model.Note != "" {
				require.Equal(t, tc.model.Note, models[0].Note)
			}
		})
	}
}

func TestProjectCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected []m.Project
	}{
		{
			name:     "Check if cached Project record by Name was updated",
			query:    m.ProjectQueryDto{Name: "_test"},
			expected: []m.Project{{Name: "_test", Title: "yes54", Flag: "__", Desc: "yes2", Note: "yes2"}},
		},
		{
			name:     "Check if cached Project record by Name 2 should be empty",
			query:    m.ProjectQueryDto{Name: "_test2"},
			expected: []m.Project{},
		},
		{
			name:     "Check if cached Project record by flag & project_id was updated",
			query:    m.ProjectQueryDto{Flag: "__", CreatedTo: createdTo, CreatedFrom: createdFrom},
			expected: []m.Project{{Name: "_test3", Title: "yes54", Flag: "__", Desc: "yes2", Note: "yes2"}, {Name: "_test", Title: "yes54", Flag: "__", Desc: "yes2", Note: "yes2"}},
		},
		{
			name:     "Read project with pagination & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 0},
			expected: []m.Project{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Title, models[i].Title)
				require.Equal(t, expected.Flag, models[i].Flag)
				require.Equal(t, expected.Desc, models[i].Desc)
				require.Equal(t, expected.Note, models[i].Note)
			}
		})
	}

}

func TestProjectDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected int
	}{
		{
			name:     "Delete Project record by Name",
			query:    m.ProjectQueryDto{Name: "_test"},
			expected: 1,
		},
		{
			name:     "Delete Project record by role & project_id",
			query:    m.ProjectQueryDto{Flag: "__", CreatedTo: createdTo, CreatedFrom: createdFrom},
			expected: 1,
		},
		{
			name:     "Delete already deleted record by name",
			query:    m.ProjectQueryDto{Name: "_test3"},
			expected: 0,
		},
		{
			name:     "Delete already deleted record by name 2",
			query:    m.ProjectQueryDto{Name: "_test2"},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := project.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestProjectState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.ProjectQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by Name",
			query:    m.ProjectQueryDto{Name: "_test"},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by Name 2",
			query:    m.ProjectQueryDto{Name: "_test2"},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by Name 3",
			query:    m.ProjectQueryDto{Name: "_test3"},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by flag & time",
			query:    m.ProjectQueryDto{Flag: "__", CreatedTo: createdTo, CreatedFrom: createdFrom},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.ProjectQueryDto{Flag: "aa", Page: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := project.Read(&tc.query)
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
		m.NewProject(),
	})

	project = *service.NewProjectService(db, redis)
}
