package domain

import "github.com/scovl/ollanta/domain/model"

type ComponentType = model.ComponentType

const (
	ComponentProject = model.ComponentProject
	ComponentModule  = model.ComponentModule
	ComponentPackage = model.ComponentPackage
	ComponentFile    = model.ComponentFile
)

type Component = model.Component
