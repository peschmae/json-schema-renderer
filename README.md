# JSON schema renderer

Minimal CLI to create a human readable documentation from a JSON schema.

The CLI currently supports outputting Markdown and Asciidoc, either to stdout
or write into a file.

## Usage

Output asciidoc to stdout
```sh
json-schema-renderer schema.json
```

Write markdown to file
```sh
json-schema-renderer schema.json -f markdown -o schema.md
```

Pipe schema to CLI
```sh
cat schema.json | json-schema-renderer -o rendered.adoc
```
This is mainly useful to combine with other tools, eg into somthing like this

```sh
ytt --data-values-schema-inspect=true --output=openapi-v3 -f schema.yaml  | openapi-to-json-schema | json-schema-renderer -o schema.adoc
```
