package config

import "github.com/illuminatingKong/kongming-kit/base/filesystem"

func (c *Config) abPath(inPath string) string {
	return filesystem.AbPath(inPath)
}
