create table if not exists providers (
    id         uuid primary key default gen_random_uuid(),
    full_name  text not null,
    email      text unique not null,
    username   text unique not null,
    timezone   text not null default 'UTC',
    bio        text default null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists clients (
    id         uuid primary key default gen_random_uuid(),
    full_name  text not null,
    email      text unique not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists slots (
    id                              uuid primary key default gen_random_uuid(),
    provider_id                     uuid not null references providers(id) on delete cascade,
    start_time                      timestamptz not null,
    end_time                        timestamptz not null,
    available_multiple_participants boolean not null default false,
    is_free_slot                    boolean not null default true,
    created_at                      timestamptz not null default now(),
    updated_at                      timestamptz not null default now(),
    unique (provider_id, start_time, end_time),
    check (end_time >= start_time + interval '30 minutes')
);

create table if not exists appointments (
    id         uuid primary key default gen_random_uuid(),
    slot_id    uuid not null unique references slots(id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists appointment_participants (
    id             uuid primary key default gen_random_uuid(),
    appointment_id uuid not null references appointments(id) on delete cascade,
    client_id      uuid not null references clients(id) on delete cascade,
    created_at     timestamptz not null default now(),
    updated_at     timestamptz not null default now(),
    unique (appointment_id, client_id)
);

create or replace function mark_slot_booked_after_appointment()
returns trigger
language plpgsql
as $$
begin
    update slots
    set is_free_slot = false,
        updated_at = now()
    where id = new.slot_id
      and is_free_slot = true;

    return new;
end;
$$;

create trigger appointments_mark_slot_booked
after insert on appointments
for each row
execute function mark_slot_booked_after_appointment();

create index on slots (provider_id);
create index on slots (provider_id, is_free_slot, start_time, id);
create index on appointments (slot_id);
create index on appointment_participants (appointment_id);
create index on appointment_participants (client_id);
