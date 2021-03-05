#!/bin/bash
docker kill surf_be
docker rm surf_be
docker-compose up -d