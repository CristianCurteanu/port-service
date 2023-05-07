# PortDomainService

A REST API service, that stores ports data, from a file upload. Also, it allows to choose whether to store all object in memory, or in a MongoDB database.

## Project setup

#### 1. Using docker

In order to run the docker container, with in memory storage, use following command:

```sh
make run-container-inmem
```

It will automatically build the docker images, and launch the server on port `8080`.

In order to start application in container, along with MongoDB, use this command:

```sh
make run-container-mongo
```

#### 2. Using manual build

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

## Endpoints

### 1. Create or Update Ports

This endpoint takes an form file as input parameter, and stores all port data

**URL** : `/ports`

**Method** : `POST`

**Auth required** : NO

**Data constraints**:

The file should contain `json`, with following structure:
```json
{
    "<PORT_CODE>": {
        "name": "<Port name>",
        "city": "<City of port location>",
        "country": "<country of port location>",
        "alias": ["<string of port name alias>"],
        "regions": [],
        "coordinates": [
            52.6126027,
            24.1915137
        ],
        "province": "<province of port location>",
        "timezone": "<time zon>",
        "unlocs": [
            "<Port code>"
        ],
        "code": "<port numerical code>"
    }
}
```

**Request example**:

```sh
curl --request POST \
  --url http://localhost:8080/ports \
  --header 'Content-Type: multipart/form-data' \
  --form ports=@/absolute/path/to/file/ports.json
```

#### Success Response

**Code** : `201 Created`

**Content example**

```json
{}
```

#### Bad File content Response

**Condition** : If the structure of JSON object in the file is wrong.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
    "code": "bad_json_file",
    "message": "Please check your json file, there might be syntax issues"
}
```

#### Data store failure Response

**Condition** : If the records, that where supposed to be stored, failed to be stored.

**Code** : `500 INTERNAL SERVER ERROR`

**Content** :

```json
{
    "code": "err_data_store",
    "message": "Error while storing the data; please contact administrator to check the reason of failure"
}
```

#### File reading error Response

**Condition** : If there are issues with file reading.

**Code** : `500 INTERNAL SERVER ERROR`

**Content** :

```json
{
    "code": "internal_error",
    "message": "Please check with the administrator"
}
```

### 2. GET Ports by Port Code

Fetch a port record by specific Port Code

**URL** : `/ports/{port_code}`

**Method** : `GET`

**Auth required** : NO

**Request example**:

```sh
curl --request GET \
  --url http://localhost:8080/ports/<port-code>
```

#### Success Response

**Code** : `200 OK`

**Content example**

```json
{
	"port_code": "<port code>",
	"name": "<port name>",
	"city": "<port city>",
	"country": "<country>",
	"code": "<port numeric code>",
	"coordinates": [
		// coordinates in float
	],
	"province": "<province of port location>",
	"unlocs": [
		"<port codes/unlocs>"
	]
}
```

#### Port record not found Response

It is returned in case record is not stored

**Code** : `404 NOT FOUND`

**Content example**

```json
{
	"code": "not_found",
	"message": "No port found with the specified port code"
}
```