CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE  IF NOT EXISTS role (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO role (name)
VALUES ('hr');

INSERT INTO role (name)
VALUES ('staff');

INSERT INTO role (name)
VALUES ('admin');


CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    role_id UUID NOT NULL REFERENCES role(id) ON DELETE CASCADE,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE, 
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
    
);

INSERT INTO users (role_id, email, password_hash)
VALUES (
    (SELECT id FROM role WHERE name = 'admin'), 
    'admin@admin.com', 
    '$2y$06$SnLcE85z.wBiHTEyehaEOu.AFoBzd3TLnoHNCqdm9kAGivuVl1YoG'
);


CREATE TABLE IF NOT EXISTS token (
    hash bytea NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) NOT NULL,
    scope TEXT NOT NULL
);

CREATE TABLE  IF NOT EXISTS  employee (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    dob DATE NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL,
    gender TEXT NOT NULL,
    job_title TEXT NOT NULL,
    department TEXT NOT NULL,
    address TEXT NOT NULL,
    is_present BOOLEAN NOT NULL DEFAULT TRUE,
    joining_date DATE NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);


INSERT INTO employee (user_id, name, dob, avatar, phone, gender, job_title, department, address, joining_date)
VALUES (
    (SELECT id FROM users WHERE email = 'admin@admin.com'), 
    'Admin Admin',
    '1990-01-01', 
    '', 
    '0654051622', 
    'Male',
    'Admin', 
    'Admin',
    'Admin Address', 
    '2024-04-29'
);




CREATE TABLE  IF NOT EXISTS  contract (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    contract_type TEXT NOT NULL,
    period INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    attachment TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);


CREATE TABLE  IF NOT EXISTS  leave (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    approved_by_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    leave_type TEXT NOT NULL DEFAULT 'annual_leave',
    description TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    leave_count SMALLINT NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    seen BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE  IF NOT EXISTS announcement (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    description TEXT NOT NULL, 
    announcement_date DATE NOT NULL,
    created_by UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

