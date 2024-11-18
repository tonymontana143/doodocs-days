### Doodocs Backend Challenge


## 1. Route: Get Archive Information

### Description
This route receives an archive file and returns detailed information about the structure of the archive, such as the file names, total size, and individual file information inside the archive.

### Request
- **Method**: `POST`
- **URL**: `/api/archive/information`
- **Content-Type**: `multipart/form-data`
- **Parameters**: The archive file must be uploaded as part of the form-data under the key `"file"`.
- **File Validation**: The server will validate that the uploaded file is a valid archive (e.g., `.zip`, `.tar`, `.tar.gz`).

### Request Example

```http
POST /api/archive/information HTTP/1.1
Content-Type: multipart/form-data; boundary=----
Content-Disposition: form-data; name="file"; filename="archive.zip"
Content-Type: application/zip
{Binary data of ZIP file}
```

### Response
- **Status Code**: `200 OK`
- **Content-Type**: `application/json`
- **Body**: A JSON object containing information about the archive and its contents.

#### Response Example

```json
{
    "filename": "test1.zip",
    "archive_size": 40194,
    "total_size": 42217,
    "total_files": 83,
    "files": [
        {
            "file_path": "a-library-for-others/",
            "size": 0,
            "mime_type": "application/octet-stream"
        },
        {
            "file_path": "a-library-for-others/example.csv",
            "size": 10,
            "mime_type": "text/csv; charset=utf-8"
        },
        {
            "file_path": "a-library-for-others/exampleWithNewLine.csv",
            "size": 15,
            "mime_type": "text/csv; charset=utf-8"
        },
        {
            "file_path": "a-library-for-others/lineWithQuotes.csv",
            "size": 43,
            "mime_type": "text/csv; charset=utf-8"
        },
        {
            "file_path": "a-library-for-others/.idea/",
            "size": 0,
            "mime_type": "application/octet-stream"
        },
        {
            "file_path": "a-library-for-others/.idea/.gitignore",
            "size": 125,
            "mime_type": "application/octet-stream"
        },
    ]
    }
```
---

## 2. Route: Create Archive

### Description
This route receives a list of valid files (based on MIME type) and creates a `.zip` archive containing those files. The archive is then returned as the response.

### Request
- **Method**: `POST`
- **URL**: `/api/archive/files`
- **Content-Type**: `multipart/form-data`
- **Parameters**: Files must be uploaded as part of the form-data under the key `"files[]"`. Only the following MIME types are allowed for the files:
  - `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (Word Document)
  - `application/xml` (XML File)
  - `image/jpeg` (JPEG Image)
  - `image/png` (PNG Image)

### Request Example

```http
POST /api/archive/files HTTP/1.1
Content-Disposition: form-data; name="files[]"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document
{Binary data of document.docx}
Content-Disposition: form-data; name="files[]"; filename="avatar.png"
Content-Type: image/png
{Binary data of avatar.png}
```


#### Response Example

```http
HTTP/1.1 200 OK
Content-Type: application/zip
{Binary data of the generated zip file}
```

## 3. Route: Send File to Multiple Emails

### Description
This route receives a file and a list of email addresses, and sends the file to all the email addresses provided. The supported MIME types for files are:
- `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (Word Document)
- `application/pdf` (PDF File)

### Request
- **Method**: `POST`
- **URL**: `/api/mail/file`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `"file"`: The file to be sent, with the supported MIME types.
  - `"emails"`: A comma-separated list of email addresses.

### Request Example

```http
POST /api/mail/file HTTP/1.1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document
{Binary data of document.docx}
Content-Disposition: form-data; name="emails"
kenes_2005@mail.ru, durov@gmail.com
```

## Usage
- **build**:
    ```
    make build
    ```
- **run**:
    ```
    make run
    ```    
- **all**:
    ```
    make all
    ```    