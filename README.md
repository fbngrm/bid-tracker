The solution implements all requirements from the task definition.

I chose to implement a `gRPC` `API`.
The gRPC API naturally uses `protocol buffers` and `HTTP/2`.
However, an `HTTP` `API` that supports `JSON` data via `HTTP/1.1` is also provided by mapping the proto definition to the `OpenAPI` standard.

The API is described in a [.proto file](./apis/auction/v1/auction.proto).
The JSON API is documented at `http://localhost:8081/openapi/`.
Enter `/openapi.json` in the box at the top of the page to browse the endpoint definitions.
Note, `TLS` encryption is not implemented, thus requests made via the Open Api UI fail.

## REST
Transpiling the gRPC proto definition into a HTTP1.1 JSON API would need much more thoughtful work to result in a "real" REST API according to a [Richardson Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html).
This is out of scope here and I simply utilized the code generation to free resources for focussing on the concurrency approach.

The `just` command runner is used to automate command execution.
A `Makefile` is provided for making build dependencies to avoid re-running time-consuming code generation steps.

## Running the app

Either just use the just task runner or copy the commands from the [Justfile](./Justfile).

```
just run-server
```

The app logs seed data and commands to interact with the API as a starting point.

## Concurrency strategy
I used pessimistic locking for synchronization.
The [itemStore.go](./pkg/item/store.go) explains the concurrency strategy in more detail in the comments.
It also points out short-comings and potential optimizations.

## Caveats
There a few places/types where I would operate on copies instead of pointers.
Tracking copying of pointers is error prone and should be avoided whenever possible.
This would be one of the first things I would optimize if there is was more time!

## Architecture
I roughly followed [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) approaches for package layout.
Dependency injection is used to achieve loose coupling, extensibility and testability (exception, pointed out in the comments).
Data structures and types are designed around domain boundaries, roughly following [Domain Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) concepts.
The above approaches are not followed strictly since this is a prototype with tight time constraints.
Though, it lays foundation for growing into a more complex codebase.

## Tests
I added only a few example unit and benchmark tests for the item store which holds the most complex code.
Due to time constranits, I didn't add more tests.
Subtests are run in parallel and tests contain programmatic concurrency.
There is a Just recipe available to run tests in parallel and test for race conditions.

```
just test
just benchmark
```

Additionally, the program prints some instructions on how to use apache-bench for load-testing.

## Notes
- I put more comments than I normally would since this is a coding challenge and I wanted to explain my reasoning.
- I worked slightly more than a day (8h) on the task but spend 3h only for research.
- I really liked the task and the volume is good for a coding challenge imo.
- I added a binary to the repo for review purposes which I normally wouldn't ship into version control. It increases the size of the repo also.

## Closing Thoughts

After implementing the prototype there are quite a few things that I would make differently when refactoring this towards a more mature version.
I would be happy to discuss my learnings, optimizations and short-comings, in depth, in a potential next interview :)

### Event Sourcing
Moving this project towards an event-driven (event-sourcing) architecture, the in-memory item store could be replaced by e.g. REDIS for storing items.
This would provide the possibility for handling and restoring internal state, allows for scaling and fault tolerance.
A message-broken would be used to communicate with other instances of the bid-tracker and services for user handling etc.
Even though the REST API brings in some restrictions towards an asynchronous communication, this could be an attractive approach, depending on the exact requirements.

### Alternative Concurrency Models
I'd like to mention these two interesting approaches I stumbled upon during research.

#### Optimistic Locking (transpiled)
This concept transpiles lock based pessimistic concurrency into optimistic concurrency.
This could be an interesting approach for posting pessimistic locking to an append only model in the context of event sourcing.
A [slide](https://www.usenix.org/system/files/atc21_slides_zhang-zhizhou.pdf) summing up the approach.
And the original [publication](https://arxiv.org/pdf/2106.01710.pdf).

#### Open Multithreaded Transactions
Nested transactions to guarantee consistent state and fault tolerance. Probably not easily scaleable.
The [paper](http://130.203.136.95/viewdoc/download?doi=10.1.1.20.1875&rep=rep1&type=pdf) describing the approach.


