package modrinth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Voxelum/minecraft-launcher-core/pkg/modrinth"
)

func setupTestServer(handler func(w http.ResponseWriter, r *http.Request)) (*modrinth.ModrinthV2Client, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	client := modrinth.NewModrinthV2Client(
		modrinth.WithBaseURL(server.URL),
		modrinth.WithHTTPClient(http.DefaultClient),
	)
	return client, server
}

func TestSearchProjects(t *testing.T) {
	tests := []struct {
		name     string
		options  modrinth.SearchProjectOptions
		expected *modrinth.SearchResult
		status   int
		body     string
		wantErr  bool
	}{
		{
			name: "Successful search",
			options: modrinth.SearchProjectOptions{
				Query:  "fabric",
				Limit:  5,
				Offset: 0,
			},
			expected: &modrinth.SearchResult{
				Hits: []modrinth.SearchResultHit{
					{
						Slug:        "fabric-api",
						ProjectID:   "P7dR8mSH",
						ProjectType: "mod",
						Title:       "Fabric API",
					},
				},
				Offset:    0,
				Limit:     5,
				TotalHits: 1,
			},
			status:  200,
			body:    `{"hits":[{"slug":"fabric-api","project_id":"P7dR8mSH","project_type":"mod","title":"Fabric API"}],"offset":0,"limit":5,"total_hits":1}`,
			wantErr: false,
		},
		{
			name:    "API error",
			options: modrinth.SearchProjectOptions{},
			status:  404,
			body:    "Not found",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			})
			defer server.Close()

			result, err := client.SearchProjects(context.Background(), tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SearchProjects() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		expected  *modrinth.Project
		status    int
		body      string
		wantErr   bool
	}{
		{
			name:      "Successful get project",
			projectID: "P7dR8mSH",
			expected: &modrinth.Project{
				ID:          "P7dR8mSH",
				Slug:        "fabric-api",
				ProjectType: "mod",
				Title:       "Fabric API",
			},
			status:  200,
			body:    `{"id":"P7dR8mSH","slug":"fabric-api","project_type":"mod","title":"Fabric API"}`,
			wantErr: false,
		},
		{
			name:      "Project ID with local- prefix",
			projectID: "local-P7dR8mSH",
			expected: &modrinth.Project{
				ID:          "P7dR8mSH",
				Slug:        "fabric-api",
				ProjectType: "mod",
				Title:       "Fabric API",
			},
			status:  200,
			body:    `{"id":"P7dR8mSH","slug":"fabric-api","project_type":"mod","title":"Fabric API"}`,
			wantErr: false,
		},
		{
			name:      "API error",
			projectID: "invalid",
			status:    404,
			body:      "Project not found",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			})
			defer server.Close()

			result, err := client.GetProject(context.Background(), tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetProject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	tests := []struct {
		name       string
		projectIDs []string
		expected   []modrinth.Project
		status     int
		body       string
		wantErr    bool
	}{
		{
			name:       "Successful get multiple projects",
			projectIDs: []string{"P7dR8mSH", "another"},
			expected: []modrinth.Project{
				{ID: "P7dR8mSH", Slug: "fabric-api"},
				{ID: "another", Slug: "other-mod"},
			},
			status:  200,
			body:    `[{"id":"P7dR8mSH","slug":"fabric-api"},{"id":"another","slug":"other-mod"}]`,
			wantErr: false,
		},
		{
			name:       "API error",
			projectIDs: []string{"invalid"},
			status:     400,
			body:       "Bad request",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			})
			defer server.Close()

			result, err := client.GetProjects(context.Background(), tt.projectIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetProjects() = %v, want %v", result, tt.expected)
			}
		})
	}
}
