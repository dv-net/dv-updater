# Quick Start

## Service Launch

To start the binary file, use the command:

```sh
make run start
```

The service will start and listen on port 8080 for incoming HTTP requests.

---

## API Endpoints

### 1. Service Update

**Method:** `POST`

**URL:** `/api/v1/update`

**Request Body:**
```json
{
    "name": "dv-merchant"
}
```

**Description:** This endpoint is used to update the service with the specified name.

---

### 2. Get Service Version

**Method:** `GET`

**URL:** `/api/v1/version/{packege_name}`

**Example Request:**
```
GET /api/v1/version/dv-merchant
```


**Description:** Returns the current version of the service by its name.

