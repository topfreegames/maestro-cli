# maestro-cli
[Maestro](https://github.com/topfreegames/maestro) command line interface

## About
Maestro-cli calls Maestro api routes. Create, delete and update schedulers.

## Quickstart
* Download the latest release
* Init configuration
```
maestro init zooba http://server.url.com
```
* Create scheduler
```
maestro create path/to/config/file.yaml
```
* Delete scheduler
```
maestro delete scheduler-name
```
* Update scheduler
```
maestro update path/to/config/file.yaml
```
* Get scheduler status
```
maestro status scheduler-name
```
* Get scheduler releases
```
maestro releases scheduler-name
```
* Rollback to previous release
```
maestro rollback scheduler-name v1
```
* Edit a scheduler and update it
```
maestro edit scheduler-name
```
* The commands update, edit and set return an operationKey. Use this operationKey to get its progress or cancel it.
```
maestro progress operation-key
```
```
maestro cancel operation-key
```
