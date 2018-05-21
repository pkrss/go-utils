
-- pkrss_user

CREATE TABLE pkrss_user (
    id bigint primary key,
    user_name text UNIQUE,
    password text NOT NULL,
    denied boolean NOT NULL DEFAULT false,
    role int NOT NULL DEFAULT 0,
    email text,
    last_error_login_time timestamp with time zone,
    last_error_login_timer int default 0,
    create_time timestamp with time zone DEFAULT now(),
    update_time timestamp with time zone DEFAULT now()
);

CREATE TABLE pkrss_user_context (
    id bigint primary key,
    user_id bigint NOT NULL,
    user_name text,
    ip text,
    role int NOT NULL DEFAULT 0,
    create_time timestamp with time zone DEFAULT now(),
    update_time timestamp with time zone DEFAULT now()
);