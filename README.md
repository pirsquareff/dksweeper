# Docker Registry Sweeper

A CLI which helps you remove a reference between obsolete Docker images and their blobs for later removing by Docker built-in garbage collector

## How to use this CLI

### Options
| Option 	| Required | Description 	|
|-----------------------------	|-----  |-------------------------------------------------------------------------------------------------------	|
| ``-u, --username <username>`` 	| ``false`` | Username to access Docker Registry 	|
| ``-p, --password <password>``	| ``false`` | Password to access Docker Registry 	|
| ``--host <host>`` 	| ``true`` | Docker Registry host with a protocol (http, https) 	|
| ``-r, --repo <repo>`` 	| ``true`` | Repository for cleanup 	|
| ``--max-age <max-age>`` 	| ``true`` | Keep images whose age is less than this value (in days). Older images will be unlinked with their blobs. 	|
| ``--keep-tag <n>`` 	| ``true`` | Number of tags to be preserved. It's to prevent deleting entire images. For example, if all the tags are obsolete (by max-age), this parameter is to make sure that ``n`` latest tags will not be unlinked. 	|
| ``--verbose`` 	| ``false`` | Print info during processing 	|

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
