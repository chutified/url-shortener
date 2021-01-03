package config_test

import (
	"testing"

	"github.com/chutified/url-shortener/config"
)

type fields struct {
	DBConn string
}

var dbConnStrTests = []struct {
	name   string
	fields fields
	want   string
	want1  string
}{
	{
		name: "foo db conn",
		fields: fields{
			DBConn: "foo",
		},
		want:  "postgres",
		want1: "foo",
	},
	{
		name: "bar db conn",
		fields: fields{
			DBConn: "bar",
		},
		want:  "postgres",
		want1: "bar",
	},
	{
		name: "no db conn",
		fields: fields{
			DBConn: "",
		},
		want:  "postgres",
		want1: "",
	},
}

func TestDB_ConnStr(t *testing.T) {
	for _, tt := range dbConnStrTests {
		t.Run(tt.name, func(t *testing.T) {
			db := &config.DB{
				DBConn: tt.fields.DBConn,
			}
			got, got1 := db.ConnStr()
			if got != tt.want {
				t.Errorf("ConnStr() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ConnStr() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
