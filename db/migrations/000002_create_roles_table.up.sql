CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  name_hash TEXT NOT NULL UNIQUE,
  description TEXT,
  is_deleted BOOLEAN DEFAULT FALSE,
  -- Ensure name is unique
  UNIQUE (name),
  -- Timestamps for record keeping
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);