#!/bin/bash
docker-compose -f services.yml --env-file vars.env up --build  --force-recreate 