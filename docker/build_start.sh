# !/bin/bash
docker-compose -f services.yml --env-file vars.env up -d --build  --force-recreate 