alter table person add constraint ch_status check ( status in ('active','inactive','blocked') );

alter table person add constraint ch_role check ( role in ('user','admin','blocked') );
