package curseforge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// CurseforgeAPIError represents an error from the Curseforge API.
type CurseforgeAPIError struct {
	URL    string
	Status int
	Body   string
}

func (e *CurseforgeAPIError) Error() string {
	return fmt.Sprintf("Fail to fetch curseforge api %s. Status=%d. %s", e.URL, e.Status, e.Body)
}

// ModAsset represents a mod asset such as a logo or screenshot.
type ModAsset struct {
	ID           int    `json:"id"`
	ModID        int    `json:"modId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	URL          string `json:"url"`
}

// ModStatus represents the status of a mod.
type ModStatus int

const (
	ModStatusNew             ModStatus = 1
	ModStatusChangesRequired ModStatus = 2
	ModStatusUnderSoftReview ModStatus = 3
	ModStatusApproved        ModStatus = 4
	ModStatusRejected        ModStatus = 5
	ModStatusChangesMade     ModStatus = 6
	ModStatusInactive        ModStatus = 7
	ModStatusAbandoned       ModStatus = 8
	ModStatusDeleted         ModStatus = 9
	ModStatusUnderReview     ModStatus = 10
)

// FileReleaseType represents the release type of a file.
type FileReleaseType int

const (
	FileReleaseTypeRelease FileReleaseType = 1
	FileReleaseTypeBeta    FileReleaseType = 2
	FileReleaseTypeAlpha   FileReleaseType = 3
)

// FileModLoaderType represents the mod loader type.
type FileModLoaderType int

const (
	FileModLoaderTypeAny        FileModLoaderType = 0
	FileModLoaderTypeForge      FileModLoaderType = 1
	FileModLoaderTypeCauldron   FileModLoaderType = 2
	FileModLoaderTypeLiteLoader FileModLoaderType = 3
	FileModLoaderTypeFabric     FileModLoaderType = 4
	FileModLoaderTypeQuilt      FileModLoaderType = 5
	FileModLoaderTypeNeoForge   FileModLoaderType = 6
)

// FileIndex represents metadata for a file's game version.
type FileIndex struct {
	GameVersion       string            `json:"gameVersion"`
	FileID            int               `json:"fileId"`
	Filename          string            `json:"filename"`
	ReleaseType       FileReleaseType   `json:"releaseType"`
	GameVersionTypeID *int              `json:"gameVersionTypeId"`
	ModLoader         FileModLoaderType `json:"modLoader"`
}

// Mod represents a Curseforge mod.
type Mod struct {
	ID                   int           `json:"id"`
	GameID               int           `json:"gameId"`
	Name                 string        `json:"name"`
	Slug                 string        `json:"slug"`
	Links                ModLinks      `json:"links"`
	Summary              string        `json:"summary"`
	Status               ModStatus     `json:"status"`
	DownloadCount        int           `json:"downloadCount"`
	IsFeatured           bool          `json:"isFeatured"`
	PrimaryCategoryID    int           `json:"primaryCategoryId"`
	Categories           []ModCategory `json:"categories"`
	ClassID              *int          `json:"classId"`
	Authors              []Author      `json:"authors"`
	Logo                 ModAsset      `json:"logo"`
	Screenshots          []ModAsset    `json:"screenshots"`
	MainFileID           int           `json:"mainFileId"`
	LatestFiles          []File        `json:"latestFiles"`
	LatestFilesIndexes   []FileIndex   `json:"latestFilesIndexes"`
	DateCreated          string        `json:"dateCreated"`
	DateModified         string        `json:"dateModified"`
	DateReleased         string        `json:"dateReleased"`
	AllowModDistribution *bool         `json:"allowModDistribution"`
	GamePopularityRank   int           `json:"gamePopularityRank"`
	IsAvailable          bool          `json:"isAvailable"`
	DefaultFileID        int           `json:"defaultFileId"`
	ThumbsUpCount        int           `json:"thumbsUpCount"`
}

// ModLinks represents links related to a mod.
type ModLinks struct {
	WebsiteURL string `json:"websiteUrl"`
	WikiURL    string `json:"wikiUrl"`
	IssuesURL  string `json:"issuesUrl"`
	SourceURL  string `json:"sourceUrl"`
}

// GameVersionLatestFile represents a latest file for a game version.
type GameVersionLatestFile struct {
	GameVersion     string `json:"gameVersion"`
	ProjectFileID   int    `json:"projectFileId"`
	ProjectFileName string `json:"projectFileName"`
	FileType        int    `json:"fileType"`
}

// CategorySection represents a category section.
type CategorySection struct {
	ID                      int     `json:"id"`
	GameID                  int     `json:"gameId"`
	Name                    string  `json:"name"`
	PackageType             int     `json:"packageType"`
	Path                    string  `json:"path"`
	InitialInclusionPattern string  `json:"initialInclusionPattern"`
	ExtraIncludePattern     *string `json:"extraIncludePattern"`
	GameCategoryID          int     `json:"gameCategoryId"`
}

// HashAlgo represents the hash algorithm used.
type HashAlgo int

const (
	HashAlgoSha1 HashAlgo = 1
	HashAlgoMd5  HashAlgo = 2
)

// FileHash represents a file hash.
type FileHash struct {
	Algo  HashAlgo `json:"algo"`
	Value string   `json:"value"`
}

// FileStatus represents the status of a file.
type FileStatus int

const (
	FileStatusProcessing         FileStatus = 1
	FileStatusChangesRequired    FileStatus = 2
	FileStatusUnderReview        FileStatus = 3
	FileStatusApproved           FileStatus = 4
	FileStatusRejected           FileStatus = 5
	FileStatusMalwareDetected    FileStatus = 6
	FileStatusDeleted            FileStatus = 7
	FileStatusArchived           FileStatus = 8
	FileStatusTesting            FileStatus = 9
	FileStatusReleased           FileStatus = 10
	FileStatusReadyForReview     FileStatus = 11
	FileStatusDeprecated         FileStatus = 12
	FileStatusBaking             FileStatus = 13
	FileStatusAwaitingPublishing FileStatus = 14
	FileStatusFailedPublishing   FileStatus = 15
)

// FileRelationType represents the type of file dependency.
type FileRelationType int

const (
	FileRelationTypeEmbeddedLibrary    FileRelationType = 1
	FileRelationTypeOptionalDependency FileRelationType = 2
	FileRelationTypeRequiredDependency FileRelationType = 3
	FileRelationTypeTool               FileRelationType = 4
	FileRelationTypeIncompatible       FileRelationType = 5
	FileRelationTypeInclude            FileRelationType = 6
)

// FileDependency represents a file dependency.
type FileDependency struct {
	ModID        int              `json:"modId"`
	RelationType FileRelationType `json:"relationType"`
}

// File represents a mod file.
type File struct {
	ID                   int                   `json:"id"`
	GameID               int                   `json:"gameId"`
	ModID                int                   `json:"modId"`
	IsAvailable          bool                  `json:"isAvailable"`
	DisplayName          string                `json:"displayName"`
	FileName             string                `json:"fileName"`
	ReleaseType          int                   `json:"releaseType"`
	FileStatus           FileStatus            `json:"fileStatus"`
	Hashes               []FileHash            `json:"hashes"`
	FileFingerprint      int                   `json:"fileFingerprint"`
	FileDate             string                `json:"fileDate"`
	FileLength           int                   `json:"fileLength"`
	DownloadCount        int                   `json:"downloadCount"`
	DownloadURL          *string               `json:"downloadUrl"`
	GameVersions         []string              `json:"gameVersions"`
	IsAlternate          bool                  `json:"isAlternate"`
	AlternateFileID      int                   `json:"alternateFileId"`
	Dependencies         []FileDependency      `json:"dependencies"`
	Modules              []Module              `json:"modules"`
	SortableGameVersions []SortableGameVersion `json:"sortableGameVersions"`
}

// SortableGameVersion represents metadata for sorting game versions.
type SortableGameVersion struct {
	GameVersionPadded      string `json:"gameVersionPadded"`
	GameVersion            string `json:"gameVersion"`
	GameVersionReleaseDate string `json:"gameVersionReleaseDate"`
	GameVersionName        string `json:"gameVersionName"`
}

// Module represents a file within a mod file.
type Module struct {
	Name        string `json:"name"`
	Fingerprint int    `json:"fingerprint"`
	Type        int    `json:"type"`
}

// Author represents a mod author.
type Author struct {
	ProjectID         int     `json:"projectId"`
	ProjectTitleID    *int    `json:"projectTitleId"`
	ProjectTitleTitle *string `json:"projectTitleTitle"`
	Name              string  `json:"name"`
	URL               string  `json:"url"`
	ID                int     `json:"id"`
	UserID            int     `json:"userId"`
	TwitchID          int     `json:"twitchId"`
}

// ModCategory represents a mod category.
type ModCategory struct {
	ID               int    `json:"id"`
	GameID           int    `json:"gameId"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	URL              string `json:"url"`
	IconURL          string `json:"iconUrl"`
	DateModified     string `json:"dateModified"`
	IsClass          *bool  `json:"isClass"`
	ClassID          *int   `json:"classId"`
	ParentCategoryID *int   `json:"parentCategoryId"`
	DisplayIndex     *int   `json:"displayIndex"`
}

