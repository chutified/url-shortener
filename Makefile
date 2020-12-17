migrate-up:
	migrate -database $(URL_SHORTENER_DBCONN) -path schema/ up

migrate-down:
	migrate -database $(URL_SHORTENER_DBCONN) -path schema/ down
