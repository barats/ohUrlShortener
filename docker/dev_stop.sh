#!/usr/bin/env bash
docker-compose -p ohurlshortener -f dev.yml  --env-file vars.env down