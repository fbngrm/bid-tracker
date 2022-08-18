
## REST
Transpiling the gRPC proto defintion into a HTTP1.1 JSON API would need much more thoughtful work to result in a "real" REST API according to a [Richardson Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html).
This is out of scope here and I simply utilized the code generation to free resources for focussing on the concurrency approach.

## Concurrency strategy
I used pessimistic locking for synchronization.
The [itemStore.go](./pkg/item/store.go) explains the concurrency strategy in more detail in the comments.
It also points out short-comings and potential optimizations.

channels vs locks...

go handlers ...

### Alternatives
I'd like to outline two interesting approaches I stumbled upon.

#### Optimistic Locking (transpiled)

#### Open Multithreaded Transactions

## Caveats
There a few places/types where I would operate on copies instead of pointers.
Tracking copying of pointers is error prone and should be avoided whenever possible.
This would be one of the first things I would optimize if there is was more time!

## Architecture
I roughly followed [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) approaches for package layout.
Dependency injection is used to achieve loose coupling, extensibility and testability.
With one exception, pointed out in the comments.
Data structures and types are designed around domain boundaries, roughly following [Domain Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) concepts.
The above approaches are not followed strictly since this is a prototype with tight time constraints.
Though, it lays foundation for growing into a more complex codebase.

## Event Soucing
Moving this project towards an event-driven (event-sourcing) architecture, the in-memory item store could be replaced by e.g. REDIS for storing items.
This would provide the possibility for handling and restoring internal state, allows for scaling and fault tolerance.
A message-broken would be used to communicate with other instances of the bid-tracker and services for user handling etc.
Though, the REST API brings in some restrictions towards an asynchronous communication, this could be an attractive approach, depending on the esact requirements.

## Notes
- I put more comments than I normally would since this is a coding challenge and I wanted to explain my reasoning.
- I worked slightly more than a day (8h) on the task but spend 3h only for research.
- I really liked the task and the volume is good for a coding challenge imo.

After implementing the prototype there are quite a few things that I would make differently when refactoring this towards a more mature version.
I would be happy to discuss my learnings, optimizations and short-comings I see, in depth, in a potential next interview :)
