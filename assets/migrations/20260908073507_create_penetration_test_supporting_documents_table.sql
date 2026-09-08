-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_supporting_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    penetration_test_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_supporting_documents_penetration_tests FOREIGN KEY (penetration_test_id)
        REFERENCES penetration_tests(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_penetration_test_supporting_documents_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_penetration_test_supporting_documents_users_2 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_supporting_documents;
-- +goose StatementEnd
