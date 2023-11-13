CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- roles
create table
  public.roles (
    id uuid default uuid_generate_v4 () not null constraint roles_pk primary key,
    created_at timestamptz default now () not null,
    updated_at timestamptz default now () not null,
    deleted_at timestamptz,
    created_by text default '' not null,
    updated_by text default '' not null,
    deleted_by text default '' not null,
    name text not null,
    description text default '' not null
  );

create unique index roles_name_uindex on public.roles (name);

-- permissions
create table
  permissions (
    id uuid default uuid_generate_v4 () not null constraint permissions_pk primary key,
    created_at timestamptz default now () not null,
    updated_at timestamptz default now () not null,
    deleted_at timestamptz,
    created_by text default '' not null,
    updated_by text default '' not null,
    deleted_by text default '' not null,
    role_id uuid not null constraint permissions_roles_id_fk references roles,
    resource text not null,
    is_create boolean default false not null,
    is_read boolean default false not null,
    is_update boolean default false not null,
    is_delete boolean default false not null
  );

create unique index permissions_role_id_resource_is_create_is_read_is_update_is_del on permissions (
  role_id,
  resource,
  is_create,
  is_read,
  is_update,
  is_delete
);

create index permissions_role_id_resource_index on public.permissions (role_id, resource);
