CREATE TYPE role_type as ENUM ('admin', 'customer', 'guest_user');
CREATE TABLE roles (
  id UUID PRIMARY KEY,
  name role_type UNIQUE,
  description TEXT
);


CREATE TYPE admin_permission as ENUM ('view', 'edit', 'full_access');

CREATE TABLE permissions (
  id UUID PRIMARY KEY,
  name admin_permission UNIQUE,
  description TEXT
);