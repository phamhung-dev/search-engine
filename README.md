![alt](https://go.dev/images/go-logo-blue.svg)

---
# SEARCH ENGINE FOR FILES ON A COMPUTER

![alt](https://img.shields.io/badge/go%20version-%3D%201.21.3-brightgreen) ![alt](https://img.shields.io/badge/platform-linux%20%7C%20windows-lightgrey) ![alt](https://img.shields.io/badge/build-passing-brightgreen) ![alt](https://img.shields.io/badge/coverage-100%25-brightgreen) [![Facebook Badge](https://img.shields.io/badge/facebook-%40phamhung-blue)](https://www.facebook.com/kunn.ngoc.5/) [![LinkedIn Badge](https://img.shields.io/badge/linkedin-%40phamhung2503-blue)](https://www.linkedin.com/in/phamhung2503/) [![Github Bagde](https://img.shields.io/badge/github-%40phamhung250301-blue)](https://github.com/pham-hung-25-03-01)

(In the spirit of the Everything application - allowing quick search of everything on a computer), we have two services: FileIndex and FileSearch, each with specific tasks:

### File Index

1. Enumerate all files on the computer, gather the following information and store it in a database (relational or non-relational): full path of the file, file name, file extension, file size, creation date, modification date, last access date, and file attributes (Read Only, Hidden). If the file is a Word, Excel, or PowerPoint file, also retrieve the content of the file (if the file is not password-protected).
2. Continuous data collection ensures that any new user actions on the file system are consistently reflected in the established database.
3. Support for Windows and Linux operating systems.

### File Search

1. Provide a Rest API for searching files based on criteria such as file name (fuzzy search), file extension (exact search), file size (range search), creation date (range search), modification date (range search), last access date (range search), and file content (fuzzy search). Support combining search conditions with AND, OR logic. The result should be a list of file paths that meet the search criteria.
2. The API should implement Basic Auth or JWT for authentication.
3. Ensure fast search speed, accurate results, and support for searching in Vietnamese.

### Deployment

1. Package the two services, FileIndex and FileSearch, and deploy them using Docker.
2. The program should not crash.
3. This deployment aims to create a robust system for indexing and searching files on a computer, ensuring data consistency, and providing a secure and efficient search API. Using Docker for deployment helps in achieving consistent environments and simplifies the deployment process.

---
# TECHNICAL STACK
Some of the technologies used in this project are:

- Clean Architecture
- Golang
- WinAPI
- LinuxAPI
- BeeGo
- Gin Gonic
- GORM
- PostgreSQL
- JWT
- Casbin
- Redis
- MinIO
- Zookeeper
- Kafka
- ElasticSearch
- Docker
- Docker Compose
- Kubernetes

(<b>Note</b>: The project is still in development, so some technologies may not be used. In the future, I will update the project to use all the above technologies and more.)