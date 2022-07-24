# drone-zip
Drone CI  plugin, use archive/zip for compressed files.

<a href="https://github.com/LunaGod/drone-zip/actions/workflows/release.yml">
  <img src="https://github.com/LunaGod/drone-zip/actions/workflows/release.yml/badge.svg?tags=latest" alt="build status">
</a>

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
      - a/**/*.js # recursive match .js file
      - a/**/* # recursive match all file
    output: release.zip
```
