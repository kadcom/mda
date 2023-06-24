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

## How to navigate this project

### `main` package

In the root of this project there's a `main` package that have `main.go`. This
is the place where you put your main program as well as runtime configuration.

### The modules 

Within the root project there are `modules` implemented as Go package. Package
is a boundaries in Go. Everything inside this module is isolated. Inside this
module there's no rule on how you organise files. However, in my project, I
usually have these:

1. __The domain object__. This defines the objects which maintain its state and
   being persisted. In this example it's in the `todo_item.go` file.

2. __The repository__. This is an abstraction where you can fetch and save your
   domain objects. I just name it 'repository' because it's _somewhat_ similar
   to repository pattern but I implemented it as with pure functions in
   `repo.go`.

3. __The storage__. This is just a place where you connect, read, and write on
   your repository. In this example, it's very simple, it's just a global `pool`
   object which represent a connection pool to postgresql instance. It's in the
   `db.go`.

4. __The service__. This is a file with functions which define a _transaction
   boundary_. This service is agnostic with the protocols. __There should be no
   protocol-related data__ in here such as HTTP Response Code. Also implemented
   in pure functions in `service.go`.

5. __The protocols__. This is a place where you put the handlers to your
   requests. I put it on separate package `handlers` because I don't want to use
   something like `CreateTodoItemHandler` and instead I can just use
   `handlers.CreateTodoItem`. The implication is that if you have many modules
   with the same name `handlers` then you still need to name it. So, you can
   also put a `handlers.go` within the directory, so that you can refer it to
   `todo.CreateItemHandler` and I think that's better than my approach here.

## Todo

- [ ] Testing by using build tags to avoid interface.

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
