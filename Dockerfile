# Dockerfile for the custom docker of the Database
# Base image
FROM postgres:17-bookworm

# Copy the sql script we wann use over
COPY ./init/init_nips_tags.sql ./docker-entrypoint-initdb.d





