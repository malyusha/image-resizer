package app

import (
	"sync"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"
)

type Preset struct {
	Name   string `json:"name"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
}

type PresetsList []Preset

// Find looks up for Preset in slice of presets and returns first matched by given name
func (p *PresetsList) Find(name string) (*Preset, bool) {
	for _, v := range *p {
		if name == v.Name {
			return &v, true
		}
	}

	return nil, false
}

// LoadPreset loads all presets from json config file file
func LoadPresets(file string) (*PresetsList, error) {
	var (
		parsingError error
		once         sync.Once
		presets      = &PresetsList{}
	)

	once.Do(func() {
		if err := jsonToSlice(file, presets); err != nil {
			fmt.Println("Error while reading json file with presets occured")
			parsingError = err

			return
		}

		parsingError = checkPresetsList(presets)
	})

	return presets, parsingError
}

// rawToPresets parses raw map to PresetsMap
func checkPresetsList(data *PresetsList) error {
	for _, preset := range *data {
		if preset.Name == "" {
			return errors.New(fmt.Sprintf("Preset must contain name. Given %v", preset))
		}

		if preset.Width == 0 && preset.Height == 0 {
			return errors.New(fmt.Sprintf("Preset must containt at least one property (width, height). Given %v", preset))
		}
	}

	return nil
}

// Retrieves json from fixture file as slice of structs
func jsonToSlice(file string, dest interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}

	return nil
}