// SearchOptions defines options for searching mods.
type SearchOptions struct {
	ClassID           *int                 `json:"-"`
	CategoryID        *int                 `json:"-"`
	GameID            *int                 `json:"-"`
	GameVersion       *string              `json:"-"`
	Index             *int                 `json:"-"`
	SortField         *ModsSearchSortField `json:"-"`
	SortOrder         *string              `json:"-"`
	ModLoaderType     *FileModLoaderType   `json:"-"`
	ModLoaderTypes    []string             `json:"-"`
	GameVersionTypeID *int                 `json:"-"`
	Slug              *string              `json:"-"`
	PageSize          *int                 `json:"-"`
	SearchFilter      *string              `json:"-"`
}

// ModsSearchSortField represents sort fields for mod search.
type ModsSearchSortField int

const (
	ModsSearchSortFieldFeatured       ModsSearchSortField = 1
	ModsSearchSortFieldPopularity     ModsSearchSortField = 2
	ModsSearchSortFieldLastUpdated    ModsSearchSortField = 3
	ModsSearchSortFieldName           ModsSearchSortField = 4
	ModsSearchSortFieldAuthor         ModsSearchSortField = 5
	ModsSearchSortFieldTotalDownloads ModsSearchSortField = 6
	ModsSearchSortFieldCategory       ModsSearchSortField = 7
	ModsSearchSortFieldGameVersion    ModsSearchSortField = 8
)

