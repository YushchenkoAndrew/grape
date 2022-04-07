package services_test

import (
	"api/client"
	"api/config"
	"api/interfaces"
	m "api/models"
	"api/service"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	subscription service.SubscriptionService

	CRON_ID  = "3073ca3c-cd0a-4036-9e5d-7bc32d205d1d"
	CRON_ID2 = "b16b056a-ccfd-41f0-a2ab-2218fe1d6950"

	TOKEN = "b16b056a-ccfd-41f0-a2ab-2218fe1d6950"
)

func TestSubscriptionCreate(t *testing.T) {
	var tests = []struct {
		name  string
		model m.Subscription
		err   bool
	}{
		{
			name:  "Create new subscription record",
			model: m.Subscription{Name: "main", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1},
		},
		{
			name:  "Create new subscription record",
			model: m.Subscription{Name: "main2", CronID: CRON_ID2, CronTime: "* * * * *", Method: "post", Path: "/", Token: CRON_ID, ProjectID: 1},
		},
		{
			name:  "Update already existed link record with 'Create' func",
			model: m.Subscription{Name: "main3", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1},
		},
		{
			name:  "Check on project_id existens",
			model: m.Subscription{Name: "main", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 10000},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var model = tc.model.Copy()
			err := subscription.Create(model)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, 0, model.ID)
			require.Equal(t, tc.model.Name, model.Name)
			require.Equal(t, tc.model.CronID, model.CronID)
			require.Equal(t, tc.model.CronTime, model.CronTime)
			require.Equal(t, tc.model.Method, model.Method)
			require.Equal(t, tc.model.Path, model.Path)
			require.Equal(t, tc.model.Token, model.Token)
		})
	}
}

