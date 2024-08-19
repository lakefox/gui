package library

import (
	"image"
)

type Shelf struct {
	Textures       map[string]*image.RGBA
	References     map[string]bool
	UnloadCallback func(string)
}

func (s *Shelf) New(key string, img *image.RGBA) string {
	if s.Textures == nil {
		s.Textures = map[string]*image.RGBA{}
	}
	if s.References == nil {
		s.References = map[string]bool{}
	}

	s.Textures[key] = img
	s.References[key] = true
	return key
}

func (s *Shelf) Get(key string) (*image.RGBA, bool) {
	a, exists := s.Textures[key]

	return a, exists
}

func (s *Shelf) Check(key string) bool {
	_, exists := s.Textures[key]
	if exists {
		s.References[key] = true
	}
	return exists
}

func (s *Shelf) Close() {
	for k, v := range s.References {
		if !v {
			s.UnloadCallback(k)
			delete(s.References, k)
			delete(s.Textures, k)
		}
		s.References[k] = false
	}
}
