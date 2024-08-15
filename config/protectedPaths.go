package config

import "strings"

type ProtectedPaths struct {
	paths map[string]struct{}
}

func NewProtectedPaths() *ProtectedPaths {
	return &ProtectedPaths{
		paths: make(map[string]struct{}),
	}
}

func (p *ProtectedPaths) Add(key string) {
	p.paths[key] = struct{}{}
}

func (p *ProtectedPaths) Contains(key string) bool {
	_, exist := p.paths[key]
	return exist
}

func (p *ProtectedPaths) String() string {
	var paths strings.Builder
	for key, _ := range p.paths {
		paths.WriteString(key + "\n")
	}
	return paths.String()
}
