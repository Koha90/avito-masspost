# Avito Masspost

```text
avito-masspost/
├─ cmd/
│  ├─ feedgen/
│  │  └─ main.go
│  ├─ tokencheck/
│  │  └─ main.go
│  └─ server/
│     └─ main.go
├─ internal/
│  ├─ avito/
│  │  ├─ auth.go
│  │  ├─ client.go
│  │  └─ report.go
│  ├─ config/
│  │  └─ config.go
│  ├─ feed/
│  │  ├─ model.go
│  │  └─ xml.go
│  ├─ listing/
│  │  ├─ model.go
│  │  ├─ repository.go
│  │  └─ service.go
│  └─ textgen/
│     ├─ title.go
│     ├─ description.go
│     └─ similarity.go
├─ migrations/
│  └─ 000001_init.sql
├─ feeds/
│  └─ avito.xml
├─ config.example.toml
├─ go.mod
├─ go.sum
└─ README.md
```

`avito-masspost` is a Go service for managing and publishing ad variants for a concrete supply business.

The project focuses on four main tasks:

1. storing domain data about products, zones and listings;
2. generating unique titles and descriptions for listing variants;
3. exporting listing data into a feed format suitable for marketplace upload;
4. synchronizing publication state with the external platform.

## Purpose

The service is designed for a business account that publishes concrete supply listings for private customers.

The system should help:

- keep listings structured and manageable;
- avoid chaotic manual editing;
- prepare non-repeating title and description variants;
- track listing state from draft to export;
- provide a stable foundation for future feed export and synchronization.

## Principles

The project follows these principles:

- idiomatic Go naming;
- no `GetSomething` methods;
- exported packages, types, functions and methods must be documented;
- domain model first, infrastructure second;
- explicit and readable code over magic;
- stable identifiers for all publishable entities;
- safe growth from simple local generation to a full production workflow.

## Scope

At the current stage the project includes:

- project conventions and repository structure;
- domain model for listings, products and zones;
- repository contracts for storage integration;
- a foundation for future feed export and text generation.

Planned next steps:

- PostgreSQL storage;
- feed generation;
- title and description generation;
- uniqueness checks for texts;
- synchronization with external publication status.

## Domain language

The project uses the following core terms.

### Product

A product is a concrete mix being offered to a customer.

Examples:

- concrete M200;
- concrete M300;
- concrete B22.5 P3 F200 W8.

A product describes the material itself, not the publication.

### Zone

A zone describes where the offer is relevant.

Examples:

- Moscow;
- Moscow region;
- a specific district or delivery area.

A zone affects text, delivery conditions and publication scope.

### Listing

A listing is a publishable ad unit.

A listing belongs to one product and one zone and contains the text and price that will be exported.

### Variant

A variant is a concrete textual form of a listing.

Different variants may describe the same product and zone, but use different titles and descriptions.

### Status

A listing changes status during its lifecycle:

- `draft`
- `ready`
- `exported`
- `archived`

## Architecture

The project follows a simple layered structure.

- `cmd/` contains application entrypoints;
- `internal/listing/` contains domain model and repository contracts;
- `internal/textgen/` contains title and description generation;
- `internal/feed/` contains feed building;
- `internal/avito/` contains external API clients;
- `internal/config/` contains application configuration.

The domain layer must not depend on:

- Avito API details;
- XML structures;
- PostgreSQL;
- HTTP transport.

Infrastructure depends on domain, not the other way around.

## Domain rules

The domain model should enforce the following rules:

- listing ID must be stable;
- title must not be empty;
- description must not be empty;
- price must be positive;
- listing cannot become `ready` without title and description;
- archived listings are not exported;
- product and zone are separate concepts and must not be mixed into one blob.

## Naming style

The project uses short and meaningful Go names.

Good examples:

- `ByID`
- `Save`
- `List`
- `Load`
- `Write`
- `Title`
- `Status`
- `Ready`

Bad examples:

- `GetListingByIdentifier`
- `GenerateUniqueAdvertisementTextData`
- `FetchAndUpdatePublicationStatus`

## Code style

Rules for this repository:

- document all exported code;
- keep package responsibilities narrow;
- keep constructors explicit;
- prefer value objects where they improve clarity;
- keep domain types readable and stable;
- avoid premature abstractions;
- write code that is easy to test.

## Development path

The expected implementation order is:

1. README and domain model;
2. PostgreSQL repositories;
3. feed export;
4. text generation;
5. uniqueness checks;
6. external synchronization.

## Status

This repository is at the foundation stage.

The first priority is to build a clean domain core that will remain stable while infrastructure grows around it.
