# monodep
Check what folders have changed between commits

Given a yaml file of something like:

```yaml
api:
  path: ./api
  deps:
    - db
db:
  path: ./db
  deps:
    - schemas
schemas:
    path: ./schemas
    deps: []
```

in a folder structure like
```
/
  /api
  /db
  /schemas
  deps.yaml
```

monodep will tell you what need to be recompiled given two git hashes.

Example: Two hashes where the only difference is that schemas/file.txt has changed
```
> monodep deps.yaml 6038f10 4862f02
api
db
schemas
```
