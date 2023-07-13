# Module Driven Go Todolist

Unlike the so-called 'standard go structure', this project is to show that to
create a software that adheres on onion architecture in go, you don't need to
create a lot of objects and structure.

## Motivation

There are people coming to Go from different backgrounds. The most common people
who want to write Go in my experience they came from these backgrounds:

- Java & Spring Boot
- PHP & Laravel
- Ruby on Rails

They have this popular frameworks with a "canonical file structure" that you
must adhere coming to Go with no 'one true way' and no 'one true frameworks' is
unfamiliar.

When I saw projects on the wild, or when I read the projects from my clients,
there are patterns that I recognised. Most of them are trying to fit Go to their
familiar framework mental models.

So, I'm sick of seeing a convoluted code that tries to shove object orientation
to Go. This is why this project use __mostly functions__.

## The Principles of building this project

This project is not a perfect example of production-grade software. I create
this object with these principle in mind:

- Easy to understand and teach.
- Idiomatic. This means that I use Go like it's intended. Use value as much as
  possible and use patterns like 'accept interface returns value' and avoid
  things like returning interface.
- High cohesion, low coupling. Things that should be together should be in the
  same place.
- Sensible dependencies. I have justification why I use the dependencies.
- Sensible defaults. Every configuration has defaults that make the development easy.

This project is as plain as vanilla. I don't use any ORMs and strive to have
very minimal dependencies. Heck, this project doesn't have any interface
defined. If you see the `go.mod` file, you'll see that these are the dependies
that I use:

```
require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/jackc/pgx/v5 v5.4.1
	github.com/oklog/ulid/v2 v2.1.0
	github.com/rs/zerolog v1.29.1
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/yaml.v3 v3.0.1
)
```

I can explain what's the purpose of the dependencies above, and why I chose
those dependencies

1. `go-chi` is a routing library, it helps parses the request parameters and
   route the request. The reason I use this library is because it has zero
   indirect dependencies. Chi is a simple routing library and it does the job
   well.

2. `pgx` is a postgresql client. I don't use standard `database/sql` because
   usually I want the feature of postgresql that's not available on the lowest
   common denominator of database driver.

3. `ulid` is a library to generate and build
   [ULID](https://github.com/ulid/spec). I need this because I want the primary
   key of the domain is lexicographically sorted but also stateless and unique.

4. `null` is a library that I usually use to avoid `nil` pointer exception and
   avoid to use pointer when I can use a value instead. This dependency is
   optional.

5. `yaml` is a library to parse yaml file. This is also optional dependencies.
   If you don't like YAML you can change it to anything you want.

## How to navigate this project

### `main` package

In the root of this project there's a `main` package that have `main.go`. This
is the place where you put your main program as well as runtime configuration.

### The modules

Within the root project there are _modules_ implemented as Go package. Package
is a boundaries in Go. Everything inside this module is isolated. Inside this
module there's __no rules__ on how you organise files. However, in my project, I
usually have these:

1. __The domain object__. This defines the objects which maintain its state and
   being persisted. In this example it's in the `todo_item.go` file. Within the
   project there's `todo_item_json.go` this is for serialisation only. I prefer
   to separate the json representation on different file.

2. __The repository__. This is an abstraction where you can fetch and save your
   domain objects. I just name it 'repository' because it's _somewhat_ similar
   to repository pattern but I implemented it as with pure functions in
   `repo.go`.

3. __The storage__. This is just a place where you connect, read, and write on
   your repository. In this example, it's very simple, it's just a global `pool`
   object which represent a connection pool to postgresql instance. It's in the
   `db.go`.

4. __The read models__. This is types and structures to represent portion or
   aggregation of data from storage. You __cannot modify__ the read models
   because as name implies, it's for reading.

5. __The service__. This is a file with functions which define a _transaction
   boundary_. This service is agnostic with the protocols. __There should be no
   protocol-related data__ in here such as HTTP Response Code. Also implemented
   in pure functions in `service.go`.

6. __The protocols__. ~~This is a place where you put the handlers to your
   requests. I put it on separate package `handlers` because I don't want to use
   something like `CreateTodoItemHandler` and instead I can just use
   `handlers.CreateTodoItem`. The implication is that if you have many modules
   with the same name `handlers` then you still need to name it. So, you can
   also put a `handlers.go` within the directory, so that you can refer it to
   `todo.CreateItemHandler` and I think that's better than my approach here.~~
   Inside the module there's `routes.go`. Every routes for this module is here.
   This will export `*chi.Mux` object which is implementing `chi.Router`
   interface. This isolates any changes to the subrouter to this module. In
   `main.go` you can just mount the router.

> **Note**
> Previously, most of the functions are exported. On this version, I've made
> functions as private as I can.

Exported functions on each module area as follows:

- `SetPool()` for setting up the global database connection pool. 
- `Router()` for exporting the `chi.Mux` object to be mounted.
- The domain object, this is optional. If you want to hide and isolate your
  domain objects, then you can just make it private

## Testing And Faking

I'm rarely uses mocks. Read the rationale
[here](https://joeblu.com/blog/2023_06_mocks/). I do use fake. The good thing
about go it allows build tags which will allow conditional compilation of files.
Because I am using files to demarcate boundaries and responsibilities, this
works well.

Take a look is the `repo.go` and `repo_fake.go` for comparison.

To build or run the program with fake implementation, you can just use `--tags`
parameter, and the fake version will be compiled instead.

```
go build --tags=fake 
```

## Configuration

### Rationale 

The early version of this program didn't have configuration file and all
parameters are being hard-coded. At first I think that's enough. However, I
think giving an example of how to implement server who adheres to [12 Factor
App](https://12factor.net/) rules are important.

### Implementation

See `config.go`. This file contains code to parse configuration files from two
sources: environment variables and configuration files. Configuration files
takes precedence. These are the environment variables, configuration file key
path and default value.

| Environment Variable  | YAML keypath  | Default value | Description          |
|-----------------------|---------------|---------------|----------------------|
| `KAD_LISTEN_HOST`     | `listen.host` | "127.0.0.1"   | Server Listen Address|
| `KAD_LISTEN_PORT`     | `listen.port` | 8080          | Server Port Address  |
| `KAD_DB_HOST`         | `db.host`     | "127.0.0.1"   | Postgres Host        |
| `KAD_DB_PORT`         | `db.port`     | 5432          | Postgres Port        |
| `KAD_DB_NAME`         | `db.db_name`  | "todo"        | Database Name        |
| `KAD_DB_SSL`          | `db.ssl_mode` | "disable"     | SSL Mode             |

The default values, if we express it in configuration file is as follows.

```yaml
listen:
  host: 127.0.0.1
  port: 8080

db:
  db_name: todo
  host: 127.0.0.1
  port: 5432 
  ssl_mode: disable
```

### Configuration file location

The program will search for `config.yaml` on current working directory, or you
can pass `-c` flag to force the program to use your own configuration file name.
For example you can run it using something like this:

```
./mda -c someconfig.yml
```

## Summary

This project is a heuristic, not a guide or a 'framework' of structure. It's to
show that you can have a sensible code structure and architecture by sticking
to the simplicity of Go.

## LICENSE

```
Copyright (c) 2023 Didiet Noor

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
