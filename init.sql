-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    role VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_name      TEXT,
    company_address   TEXT,
    drive             TEXT,
    type_of_drive     TEXT,
    follow_up         TEXT,
    is_contacted      BOOLEAN DEFAULT false,
    remarks           TEXT,
    contact_details   TEXT,
    hr1_details       TEXT,
    hr2_details       TEXT,
    package           TEXT,
    assigned_officer  TEXT[] DEFAULT '{}',
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS companies_temp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id        UUID REFERENCES companies(id),
    company_name      TEXT,
    company_address   TEXT,
    drive             TEXT,
    type_of_drive     TEXT,
    follow_up         TEXT,
    is_contacted      BOOLEAN DEFAULT false,
    remarks           TEXT,
    contact_details   TEXT,
    hr1_details       TEXT,
    hr2_details       TEXT,
    package           TEXT,
    assigned_officer  TEXT[] DEFAULT '{}',
    status            TEXT DEFAULT 'pending',
    created_by        TEXT,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for companies table
CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(company_name);
CREATE INDEX IF NOT EXISTS idx_companies_drive ON companies(drive);
CREATE INDEX IF NOT EXISTS idx_companies_is_contacted ON companies(is_contacted);

-- Create indexes for events table
CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);

-- Insert initial sample data
INSERT INTO users (username, email, role, password, created_at) VALUES 
    ('admin', 'admin@company.com', 'Admin', 'password', '2024-01-01T00:00:00Z'),
    ('manager', 'manager@company.com', 'Manager', 'password', '2024-01-01T00:00:00Z'),
    ('officer', 'officer@company.com', 'Officer', 'password', '2024-01-01T00:00:00Z')
ON CONFLICT (username) DO NOTHING;
