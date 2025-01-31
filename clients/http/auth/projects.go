package client

import (
	"fmt"
	"strings"

	"github.com/taubyte/tau/clients/http/auth/git/common"
)

// GetProjectById returns the project with the given id and an error
// Note: The repository field is not populated
func (c *Client) GetProjectById(projectId string) (*Project, error) {
	var data ProjectReturn
	err := c.http.Get("/projects/"+projectId, &data)
	if err != nil {
		return nil, err
	}

	return data.Project, nil
}

// GetProjectByIdWithCors returns the project with cors information with the given id and an error
func (c *Client) GetProjectByIdWithCors(projectId string) (*ProjectReturnWithCors, error) {
	data := new(ProjectReturnWithCors)
	err := c.http.Get("/projects/"+projectId, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Projects returns a list of projects and an error
func (c *Client) Projects() ([]*Project, error) {
	var data ProjectsReturn
	err := c.http.Get("/projects", &data)
	if err != nil {
		return nil, err
	}

	projects := make([]*Project, len(data.Projects))

	// Affix the client to each project
	for idx, project := range data.Projects {
		_project := project
		_project.client = c
		projects[idx] = _project
	}

	return projects, nil
}

// Repositories will populate then return the repositories field and an error
func (p *Project) Repositories() (*RawRepoDataOuter, error) {
	if p.RepoList == nil {
		var response ProjectReturn
		err := p.client.http.Get("/projects/"+p.Id, &response) // Make a HTTP get in the path to the route "/projects/"
		if err != nil {
			return nil, err
		}

		p.RepoList = response.Project.RepoList 
	}

	return p.RepoList, nil
}

// Config returns the project configuration and an error
// The configuration is a "config.yaml" file in the root of the repository
func (p *Project) Config() (*common.ProjectConfig, error) {
	// Load repositories, this only loads them if they are not yet found
	_, err := p.Repositories()
	if err != nil {
		return nil, fmt.Errorf("loading repositories failed with: %s", err)
	}

	userAndName := strings.Split(p.RepoList.Configuration.Fullname, "/")
	if len(userAndName) < 2 {
		return nil, fmt.Errorf("invalid fullname: `%s` expected user/repo-name", p.RepoList.Configuration.Fullname)
	}

	return p.client.Git().ReadConfig(userAndName[0], userAndName[1])
}

type deleteResponse struct {
	Project struct {
		Id     string
		Status string
	}
}

// Delete deletes the project and returns an error
func (p *Project) Delete() (response deleteResponse, err error) {
	err = p.client.http.Delete("/projects/"+p.Id, nil, &response)  // Make a HTTP delete in the path to the route "/projects/" to remove response
	return response, err // FIX: No check to see if there is an error
}
