package model

import (
	"fmt"
	"golang.org/x/exp/utf8string"
	"strings"
	"unicode/utf8"
)

// models used in requests

type Validable interface {
	Validate() error
}

type PlayerRegister struct {
	Name string `json:"name"`
}

func (p PlayerRegister) Validate() error {
	if utf8.RuneCountInString(p.Name) > 64 || len(p.Name) > 64 {
		return fmt.Errorf("name too long")
	}
	if p.Name == "" {
		return fmt.Errorf("name empty")
	}
	return nil
}

type PlayerObtainCard struct {
	Key string `json:"key"`
}

func (p PlayerObtainCard) Validate() error {
	if p.Key == "" {
		return fmt.Errorf("key empty")
	}
	if !utf8string.NewString(p.Key).IsASCII() {
		return fmt.Errorf("non-ascii characters")
	}
	if len(p.Key) < 9 || len(p.Key) > 10 {
		return fmt.Errorf("invalid length")
	}
	return nil
}

type PatchPlayerParams struct {
	IsAdmin bool `json:"is_admin"`
}

func (p PatchPlayerParams) Validate() error {
	return nil
}

type CreateCardBaseParams struct {
	Key    *string `json:"key"` // key empty for unobtainable system cards
	Name   string  `json:"name"`
	Source *string `json:"source"`
	Place  int16   `json:"place"` // surprisingly negative numbers are valid for place
}

func (p CreateCardBaseParams) Validate() error {
	if p.Key != nil {
		if !utf8string.NewString(*p.Key).IsASCII() {
			return fmt.Errorf("non-ascii characters in key")
		}
		if len(*p.Key) != 9 {
			return fmt.Errorf("invalid key length")
		}
		if strings.ToUpper(*p.Key) != *p.Key {
			return fmt.Errorf("not all upper-case")
		}
	}
	if p.Name == "" {
		return fmt.Errorf("name empty")
	}
	if utf8.RuneCountInString(p.Name) > 64 || len(p.Name) > 64 {
		return fmt.Errorf("name too long")
	}
	if p.Source != nil {
		if utf8.RuneCountInString(*p.Source) > 255 || len(*p.Source) > 255 {
			return fmt.Errorf("source too long")
		}
	}
	return nil
}

type AssignWearLevelParams struct {
	ImgUrl string `json:"img_url"`
}

func (p AssignWearLevelParams) Validate() error {
	if p.ImgUrl == "" {
		return fmt.Errorf("missing img url")
	}
	return nil
}
