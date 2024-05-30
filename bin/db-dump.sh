#!/usr/bin/env fish

pg_dump --schema-only --dbname=scrapygo_development > ./db/structure.sql
