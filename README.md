# Carlang

Carlang is an unofficial Invoker-inspired esoteric programming language.

Programs combine **Quas**, **Wex**, and **Exort** into three-orb recipes, invoke those recipes with `R`, and cast from two spell slots with `D` and `F`.

> This is a fan project and is not affiliated with Valve, Dota 2, or any official Dota property.

## Quick Start

```bash
go test ./...
go run ./cmd/carl run examples/hello_world.carl
```

Expected output:

```txt
Hello, World!
```

Build the CLI:

```bash
go build -o bin/carl ./cmd/carl
```

Common commands:

```bash
bin/carl run examples/hello_world.carl
bin/carl compile examples/hello_world.carl -o dist/hello_world.cbc
bin/carl vm dist/hello_world.cbc
bin/carl dis examples/hello_world.carl
bin/carl repl
```

## Docs

- [Language guide](docs/LANGUAGE.md): syntax, runtime model, spells, combos, and Hello World notes.
- [Language specification](docs/SPEC.md): grammar-oriented reference.
- [Compiler notes](docs/COMPILER.md): bytecode pipeline, compile-time expansion, and VM opcodes.
- [Project layout](docs/PROJECT.md): package map for the repository.
- [Status](docs/STATUS.md): implemented features and current compiler limits.

## License

MIT. See [LICENSE](LICENSE).
