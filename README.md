# Docker Registry Sweeper

A CLI tool which helps you clean up your Docker Registry by deleting outdated images

Deep diving in how this tool works, it uses the registry REST API to retrieve all the tags along with their creation date. After that, it makes a DELETE request to your Docker Registry for all the tags met with requirement specified. This tool provides two adjustable parameters for filtering images to be deleted: ``--older-than`` and ``--keep-tag``. See [Options section](#options) below for their description.

This tool was implemented using the Go programing language to utilize its concurrency feature. Doing so, it makes fetching and deleting tags faster and more efficiently.

Please note that it only removes references to images and makes them eligible for garbage collection. To remove them completely from the filesystem, you should explicitly run Docker Registry built-in garbage collection. See [How to run garbage collection in Docker Registry section](#how-to-run-garbage-collection-in-docker-registry).

## How to use this CLI
### Run from Docker Image
The official image was published to Docker Hub. You can find it at https://hub.docker.com/r/pirsquareff/dksweeper.

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
| ``--keep-tag <integer n>`` | ``KEEP_TAG`` | ``true`` | Number of tags to be preserved. It's to prevent deleting entire images in a repo. For example, if all the tags are condisered obsolete (by --older-than), this parameter is to make sure that ``n`` latest tags will not be deleted. |
| ``--verbose`` | ``VERBOSE`` | ``false`` | Print info during processing |


## How to run garbage collection in Docker Registry
Run the following command in your Docker Registry server

```bash
bin/registry garbage-collect /etc/docker/registry/config.yml
```

Change a path to a config.yml file to match with your configuration. For further information, please consult the Docker official document [here](https://docs.docker.com/registry/garbage-collection/).

> Note: You should ensure that the registry is in read-only mode or not running at all. If you were to upload an image while garbage collection is running, there is the risk that the image’s layers are mistakenly deleted leading to a corrupted image.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/pirsquareff/dksweeper/blob/master/LICENSE) file for details
