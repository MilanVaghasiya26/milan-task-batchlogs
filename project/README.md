### API's

#### 1: Platform Routes

###### Platform Public Routes

- Logs Module

1. POST "Create Batch Logs" => /platform/api/v1/ingest

   Request Body :-

```
{
    "logs": [
            "2023-10-11T10:31:00Z INFO [apache] Received GET request from 192.168.0.1 for /index.html",
            "2023-10-11T10:32:15Z INFO [apache] Request from 10.0.0.2 failed with status code 404 for /page-not-found.html",
            "2023-10-11T11:33:30Z WARN [nginx] Received POST request from 192.168.0.3 for /submit-form"
        ]
}
```

2. GET "List Batch Logs" => /platform/api/v1/query

   Response Body :-

```
[
    {
        "timestamp": "2023-10-11T10:31:00Z",
        "body": "Received GET request from 192.168.0.1 for /index.html",
        "service": "apache",
        "severity": "INFO"
    },
    {
        "timestamp": "2023-10-11T10:32:15Z",
        "body": "Request from 10.0.0.2 failed with status code 404 for /page-not-found.html",
        "service": "apache",
        "severity": "INFO"
    },
    {
        "timestamp": "2023-10-11T11:33:30Z",
        "body": "Received POST request from 192.168.0.3 for /submit-form",
        "service": "nginx",
        "severity": "WARN"
    }
]{
    "message": "Batch logs fetched successfully."
}
```

## Generate Swagger Docs

```
cd project
swag init --pd
```

- To Run Swagger Doc

  - Run the project first and then hit

```
http://localhost:3001/project/docs
```
