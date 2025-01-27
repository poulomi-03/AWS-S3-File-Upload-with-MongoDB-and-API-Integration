# AWS S3 File Upload with MongoDB and API Integration
## Poulomi Bhattacharya

This project demonstrates how to upload static files (e.g., images) to AWS S3, store the file URLs in a MongoDB collection, and display the data on the frontend using API endpoints. It includes a backend (GoLang with Gin framework) and a React-based frontend.

---

## Project Features
1. **File Upload**: Users can upload static files (e.g., images) through a form.
2. **AWS S3 Integration**: Uploaded files are stored in an AWS S3 bucket with a public URL.
3. **MongoDB Integration**: The file metadata (e.g., name, email, file URL) is stored in a MongoDB collection.
4. **API Endpoints**: RESTful APIs to handle file upload and fetch data.
5. **Frontend Display**: A React.js application displays uploaded data.

---

## Prerequisites
1. **Backend**:
   - Go installed on your machine.
   - MongoDB database (local or hosted).
   - AWS S3 bucket and credentials.
2. **Frontend**:
   - Node.js and npm installed.

---

## Environment Variables
Create a `.env` file in the root directory with the following variables:

AWS_ACCESS_KEY=your_aws_access_key
AWS_SECRET_KEY=your_aws_secret_key
AWS_REGION=your_aws_region
AWS_BUCKET=your_s3_bucket_name
MONGODB_CONN_URI=your_mongodb_connection_uri
MONGODB_DB_NAME=your_database_name
COLLECTION_NAME=your_collection_name

---

## Project Setup
### Backend Setup
1. **Clone the Repository**:
   git clone https://github.com/your-repo.git
   cd backend

2. **Install Dependencies**:
   go mod tidy

3. **Run the Backend**:
   go run main.go

---

### Frontend Setup
1. **Navigate to Frontend Directory**:
   cd frontend

2. **Install Dependencies**:
   npm install

3. **Run the Frontend**:
   npm start

---

## API Endpoints
1. **POST /upload**:
   - **Description**: Uploads a file to AWS S3 and stores metadata in MongoDB.
   - **Request Body**:
     - `name` (string): User's name.
     - `email` (string): User's email.
     - `file` (binary): File to upload.
   - **Response**:
     - `message`: Upload success or failure.

2. **GET /files**:
   - **Description**: Fetches all uploaded files and metadata.
   - **Response**:
     - `data`: List of file metadata.

---

## File Structure
- **Backend**:
  - `main.go`: Main entry point.
  - `routes/`: Contains route definitions.
  - `services/`: Handles AWS S3 and MongoDB logic.
- **Frontend**:
  - `src/App.js`: Main React component.
  - `src/components/`: Contains UI components.

---

## Technologies Used
1. **Backend**:
   - GoLang
   - Gin framework
   - AWS SDK
   - MongoDB driver for Go
2. **Frontend**:
   - React.js
   - Axios for API calls
3. **Database**:
   - MongoDB (NoSQL)
4. **Storage**:
   - AWS S3 (Simple Storage Service)

---

## Deployment
1. **Backend**:
   - Use Docker to containerize the app.
   - Deploy on AWS EC2 or any cloud VM.
2. **Frontend**:
   - Deploy using Vercel, Netlify, or AWS Amplify.
3. **Database**:
   - Use a cloud MongoDB provider like MongoDB Atlas.
4. **Environment Variables**:
   - Secure sensitive data using `.env` files or cloud secrets managers.

---

## Notes
- Ensure proper CORS handling between backend and frontend.
- Validate file types and sizes before uploading.
- Secure your AWS credentials using IAM roles and least privilege principles.
