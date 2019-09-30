# Docker Registry Sweeper

A CLI which helps you remove a reference between obsolete Docker images and their blobs for later removing by Docker built-in garbage collector

## How to use this CLI
### Run from Docker Image
```bash
docker run pirsquareff/dksweeper <options>
```

### Options
| Option 	| Env | Required | Description 	|
|-----------------------------	|----- |-----  |-------------------------------------------------------------------------------------------------------	|
| ``-u, --username <username>`` 	| ``DOCKER_USERNAME`` | ``false`` | Username to access Docker Registry 	|
| ``-p, --password <password>``	| ``DOCKER_PASSWORD`` | ``false`` | Password to access Docker Registry 	|
| ``--host <host>`` 	| ``REGISTRY_HOST`` | ``true`` | Docker Registry host with a protocol (http, https) 	|
| ``-r, --repo <repo>`` 	| ``REPOSITORY`` | ``true`` | Repository for cleanup 	|
| ``--max-age <max-age>`` 	| ``MAX_AGE`` | ``true`` | Keep images whose age is less than this value (in days). Older images will be unlinked with their blobs. 	|
| ``--keep-tag <n>`` 	| ``KEEP_TAG`` | ``true`` | Number of tags to be preserved. It's to prevent deleting entire images. For example, if all the tags are obsolete (by max-age), this parameter is to make sure that ``n`` latest tags will not be unlinked. 	|
| ``--verbose`` 	| ``VERBOSE`` | ``false`` | Print info during processing 	|

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
