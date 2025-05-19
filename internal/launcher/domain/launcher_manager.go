package domain

import (
	"context"
	"fmt"
)

type LauncherManager struct {
}

func NewLauncherManager() *LauncherManager {
	return &LauncherManager{}
}

func (lm *LauncherManager) Run(ctx context.Context) error {
	fmt.Println("===========> LAUNCH RUN")
	return nil
}

func (lm *LauncherManager) Stop(ctx context.Context) error {
	fmt.Println("===========> LAUNCH STOPED")
	return nil
}
