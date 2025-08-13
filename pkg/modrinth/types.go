package modrinth

import (
	"fmt"
	"net/http"
)

// ModrinthAPIError represents an error from the Modrinth API.
type ModrinthAPIError struct {
	URL    string
	Status int
	Body   string
}

func (e *ModrinthAPIError) Error() string {
	return fmt.Sprintf("Fail to fetch modrinth api %s. Status=%d. %s", e.URL, e.Status, e.Body)
}

// SearchResultHit represents a single hit in search results.
type SearchResultHit struct {
	Slug               string   `json:"slug"`
	ProjectID          string   `json:"project_id"`
	ProjectType        string   `json:"project_type"`
	Author             string   `json:"author"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Categories         []string `json:"categories"`
	Versions           []string `json:"versions"`
	Downloads          int      `json:"downloads"`
	Follows            int      `json:"follows"`
	PageURL            string   `json:"page_url"`
	IconURL            string   `json:"icon_url"`
	AuthorURL          string   `json:"author_url"`
	DateCreated        string   `json:"date_created"`
	DateModified       string   `json:"date_modified"`
	LatestVersion      string   `json:"latest_version"`
	License            string   `json:"license"`
	ClientSide         string   `json:"client_side"`
	ServerSide         string   `json:"server_side"`
	Host               string   `json:"host"`
	Gallery            []string `json:"gallery"`
	FeaturedGallery    string   `json:"featured_gallery"`
	MonetizationStatus string   `json:"monetization_status"`
}

// SearchProjectOptions defines options for searching projects.
type SearchProjectOptions struct {
	Query  string `json:"-"`
	Facets string `json:"-"`
	Filter string `json:"-"`
	Index  string `json:"-"`
	Offset int    `json:"-"`
	Limit  int    `json:"-"`
}

// SearchResult represents the result of a project search.
type SearchResult struct {
	Hits      []SearchResultHit `json:"hits"`
	Offset    int               `json:"offset"`
	Limit     int               `json:"limit"`
	TotalHits int               `json:"total_hits"`
}

// GetProjectVersionsOptions defines options for fetching project versions.
type GetProjectVersionsOptions struct {
	Loaders      []string `json:"-"`
	GameVersions []string `json:"-"`
	Featured     *bool    `json:"-"`
}

// Project represents a Modrinth project.
type Project struct {
	ID                 string            `json:"id"`
	Slug               string            `json:"slug"`
	ProjectType        string            `json:"project_type"`
	Team               string            `json:"team"`
	Title              string            `json:"title"`
	Description        string            `json:"description"`
	Body               string            `json:"body"`
	Published          string            `json:"published"`
	Updated            string            `json:"updated"`
	Approved           string            `json:"approved,omitempty"`
	Status             string            `json:"status"`
	RequestedStatus    string            `json:"requested_status,omitempty"`
	ModeratorMessage   *ModeratorMessage `json:"moderator_message,omitempty"`
	License            License           `json:"license"`
	ClientSide         string            `json:"client_side"`
	ServerSide         string            `json:"server_side"`
	Downloads          int               `json:"downloads"`
	Followers          int               `json:"followers"`
	Categories         []string          `json:"categories"`
	Versions           []string          `json:"versions"`
	IconURL            string            `json:"icon_url,omitempty"`
	Color              *int              `json:"color,omitempty"`
	ThreadID           string            `json:"thread_id,omitempty"`
	MonetizationStatus string            `json:"monetization_status"`
	IssuesURL          string            `json:"issues_url,omitempty"`
	SourceURL          string            `json:"source_url,omitempty"`
	WikiURL            string            `json:"wiki_url,omitempty"`
	DiscordURL         string            `json:"discord_url,omitempty"`
	DonationURLs       []DonationURL     `json:"donation_urls,omitempty"`
	Gallery            []GalleryItem     `json:"gallery,omitempty"`
}

// ModeratorMessage represents a moderator message for a project.
type ModeratorMessage struct {
	Message string `json:"message"`
	Body    string `json:"body"`
}

// DonationURL represents a donation link for a project.
type DonationURL struct {
	ID       string `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

// GalleryItem represents an item in a project's gallery.
type GalleryItem struct {
	URL         string `json:"url"`
	Featured    bool   `json:"featured"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Created     string `json:"created"`
	Ordering    int    `json:"ordering"`
}

// ProjectVersion represents a version of a Modrinth project.
type ProjectVersion struct {
	ID              string              `json:"id"`
	ProjectID       string              `json:"project_id"`
	AuthorID        string              `json:"author_id"`
	Name            string              `json:"name"`
	VersionNumber   string              `json:"version_number"`
	Changelog       string              `json:"changelog,omitempty"`
	DatePublished   string              `json:"date_published"`
	Downloads       int                 `json:"downloads"`
	VersionType     string              `json:"version_type"`
	Status          string              `json:"status"`
	RequestedStatus string              `json:"requested_status,omitempty"`
	Files           []VersionFile       `json:"files"`
	Dependencies    []VersionDependency `json:"dependencies"`
	GameVersions    []string            `json:"game_versions"`
	Loaders         []string            `json:"loaders"`
	Featured        bool                `json:"featured"`
}

// VersionFile represents a file in a project version.
type VersionFile struct {
	Hashes   map[string]string `json:"hashes"`
	URL      string            `json:"url"`
	Filename string            `json:"filename"`
	Primary  bool              `json:"primary"`
	Size     int               `json:"size"`
	FileType string            `json:"file_type,omitempty"`
}

// VersionDependency represents a dependency in a project version.
type VersionDependency struct {
	VersionID      string `json:"version_id,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
	FileName       string `json:"file_name,omitempty"`
	DependencyType string `json:"dependency_type"`
}

// Category represents a Modrinth category tag.
type Category struct {
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	ProjectType string `json:"project_type"`
	Header      string `json:"header"`
}

// License represents a Modrinth license tag.
type License struct {
	ID   string `json:"short"`
	Name string `json:"name"`
	URL  string `json:"link,omitempty"`
}

// GameVersion represents a Modrinth game version tag.
type GameVersion struct {
	Version     string `json:"version"`
	VersionType string `json:"version_type"`
	Date        string `json:"date"`
	Major       bool   `json:"major"`
}

// Loader represents a Modrinth loader tag.
type Loader struct {
	Icon                  string   `json:"icon"`
	Name                  string   `json:"name"`
	SupportedProjectTypes []string `json:"supported_project_types"`
}

// User represents a Modrinth user.
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Bio       string `json:"bio,omitempty"`
	AvatarURL string `json:"avatar_url"`
	Created   string `json:"created"`
	Role      string `json:"role"`
	Badges    int    `json:"badges"`
	GithubID  int    `json:"github_id,omitempty"`
}

// TeamMember represents a member of a project team.
type TeamMember struct {
	TeamID       string  `json:"team_id"`
	User         User    `json:"user"`
	Role         string  `json:"role"`
	Permissions  int     `json:"permissions"`
	Accepted     bool    `json:"accepted"`
	PayoutsShare float64 `json:"payouts_share,omitempty"`
	Ordering     int     `json:"ordering"`
}

// Collection represents a Modrinth collection (inferred from API usage).
type Collection struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Projects    []string `json:"projects"`
	IconURL     string   `json:"icon_url,omitempty"`
}

// ModrinthV2Client is a client for the Modrinth V2 API.
type ModrinthV2Client struct {
	baseURL    string
	headers    map[string]string
	httpClient *http.Client
}
