# PortDomainService

A REST API service, that stores ports data, from a file upload. Also, it allows to choose whether to store all object in memory, or in a MongoDB database.

### Project setup

##### 1. Using docker

In order to run the docker container, with in memory storage, use following command:

```sh
make run-container-inmem
```

It will automatically build the docker images, and launch the server on port `8080`.

In order to start application in container, along with MongoDB, use this command:

```sh
make run-container-mongo
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
