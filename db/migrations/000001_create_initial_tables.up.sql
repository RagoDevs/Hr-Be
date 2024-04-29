CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE  IF NOT EXISTS role (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    role_id UUID NOT NULL REFERENCES role(id) ON DELETE CASCADE,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
    
);


CREATE TABLE  IF NOT EXISTS  employee (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL,
    avatar TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL,
    gender TEXT NOT NULL,
    job_title TEXT NOT NULL,
    department TEXT NOT NULL,
    address TEXT NOT NULL,
    joining_date DATE NOT NULL
);



CREATE TABLE  IF NOT EXISTS  contract (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    contract_type TEXT NOT NULL,
    period INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    attachment TEXT NOT NULL DEFAULT ''
);


CREATE TABLE  IF NOT EXISTS  leave (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    employee_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    approved_by UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    approved BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    leave_count SMALLINT NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    seen BOOLEAN NOT NULL DEFAULT FALSE
);



CREATE TABLE IF NOT EXISTS token (
    hash bytea NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) NOT NULL
);

