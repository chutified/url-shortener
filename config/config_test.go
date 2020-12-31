package config_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/chutified/url-shortener/config"
	"github.com/stretchr/testify/assert"
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
			if tc.err != nil {
				fmt.Println(err)
				if assert.NotNil(t, err) {
					assert.EqualError(t, err, tc.err.Error())
				}
			}
			assert.Equal(t, tc.cfg, cfg)
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
		cfg   *config.Config
		noErr bool
	}{
		{
			name: "ok",
			file: "settings_1.json",
			cfg: &config.Config{
				SrvPort:    8080,
				SrvTimeOut: "10s",
				DB: &config.DB{
					Driver: "postgres",
				},
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

			// check config values
			if tc.cfg != nil {
				if assert.NotNil(t, cfg) {

					if assert.NotNil(t, cfg) {
						assert.Equal(t, tc.cfg.SrvPort, cfg.SrvPort)
					}
					if tc.cfg.DB != nil {
						if assert.NotNil(t, cfg.DB) {
							assert.Equal(t, tc.cfg.DB.Driver, cfg.DB.Driver)
						}
					}
					assert.Equal(t, tc.cfg.SrvTimeOut, cfg.SrvTimeOut)
				}
			}
		})
	}
}
