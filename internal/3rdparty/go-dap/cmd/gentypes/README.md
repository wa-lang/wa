The schema for DAP messages is defined in JSON at
https://github.com/microsoft/vscode-debugadapter-node/blob/main/debugProtocol.json

The auto-generated TypeScript representation of the schema is at
https://github.com/microsoft/vscode-debugadapter-node/blob/main/protocol/src/debugProtocol.ts

----

In this directory we have a copy of the schema, which is licensed by Microsoft
with a [MIT
License](https://github.com/microsoft/vscode-debugadapter-node/blob/main/License.txt).
This copy must be updated whenever the schema changes.

To generate Go types from the schema, run:

```
$ go run cmd/gentypes/gentypes.go cmd/gentypes/debugProtocol.json > schematypes.go
```

The generated ``schematypes.go`` is also checked in, so there is no need to
regenerate it unless the schema changes.
