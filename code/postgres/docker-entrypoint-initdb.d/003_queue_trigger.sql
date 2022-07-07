create or replace function audit_insert_trigger_fnc()
  returns trigger as $$
    begin
        insert into 
            app.queue_audit_log ( 
             action
            ,user_id
            ,content_item_id
            ,create_date
            )
        values(
             new."action"
            ,new."user_id"
            ,new."content_item_id"
            ,new."create_date"
        );

        return new;
    end;
$$ language 'plpgsql';


create trigger audit_insert_trigger
  after insert on app.audit_log
  for each row
  execute procedure audit_insert_trigger_fnc();