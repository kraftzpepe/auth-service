#### Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- Automatically generate a UUID
    username VARCHAR(255) UNIQUE NOT NULL,                -- Username of the user
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

#### Password reset tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

#### Run server
go run server/main.go

## CLI

### Signup
go run cmd/main.go signup --username user1 --email user2@email.com --password "Password1@"

### Login
go run cmd/main.go login --email user4@email.com --password "Password1@"

### Query User
go run cmd/main.go query-user --email user1@email.com
