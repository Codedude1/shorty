
# Shorty - cut it short
### Project Description
A simple URL Shortener Service that allows users to convert long URLs into shortened, unique URLs, similar to services like Bitly. The service ensures that each unique long URL generates a unique short URL, reusing the same short URL if the same long URL is submitted again. It also provides redirection to the original URL when the shortened URL is accessed. Other features include access statistics and URL expiration (TTL).

### Features
* Shorten URLs: Submit a long URL and receive a shortened version.

* Redirection: Accessing the shortened URL redirects to the original long URL.

* Unique URLs: Each unique long URL generates a unique short URL. Duplicate submissions reuse the same short URL.

* Validation: Validates input to ensure the URL is valid.

* Access Statistics: Tracks and displays the number of times a shortened URL has been accessed.

* Time-to-Live (TTL): Allows URLs to expire after a specified duration, with appropriate cleanup.
### Architecture and Design Decisions
1. Overall Architecture
   
   The application follows a modular design with separation of concerns, enhancing maintainability and scalability. The key components are:
   
      * API Layer: Handles HTTP requests and responses.
      * Services Layer: Contains business logic for URL shortening, validation, redirection, and optional features.
      * Data Storage Layer: In-memory storage system to store mappings between long URLs and short URLs.
      * Concurrency Handling: Utilizes thread-safe data structures and synchronization mechanisms to handle multiple requests simultaneously.
3. Design Decisions
* Programming Language and Framework: Implemented using Golang with the Gin web framework, chosen for its performance, simplicity, and built-in concurrency support.

* URL Shortening Algorithm: 
    * Unique ID Generation: Uses an auto-incrementing integer as a unique identifier for each URL.
    * Encoding: Applies Base62 encoding to the unique ID to generate a short string consisting of alphanumeric characters.
* In-Memory Data Storage:
    * Rationale: Satisfies the assignment's constraints and provides quick read/write operations.
    * Scalability Consideration: Abstracted the data storage layer to facilitate future migration to a persistent database.
* Validation: 
    * URL Format Validation: Utilizes Go's net/url package to ensure the URL is properly formatted.
    * Reachability Check: Not implemented by default to maintain performance, but can be added if needed.
* Concurrency Handling:
    * Thread Safety: Implements sync.RWMutex to synchronize access to shared data structures.
    * Concurrent Requests: Designed to handle multiple simultaneous requests efficiently.
* Additional Features: 
    * Access Statistics: Stores an access count for each URL and increments it atomically upon each redirection.
    * Time-to-Live (TTL): Stores an expiration timestamp for each URL and checks it upon access.

### Setup and Installation

#### Prerequisites
    Golang version 1.16 or higher.
#### Installation Steps
* Clone the repository

        git clone https://github.com/Codedude1/shorty.git
        cd shorty
* Initialize go modules
    
        go mod tidy

* Run the application
    
        go run main.go

The server will start on http://localhost:8081.

### Configuration Options
Shorty can be customized using environment variables. Below are the configurable options:

* Server Port:

    * Environment Variable: PORT
    * Default: 8081
    * Description: Specifies the port on which the server listens.


* Default TTL:

    * Environment Variable: DEFAULT_TTL
    * Default: 24h
    * Description: Sets the default time-to-live for shortened URLs.

* Storage Type:

    * Environment Variable: STORAGE_TYPE
    * Options: memory, redis, postgresql
    * Default: memory
    * Description: Determines the type of storage backend used.

To set environment variables, you can create a .env file in the project root:
```env   
   PORT=8081
   DEFAULT_TTL=24h
   STORAGE_TYPE=memory
```



And load them using a package like github.com/joho/godotenv:



    import (
        "github.com/joho/godotenv"
    )

    func main() {
        err := godotenv.Load()
        if err != nil {
            log.Println("No .env file found")
        }
    // Rest of the code
    }

### Usage
* Shorten a URL

    
    Endpoint: POST /shorten
    
    Request Body:
        
        {"url": "https://www.example.com"}
    Example using cURL:

        curl -X POST -H "Content-Type: application/json" -d '{"url":"https://www.example.com"}' http://localhost:8081/shorten
    Response:

        {"short_url": "http://localhost:8081/abc123"}
* Shorten a URL with TTL


   Endpoint: POST /shorten
   
   Request Body:

      {"url": "https://www.example.com", "expiry_in_sec": 30}
   Example using cURL:

        curl -X POST -H "Content-Type: application/json" -d '{"url":"https://www.example.com", "expiry_in_sec": 30}' http://localhost:8081/shorten
   Response:

        {"short_url": "http://localhost:8081/abc123"}

*  Redirect to Original URL

    Access the shortened URL in a web browser or via an HTTP       
    GET request:

        curl -L http://localhost:8081/abc123


This will redirect you to https://www.example.com.

