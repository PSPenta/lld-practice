# Low Level System Design (LLD) Solutions

This repository contains solutions to various Low Level System Design questions. Each folder in this repository represents a separate LLD problem and its code solution.

## Table of Contents

- [Overview](#overview)
- [Folders and Files](#folders-and-files)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Overview

This repository showcases solutions to various Low Level System Design problems. Each folder in this repository represents a separate LLD problem and its code solution. The problems covered include:

- SplitWise: A system for splitting bills among a group of people.
- LRU Cache: An implementation of an LRU cache data structure.
- Parking Lot: A system for managing a parking lot, including parking, unparking, and showing availability.
- Polling System: A system for creating and managing polls with multiple options and users.
- Rate Limiter: An implementation of a rate limiter using the token bucket algorithm.
- Search Engine: An implementation of a search engine, including tokenizing documents, creating inverted indexes, and ranking search results.
- URL Shortener: A system for shortening URLs using a custom algorithm.
- Database: An implementation of a simple key-value store using a Map.
- Pub-Sub: An implementation of a publish-subscribe system using an event emitter.

## Folders and Files

The following folders and files are included in this repository:

- `SplitWise`: Contains code for managing bill splitting among a group of people.
- `LRU`: Contains code for implementing an LRU cache.
- `ParkingLot2`: Contains code for managing a parking lot, including parking, unparking, and showing availability.
- `PollingSystem`: Contains code for managing polls, including creating polls, submitting votes, and showing statistics.
- `Ratelimiter`: Contains code for implementing a rate limiter using the token bucket algorithm.
- `SearchEngine`: Contains code for implementing a search engine, including tokenizing documents, creating inverted indexes, and ranking search results.
- `UrlShortener`: Contains code for shortening URLs using a custom algorithm.
- `Database`: Contains code for implementing a simple key-value store using a Map.
- `Pub-Sub`: Contains code for implementing a publish-subscribe system using an event emitter.

## Usage

To use the SplitWise system, create a `SplitWise` instance and use the methods provided to manage bill splitting among a group of people.

To use the LRU cache, create an `LRUCache` instance and use the `set` and `get` methods to store and retrieve data.

To use the parking lot system, create a `ParkingLot` instance and use the methods provided to manage parking, unparking, and showing availability.

To use the polling system, create a `Poll` instance and add it to the `Polls` module.

To use the rate limiter, create a `RateLimiter` instance and call the `limitConcurrency` method with a list of tasks and the desired concurrency.

To use the search engine, create a `SearchEngine` instance and call the `search` method with a query string.

To use the URL shortener, create a `UrlShortener` instance and use the `shorten` and `expand` methods to shorten and expand URLs.

To use the database, create a `Database` instance and use the `set` and `get` methods to store and retrieve data.

To use the pub-sub system, create an `EventEmitter` instance and use the `emit` and `on` methods to publish and subscribe to events.

## Contributing

Contributions are welcome! Please follow the guidelines in CONTRIBUTING.md.

## License

This project is licensed under the MIT License.
