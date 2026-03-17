-- name: GetProviderByUsername :one
select
    p.id,
    p.full_name,
    p.email,
    p.username,
    p.timezone,
    p.bio
from providers p
where p.username = sqlc.arg('username');

-- name: GetSlotsByProviderUsernameAndDate :many
select
    s.id,
    timezone(p.timezone, s.start_time)::date as date_of_day,
    s.start_time,
    s.end_time,
    timezone(p.timezone, s.start_time) as slot_start_at,
    timezone(p.timezone, s.end_time) as slot_end_at,
    s.available_multiple_participants
from slots s
join providers p on p.id = s.provider_id
where p.username = sqlc.arg('username')
  and timezone(p.timezone, s.start_time)::date = sqlc.arg('date_of_day')
  and s.is_free_slot = true
  and (
      sqlc.narg('cursor_start_time')::timestamptz is null
      or s.start_time > sqlc.narg('cursor_start_time')::timestamptz
      or (
          s.start_time = sqlc.narg('cursor_start_time')::timestamptz
          and s.id > sqlc.narg('cursor_id')::uuid
      )
  )
order by s.start_time, s.id
limit sqlc.arg('limit');

-- name: BookSlot :one
insert into appointments (slot_id)
select s.id
from slots s
where s.id = sqlc.arg('slot_id')
  and s.is_free_slot = true
  and not exists (
      select 1
      from appointments a
      where a.slot_id = s.id
  )
returning id, slot_id, created_at;

-- name: UpsertClientByEmail :one
insert into clients (full_name, email)
values (sqlc.arg('full_name'), sqlc.arg('email'))
on conflict (email)
do update
set full_name = excluded.full_name,
    updated_at = now()
returning id, full_name, email;

-- name: AddParticipantToAppointment :exec
insert into appointment_participants (appointment_id, client_id)
values (sqlc.arg('appointment_id'), sqlc.arg('client_id'));

-- name: GetAppointmentWithDetails :one
select
    a.id as appointment_id,
    p.full_name as provider_name,
    p.timezone as provider_timezone,
    timezone(p.timezone, s.start_time)::date as date_of_day,
    s.start_time,
    s.end_time,
    timezone(p.timezone, s.start_time) as slot_start_at,
    timezone(p.timezone, s.end_time) as slot_end_at,
    s.available_multiple_participants,
    json_agg(
        json_build_object(
            'clientId', c.id,
            'fullName', c.full_name,
            'email', c.email
        )
        order by ap.created_at
    ) as participants
from appointments a
join slots s on s.id = a.slot_id
join providers p on p.id = s.provider_id
join appointment_participants ap on ap.appointment_id = a.id
join clients c on c.id = ap.client_id
where a.id = sqlc.arg('appointment_id')
group by a.id, p.full_name, p.timezone, s.start_time, s.end_time, s.available_multiple_participants;
