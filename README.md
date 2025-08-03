# Resume Parser & Job Matching API

A comprehensive backend service for parsing resumes, managing job postings, and providing AI-powered job recommendations.

## Features

### 🔍 Resume Parsing
- Upload and parse resume files
- Extract structured data from resumes
- Profile image upload support

### 👤 User Management
- User registration and authentication
- Login/logout functionality
- Password reset capabilities
- User profile management

### 💼 Job Management
- Create and manage job postings
- Retrieve jobs by various criteria
- Recruiter-specific job listings

### 🤖 AI-Powered Recommendations
- AI job suggestions for users
- Job matching based on resume data
- Personalized recommendations

### 💳 Payment Integration
- Payment verification system
- User credit management

## API Endpoints

### Authentication & User Management
```
POST   /signup              - User registration
POST   /login               - User login
POST   /reset-password      - Reset user password
GET    /user/{id}           - Get user by ID
```

### Resume Processing
```
POST   /upload              - Upload and parse resume
POST   /upload/profile-image - Upload profile image
```

### Job Management
```
POST   /jobs                - Create new job posting
GET    /jobs                - List all jobs
GET    /jobs/id/{jobID}     - Get specific job by ID
GET    /jobs/recruiter/{recruiterID} - Get jobs by recruiter
```

### AI Recommendations
```
GET    /jobs/{jobID}/suggestions     - Get AI suggestions for a job
GET    /users/{userID}/suggestions   - Get personalized job suggestions for user
```

### Credits & Payments
```
GET    /credit/{userID}     - Get user credit information
POST   /api/verify-payment  - Verify payment transactions
```

## Project Structure

```
resume-parser/
├── cmd/server/          # Application entry point
├── config/             # Configuration files
├── internal/           # Internal application code
│   ├── controller/     # HTTP request handlers
│   ├── handler/        # Specific route handlers
│   ├── middleware/     # HTTP middleware (CORS, auth, etc.)
│   └── service/        # Business logic layer
├── storage/            # File storage utilities
├── .env               # Environment variables
├── .gitignore         # Git ignore rules
├── README.md          # Project documentation
├── go.mod             # Go module definition
└── go.sum             # Go module checksums
```

## Technologies Used

- **Language**: Go
- **Database**: PostgreSQL with GORM
- **HTTP Router**: Native Go HTTP multiplexer
- **Authentication**: JWT tokens
- **File Upload**: Multipart form handling
- **AI Integration**: Custom AI service for job matching

## Getting Started

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- Environment variables configured

### Installation

1. Clone the repository:
```bash
git clone https://github.com/satyam-svg/resume-parser.git
cd resume-parser
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations:
```bash
go run cmd/server/main.go migrate
```

5. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on the configured port (default: 8080).

## Environment Variables

Create a `.env` file with the following variables:

```env
DATABASE_URL=postgresql://username:password@localhost:5432/resume_parser
JWT_SECRET=your_jwt_secret_key
PORT=8080
UPLOAD_PATH=./uploads
AI_SERVICE_URL=your_ai_service_endpoint
PAYMENT_GATEWAY_KEY=your_payment_key
```

## API Usage Examples

### Upload Resume
```bash
curl -X POST http://localhost:8080/upload \
  -F "resume=@path/to/resume.pdf" \
  -F "user_id=123"
```

### Create Job Posting
```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Software Engineer",
    "company": "Tech Corp",
    "description": "We are looking for...",
    "requirements": ["Go", "PostgreSQL", "REST APIs"],
    "recruiter_id": "456"
  }'
```

### Get AI Suggestions
```bash
curl -X GET http://localhost:8080/users/123/suggestions \
  -H "Authorization: Bearer your_jwt_token"
```

## Response Format

All API responses follow this structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data here
  },
  "error": null
}
```

Error responses:
```json
{
  "success": false,
  "message": "Error description",
  "data": null,
  "error": "Detailed error information"
}
```

## Features in Detail

### Resume Parsing
- Supports PDF, DOC, and DOCX formats
- Extracts: personal info, skills, experience, education
- Stores parsed data in structured format

### Job Matching Algorithm
- Analyzes resume skills vs job requirements
- Considers experience level and location preferences
- Provides match percentage and recommendations

### Credit System
- Users earn/spend credits for premium features
- Payment integration for credit purchases
- Usage tracking and limits

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

- JWT-based authentication
- CORS protection enabled
- Input validation and sanitization
- Secure file upload handling

## Performance

- Database connection pooling
- Efficient query optimization with GORM
- Caching for frequently accessed data
- Asynchronous processing for heavy operations

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please contact:
- Email: support@resumeparser.com
- GitHub Issues: Create an issue in this repository

## Changelog

### v1.0.0 (Current)
- Initial release
- Resume parsing functionality
- Job management system
- AI-powered recommendations
- Payment integration
- User authentication system
