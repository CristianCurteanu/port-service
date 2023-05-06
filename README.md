# PortDomainService

A REST API service, that stores ports data, from a file upload. Also, it allows to choose whether to store all object in memory, or in a MongoDB database.

### Project setup

##### 1. Using docker

In order to set the project up, it is enough to run:

```sh
make run-http
```

It will automatically build the docker images, and set up mongo along with rest api server.

By default, it will store the data in memory. In order to store data in database, make sure to add `mongo-db-name` and `mongo-db-uri` arguments to the entrypoint, in docker-compose file:

```yaml
    entrypoint:
      ...
      - "--mongo-db-uri=mongodb://root:S3cret@mongodb:27017/?maxPoolSize=20&w=majority"
      - "--mongo-db-name=port"
      ...
```

and run the make command again:

```sh
make run-http
```

##### 2. Using manual build

There is also possibility to build it manually, by using:

```sh
make build
```

this will create a binary file `bin/server`; it will be used to run the application, using the command:

```sh
make run
```

Additionaly, for running `make run` command, there could be added following environment variables:
- `PORT` - in order to set a port different than default one, ie. `8080`
- `MONGO_DB_URI` - The MongoDB URI, in order to define the MongoDB connection string
- `MONGO_DB_NAME` - The MongoDB database name, where the data will be stored

Please note, that the MongoDB storage is enabled only if both `URI` and `DB_NAME` are set.
