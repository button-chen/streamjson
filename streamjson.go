package streamjson

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type StreamJson struct {
	cbs    map[string]func(any, error)
	called map[string]bool
	dec    *json.Decoder
}

func NewStreamJson() *StreamJson {
	s := &StreamJson{
		cbs:    make(map[string]func(any, error)),
		called: make(map[string]bool),
	}
	return s
}

func (s *StreamJson) AddMonitor(pattern string, cb func(any, error)) {
	s.cbs[pattern] = cb
	s.called[pattern] = false
}

func (s *StreamJson) ProcessStream(r io.Reader) error {
	s.dec = json.NewDecoder(r)
	keys := []string{}
	err := s.process(&keys)
	if err != nil {
		return err
	}
	for pattern, ok := range s.called {
		if !ok {
			s.cbs[pattern](nil, fmt.Errorf("not find pattern: "+pattern))
		}
	}
	return nil
}

func (s *StreamJson) process(keys *[]string) error {
	var key string
	for {
		tok, err := s.dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		f, ok := tok.(json.Delim)
		if ok && f == '{' {
			if key != "" {
				*keys = append(*keys, key)
				key = ""
			}
			err := s.process(keys)
			if err != nil {
				return err
			}
			continue
		}
		if ok && f == '}' {
			if len(*keys) > 0 && (*keys)[len(*keys)-1] != "*" {
				*keys = (*keys)[:len(*keys)-1]
			}
			break
		}

		if key == "" {
			key = tok.(string)
			continue
		}
		newKey := key
		if len(*keys) > 0 {
			newKey = strings.Join(append(*keys, key), ".")
		}
		if f, ok := tok.(json.Delim); ok && f == '[' {
			*keys = append(*keys, key)
			*keys = append(*keys, "*")
			s.array(keys, newKey)
			if len(*keys) > 0 {
				*keys = (*keys)[:len(*keys)-1]
			}
			key = ""
			continue
		}
		if cb, ok := s.cbs[newKey]; ok {
			cb(tok, nil)
			s.called[newKey] = true
		}
		key = ""
	}
	return nil
}

func (s *StreamJson) array(keys *[]string, key string) error {
	key = key + ".*"

	for {
		tok, err := s.dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		f, ok := tok.(json.Delim)
		if ok && f == '[' {
			*keys = append(*keys, "*")
			s.array(keys, key)
			continue
		}
		if ok && f == ']' {
			if len(*keys) > 0 {
				*keys = (*keys)[:len(*keys)-1]
			}
			break
		}
		if ok && f == '{' {
			s.process(keys)
			continue
		}

		if cb, ok := s.cbs[key]; ok {
			cb(tok, nil)
			s.called[key] = true
		}
	}
	return nil
}
