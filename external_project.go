package main

type ExternalProject struct {
	Owner string `yaml:"Owner"`
	Repo  string `yaml:"Repo"`
	Path  string `yaml:"Path"`
}

func (c *ExternalProject) GetProject() *Project {
	ctx, githubClient := githubClient()
	repoConfig, _, _, err := githubClient.Repositories.GetContents(ctx, c.Owner, c.Repo, c.Path, nil)
	if err != nil {
		return nil
	}

	repoConfigContent, err := repoConfig.GetContent()
	if err != nil {
		return nil
	}

	if isYaml(c.Path) {
		return NewProjectFromYaml([]byte(repoConfigContent))
	}

	return NewProject([]byte(repoConfigContent))
}
