-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_result_performance_measures (
    id INT PRIMARY KEY IDENTITY(1,1),
    test_result_id INT NOT NULL,
    criterion NVARCHAR(MAX),
    target NVARCHAR(MAX),
    achieved_result NVARCHAR(MAX),
    remarks NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_test_result_performance_measures_test_result
        FOREIGN KEY (test_result_id) 
        REFERENCES test_and_exercise_test_results(id)
        ON DELETE NO ACTION 
        ON UPDATE NO ACTION,
    CONSTRAINT FK_test_result_performance_measures_created_by
        FOREIGN KEY (created_by_id) 
        REFERENCES users(id)
        ON DELETE NO ACTION 
        ON UPDATE NO ACTION,
    CONSTRAINT FK_test_result_performance_measures_updated_by
        FOREIGN KEY (updated_by_id) 
        REFERENCES users(id)
        ON DELETE NO ACTION 
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test_result_performance_measures;
-- +goose StatementEnd
