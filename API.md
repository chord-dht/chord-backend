# Backend API Documentation

## Basic Information

- Base Path: `/api`
- All endpoints return data in JSON format

## Endpoints

### 1. Get Node Status

- **URL**: `/api/nodestatus`
- **Method**: `GET`
- **Description**: Check if the node exists

**Response**:

- `200 OK`

```json
{
    "exists": true | false
}
```

### 2. New Node

Combine `/api/new` and `/api/initialize`.

#### 2.1 Create New Node

- **URL**: `/api/new`
- **Method**: `POST`
- **Description**: Create a new node
- **Request Body**: JSON formatted configuration data

example:

```json
{
    "IdentifierLength": 10,
    "SuccessorsLength": 4,
    "IpAddress": "127.0.0.1",
    "Port": "7000",
    "Mode": "create",
    "JoinAddress": "",
    "JoinPort": "",
    "StabilizeTime": 3000,
    "FixFingersTime": 1000,
    "CheckPredecessorTime": 3000,
    "StorageDir": "./storage",
    "BackupDir": "./backup",
    "AESBool": false,
    "AESKeyPath": "",
    "TLSBool": false,
    "CaCert": "",
    "ServerCert": "",
    "ServerKey": ""
}
```

**Response**:

- `201 Created`

```json
{
    "status": "success",
    "message": "new node succeeded"
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "NODE_EXISTS_ERROR" | "BIND_JSON_ERROR" | "PARSE_JSON_ERROR" | "VALIDATE_CONFIG_ERROR",
    "error_message": "..."
}
```

- `500 Internal Server Error`

```json
{
    "status": "error",
    "error_code": "CREATE_NODE_ERROR",
    "error_message": "..."
}
```

#### 2.2 Initialize Node

- **URL**: `/api/initialize`
- **Method**: `GET`
- **Description**: Initialize the node

**Response**:

- `200 OK`

```json
{
    "status": "success",
    "message": "initialize node succeeded"
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "NODE_NOT_EXISTS_ERROR",
    "error_message": "node not created: Please create a node first"
}
```

- `500 Internal Server Error`

```json
{
    "status": "error",
    "error_code": "INITIALIZE_ERROR",
    "error_message": "failed to initialize node: ..."
}
```

### 3. Quit Node

- **URL**: `/api/quit`
- **Method**: `GET`
- **Description**: Quit the current node

**Response**:

- `200 OK`

```json
{
    "status": "success",
    "message": "quit node succeeded"
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "NODE_NOT_EXISTS_ERROR",
    "error_message": "node not created: Please create a node first"
}
```

### 4. Get Node State

- **URL**: `/api/printstate`
- **Method**: `GET`
- **Description**: Get the current node's state information

**Response**:

- `200 OK`

```json
{
    "status": "success",
    "data": {
    "node_state": "..."
    }
}
```

node state example:

```json
{
    "data": {
        "node_state": {
            "info": {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
            "predecessor": {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
            "successors": [
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"}
            ],
            "fingerTable": [
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"},
                {"Identifier": 308, "IpAddress": "127.0.0.1", "Port": "7000"}
            ],
            "fingerIndex": [309, 310, 312, 316, 324, 340, 372, 436, 564, 820],
            "localStorageName": [],
            "backupStoragesName": [
                [], [], [], []
            ]
        }
    },
    "status": "success"
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "NODE_NOT_EXISTS_ERROR",
    "error_message": "node not created: Please create a node first"
}
```

### 5. Store File

- **URL**: `/api/storefile`
- **Method**: `POST`
- **Description**: Store a file to the node
- **Request Body**: Form data including the file

**Response**:

- `200 OK`

```json
{
    "status": "success",
    "data": {
        "file_identifier": "...",
        "target_node": "..."
    }
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "FORM_FILE_ERROR" | "NODE_NOT_EXISTS_ERROR",
    "error_message": "..."
}
```

- `500 Internal Server Error`

```json
{
    "status": "error",
    "error_code": "OPEN_FILE_ERROR" | "READ_FILE_ERROR" | "FIND_ERROR" | "ENCRYPT_ERROR" | "ACCESS_ERROR" | "STORE_DENIED_ERROR",
    "error_message": "...",
    "data": {
        "file_identifier": "...", // if has
        "target_node": "..." // if has
    }
}
```

### 6. Get File

Combine `/api/getfile` and `/api/downloadfile`.

#### 6.1 Get File

- **URL**: `/api/getfile`
- **Method**: `POST`
- **Description**: Retrieve a file from the node
- **Request Body**: JSON formatted data including the file name

**Response**:

- `200 OK`

```json
{
    "status": "success",
    "data": {
        "file_identifier": "...",
        "target_node": "..."
    }
}
```

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "BIND_JSON_ERROR" | "PARSE_JSON_ERROR" | "NODE_NOT_EXISTS_ERROR",
    "error_message": "node not created: Please create a node first"
}
```

- `500 Internal Server Error`

```json
{
    "status": "error",
    "error_code": "FIND_ERROR" | "ACCESS_ERROR" | "NON_FILE_ERROR" | "DECRYPT_ERROR" | "TEMP_ERROR",
    "error_message": "...",
    "data": {
        "file_identifier": "...", // if has
        "target_node": "..." // if has
    }
}
```

#### 6.2 Download File

- **URL**: `/api/downloadfile`
- **Method**: `POST`
- **Description**: Download a file
- **Request Body**: JSON formatted data including the file name

- **Response**:

- `200 OK`: File download

- `400 Bad Request`

```json
{
    "status": "error",
    "error_code": "BIND_JSON_ERROR" | "PARSE_JSON_ERROR",
    "error_message": "node not created: Please create a node first"
}
```

- `404 Not Found`

```json
{
    "status": "error",
    "error_code": "FILE_NOT_FIND_ERROR",
    "error_message": "file not found"
}
```
