#!/usr/bin/env bash
docker-compose -p ohurlshortener -f pull_build.yml  --env-file vars.env down