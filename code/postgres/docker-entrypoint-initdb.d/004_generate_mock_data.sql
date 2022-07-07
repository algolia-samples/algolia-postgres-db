-- set random seed for repeatable random data generation
SELECT setseed(0.8);
DO $$
    DECLARE
        -- configurable parameters for data generation
        nr_lines integer := 20;
        user_min integer := 10;
        user_max integer := 20;
        citm_min integer := 1500;
        citm_max integer := 2300;
        actn_min integer := 1;
        actn_max integer := 3;
    BEGIN
        with
            -- generate user_ids
            users as (
                select generate_series(user_min, user_max) as user_id
            )
            -- generate content_ids
           ,content as (
               select generate_series(citm_min, citm_max) as content_id
            )
            -- generate action_ids
           ,actions as (
               select generate_series(actn_min, actn_max) as action_id
            )
            -- get the cartesian product of the above in a random sort
           ,limited_data as (
               select
                 random() randomizer
                 ,* 
               from users, content, actions 
               order by randomizer
               limit nr_lines
            )
        insert 
            into app.audit_log (
                action
                ,user_id
                ,content_item_id
            )
            select
                 action_id
                ,user_id
                ,content_id
            from limited_data
        ;
END $$
;

-- select * from audit_log order by content_item_id, user_id, action;