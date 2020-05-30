TAIGO
-----

[![Travis CI Build Status](https://travis-ci.org/theriverman/taigo.svg?branch=master)](https://travis-ci.org/theriverman/taigo)
[![GoDoc](https://godoc.org/github.com/theriverman/taigo?status.svg)](https://pkg.go.dev/github.com/theriverman/taigo?tab=doc "Docs @ pkg.go.dev")
[![Go Report Card](https://goreportcard.com/badge/github.com/theriverman/taigo)](https://goreportcard.com/report/github.com/theriverman/taigo)
![Taigo Banner](assets/banner_627x300.png "TAIGO Banner")



# Description
**Taigo** is a driver written in GO for [Taiga](https://github.com/taigaio) targeting to implement all publicly available API v1 endpoints. <br>
Taiga is an Agile, Free and Open Source Project Management Tool.

Should you have any ideas or recommendations, feel free to report it as an issue or open a pull request.

# Integration
Download **Taigo**:
```bash
go get github.com/theriverman/taigo
```
Import it into your project: <br>
```go
import (
    taiga "github.com/theriverman/taigo"
)
```

# Documentation
Documentation is located at [pkg.go.dev/Taigo](https://pkg.go.dev/github.com/theriverman/taigo?tab=doc).

# Architecture
To use **Taigo**, you must instantiate a `*Client`, and authenticate against Taiga to either get a new token or validate the one given. To authenticate, use one of the following methods:
  * `client.AuthByCredentials`
  * `client.AuthByToken`

All Taiga objects, such as, Epic, User Story, Issue, Task, Sprint, etc.. are represented as a struct.<br>

## Meta System
Since the API endpoints of Taiga return various payloads for the same object types, a meta system has been added 
to simplify the interaction with the **Taigo** APIs. <br>

For example, in Taiga the `http://localhost:8000/api/v1/userstories` endpoint can return the following payload types:
  - User story detail
  - User story detail (GET)
  - User story detail (LIST)

The returned payload type depends on the executed action: list, create, edit, get, etc..

Since Go does not support generics and casting between different types is not a trivial operation, a generic type has been introduced in **Taigo** for each implemented basic object type (Epic, Milestone, UserStory, Task, etc..) found in Taiga.

These generic types come with meta fields (pointers) to the originally returned payload.
In case you need access to a field *not* represented in the generic type provided by **Taigo**, you can use the appropriate meta field.

For example, struct `Epic` has the following meta fields:

```go
type Epic struct {
    ID                int
    // ...
    EpicDetail        *EpicDetail
    EpicDetailGET     *EpicDetailGET
    EpicDetailLIST    *EpicDetailLIST
}
```

The available meta field is always stated in the operation's docstring as `Available Meta`, for example:
```go
// CreateEpic => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create
//
// Available Meta: *EpicDetail
func (s *EpicService) Create(epic Epic) (*Epic, error) {}
```

## Docs-As-Code ( godoc / docstrings )
To avoid creating and maintaining a redundant API documentation for each covered Taiga API resource, 
the API behaviours and fields are **not** documented in **Taigo**.
Instead, the docstring comment provides a direct URL to each appropriate online resource.<br>
**Example:**<br>
```go
// CreateEpic => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create
func (s *EpicService) Create(epic Epic) (*Epic, error) {}
```

Upon opening that link, you can find out what fields are available through `*Epic`:
  > - project (required): project id <br>
  > - subject (required) 

To find out which operation requires what fields and what they return, always refer to the official 
[Taiga REST API](https://taigaio.github.io/taiga-doc/dist/api.html) documentation.

## Pagination
It is rarely useful to receive results paginated in a driver like **Taigo**, so **pagination is disabled** by default. <br>

To enable pagination (for all future requests) set `Client.DisablePagination` to `false`: <br>
```go
// Enable pagination (disabled by default)
client.DisablePagination(false)
```

If you enable pagination, you can collect pagination details by accessing the returned `*http.Response`, <br>
and extracting the relevant headers into a taigo.Pagination struct the following way:
```go
// Collect returned Pagination headers
p := Pagination{}
p.LoadFromHeaders(client, response)
```

**Note:** See `type Pagination struct` in [TAIGO/common.models.go](common.models.go) for more details.

For details on Taiga's pagination, see [Taiga REST API / Pagination](https://taigaio.github.io/taiga-doc/dist/api.html#_pagination).

# Usage

## Client Instantiation
```go
package main

import (
  "fmt"
  "net/http"

  taiga "github.com/theriverman/taigo/gotaiga"
)
func main() {
	// Create client
	client := taiga.Client{
		BaseURL:    "https://api.taiga.io",
		HTTPClient: &http.Client{},
	}
	
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		log.Fatalln(err)
	}
	
	// Authenticate (get/set Token)
	client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: "my_pretty_username",
		Password: "123123",
	})

	// Set default project (optional. it is for convenience)
	client.SetDefaultProjectBySlug("theriverman-test-1337")

	// Get /users/me
	me, _ := client.User.Me()
	fmt.Println("Me: (ID, Username, FullName)", me.ID, me.Username, me.FullName)

  	// Get Project [ theriverman-api-dev-testing-1337 ]
	slug := "theriverman-api-dev-testing-1337"
	fmt.Printf("Getting Project (slug=%s)..\n", slug)
	project, err := client.Project.GetBySlug(slug)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Project name: %s \n\n", project.Name)
}
```

## Basic Operations
```go
// get all projects (no filtering, no ordering)
projectsList, _ := client.Project.List(nil)
for _, proj := range *projectsList {
    fmt.Println("Project ID:", proj.ID, "Project Name:", proj.Name)
}
/*
  ListProjects accepts a `*ProjectsQueryParameters` as an argument, but if you don't need any filtering, you can pass in a nil.
*/

// get all projects for user ID 1337
queryParams := taiga.Project.QueryParameters{
    Member: 1337,
}
queryParams.TotalFansLastMonth() // results are ordered by TotalFansLastMonth
projectsList, _ := client.Project.List(&queryParams)
for _, proj := range *projectsList {
    fmt.Println("Project ID:", proj.ID, "Project Name:", proj.Name)
}

// upload a file (Create an Epic attachment)
newAttachment, err := client.Epic.CreateAttachment(&taiga.Attachment{ObjectID: 1337, Project: 7331}, "C:/Users/theriverman/Pictures/nghfb.jpg")
```

## Non-Standard Operations (Non-Standard)
1. Clone Epic (with UserStories)
2. Clone Epic (without UserStories)
3. Clone UserStory (with sub-tasks)
4. Clone UserStory (without sub-tasks)
5. Clone Sub-Task
6. Copy UserStory to another project [will lose comments and attachments]

## Advanced Operations
Do you need access to a non yet implemented or special API endpoint? No problem! <br>
It is possible to make requests to custom API endpoints by leveraging `*RequestService`.

HTTP operations (`GET`, `POST`, etc..) should be executed through `RequestService` which provides a managed environment for communicating with Taiga.

For example, let's try accessing the `epic-custom-attributes` endpoint:

First, a model struct must be created to represent the data returned by the endpoint. Refer to the official Taiga API documentation for response examples:
```go
// EpicCustomAttributeDetail -> https://taigaio.github.io/taiga-doc/dist/api.html#object-epic-custom-attribute-detail
// Converted via https://mholt.github.io/json-to-go/
type EpicCustomAttributeDetail struct {
	CreatedDate  time.Time   `json:"created_date"`
	Description  string      `json:"description"`
	Extra        interface{} `json:"extra"`
	ID           int         `json:"id"`
	ModifiedDate time.Time   `json:"modified_date"`
	Name         string      `json:"name"`
	Order        int         `json:"order"`
	Project      int         `json:"project"`
	Type         string      `json:"type"`
}
```

Then initialise an instance which is going to store the reponse:
```go
epicCustomAttributes := []EpicCustomAttributeDetail{}
```

Using `client.MakeURL` the absolute URL endpoint is built for the request. <br>
Using `client.Request.Get` an HTTP GET request is sent to the API endpoint URL. <br>
The `client.Request.Get` call must receive a URL and a pointer to the model struct:
```go
// Final URL should be https://api.taiga.io/api/v1/epic-custom-attributes
resp, err := client.Request.Get(client.MakeURL("epic-custom-attributes"), &epicCustomAttributes)
if err != nil {
	fmt.Println(err)
	fmt.Println(resp)
} else {
	for i := 0; i < 3; i++ {
		ca := epicCustomAttributes[i]
		fmt.Println("  * ", ca.ID, ca.Name)
	}
}
```

Such raw requests return `(*http.Response, error)`.

# Contribution
You're contribution would be much appreciated! <br>
Feel free to open Issue tickets or create Pull requests. <br>

Should you have an idea which would cause non-backward compatibility, please open an issue first to discuss your proposal!

For more details, see [TAIGO Contribution](CONTRIBUTION.md).

## Branching Strategy
This package conforms to the stable HEAD philosophy.
  * **master** is <i>always</i> stable
  * **develop** is *not* <i>always</i> stable

# Licenses & Recognitions
|							|			|															|
|---------------------------|-----------|-----------------------------------------------------------|
| **Go Gopher**             | CC3.0     | [Renee French](https://www.instagram.com/reneefrench/) 	|
| **Gopher Konstructor**    | CC0-1.0   | [quasilyte](https://github.com/quasilyte/gopherkon)    	|
| **TAIGA**                 | AGPL-3.0  | [Taiga.io](https://github.com/taigaio)                 	|
| **Taigo**                 | MIT       | [theriverman](https://github.com/theriverman)          	|
