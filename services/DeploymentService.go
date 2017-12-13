package services

import (
	"github.com/softleader/deployer/datamodels"
	"fmt"
	"github.com/softleader/deployer/cmd"
)

type DeploymentService struct {
	DockerStack cmd.DockerStack
}

func (s *DeploymentService) GetAll() string {
	return s.DockerStack.Ls()
}

func (s *DeploymentService) Deploy(stack string, d datamodels.Deployment) string {
	fmt.Printf("deploy %v with [%+v]...", stack, d)
	return s.GetAll()
}

func (s *DeploymentService) Delete(stack string) string {
	return s.DockerStack.Rm(stack)
}
