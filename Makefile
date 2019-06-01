SMTP_PASS = 
DB_HOST = db
DB_USER = 
DB_PORT = 3306
DB_PASS = 
TEST_SECRET = 
APP_NAME = anon-solicitor
SEND_EMAILS = false

.PHONY: test
test:
	@echo "Running tests"
	@GO_ENV=test APP_NAME=$(APP_NAME) go test ./...

.PHONY: dev
dev:
	@echo "Running dev server"
	@GO_ENV=development \
		DB_HOST=$(DB_HOST) \
		DB_USER=$(DB_USER) \
		DB_PORT=$(DB_PORT) \
		DB_PASS=$(DB_PASS) \
		ANON_SOLICITOR_SECRET=$(TEST_SECRET) \
		APP_NAME=$(APP_NAME) \
		SEND_EMAILS=$(SEND_EMAILS) \
		SMTP_PASS=$(SMTP_PASS) \
		go run ./cmd/app/main.go