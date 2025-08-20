<p align="center">
  <img src="identity/logo_insights.svg" />
</p>

# Artio Insights

Artio Insights is a aggregator of Nostr data pulling from external sources. It has a frontend component that illustrates the data from the external sources and the [Artio Relay](https://github.com/SEG-UNIBE/artio-relay) into graphs.
Another component is the DV (Data Vault) that allows the aggregation of data from multiple sources into a single data model.
## Development 

The project is organized into different modules. Each one of them is a folder with this nomenclature : `{data-source}_{data-type}`. For example, `nostr-protocol-nips_tags` pulls NIP *tags* data from the github repo [*nostr-protocol/nips*](https://github.com/nostr-protocol/nips).

Each module is and should be independent from the others, only sharing the database docker container, the Grafana instance and the `.env` file. The modules can be written in any language of choice.

To add another Go module, the easiest way would be to copy the current module, remove the specific code and keep the file structure. This way is best because the database connection is already coded properly, the sql file is ran and many checks are already made.

## Running

### Prerequisites

- Docker
- Go (if running the module bare bones)

### Instructions

#### General

1. Copy the `.env-example` into a `.env` file and fill out the information in the `.env` file.
2. Build and run the docker containers : `docker compose build & docker compose up`

#### nostr-protocol-nips_tags

The database must be running for this module to work properly.

##### Bare Bones

1. Move to the corresponding folder: `cd nostr-protocol-nips_nip-tags`
2. Download the packages: `go mod download`
3. Run the app: `go run cmd/main.go`

##### In Docker

Use the `docker-compose.yml` file.
