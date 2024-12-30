#### Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- Automatically generate a UUID
    username VARCHAR(255) NOT NULL,                -- Username of the user
    email VARCHAR(255) UNIQUE NOT NULL,            -- Unique email address
    password VARCHAR(255) NOT NULL,                -- Hashed password
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Automatically set creation time
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Automatically set update time
);

#### Refresh Tokens Table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),      -- Automatically generate a UUID
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Foreign key to users table
    token VARCHAR(512) NOT NULL,                        -- Refresh token
    expires_at TIMESTAMP NOT NULL,                      -- Expiration time for the token
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Automatically set creation time
);

#### Run server
go run server/main.go

## CLI

### Signup
go run cmd/main.go signup --username user1 --email user1@email.com --password "Password1@"

### Login
go run cmd/main.go login --email user1@email.com --password "Password1@"

### Query User
go run cmd/main.go query-user --email user1@email.com
