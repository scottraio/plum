package plum

import (
	"context"
)

type Job struct {
	App  AppConfig
	Ctx  context.Context
	Name string
	Run  func() error
}