func TestSubscriptionRead(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.SubscribeQueryDto
		expected []m.Subscription
	}{
		{
			name:     "Read subscription record by cron_id",
			query:    m.SubscribeQueryDto{CronID: CRON_ID},
			expected: []m.Subscription{{Name: "main3", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Read subscription record by name & project_id",
			query:    m.SubscribeQueryDto{Name: "main3", ProjectID: 1},
			expected: []m.Subscription{{Name: "main3", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Read subscription with pagination & query filter",
			query:    m.SubscribeQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Subscription{{Name: "main2", CronID: CRON_ID2, CronTime: "* * * * *", Method: "post", Path: "/", Token: CRON_ID, ProjectID: 1}, {Name: "main3", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Read subscription with pagination (outside borders) & query filter",
			query:    m.SubscribeQueryDto{Name: "main2", ProjectID: 1, Page: 1},
			expected: []m.Subscription{},
		},
		{
			name:     "Read subscription record by project_id",
			query:    m.SubscribeQueryDto{ProjectID: 1},
			expected: []m.Subscription{{Name: "main2", CronID: CRON_ID2, CronTime: "* * * * *", Method: "post", Path: "/", Token: CRON_ID, ProjectID: 1}, {Name: "main3", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Read not existed values",
			query:    m.SubscribeQueryDto{Name: "main2", ProjectID: 10000},
			expected: []m.Subscription{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := subscription.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.CronID, models[i].CronID)
				require.Equal(t, expected.CronTime, models[i].CronTime)
				require.Equal(t, expected.Method, models[i].Method)
				require.Equal(t, expected.Path, models[i].Path)
				require.Equal(t, expected.Token, models[i].Token)
			}
		})
	}
}

func TestSubscriptionUpdate(t *testing.T) {
	var tests = []struct {
		name  string
		query m.SubscribeQueryDto
		model m.Subscription
		err   bool
	}{
		{
			name:  "Update subscription record by cron_id",
			query: m.SubscribeQueryDto{CronID: CRON_ID},
			model: m.Subscription{Name: "main2"},
		},
		{
			name:  "Update multiple records at once",
			query: m.SubscribeQueryDto{Name: "main2", ProjectID: 1},
			model: m.Subscription{Name: "main4"},
		},
		{
			name:  "Update not existed record",
			query: m.SubscribeQueryDto{ProjectID: 10000},
			model: m.Subscription{Name: "main"},
			err:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := subscription.Update(&tc.query, tc.model.Copy())
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tc.model.Name != "" {
				require.Equal(t, tc.model.Name, models[0].Name)
			}
		})
	}
}

func TestSubscriptionCheckCache(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.SubscribeQueryDto
		expected []m.Subscription
	}{
		{
			name:     "Check if cached Subscription record by cron_id was updated",
			query:    m.SubscribeQueryDto{CronID: CRON_ID},
			expected: []m.Subscription{{Name: "main4", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Check if cached subscription record by project_id was updated",
			query:    m.SubscribeQueryDto{ProjectID: 1},
			expected: []m.Subscription{{Name: "main4", CronID: CRON_ID2, CronTime: "* * * * *", Method: "post", Path: "/", Token: CRON_ID, ProjectID: 1}, {Name: "main4", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
		{
			name:     "Check if cached subscription record with pagination & project_id was updated",
			query:    m.SubscribeQueryDto{ProjectID: 1, Page: 0},
			expected: []m.Subscription{{Name: "main4", CronID: CRON_ID2, CronTime: "* * * * *", Method: "post", Path: "/", Token: CRON_ID, ProjectID: 1}, {Name: "main4", CronID: CRON_ID, CronTime: "* * * * *", Method: "post", Path: "/", Token: TOKEN, ProjectID: 1}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := subscription.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(models))

			for i, expected := range tc.expected {
				require.Equal(t, expected.Name, models[i].Name)
				require.Equal(t, expected.CronID, models[i].CronID)
				require.Equal(t, expected.CronTime, models[i].CronTime)
				require.Equal(t, expected.Method, models[i].Method)
				require.Equal(t, expected.Path, models[i].Path)
				require.Equal(t, expected.Token, models[i].Token)
			}
		})
	}

}

func TestSubscriptionDelete(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.SubscribeQueryDto
		expected int
	}{
		{
			name:     "Delete subscription record by ID",
			query:    m.SubscribeQueryDto{CronID: CRON_ID},
			expected: 1,
		},
		{
			name:     "Delete subscription record by project_id",
			query:    m.SubscribeQueryDto{ProjectID: 1},
			expected: 1,
		},
		{
			name:     "Delete subscription record by name & project_id",
			query:    m.SubscribeQueryDto{Name: "main4", ProjectID: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := subscription.Delete(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, items)
		})
	}
}

func TestSubscriptionFinalCacheState(t *testing.T) {
	var tests = []struct {
		name     string
		query    m.SubscribeQueryDto
		expected int
	}{
		{
			name:     "Check if values was deleted correctly by ID",
			query:    m.SubscribeQueryDto{CronID: CRON_ID},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by name & project_id",
			query:    m.SubscribeQueryDto{Name: "main4", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination & query filter",
			query:    m.SubscribeQueryDto{ProjectID: 1, Page: 0},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly 2",
			query:    m.SubscribeQueryDto{ProjectID: 1, Page: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly with pagination (outside border) & query filter",
			query:    m.SubscribeQueryDto{Name: "main", ProjectID: 1},
			expected: 0,
		},
		{
			name:     "Check if values was deleted correctly by project_id",
			query:    m.SubscribeQueryDto{ProjectID: 1},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			models, err := subscription.Read(&tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.expected, len(models))
		})
	}
}

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./", ""),
	}).Init()

	db, client, _, _ := client.Init([]interfaces.Table{
		m.NewSubscription(),
		m.NewProject(),
	})

	subscription = *service.NewSubscriptionService(db, client)

	var project = *service.NewProjectService(db, client)
	project.Create(&m.Project{ID: 1, Name: "yes", Title: "js", Flag: "js", Desc: "js", Note: "js"})
}