// GetModFilesOptions defines options for fetching mod files.
type GetModFilesOptions struct {
	ModID             int                `json:"modId"`
	GameVersion       *string            `json:"gameVersion"`
	ModLoaderType     *FileModLoaderType `json:"modLoaderType"`
	GameVersionTypeID *int               `json:"gameVersionTypeId"`
	Index             *int               `json:"index"`
	PageSize          *int               `json:"pageSize"`
}

// Pagination represents pagination metadata.
type Pagination struct {
	Index       int `json:"index"`
	PageSize    int `json:"pageSize"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}

// FingerprintMatch represents a fingerprint match result.
type FingerprintMatch struct {
	ID          int    `json:"id"`
	File        File   `json:"file"`
	LatestFiles []File `json:"latestFiles"`
}

// FingerprintsMatchesResult represents the result of a fingerprints match.
type FingerprintsMatchesResult struct {
	IsCacheBuilt          bool                   `json:"isCacheBuilt"`
	ExactMatches          []FingerprintMatch     `json:"exactMatches"`
	ExactFingerprints     []int                  `json:"exactFingerprints"`
	PartialMatches        []FingerprintMatch     `json:"partialMatches"`
	PartialFingerprints   map[string]interface{} `json:"partialFingerprints"`
	UnmatchedFingerprints []int                  `json:"unmatchedFingerprints"`
}

// FingerprintFuzzyMatch represents a fuzzy fingerprint match.
type FingerprintFuzzyMatch struct {
	ID           int    `json:"id"`
	File         File   `json:"file"`
	LatestFiles  []File `json:"latestFiles"`
	Fingerprints []int  `json:"fingerprints"`
}

// FingerprintFuzzyMatchResult represents the result of a fuzzy fingerprints match.
type FingerprintFuzzyMatchResult struct {
	FuzzyMatches []FingerprintFuzzyMatch `json:"fuzzyMatches"`
}

// CurseforgeV1Client is a client for the Curseforge V1 API.
type CurseforgeV1Client struct {
	baseURL    string
	headers    map[string]string
	httpClient *http.Client
}

// CurseforgeClientOption is a functional option for configuring the client.
type CurseforgeClientOption func(*CurseforgeV1Client)

// WithBaseURL sets the base URL for the client.
func WithBaseURL(baseURL string) CurseforgeClientOption {
	return func(c *CurseforgeV1Client) {
		c.baseURL = baseURL
	}
}

// WithHeaders sets additional headers for the client.
func WithHeaders(headers map[string]string) CurseforgeClientOption {
	return func(c *CurseforgeV1Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) CurseforgeClientOption {
	return func(c *CurseforgeV1Client) {
		c.httpClient = client
	}
}

// NewCurseforgeV1Client creates a new Curseforge V1 API client.
func NewCurseforgeV1Client(apiKey string, options ...CurseforgeClientOption) *CurseforgeV1Client {
	c := &CurseforgeV1Client{
		baseURL:    "https://api.curseforge.com",
		headers:    map[string]string{"x-api-key": apiKey, "accept": "application/json"},
		httpClient: http.DefaultClient,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *CurseforgeV1Client) request(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Response, error) {
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

func (c *CurseforgeV1Client) doJSON(ctx context.Context, method, path string, body any, result any) error {
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
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return &CurseforgeAPIError{
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

// GetCategories fetches available categories for a game.
func (c *CurseforgeV1Client) GetCategories(ctx context.Context) ([]ModCategory, error) {
	v := url.Values{}
	v.Add("gameId", "432")
	path := "/v1/categories?" + v.Encode()
	var result struct {
		Data []ModCategory `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return result.Data, err
}

