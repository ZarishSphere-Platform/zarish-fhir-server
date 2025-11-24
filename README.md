# Zarish FHIR Server

**Part of the ZarishSphere Platform - A No-Code FHIR Healthcare Data Management System**

The Zarish FHIR Server is a Go-based implementation of the HL7 FHIR (Fast Healthcare Interoperability Resources) standard. It provides RESTful APIs for managing healthcare data including patients, observations, medications, and other FHIR resources.

## 🚀 Technology Stack

- **Language**: Go 1.23+
- **Web Framework**: Gin
- **Database ORM**: GORM
- **Database**: PostgreSQL
- **Search Engine**: Elasticsearch (optional)
- **Authentication**: Keycloak integration
- **Containerization**: Docker

## 📋 Prerequisites for Local Development

Before you begin, ensure you have the following installed:

- **Go**: Version 1.23 or higher ([Download Go](https://go.dev/dl/))
- **PostgreSQL**: Version 14 or higher ([Download PostgreSQL](https://www.postgresql.org/download/))
- **Docker** (optional): For containerized deployment ([Download Docker](https://www.docker.com/))
- **Git**: For version control ([Download Git](https://git-scm.com/))
- **Code Editor**: VS Code with Go extension recommended

### Checking Your Installation

```bash
go version      # Should show go1.23 or higher
psql --version  # Should show PostgreSQL 14.x or higher
docker --version # Should show Docker version 20.x or higher
git --version   # Should show git version 2.x.x
```

## 🛠️ Step-by-Step Development Setup

### Step 1: Clone the Repository

```bash
# Navigate to your desired directory
cd ~/Desktop

# Clone the repository
git clone https://github.com/ZarishSphere-Platform/zarish-fhir-server.git

# Navigate into the project
cd zarish-fhir-server
```

### Step 2: Install Go Dependencies

```bash
# Download and install all Go modules
go mod download

# Verify dependencies
go mod tidy
```

### Step 3: Set Up PostgreSQL Database

#### Option A: Using Local PostgreSQL

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE zarish_fhir;

# Create user (optional)
CREATE USER zarish WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE zarish_fhir TO zarish;

# Exit psql
\q
```

#### Option B: Using Docker

```bash
# Run PostgreSQL in Docker
docker run --name zarish-postgres \
  -e POSTGRES_DB=zarish_fhir \
  -e POSTGRES_USER=zarish \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 \
  -d postgres:14
```

### Step 4: Configure Environment Variables

Create a `.env` file in the project root:

```bash
# Create the environment file
touch .env
```

Add the following configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=zarish
DB_PASSWORD=your_password
DB_NAME=zarish_fhir
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8081
SERVER_HOST=0.0.0.0

# Keycloak Configuration
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=zarish
KEYCLOAK_CLIENT_ID=zarish-fhir-server
KEYCLOAK_CLIENT_SECRET=your_client_secret

# Elasticsearch (optional)
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_ENABLED=false
```

### Step 5: Run Database Migrations

```bash
# The application will automatically create tables on first run
# Or you can run migrations manually if implemented
go run cmd/server/main.go migrate
```

### Step 6: Start the Development Server

```bash
# Run the server
go run cmd/server/main.go

# Or build and run
go build -o zarish-fhir-server cmd/server/main.go
./zarish-fhir-server
```

The server will start at `http://localhost:8081`

### Step 7: Test the API

```bash
# Check server health
curl http://localhost:8081/health

# Get FHIR metadata
curl http://localhost:8081/fhir/metadata

# Create a patient (example)
curl -X POST http://localhost:8081/fhir/Patient \
  -H "Content-Type: application/fhir+json" \
  -d '{
    "resourceType": "Patient",
    "name": [{"family": "Doe", "given": ["John"]}],
    "gender": "male"
  }'
```

## 🔧 Available Commands

| Command | Description |
|---------|-------------|
| `go run cmd/server/main.go` | Start development server |
| `go build -o zarish-fhir-server cmd/server/main.go` | Build production binary |
| `go test ./...` | Run all tests |
| `go mod tidy` | Clean up dependencies |
| `go fmt ./...` | Format code |
| `go vet ./...` | Run Go vet for code analysis |

## 📁 Project Structure

```
zarish-fhir-server/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routes
│   ├── auth/              # Authentication middleware
│   ├── database/          # Database connection and migrations
│   ├── models/            # FHIR resource models
│   └── search/            # Search functionality
├── Dockerfile             # Docker configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
└── README.md              # This file
```

## 🧪 Testing

### Run All Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Manual API Testing

Use tools like:
- **curl**: Command-line HTTP client
- **Postman**: GUI for API testing
- **Insomnia**: Alternative to Postman

## 🐳 Docker Deployment

### Build Docker Image

```bash
# Build the image
docker build -t zarish-fhir-server .

# Run the container
docker run -p 8081:8081 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=zarish \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=zarish_fhir \
  zarish-fhir-server
```

### Using Docker Compose

See the [zarish-fhir-infra](https://github.com/ZarishSphere-Platform/zarish-fhir-infra) repository for complete Docker Compose setup.

## 📚 FHIR Resources Supported

The server currently supports the following FHIR resources:

- **Patient**: Patient demographics and information
- **Observation**: Clinical observations and measurements
- **Medication**: Medication information
- **MedicationRequest**: Medication prescriptions
- **Condition**: Patient conditions and diagnoses
- **Encounter**: Patient encounters and visits
- **Practitioner**: Healthcare practitioners
- **Organization**: Healthcare organizations

## 🔍 Search Capabilities

The server supports FHIR search parameters:

```bash
# Search patients by name
GET /fhir/Patient?name=John

# Search observations by patient
GET /fhir/Observation?patient=Patient/123

# Search with multiple parameters
GET /fhir/Patient?name=John&gender=male
```

## 🐛 Troubleshooting

### Database Connection Fails

```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Check database exists
psql -U postgres -l | grep zarish_fhir

# Test connection
psql -U zarish -d zarish_fhir -h localhost
```

### Port Already in Use

```bash
# Find process using port 8081
lsof -i :8081

# Kill the process
kill -9 <PID>

# Or change the port in .env
SERVER_PORT=8082
```

### Module Download Fails

```bash
# Clear Go module cache
go clean -modcache

# Re-download modules
go mod download
```

## 📚 Learning Resources

- [FHIR Specification](https://www.hl7.org/fhir/)
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Write/update tests
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

This project is part of the ZarishSphere Platform.

## 🆘 Getting Help

- Check [Issues](https://github.com/ZarishSphere-Platform/zarish-fhir-server/issues)
- Create a new issue with details
- Include error logs and steps to reproduce

## 🔗 Related Repositories

- [zarish-frontend-shell](https://github.com/ZarishSphere-Platform/zarish-frontend-shell) - Frontend Application
- [zarish-terminology-server](https://github.com/ZarishSphere-Platform/zarish-terminology-server) - Terminology Server
- [zarish-fhir-infra](https://github.com/ZarishSphere-Platform/zarish-fhir-infra) - Infrastructure & Deployment
- [zarish-fhir-data](https://github.com/ZarishSphere-Platform/zarish-fhir-data) - FHIR Data Resources
