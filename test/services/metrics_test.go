package services_test

import (
	"api/client"
	"api/config"
	"api/interfaces"
	m "api/models"
	"api/service"
	"api/service/k3s/pods"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const METRICS_ID = 305

var (
	metrics pods.MetricsService

	m_createdTo   = time.Now().Add(time.Hour)
	m_createdFrom = time.Now().Add(-1 * time.Hour)
)

func TestMetricsCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Metrics
		err   bool
	}{
		{
			name:  "Create new Metrics record",
			model: m.Metrics{ID: METRICS_ID, Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1},
		},
		{
			name:  "Create new Metrics record",
			model: m.Metrics{Name: "test", Namespace: "test2", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1},
		},
		{
			name:  "Check what will happend if create already existed Metrics",
			model: m.Metrics{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1},
		},
		{
			name:  "Check on project_id existens",
			model: m.Metrics{Name: "test", Namespace: "test", ProjectID: 10000},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := metrics.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.Namespace, model.Namespace)
			require.Equal(t, tc.model.ContainerName, model.ContainerName)
			require.Equal(t, tc.model.CPU, model.CPU)
			require.Equal(t, tc.model.Memory, model.Memory)
		})
	}
}

func TestMetricsRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.MetricsQueryDto
		expected []m.Metrics
	}{
		{
			name:     "Read metrics record by ID",
			query:    m.MetricsQueryDto{ID: METRICS_ID},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}},
		},
		{
			name:     "Read metrics record by name & namespace",
			query:    m.MetricsQueryDto{Name: "test", Namespace: "test"},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}},
		},
		{
			name:     "Read project record by role & project_id",
			query:    m.MetricsQueryDto{Name: "test", CreatedTo: m_createdTo, CreatedFrom: m_createdFrom},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}},
		},
		{
			name:     "Read metrics with pagination & query filter",
			query:    m.MetricsQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}},
		},
		{
			name:     "Read file metrics pagination (outside borders) & query filter",
			query:    m.MetricsQueryDto{Name: "test", ProjectID: 1, Page: 1},
			expected: []m.Metrics{},
		},
		{
			name:     "Read metrics record by project_id",
			query:    m.MetricsQueryDto{ProjectID: 1},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 1, Memory: 2, ProjectID: 1}},
		},
		{
			name:     "Read not existed values",
			query:    m.MetricsQueryDto{Name: "test", ProjectID: 10000},
			expected: []m.Metrics{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := metrics.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Namespace, models[i].Namespace)
				require.Equal(t, expected.ContainerName, models[i].ContainerName)
				require.Equal(t, expected.CPU, models[i].CPU)
				require.Equal(t, expected.Memory, models[i].Memory)
			}
		})
	}
}

func TestMetricsUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.MetricsQueryDto
		model m.Metrics
		err   bool
	}{
		{
			name:  "Update metrics record by ID",
			query: m.MetricsQueryDto{ID: METRICS_ID},
			model: m.Metrics{CPU: 2, Memory: 3},
		},
		{
			name:  "Update multiple metrics record by query",
			query: m.MetricsQueryDto{Namespace: "test2", ProjectID: 1},
			model: m.Metrics{CPU: 1, Memory: 5},
		},
		{
			name:  "Update multiple metrics record by query",
			query: m.MetricsQueryDto{ProjectID: 1},
			model: m.Metrics{CPU: 3, Memory: 3},
		},
		{
			name:  "Update not existed record",
			query: m.MetricsQueryDto{Name: "test", ProjectID: 10000},
			model: m.Metrics{CPU: 1, Memory: 1},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := metrics.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.model.CPU, models[0].CPU)
			require.Equal(t, tc.model.Memory, models[0].Memory)
		})
	}
}

func TestMetricsCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.MetricsQueryDto
		expected []m.Metrics
	}{
		{
			name:     "Check if cached metrics record by ID was updated",
			query:    m.MetricsQueryDto{ID: METRICS_ID},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}},
		},
		{
			name:     "Read metrics record by name & namespace",
			query:    m.MetricsQueryDto{Name: "test", Namespace: "test"},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}},
		},
		{
			name:     "Read project record by role & project_id",
			query:    m.MetricsQueryDto{Name: "test", CreatedTo: m_createdTo, CreatedFrom: m_createdFrom},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}},
		},
		{
			name:     "Read metrics with pagination & query filter",
			query:    m.MetricsQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}},
		},
		{
			name:     "Read metrics with pagination & query filter",
			query:    m.MetricsQueryDto{ProjectID: 1},
			expected: []m.Metrics{{Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test2", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}, {Name: "test", Namespace: "test", ContainerName: "void", CPU: 3, Memory: 3, ProjectID: 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := metrics.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.Namespace, models[i].Namespace)
				require.Equal(t, expected.ContainerName, models[i].ContainerName)
				require.Equal(t, expected.CPU, models[i].CPU)
				require.Equal(t, expected.Memory, models[i].Memory)
			}
		})
	}

}

func TestMetricsDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.MetricsQueryDto
		expected int
	}{
		{
			name:     "Delete Metrics record by ID",
			query:    m.MetricsQueryDto{ID: METRICS_ID},
			expected: 1,
		},
		{
			name:     "Delete Metrics record by role & project_id",
			query:    m.MetricsQueryDto{Namespace: "test", ProjectID: 1},
			expected: 1,
		},
		{
			name:     "Delete Metrics record by project_id",
			query:    m.MetricsQueryDto{ProjectID: 1},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := metrics.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestMetricsState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.MetricsQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by ID",
			query:    m.MetricsQueryDto{ID: METRICS_ID},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by role & project_id",
			query:    m.MetricsQueryDto{Name: "test", Namespace: "test"},
			expected: 0,
		},
		{
			name:     "Read project record by role & project_id",
			query:    m.MetricsQueryDto{Name: "test", CreatedTo: m_createdTo, CreatedFrom: m_createdFrom},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination & query filter",
			query:    m.MetricsQueryDto{ProjectID: 1, Page: 0},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.MetricsQueryDto{Name: "test", ProjectID: 1, Page: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by project_id",
			query:    m.MetricsQueryDto{ProjectID: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := metrics.Read(&tc.query)
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
		m.NewMetrics(),
		m.NewProject(),
	})

	metrics = *pods.NewMetricsService(db, redis)

	var project = *service.NewProjectService(db, redis)
	project.Create(&m.Project{ID: 1, Name: "yes", Title: "js", Flag: "js", Desc: "js", Note: "js"})
}
