#!/bin/bash

# Variables
MIGRATIONS_DIR="./deployments/migration/"

read -p "Reason for migration:" args;

if [ -z "$args" ]; then
  echo "reason for migration cannot be empty"
  exit 1
fi

FILE_NAME=$(date +%s_)$(echo ${args} | tr " " _ )
FILE_PREFIX=${MIGRATIONS_DIR}/${FILE_NAME}

if [ -f ${FILE_PREFIX}.up.sql ]; then
    echo "${FILE_PREFIX}.up.sql exists."
    exit 1
fi

if [ -f ${FILE_PREFIX}.down.sql ]; then
    echo "${FILE_PREFIX}.down.sql exists."
    exit 1
fi

touch ${FILE_PREFIX}.up.sql
touch ${FILE_PREFIX}.down.sql

echo

echo "following files has been created in: ${MIGRATIONS_DIR} please add migrations"
echo ${FILE_NAME}.up.sql; echo ${FILE_NAME}.down.sql