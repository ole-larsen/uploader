This is a Go service built using goswagger.io for handling file uploads. The service allows files to be uploaded through secured provider routes and provides a public route for retrieving files. The service uploads files to the uploads directory and stores file metadata in a PostgreSQL database.

```
swagger generate server -A uploader -f ./schema/swagger.yml -P models.Principal
```
# Docker Build Command
To build the Docker image for the service, run the following command:
```
sudo docker build -t olelarsen/uploader \
  --build-arg NODE_ENV=production \
  --build-arg APP_NAME=uploader \
  --build-arg PORT=1234 \
  --build-arg SESSION_SECRET=12345abscdefg \
  --build-arg X_TOKEN=abcdefg12345 \
  --build-arg USE_HASH=true \
  --build-arg DB_SQL_HOST=localhost \
  --build-arg DB_SQL_PORT=5432 \
  --build-arg DB_SQL_USERNAME=postgres \
  --build-arg DB_SQL_PASSWORD=postgres \
  --build-arg DB_SQL_DATABASE=files \
  --build-arg USE_DB=true .
```
# Provider Routes
The service provides the following secured provider routes:

GET /files: Retrieves a list of files.
POST /files: Uploads a new file.
PUT /files/{id}: Updates an existing file.

# Public Route
The service includes a public route for retrieving and resizing files. The URL format for the public route is:
```
GET /files/{size}:{hash}.{extension}
{size}: The desired dimensions for resizing the image in the format widthxheight (e.g., 400x400).
{hash}: The hashed name of the file.
{extension}: The file extension (e.g., jpg, png).
```
When a request is made to the public route with the specified size, hash, and extension, the service will return the resized image.


# The files table stores the following information for each file:
PostgreSQL Table Structure
The service uses a PostgreSQL database to store file metadata. The table structure is defined as follows:
```
id: Unique identifier of the file.
name: Name of the file (must be unique).
alt: Alternate text for the file.
caption: Caption for the file.
width: Width of the file (in pixels).
height: Height of the file (in pixels).
formats: JSON object storing additional formats of the file.
hash: Hash value of the file.
ext: Extension of the file.
mime: MIME type of the file.
size: Size of the file (in bytes).
url: URL of the file.
preview_url: URL of the file's preview.
provider: Provider of the file.
provider_metadata: JSON object storing provider-specific metadata.
created_by_id: ID of the user who created the file.
updated_by_id: ID of the user who last updated the file.
created: Timestamp indicating the creation time of the file.
updated: Timestamp indicating the last update time of the file.
deleted: Timestamp indicating the deletion time of the file (nullable).
```