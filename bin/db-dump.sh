#!/usr/bin/env fish

pg_dump --clean --schema-only --dbname=scrapygo_development > ./db/structure.sql
