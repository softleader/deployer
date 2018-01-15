package services

import (
	"github.com/softleader/deployer/cmd"
	"github.com/softleader/deployer/models"
	"strings"
)

type PracticeService struct {
	cmd.Workspace
}

func (ps *PracticeService) GetAll() (p models.Practice, err error) {
	return models.ReadFromFile(ps.Workspace.Path())
}

func (ps *PracticeService) Add(content string) (err error) {
	content = strings.TrimSpace(content)
	if len(content) <= 0 {
		return nil
	}
	p, err := models.ReadFromFile(ps.Workspace.Path())
	if err != nil {
		return err
	}
	p = append(p, content)
	p.SaveToFile(ps.Workspace.Path())
	return nil
}

func (ps *PracticeService) Delete(idx int) (err error) {
	p, err := models.ReadFromFile(ps.Workspace.Path())
	if err != nil {
		return err
	}
	p = append(p[:idx], p[idx+1:]...)
	p.SaveToFile(ps.Workspace.Path())
	return nil
}
