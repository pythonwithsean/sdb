# sdb — a learning database

sdb is a small, experimental relational database built from scratch in Go.
It is a learning project: the goal is to understand how databases work
internally, not to build a production-ready system.

## Goal

A simple database with a simple notation, supporting the fundamentals first:

- create a database
- create a table
- select rows
- insert rows
- update rows
- delete rows

Later milestones will explore durability and performance optimizations —
things like pages, records, indexes, and crash recovery — one concept at a
time.

## Layout

- `QueryProcessor` — parsing and executing the simple query notation
- `StorageEngine` — how data is stored on disk and in memory
- `Transport` — how commands reach the database (e.g. a client/server interface)

## Versioning

Each version is a deliberate learning milestone. New features are added only
as the existing ones are understood, so the project grows step by step.
