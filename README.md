# drone-zip
Woodpecker CI (or Drone CI)  plugin, use [compress](https://github.com/klauspost/compress) for compressed files.

<a href="https://github.com/loosheng/drone-zip/actions/workflows/release.yml">
  <img src="https://github.com/loosheng/drone-zip/actions/workflows/release.yml/badge.svg?tags=latest" alt="build status">
</a>
<a href="https://github.com/loosheng/drone-zip/actions/workflows/ci.yml">
  <img src="https://github.com/loosheng/drone-zip/actions/workflows/ci.yml/badge.svg?tags=latest" alt="build status">
</a>

  <a href="https://hub.docker.com/r/lunagod/drone-zip" title="Docker pulls">
    <img src="https://img.shields.io/docker/pulls/lunagod/drone-zip">
  </a>

## Use in Woodpecker-CI
```yaml
pipeline:
  zip:
    image: lunagod/drone-zip
    settings:
      input: 
      - a.txt
      - a/*.js # globs are allowed
      - a/**/*.js # recursive match .js file
      - a/**/* # recursive match all file
      - ./a # recursively compress the a folder
      output: release.zip
```
## Use in Drone-CI

Drone CI version `1.x` or `2`

```yaml
steps:
- name: zip
  image: lunagod/drone-zip
  settings:
    input: 
      - a.txt
      - a/*.js # globs are allowed
      - a/**/*.js # recursive match .js file
      - a/**/* # recursive match all file
      - ./a # recursively compress the a folder
    output: release.zip
```
