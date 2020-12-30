package config_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/chutified/url-shortener/config"
)

var testDir = "test-settings/"

func TestOpenConfig(t *testing.T) {

	// define test table
	tt := []struct {
		name string
		file string
		cfg  config.Config
		err  error
	}{
		{
			name: "ok",
			file: "settings_1.json",
			cfg: config.Config{
				SrvPort:    8080,
				SrvTimeOut: "10s",
				DB: &config.DB{
					Driver: "postgres",
				},
			},
			err: nil,
		},
		{
			name: "invalid extension",
			file: "settings_2.yaml",
			cfg:  config.Config{},
			err:  config.ErrInvalidFileFormat,
		},
		{
			name: "corrupted content",
			file: "settings_3.json",
			cfg:  config.Config{},
			err:  config.ErrInvalidJSONFile,
		},
		{
			name: "invalid driver",
			file: "settings_4.json",
			cfg:  config.Config{},
			err:  config.ErrDriverNotSupported,
		},
		{
			name: "invalid time format",
			file: "settings_5.json",
			cfg:  config.Config{},
			err:  config.ErrInvalidTimeFormat,
		},
		{
			name: "file not found",
			file: "settings_6.json",
			cfg:  config.Config{},
			err:  config.ErrFileNotFound,
		},
	}

	// iterate over test cases
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cfg, err := config.OpenConfig(testDir + tc.file)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected %#v; got %#v", tc.err, err)
			}
			if !reflect.DeepEqual(cfg, tc.cfg) {
				t.Errorf("expected %#v; got %#v", tc.cfg, cfg)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {

	// define test table
	tt := []struct {
		name  string
		file  string
		cfg   *config.DB
		noErr bool
	}{
		{
			name: "ok",
			file: "settings_1.json",
			cfg: &config.DB{
				Driver: "postgres",
			},
			noErr: true,
		},
	}

	// iterate over test table
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cfg, err := config.GetConfig(testDir + tc.file)
			if (err == nil) != tc.noErr {
				t.Errorf("expected no error; got %#v", err)
			}
			if !reflect.DeepEqual(cfg, tc.cfg) {
				t.Errorf("expected %#v; got %#v", tc.cfg, cfg)
			}
		})
	}
}
