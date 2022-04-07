#!/usr/bin/env bash
docker-compose -p ohurlshortener -f services.yml --env-file vars.env up -d --build --force-recreate
