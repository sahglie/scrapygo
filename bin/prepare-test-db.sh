#!/usr/bin/env fish

GOOSE_MIGRATION_DIR="./db/schema" goose postgres "postgres://app_scrapygo@localhost:5432/scrapygo_test?sslmode=disable" up


