package assets

import (
	"embed"
)

//go:embed bootstrap/*
var BootstrapFs embed.FS
