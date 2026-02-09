# GO Blog Aggregator

A Go-powered blog aggregator that lets users follow feeds, browse posts, and manage subscriptions. Built with Drizzle ORM and designed for extensibility.

## Features

- **Feed Management:** Add, list, and remove RSS/Atom feeds.
- **User Subscriptions:** Users can follow or unfollow feeds.
- **Post Aggregation:** Bulk import posts from feeds, with upsert support for duplicates.
- **Browse Posts:** View recent posts from feeds a user follows.
- **Command-Line Interface:** Interact with the aggregator using simple commands.

## Getting Started

### Prerequisites

- Node.js (v18+ recommended)
- PostgreSQL (or compatible SQL database)

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/go-blog-aggregator.git
   cd go-blog-aggregator
   ```

2. Install dependencies:

   ```sh
   npm install
   ```

3. Set up your database:
   - Create a .gatorconfig.json file at home-dir with your database connection string:
     ```
     "db_url": "postgres://example"
     ```
   - Run migrations:
     ```sh
     goose postgres <database_url> up
     ```

### Usage

#### Build the project

```sh
go build
```

#### Register a User

```sh
 ./go-blog-aggregator register <username>
```

#### Login

```sh
 ./go-blog-aggregator login <username>
```

#### Reset the database

```sh
 ./go-blog-aggregator reset
```

#### List Users

```sh
 ./go-blog-aggregator users
```

#### Aggregate Posts

```sh
 ./go-blog-aggregator agg <time_between_reqs>
```

#### Add a Feed

```sh
 ./go-blog-aggregator addfeed <feed-name> <url>
```

#### List Feeds

```sh
 ./go-blog-aggregator feeds
```

#### Follow a Feed

```sh
 ./go-blog-aggregator follow <url>
```

#### Unfollow a Feed

```sh
 ./go-blog-aggregator unfollow <url>
```

#### List Followed Feeds

```sh
 ./go-blog-aggregator following
```

#### Browse Posts

```sh
 ./go-blog-aggregator browse
```

Lists recent posts from feeds the user follows.

## Development

- All source code is in Go.
- Database schema and migrations are managed with Goose.
- Queries are generated using sqlc
- Extend commands or add new features by editing the CLI and schema files.
