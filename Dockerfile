# Dockerfile for the custom docker of the Database
# Base image
FROM postgres:17-bookworm

# Copy the sql script we want use over
COPY ./nostr-protocol-nips_nip-tags/init/init.sql ./docker-entrypoint-initdb.d
