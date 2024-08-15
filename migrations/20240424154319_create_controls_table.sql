-- +goose Up
-- +goose StatementBegin
CREATE TABLE controls
(
    id INT PRIMARY KEY IDENTITY(1,1),
    control_type_id INT NOT NULL,
    control_category_id INT NOT NULL,
    owner_unit_id INT,
    owner_id INT,
    control_implementation_status_id INT,
    control_efficiency_status_id INT,
    efficiency_score INT DEFAULT 0,
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    code VARCHAR(45),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX),
    objective NVARCHAR(MAX),
    review_frequency_id INT,
    last_review_date DATE,
    next_review_date DATE,
    last_review_note NVARCHAR(MAX),
    test_responsible_user_id INT NOT NULL,
    test_frequency_id INT,
    last_test_date DATE,
    next_test_date DATE,
    implemented_at DATE,
    marked_as_auditable_item BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    lifecycle_status_id INT NOT NULL,
    last_reviewed_user_id INT NULL,
    control_review_status_id INT NULL,

    CONSTRAINT FK_controls_control_types FOREIGN KEY (control_type_id) REFERENCES control_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_control_categories FOREIGN KEY (control_category_id) REFERENCES control_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_organization_units FOREIGN KEY (owner_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_user_1 FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_control_implementation_statuses FOREIGN KEY (control_implementation_status_id) REFERENCES control_implementation_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_control_efficiency_statuses FOREIGN KEY (control_efficiency_status_id) REFERENCES control_efficiency_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_review_frequencies1 FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_user_2 FOREIGN KEY (test_responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_review_frequencies2 FOREIGN KEY (test_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_users_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_users_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_users_5 FOREIGN KEY (last_reviewed_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_lifecycle_statuses FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_controls_review_statuses FOREIGN KEY (control_review_status_id) REFERENCES control_review_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS controls;
-- +goose StatementEnd
