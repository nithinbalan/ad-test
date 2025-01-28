-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_processes
ADD is_recovered TINYINT NOT NULL DEFAULT 0;

ALTER TABLE test_and_exercise_processes 
ADD timeliness_id INT NULL;


ALTER TABLE test_and_exercise_processes
ADD CONSTRAINT FK_test_and_exercise_processes_timeliness FOREIGN KEY (timeliness_id) REFERENCES timeliness(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_processes
DROP CONSTRAINT FK_test_and_exercise_processes_timeliness;

ALTER TABLE test_and_exercise_processes
DROP COLUMN timeliness_id;

ALTER TABLE test_and_exercise_processes
DROP COLUMN is_recovered;
-- +goose StatementEnd
