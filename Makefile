.PHONY: migrate-up migrate-down

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

migrate-up:
	migrate -database $(URL_SHORTENER_DBCONN)?sslmode=disable -path sql/schema/ up

migrate-down:
	migrate -database $(URL_SHORTENER_DBCONN)?sslmode=disable -path sql/schema/ down
