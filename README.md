# maestro-cli
[Maestro](https://github.com/topfreegames/maestro) command line interface

[![codecov](https://codecov.io/gh/topfreegames/maestro-cli/branch/next/graph/badge.svg?token=IJIA498X2D)](https://codecov.io/gh/topfreegames/maestro-cli)

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
* Create scheduler version
```
maestro create path/to/config/file.yaml
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
* Switch active schedule version
```
maestro-cli switch active-version scheduler-name v2.0.0
```
* The commands update, edit and set return an operationKey. Use this operationKey to get its progress or cancel it.
```
maestro progress operation-key
```
```
maestro cancel operation-key
```
