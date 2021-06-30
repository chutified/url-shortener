package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/chutommy/url-shortener/config"
	"github.com/stretchr/testify/assert"
)

var testDir = "test-settings/"

var openConfigTests = []struct {
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

func TestOpenConfig(t *testing.T) {
	for _, tc := range openConfigTests {
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

var getConfigTests = []struct {
	name       string
	file       string
	cfg        *config.Config
	noErr      bool
	preAction  func() error
	postAction func() error
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
		preAction: func() error {
			return os.Setenv("URL_SHORTENER_DBCONN", "url_shortener_dbconn")
		},
		postAction: func() error {
			return os.Unsetenv("URL_SHORTENER_DBCONN")
		},
	},
	{
		name:  "no dbconn env variable",
		file:  "settings_1.json",
		cfg:   nil,
		noErr: false,
		preAction: func() error {
			return nil
		},
		postAction: func() error {
			return nil
		},
	},
	{
		name:  "fail to open file",
		file:  "settings_2.json",
		cfg:   nil,
		noErr: false,
		preAction: func() error {
			return nil
		},
		postAction: func() error {
			return nil
		},
	},
}

func TestGetConfig(t *testing.T) {
	for _, tc := range getConfigTests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Nil(t, tc.preAction())

			// test function
			cfg, err := config.GetConfig(testDir + tc.file)
			if tc.noErr {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}

			// check config values
			if tc.cfg == nil {
				assert.Nil(t, cfg)
			} else if assert.NotNil(t, cfg) {
				assert.Equal(t, tc.cfg.SrvPort, cfg.SrvPort)
				assert.Equal(t, tc.cfg.SrvTimeOut, cfg.SrvTimeOut)
				if tc.cfg.DB != nil {
					if assert.NotNil(t, cfg.DB) {
						assert.Equal(t, tc.cfg.DB.Driver, cfg.DB.Driver)
					}
				}
			}

			assert.Nil(t, tc.postAction())
		})
	}
}

var addrTests = []struct {
	name string
	cfg  *config.Config
	port string
}{
	{
		name: "port 80",
		cfg: &config.Config{
			SrvPort: 80,
		},
		port: ":80",
	},
	{
		name: "port 8080",
		cfg: &config.Config{
			SrvPort: 8080,
		},
		port: ":8080",
	},
	{
		name: "port 1234",
		cfg: &config.Config{
			SrvPort: 1234,
		},
		port: ":1234",
	},
	{
		name: "no port",
		cfg:  &config.Config{},
		port: ":0",
	},
	{
		name: "no config",
		cfg:  nil,
		port: "",
	},
}

func TestConfig_Addr(t *testing.T) {
	for _, tc := range addrTests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.port, tc.cfg.Addr())
		})
	}
}
