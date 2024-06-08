#!/usr/bin/env fish

pg_dump -T goose_db_version -T goose_db_version_id_seq --clean --schema-only --dbname=scrapygo_development > ./db/structure.sql
