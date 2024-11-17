#!/bin/bash

duckdb -c "create or replace table first_names as select #1 as name from 'sources/first_names.csv'" sources.db
duckdb -c "create or replace table last_names as select #1 as name from 'sources/last_names.csv'" sources.db
