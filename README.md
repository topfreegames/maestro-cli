# maestro-cli
[Maestro](https://github.com/topfreegames/maestro) command line interface

## About
Maestro-cli calls Maestro api routes. Create, delete and update schedulers. All calls are authenticated with Google Oauth2, so it's necessary to login before. 

## Quickstart
* Download the latest release
* Login
```
maestro login http://server.url.com
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
