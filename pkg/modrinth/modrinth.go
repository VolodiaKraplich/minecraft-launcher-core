package modrinth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ModrinthClientOption is a functional option for configuring the client.
type ModrinthClientOption func(*ModrinthV2Client)

// WithBaseURL sets the base URL for the client.
func WithBaseURL(baseURL string) ModrinthClientOption {
	return func(c *ModrinthV2Client) {
		c.baseURL = baseURL
	}
}

// WithHeaders sets additional headers for the client.
func WithHeaders(headers map[string]string) ModrinthClientOption {
	return func(c *ModrinthV2Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ModrinthClientOption {
	return func(c *ModrinthV2Client) {
		c.httpClient = client
	}
}

// NewModrinthV2Client creates a new Modrinth V2 API client.
func NewModrinthV2Client(options ...ModrinthClientOption) *ModrinthV2Client {
	c := &ModrinthV2Client{
		baseURL:    "https://api.modrinth.com",
		headers:    make(map[string]string),
		httpClient: http.DefaultClient,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *ModrinthV2Client) request(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Response, error) {
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.httpClient.Do(req)
}

func (c *ModrinthV2Client) doJSON(ctx context.Context, method, path string, body any, result any) error {
	var r io.Reader
	var ct string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
		ct = "application/json"
	}
	resp, err := c.request(ctx, method, path, r, ct)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &ModrinthAPIError{
			URL:    resp.Request.URL.String(),
			Status: resp.StatusCode,
			Body:   string(b),
		}
	}
	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// SearchProjects searches for projects on Modrinth.
func (c *ModrinthV2Client) SearchProjects(ctx context.Context, options SearchProjectOptions) (*SearchResult, error) {
	v := url.Values{}
	v.Add("query", options.Query)
	v.Add("filter", options.Filter)
	index := options.Index
	if index == "" {
		if options.Query != "" {
			index = "relevance"
		} else {
			index = "downloads"
		}
	}
	v.Add("index", index)
	v.Add("offset", fmt.Sprint(options.Offset))
	v.Add("limit", fmt.Sprint(options.Limit))
	if options.Facets != "" {
		v.Add("facets", options.Facets)
	}
	path := "/v2/search?" + v.Encode()
	var result SearchResult
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return &result, err
}

// GetProject fetches a single project by ID.
func (c *ModrinthV2Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	projectID = strings.TrimPrefix(projectID, "local-")
	path := "/v2/project/" + projectID
	var project Project
	err := c.doJSON(ctx, http.MethodGet, path, nil, &project)
	return &project, err
}

// GetProjects fetches multiple projects by IDs.
func (c *ModrinthV2Client) GetProjects(ctx context.Context, projectIDs []string) ([]Project, error) {
	v := url.Values{}
	b, _ := json.Marshal(projectIDs)
	v.Add("ids", string(b))
	path := "/v2/projects?" + v.Encode()
	var projects []Project
	err := c.doJSON(ctx, http.MethodGet, path, nil, &projects)
	return projects, err
}

// GetProjectVersions fetches versions for a project.
func (c *ModrinthV2Client) GetProjectVersions(ctx context.Context, projectID string, options GetProjectVersionsOptions) ([]ProjectVersion, error) {
	v := url.Values{}
	if len(options.Loaders) > 0 {
		b, err := json.Marshal(options.Loaders)
		if err != nil {
			return nil, err
		}
		v.Add("loaders", string(b))
	}
	if len(options.GameVersions) > 0 {
		b, err := json.Marshal(options.GameVersions)
		if err != nil {
			return nil, err
		}
		v.Add("game_versions", string(b))
	}
	if options.Featured != nil {
		v.Add("featured", fmt.Sprintf("%t", *options.Featured))
	}
	path := "/v2/project/" + projectID + "/version?" + v.Encode()
	var versions []ProjectVersion
	err := c.doJSON(ctx, http.MethodGet, path, nil, &versions)
	return versions, err
}

// GetProjectVersion fetches a single project version by ID.
func (c *ModrinthV2Client) GetProjectVersion(ctx context.Context, versionID string) (*ProjectVersion, error) {
	path := "/v2/version/" + versionID
	var version ProjectVersion
	err := c.doJSON(ctx, http.MethodGet, path, nil, &version)
	return &version, err
}

// GetProjectVersionsByID fetches multiple project versions by IDs.
func (c *ModrinthV2Client) GetProjectVersionsByID(ctx context.Context, ids []string) ([]ProjectVersion, error) {
	v := url.Values{}
	b, _ := json.Marshal(ids)
	v.Add("ids", string(b))
	path := "/v2/versions?" + v.Encode()
	var versions []ProjectVersion
	err := c.doJSON(ctx, http.MethodGet, path, nil, &versions)
	return versions, err
}

// GetProjectVersionsByHash fetches project versions by file hashes.
func (c *ModrinthV2Client) GetProjectVersionsByHash(ctx context.Context, hashes []string, algorithm string) (map[string]ProjectVersion, error) {
	if algorithm == "" {
		algorithm = "sha1"
	}
	body := map[string]any{"hashes": hashes, "algorithm": algorithm}
	path := "/v2/version_files"
	var versions map[string]ProjectVersion
	err := c.doJSON(ctx, http.MethodPost, path, body, &versions)
	return versions, err
}

// GetLatestVersionsFromHashes fetches the latest versions matching the hashes and filters.
func (c *ModrinthV2Client) GetLatestVersionsFromHashes(ctx context.Context, hashes []string, algorithm string, loaders []string, gameVersions []string) (map[string]ProjectVersion, error) {
	body := map[string]any{
		"hashes":        hashes,
		"algorithm":     algorithm,
		"loaders":       loaders,
		"game_versions": gameVersions,
	}
	path := "/v2/version_files/update"
	var versions map[string]ProjectVersion
	err := c.doJSON(ctx, http.MethodPost, path, body, &versions)
	return versions, err
}

// GetLatestProjectVersion fetches the latest version matching the hash and filters.
func (c *ModrinthV2Client) GetLatestProjectVersion(ctx context.Context, sha1 string, algorithm string, loaders []string, gameVersions []string) (*ProjectVersion, error) {
	if algorithm == "" {
		algorithm = "sha1"
	}
	v := url.Values{}
	v.Add("algorithm", algorithm)
	path := "/v2/version_file/" + sha1 + "/update?" + v.Encode()
	body := map[string]any{
		"loaders":       loaders,
		"game_versions": gameVersions,
	}
	var version ProjectVersion
	err := c.doJSON(ctx, http.MethodPost, path, body, &version)
	return &version, err
}

// GetLicenseTags fetches available license tags.
func (c *ModrinthV2Client) GetLicenseTags(ctx context.Context) ([]License, error) {
	path := "/v2/tag/license"
	var licenses []License
	err := c.doJSON(ctx, http.MethodGet, path, nil, &licenses)
	return licenses, err
}

// GetCategoryTags fetches available category tags.
func (c *ModrinthV2Client) GetCategoryTags(ctx context.Context) ([]Category, error) {
	path := "/v2/tag/category"
	var categories []Category
	err := c.doJSON(ctx, http.MethodGet, path, nil, &categories)
	return categories, err
}

// GetGameVersionTags fetches available game version tags.
func (c *ModrinthV2Client) GetGameVersionTags(ctx context.Context) ([]GameVersion, error) {
	path := "/v2/tag/game_version"
	var versions []GameVersion
	err := c.doJSON(ctx, http.MethodGet, path, nil, &versions)
	return versions, err
}

// GetLoaderTags fetches available loader tags.
func (c *ModrinthV2Client) GetLoaderTags(ctx context.Context) ([]Loader, error) {
	path := "/v2/tag/loader"
	var loaders []Loader
	err := c.doJSON(ctx, http.MethodGet, path, nil, &loaders)
	return loaders, err
}

// GetCollections fetches collections for a user.
func (c *ModrinthV2Client) GetCollections(ctx context.Context, userID string) ([]Collection, error) {
	path := "/v3/user/" + userID + "/collections"
	var collections []Collection
	err := c.doJSON(ctx, http.MethodGet, path, nil, &collections)
	return collections, err
}

// UpdateCollectionIcon updates the icon for a collection.
func (c *ModrinthV2Client) UpdateCollectionIcon(ctx context.Context, collectionID string, iconData []byte, mimeType string) error {
	extParts := strings.Split(mimeType, "/")
	if len(extParts) < 2 {
		return errors.New("invalid mime type")
	}
	ext := extParts[1]
	path := "/v3/collection/" + collectionID + "/icon?ext=" + ext
	resp, err := c.request(ctx, http.MethodPatch, path, bytes.NewReader(iconData), mimeType)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &ModrinthAPIError{
			URL:    resp.Request.URL.String(),
			Status: resp.StatusCode,
			Body:   string(b),
		}
	}
	return nil
}

// CreateCollection creates a new collection.
func (c *ModrinthV2Client) CreateCollection(ctx context.Context, name, description string, projectIDs []string) (*Collection, error) {
	body := map[string]any{
		"name":        name,
		"description": description,
		"projects":    projectIDs,
	}
	path := "/v3/collection"
	var collection Collection
	err := c.doJSON(ctx, http.MethodPost, path, body, &collection)
	return &collection, err
}

// UpdateCollection updates a collection with new projects.
func (c *ModrinthV2Client) UpdateCollection(ctx context.Context, collectionID string, projectIDs []string) error {
	body := map[string][]string{"new_projects": projectIDs}
	path := "/v3/collection/" + collectionID
	return c.doJSON(ctx, http.MethodPatch, path, body, nil)
}

// GetAuthenticatedUser fetches the authenticated user.
func (c *ModrinthV2Client) GetAuthenticatedUser(ctx context.Context) (*User, error) {
	path := "/v2/user"
	var user User
	err := c.doJSON(ctx, http.MethodGet, path, nil, &user)
	return &user, err
}

// FollowProject follows a project.
func (c *ModrinthV2Client) FollowProject(ctx context.Context, id string) error {
	path := "/v2/project/" + id + "/follow"
	return c.doJSON(ctx, http.MethodPost, path, nil, nil)
}

// UnfollowProject unfollows a project.
func (c *ModrinthV2Client) UnfollowProject(ctx context.Context, id string) error {
	path := "/v2/project/" + id + "/follow"
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}

// GetUserFollowedProjects fetches projects followed by a user.
func (c *ModrinthV2Client) GetUserFollowedProjects(ctx context.Context, userID string) ([]Project, error) {
	path := "/v2/user/" + userID + "/follows"
	var projects []Project
	err := c.doJSON(ctx, http.MethodGet, path, nil, &projects)
	return projects, err
}

// GetProjectTeamMembers fetches team members for a project.
func (c *ModrinthV2Client) GetProjectTeamMembers(ctx context.Context, projectID string) ([]TeamMember, error) {
	path := "/v2/project/" + projectID + "/members"
	var members []TeamMember
	err := c.doJSON(ctx, http.MethodGet, path, nil, &members)
	return members, err
}

// GetUser fetches a user by ID.
func (c *ModrinthV2Client) GetUser(ctx context.Context, id string) (*User, error) {
	path := "/v2/user/" + id
	var user User
	err := c.doJSON(ctx, http.MethodGet, path, nil, &user)
	return &user, err
}

// GetUserProjects fetches projects for a user.
func (c *ModrinthV2Client) GetUserProjects(ctx context.Context, id string) ([]Project, error) {
	path := "/v2/user/" + id + "/projects"
	var projects []Project
	err := c.doJSON(ctx, http.MethodGet, path, nil, &projects)
	return projects, err
}
