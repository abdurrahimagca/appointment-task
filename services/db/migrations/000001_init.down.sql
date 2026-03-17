drop trigger if exists appointments_mark_slot_booked on appointments;

drop function if exists mark_slot_booked_after_appointment();

drop table if exists appointment_participants;
drop table if exists appointments;
drop table if exists slots;
drop table if exists clients;
drop table if exists providers;
