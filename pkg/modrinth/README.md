# Modrinth API

Provides a client for the Modrinth API as described in <https://docs.modrinth.com/api/>.

## Usage

This package uses the standard `net/http` package for HTTP requests and has no external dependencies.

### Search Projects in Modrinth

You can use a keyword to search.

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/Voxelum/minecraft-launcher-core/pkg/modrinth"
)

func main() {
 client := modrinth.NewModrinthV2Client()
 searchOptions := modrinth.SearchProjectOptions{
  Query: "shader", // searching shader
 }
 result, err := client.SearchProjects(context.Background(), searchOptions)
 if err != nil {
  log.Fatal(err)
 }
 totalProjectCounts := result.TotalHits
 for _, project := range result.Hits {
  fmt.Printf("%s %s %s\n", project.ProjectID, project.Title, project.Description) // print project info
 }
}
```

### Get Project in Modrinth

You can get project detail info via project ID, including the download URL.

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/Voxelum/minecraft-launcher-core/pkg/modrinth"
)

func main() {
 client := modrinth.NewModrinthV2Client()
 projectID := "some-project-id" // you can get this id from SearchProjects
 project, err := client.GetProject(context.Background(), projectID) // project details
 if err != nil {
  log.Fatal(err)
 }
 versions := project.Versions
 oneVersion := versions[0]

 modVersion, err := client.GetProjectVersion(context.Background(), oneVersion)
 if err != nil {
  log.Fatal(err)
 }

 files := modVersion.Files
 file := files[0]
 url := file.URL
 name := file.Filename
 hashes := file.Hashes
 // now you can get file name, file hashes and download url of the file
 fmt.Printf("URL: %s, Name: %s, Hashes: %v\n", url, name, hashes)
}
```
