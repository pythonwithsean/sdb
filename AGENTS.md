# AGENTS.md

## Project purpose

This is a learning project for building a relational database from scratch.
The primary goal is to understand how databases work, not to reach production
readiness as quickly as possible.

Version 1 is intentionally a small, experimental implementation. Each future
version should grow from the previous one as the project owner gains a deeper
understanding of database concepts.

## Working with the project owner

- Only build exactly what the project owner explicitly requests.
- Do not work on changes that were not requested, even if they appear useful,
  necessary for completeness, or like obvious next steps.
- Do not add features, abstractions, optimizations, dependencies, tests,
  documentation, or refactors without an explicit request.
- If a request is ambiguous or has multiple meaningful design choices, explain
  the choices and ask before implementing.
- Prefer a simple implementation that is easy to read and learn from over a
  clever or highly optimized implementation.
- The project owner is building and learning the project. The AI's role is
  only to write code for the exact component the owner requests; it must not
  take ownership of the project or decide what gets built next.
- Keep every component explicit and understandable so the project owner can
  study it and build it themselves.
- Do not hide important database behavior behind external libraries unless the
  project owner explicitly asks to use one.

## Learning-first implementation style

When the project owner requests a change:

1. Explain the database concept being introduced.
2. Describe the smallest reasonable implementation approach.
3. Confirm the exact requested scope when it is ambiguous.
4. Implement only that requested scope.
5. Point out important trade-offs and limitations without implementing
   unrequested follow-up work.

Teach the underlying ideas, including where relevant:

- pages, records, schemas, and serialization;
- tables, rows, columns, and data types;
- primary keys, indexes, and constraints;
- scanning, lookup, insertion, update, and deletion;
- parsing and executing queries;
- transactions, durability, recovery, and concurrency;
- relational design and query planning.

Do not imply that a feature is production-safe when it is only a learning
implementation. Clearly label incomplete behavior and known limitations.

## Versioning expectations

- Treat each version as a deliberate learning milestone.
- Avoid implementing later-version features early.
- Keep changes small enough that the project owner can understand the full
  implementation.
- When a new feature depends on a concept not yet implemented, explain the
  dependency instead of silently adding the prerequisite system.
- Preserve earlier behavior unless the project owner explicitly approves a
  breaking change.

## Go project conventions

- Follow standard Go formatting and idioms.
- Prefer the standard library while the project is learning database
  fundamentals.
- Keep the code explicit and easy to trace.
- Use clear names for database concepts and avoid unnecessary indirection.
- Add tests only when the project owner explicitly asks for them.
- Run validation only when requested, or when it is required to confirm that
  the explicitly requested change works.

## Scope control

Before changing anything, identify:

- what the owner explicitly requested;
- which files need to change;
- which behavior will be added or changed;
- what is intentionally out of scope.

If a useful improvement is outside the request, do not implement it. It may be
mentioned only when relevant to explaining the requested change, and must
remain separate from the implementation.