* Access Statistics 

    Endpoint: GET /stats/{shortURL}
    
    Example:

        curl http://localhost:8081/stats/abc123
    
    Response:

        {"long_url": "https://www.example.com", "access_count": 42}
### Error Handling & Validation
Shorty handles various error scenarios to ensure robust and reliable operation:

* Invalid URLs:

    * Scenario: Submitting a malformed or invalid URL.

    * Response: 400 Bad Request

    * Example Response:

        ```json
        {
            "error": "Invalid URL format."
        }
* URL Expired:

    * Scenario: Accessing a shortened URL that has expired.

    * Response: 410 Gone

    * Example Response:

        ```json
        {
            "error": "This short URL has expired."
        }
* Non-Existent Short URL:

    * Scenario: Accessing a short URL that does not exist.

    * Response: 404 Not Found

    * Example Response:

        ```json
        {
            "error": "Short URL not found."
        }
* Duplicate URL Submission:

    * Scenario: Submitting the same long URL multiple times.

    * Response: Returns the existing short URL without creating duplicates.

    * Example Response:

        ```json
        {
            "short_url": "http://localhost:8081/abc123"
        }
* Server Errors:

    * Scenario: Unexpected server-side issues.

    * Response: 500 Internal Server Error

    * Example Response:

    ```json
    {
        "error": "An unexpected error occurred. Please try again later."
  }

### Challenges Faced 
* Ensuring Thread Safety
    
    Problem: Concurrent access to shared data structures could lead to race conditions.

    Solution: Used sync.RWMutex to synchronize read and write operations on the in-memory storage.
* URL Validation Accuracy

    Problem: Accurately validating URLs without rejecting valid ones or accepting invalid ones.

    Solution: Utilized Go's net/url package for robust URL parsing and validation.
* Handling URL Expiration

    Problem: Efficiently managing expired URLs without degrading performance.

    Solution: Implemented TTL checks during access and a periodic cleanup routine running in a separate goroutine.
### Future Improvements
* Persistent Storage: Migrate to a persistent database like Redis or PostgreSQL for data durability across restarts.
* Custom Aliases: Allow users to specify custom aliases for their short URLs.
* User Authentication: Implement user accounts to manage personal URL mappings.
* Analytics Dashboard: Provide a web interface to view access statistics and manage URLs.
* Enhanced Validation: Add checks for malicious URLs or phishing attempts.
* Rate Limiting: Implement rate limiting to prevent abuse of the service.
### Appendix
#### Workflow Diagrams
    
* URL Shortening Workflow

   ```mermaid
      graph TD
    A[User Submits Long URL via POST /shorten] --> B[Validate URL Format]
    B -->|Valid| C[Check if URL Already Exists]
    C -->|Exists| D[Return Existing Short URL]
    C -->|Does Not Exist| E[Generate Unique ID]
    E --> F[Encode ID using Base62]
    F --> G[Store Mapping in In-Memory Storage]
    G --> H[Return Shortened URL to User]
    B -->|Invalid| I[Return 400 Bad Request]



* URL Redirection Workflow

   ```mermaid
      graph TD
    A["User Accesses Shortened URL via GET /{shortURL}"] --> B["Retrieve Original URL from Storage"]
    B --> C{Check if URL is Expired}
    C -->|Expired| D["Return 410 Gone"]
    C -->|Not Expired| E["Increment Access Count"]
    E --> F["Redirect to Original URL"]
    B -->|Not Found| G["Return 404 Not Found"]




#### Database Schema

While the service uses in-memory storage, the following database schemas outline how the data would be structured in a relational database for future scalability.

* URLMappings Table 
    
        CREATE TABLE URLMappings (id INTEGER PRIMARY KEY AUTOINCREMENT,
        long_url TEXT NOT NULL UNIQUE,
        short_code VARCHAR(10) NOT NULL UNIQUE,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
        access_count INTEGER NOT NULL DEFAULT 0,
        expires_at DATETIME
        );
    * id: Auto-incrementing unique identifier.
    * long_url: The original long URL provided by the user.
    * short_code: The encoded string used in the shortened URL.
    * created_at: Timestamp of when the mapping was created.
    * access_count: Number of times the short URL has been accessed.
    * expires_at: Optional expiration date and time for the short URL.
* AccessLogs Table 
        
        CREATE TABLE AccessLogs (id INTEGER PRIMARY KEY AUTOINCREMENT,
        short_code VARCHAR(10) NOT NULL,
        accessed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_agent TEXT,
        ip_address VARCHAR(45),
        FOREIGN KEY (short_code) REFERENCES URLMappings(short_code));



    
    * id: Auto-incrementing unique identifier.
    * short_code: The short URL code accessed.
    * accessed_at: Timestamp of the access event.
    * user_agent: The user agent string from the request header (optional).
    * ip_address: IP address of the client making the request (optional).


Contact Information
For any questions or suggestions, please contact:

Name: Yash Mishra
Email: [mishra.y19@gmail.com]
GitHub: [https://github.com/Codedude1]

