# go-qa-service
A simple API microservice for monitoring the test coverage of Go services

This service uses DynamoDB for persistant storage.

## Configuration
The server uses the [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) will load configuration from environment variables, 
AWS shared configuration file (\~/.aws/config), and AWS shared credentials file (\~/.aws/credentials).

The DynamoDB table must be specifide with the environment variable `QA_TABLE_NAME` corresponding to the table name in the configured default AWS region.

The server must be launched with the specified port as an arguement.
```console
./go-qa-service 8080
```
Ensure that the credentials of the server enable reading and writing items to DynamoDB.

## Routes
The api has the following routes:  
`/` This is used for updating service coverage. The request body must be a JSON in the following format:
```json
{
  "payload": { 
    "service_name": "test-service", 
    "coverage": 80}
}
```
`/api/stats` This returns a json formatted response of all the services and their coverage.  
`/stats` This returns an HTML response of all the services and thier coverage.

## Events
Calls to update the service coverage creates an AMPQ message that are published to the `qa.events` exchange on `amqp://guest:guest@localhost:5672/`.  
The server by default creates a consumer that logs events to stdout.

## Exaple pipeline
Here is a [github action workflow](https://github.com/devinmarder/go-test/blob/main/.github/workflows/publish-coverage.yml) 
that invokes an ec2 instance running a server that publishs the test coverage of a Go application.
