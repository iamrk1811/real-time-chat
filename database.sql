    -- Create Users table
    CREATE TABLE Users (
        user_id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Create Groups table
    CREATE TABLE Groups (
        group_id SERIAL PRIMARY KEY,
        group_name VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Create UserGroups table to manage many-to-many relationships between users and groups
    CREATE TABLE UserGroups (
        user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
        group_id INT REFERENCES Groups(group_id) ON DELETE CASCADE,
        joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (user_id, group_id)
    );

    -- Create Messages table
    CREATE TABLE Messages (
        message_id SERIAL PRIMARY KEY,
        sender_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
        group_id INT REFERENCES Groups(group_id) ON DELETE CASCADE, -- Can be NULL for user-to-user messages
        receiver_id INT REFERENCES Users(user_id) ON DELETE CASCADE, -- Can be NULL for group messages
        content TEXT NOT NULL,
        sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Create indexes for commonly queried fields
    CREATE INDEX idx_messages_sender_id ON Messages(sender_id);
    CREATE INDEX idx_messages_receiver_id ON Messages(receiver_id);
    CREATE INDEX idx_messages_group_id ON Messages(group_id);
    CREATE INDEX idx_messages_sent_at ON Messages(sent_at);


CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY, -- Unique session ID
    user_id INT NOT NULL REFERENCES Users(user_id), -- User ID associated with the session
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the session was created
    expires_at TIMESTAMPTZ NOT NULL -- Timestamp when the session expires
);