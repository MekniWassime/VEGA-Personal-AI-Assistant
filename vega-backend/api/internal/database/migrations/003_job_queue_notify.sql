CREATE FUNCTION notify_job_created() RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('job_created', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER job_created_notify
AFTER INSERT ON job_queue
FOR EACH ROW
EXECUTE FUNCTION notify_job_created();

---- create above / drop below ----

DROP TRIGGER job_created_notify ON job_queue;
DROP FUNCTION notify_job_created;
