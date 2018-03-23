delete from sym_trigger_router;
delete from sym_trigger;
delete from sym_router;
delete from sym_channel where channel_id in ('discovery');
delete from sym_node_group_link;
delete from sym_node_group;
delete from sym_node_host;
delete from sym_node_identity;
delete from sym_node_security;
delete from sym_node;

insert into sym_channel
(channel_id, processing_order, max_batch_size, enabled, description)
values('discovery', 1, 100000, 1, 'patient and location data');

insert into sym_node_group (node_group_id) values ('cloud');
insert into sym_node_group (node_group_id) values ('local');

insert into sym_node_group_link (source_node_group_id, target_node_group_id, data_event_action) values ('cloud', 'local', 'W');
insert into sym_node_group_link (source_node_group_id, target_node_group_id, data_event_action) values ('local', 'cloud', 'P');

insert into sym_trigger
(trigger_id,source_table_name,channel_id,last_update_time,create_time)
values('patients','patients','discovery',current_timestamp,current_timestamp);

insert into sym_trigger
(trigger_id,source_table_name,channel_id,last_update_time,create_time)
values('locations','locations','discovery',current_timestamp,current_timestamp);

insert into sym_trigger
(trigger_id,source_table_name,channel_id,last_update_time,create_time)
values('connections','connections','discovery',current_timestamp,current_timestamp);

insert into sym_trigger
(trigger_id,source_table_name,channel_id, sync_on_insert, sync_on_update, sync_on_delete,last_update_time,create_time)
values('patients_cloud','patients','discovery',0,0,0,current_timestamp,current_timestamp);

insert into sym_trigger
(trigger_id,source_table_name,channel_id, sync_on_insert, sync_on_update, sync_on_delete,last_update_time,create_time)
values('locations_cloud','locations','discovery',0,0,0,current_timestamp,current_timestamp);

insert into sym_trigger
(trigger_id,source_table_name,channel_id, sync_on_insert, sync_on_update, sync_on_delete,last_update_time,create_time)
values('connections_cloud','connections','discovery',0,0,0,current_timestamp,current_timestamp);

insert into sym_router
(router_id,source_node_group_id,target_node_group_id,router_type,create_time,last_update_time)
values('cloud_2_local', 'cloud', 'local', 'default',current_timestamp, current_timestamp);

insert into sym_router
(router_id,source_node_group_id,target_node_group_id,router_type,create_time,last_update_time)
values('local_2_cloud', 'local', 'cloud', 'default',current_timestamp, current_timestamp);

insert into sym_router
(router_id,source_node_group_id,target_node_group_id,router_type,router_expression,create_time,last_update_time)
values('cloud_2_select_local','cloud', 'local', 'lookuptable',
       'LOOKUP_TABLE=LOCATIONS KEY_COLUMN=PATIENT_ID LOOKUP_KEY_COLUMN=PATIENT_ID EXTERNAL_ID_COLUMN=LOCATION_ID', current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('patients','cloud_2_select_local', 100, current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('locations','cloud_2_select_local', 100, current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('connections','cloud_2_select_local', 100, current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('patients','local_2_cloud', 200, current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('locations','local_2_cloud', 200, current_timestamp, current_timestamp);

insert into sym_trigger_router
(trigger_id,router_id,initial_load_order,last_update_time,create_time)
values('connections','local_2_cloud', 200, current_timestamp, current_timestamp);

insert into sym_node (node_id,node_group_id,external_id,sync_enabled,sync_url,schema_version,symmetric_version,database_type,database_version,heartbeat_time,timezone_offset,batch_to_send_count,batch_in_error_count,created_at_node_id)
 values ('cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326','cloud','f7e41e48-ec79-4c78-9db6-37c0c4f78326',1,null,null,null,null,null,current_timestamp,null,0,0,'cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326');
insert into sym_node (node_id,node_group_id,external_id,sync_enabled,sync_url,schema_version,symmetric_version,database_type,database_version,heartbeat_time,timezone_offset,batch_to_send_count,batch_in_error_count,created_at_node_id)
 values ('local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa','local','2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa',1,null,null,null,null,null,current_timestamp,null,0,0,'cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326');

insert into sym_node_security (node_id,node_password,registration_enabled,registration_time,initial_load_enabled,initial_load_time,created_at_node_id)
 values ('cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326','3fdf43e0e03869a235adb28e6e7512b7',0,current_timestamp,0,current_timestamp,'cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326');
insert into sym_node_security (node_id,node_password,registration_enabled,registration_time,initial_load_enabled,initial_load_time,created_at_node_id)
 values ('local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa','3fdf43e0e03869a235adb28e6e7512b7',1,null,1,null,'cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326');

insert into sym_node_identity values ('cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326');
