BEGIN;

CREATE TABLE users(
		user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(64) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
);

CREATE TABLE rooms(
        room_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        room_name VARCHAR(64) UNIQUE NOT NULL
);

CREATE TABLE messages(
        message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        room_id UUID NOT NULL REFERENCES rooms(room_id),
        content VARCHAR(64) NOT NULL
);

COMMIT;