// GetMod fetches a mod by its ID.
func (c *CurseforgeV1Client) GetMod(ctx context.Context, modID int) (*Mod, error) {
	path := fmt.Sprintf("/v1/mods/%d", modID)
	var result struct {
		Data Mod `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return &result.Data, err
}

// GetModDescription fetches the description of a mod.
func (c *CurseforgeV1Client) GetModDescription(ctx context.Context, modID int) (string, error) {
	path := fmt.Sprintf("/v1/mods/%d/description", modID)
	var result struct {
		Data string `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return result.Data, err
}

// ModFilesResponse represents the response structure for mod files
type ModFilesResponse struct {
	Data       []File     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// GetModFiles fetches files for a mod with specified options.
func (c *CurseforgeV1Client) GetModFiles(ctx context.Context, options GetModFilesOptions) (ModFilesResponse, error) {
	v := url.Values{}
	if options.GameVersion != nil {
		v.Add("gameVersion", *options.GameVersion)
	}
	if options.ModLoaderType != nil {
		v.Add("modLoaderType", fmt.Sprint(*options.ModLoaderType))
	}
	if options.GameVersionTypeID != nil {
		v.Add("gameVersionTypeId", fmt.Sprint(*options.GameVersionTypeID))
	}
	if options.Index != nil {
		v.Add("index", fmt.Sprint(*options.Index))
	}
	if options.PageSize != nil {
		v.Add("pageSize", fmt.Sprint(*options.PageSize))
	}
	path := fmt.Sprintf("/v1/mods/%d/files?%s", options.ModID, v.Encode())
	var result struct {
		Data       []File     `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

// GetModFile fetches a specific file for a mod.
func (c *CurseforgeV1Client) GetModFile(ctx context.Context, modID, fileID int) (*File, error) {
	path := fmt.Sprintf("/v1/mods/%d/files/%d", modID, fileID)
	var result struct {
		Data File `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return &result.Data, err
}

// GetMods fetches multiple mods by their IDs.
func (c *CurseforgeV1Client) GetMods(ctx context.Context, modIDs []int) ([]Mod, error) {
	path := "/v1/mods"
	body := map[string][]int{"modIds": modIDs}
	var result struct {
		Data []Mod `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodPost, path, body, &result)
	return result.Data, err
}

// GetFiles fetches multiple files by their IDs.
func (c *CurseforgeV1Client) GetFiles(ctx context.Context, fileIDs []int) ([]File, error) {
	path := "/v1/mods/files"
	body := map[string][]int{"fileIds": fileIDs}
	var result struct {
		Data []File `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodPost, path, body, &result)
	return result.Data, err
}

// SearchMods searches for mods based on the provided options.
func (c *CurseforgeV1Client) SearchMods(ctx context.Context, options SearchOptions) (struct {
	Data       []Mod      `json:"data"`
	Pagination Pagination `json:"pagination"`
}, error) {
	v := url.Values{}
	v.Add("gameId", "432")
	if options.ClassID != nil {
		v.Add("classId", fmt.Sprint(*options.ClassID))
	}
	if options.CategoryID != nil {
		v.Add("categoryId", fmt.Sprint(*options.CategoryID))
	}
	if options.GameVersion != nil {
		v.Add("gameVersion", *options.GameVersion)
	}
	if options.SearchFilter != nil {
		v.Add("searchFilter", *options.SearchFilter)
	}
	sortField := ModsSearchSortFieldPopularity
	if options.SortField != nil {
		sortField = *options.SortField
	}
	v.Add("sortField", fmt.Sprint(sortField))
	sortOrder := "desc"
	if options.SortOrder != nil {
		sortOrder = *options.SortOrder
	}
	v.Add("sortOrder", sortOrder)
	if options.ModLoaderType != nil {
		v.Add("modLoaderType", fmt.Sprint(*options.ModLoaderType))
	}
	if len(options.ModLoaderTypes) > 0 {
		v.Add("modLoaderTypes", "["+strings.Join(options.ModLoaderTypes, ",")+"]")
	}
	if options.GameVersionTypeID != nil {
		v.Add("gameVersionTypeId", fmt.Sprint(*options.GameVersionTypeID))
	}
	if options.Slug != nil {
		v.Add("slug", *options.Slug)
	}
	index := 0
	if options.Index != nil {
		index = *options.Index
	}
	v.Add("index", fmt.Sprint(index))
	pageSize := 25
	if options.PageSize != nil {
		pageSize = *options.PageSize
	}
	v.Add("pageSize", fmt.Sprint(pageSize))
	path := "/v1/mods/search?" + v.Encode()
	var result struct {
		Data       []Mod      `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

// GetModFileChangelog fetches the changelog for a specific mod file.
func (c *CurseforgeV1Client) GetModFileChangelog(ctx context.Context, modID, fileID int) (string, error) {
	path := fmt.Sprintf("/v1/mods/%d/files/%d/changelog", modID, fileID)
	var result struct {
		Data string `json:"data"`
	}
	err := c.doJSON(ctx, http.MethodGet, path, nil, &result)
	return result.Data, err
}

// GetFingerprintsMatchesByGameID fetches fingerprint matches for a game.
func (c *CurseforgeV1Client) GetFingerprintsMatchesByGameID(ctx context.Context, gameID int, fingerprints []int) (*FingerprintsMatchesResult, error) {
	path := fmt.Sprintf("/v1/fingerprints/%d", gameID)
	body := map[string][]int{"fingerprints": fingerprints}
	var result FingerprintsMatchesResult
	err := c.doJSON(ctx, http.MethodPost, path, body, &result)
	return &result, err
}

// GetFingerprintsFuzzyMatchesByGameID fetches fuzzy fingerprint matches for a game.
func (c *CurseforgeV1Client) GetFingerprintsFuzzyMatchesByGameID(ctx context.Context, gameID int, fingerprints []int) (*FingerprintFuzzyMatchResult, error) {
	path := fmt.Sprintf("/v1/fingerprints/fuzzy/%d", gameID)
	body := map[string][]int{"fingerprints": fingerprints}
	var result FingerprintFuzzyMatchResult
	err := c.doJSON(ctx, http.MethodPost, path, body, &result)
	return &result, err
}
