package services

import (
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/models"
	"strings"
)

type PracticeService struct {
	cmd.Workspace
}

func (ps *PracticeService) Get() (content string, err error) {
	return models.ReadPractices(ps.Workspace.Path())
}

func (ps *PracticeService) Save(content string) (err error) {
	content = strings.TrimSpace(content)
	if len(content) <= 0 {
		return nil
	}
	return models.SavePractices(ps.Workspace.Path(), content)
}
