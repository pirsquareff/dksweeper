# Docker Registry Sweeper

A CLI which helps you remove a reference between obsolete Docker images and their blobs for later removing by Docker built-in garbage collection

## How to use this CLI
### Run from Docker Image
```bash
docker run pirsquareff/dksweeper <options>
```

### Options
| Option | Env | Required | Description |
|--------|-----|----------|-------------|
| ``-u, --username <string s>`` | ``DOCKER_USERNAME`` | ``false`` | Username to access Docker Registry |
| ``-p, --password <string s>``	| ``DOCKER_PASSWORD`` | ``false`` | Password to access Docker Registry |
| ``--host <string s>`` | ``REGISTRY_HOST`` | ``true`` | Docker Registry host with a protocol (http, https) |
| ``-r, --repo <string s>`` | ``REPOSITORY`` | ``true`` | Repository for cleanup |
| ``--older-than <integer n>`` | ``OLDER_THAN`` | ``true`` | Delete images older than this value (in days) |
| ``--keep-tag <integer n>`` | ``KEEP_TAG`` | ``true`` | Number of tags to be preserved. It's to prevent deleting entire images in a repo. For example, if all the tags are obsolete (by --older-than), this parameter is to make sure that ``n`` latest tags will not be unlinked. |
| ``--verbose`` | ``VERBOSE`` | ``false`` | Print info during processing |


## How to run garbage collection in Docker Registry
```bash
bin/registry garbage-collect /etc/docker/registry/config.yml
```

Change a path to a config.yml file to match with your configuration. For further information, please consult an official documentation [here](https://docs.docker.com/registry/garbage-collection/).

> Note: You should ensure that the registry is in read-only mode or not running at all. If you were to upload an image while garbage collection is running, there is the risk that the image’s layers are mistakenly deleted leading to a corrupted image.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/pirsquareff/dksweeper/blob/master/LICENSE) file for details
