# drone-zip
Drone CI  plugin, use archive/zip for compressed files.

![Github Action](https://github.com/LunaGod/drone-zip/actions?query=workflow%3A%22Publish+release%22)

### Drone CI Plugin Config

`1.x` or `2`
```yaml
steps:
...
- name: zip
  image: lunagod/drone-zip
  settings:
    input: 
      - a.txt
      - a/*.js # globs are allowed
    output: release.zip
```
