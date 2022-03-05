#!/bin/bash

docker-compose -f services.yml  --env-file vars.env down && docker volume prune -f