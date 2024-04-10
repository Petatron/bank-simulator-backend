-- This statement creates a new table named "sessions".
CREATE TABLE "sessions" (
    -- "id" is a UUID that serves as the primary key for the "sessions" table.
                            "id" uuid PRIMARY KEY,

    -- "username" is a string that represents the username of the user. It cannot be null.
                            "username" varchar NOT NULL,

    -- "refresh_token" is a string that represents the refresh token of the user. It cannot be null.
                            "refresh_token" varchar NOT NULL,

    -- "user_agent" is a string that represents the user agent of the user's device. It cannot be null.
                            "user_agent" varchar NOT NULL,

    -- "client_ip" is a string that represents the IP address of the user's device. It cannot be null.
                            "client_ip" varchar NOT NULL,

    -- "is_blocked" is a boolean that indicates whether the user is blocked. It cannot be null and defaults to false.
                            "is_blocked" boolean NOT NULL DEFAULT false,

    -- "expires_at" is a timestamp that represents when the session expires. It cannot be null.
                            "expires_at" timestamptz NOT NULL,

    -- "created_at" is a timestamp that represents when the session was created. It cannot be null and defaults to the current time.
                            "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- This statement adds a foreign key constraint to the "sessions" table.
-- The "username" column in the "sessions" table references the "username" column in the "users" table.
ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
