-- +goose Up
-- +goose StatementBegin
INSERT INTO strategic_profiles (title_en, title_ar, address, location, email, industry, phone, linkedin, x, facebook, instagram, youtube, about_en, about_ar, vision_en, vision_ar, mission_en, mission_ar, created_at, created_by)
VALUES
('Demo organization', 'Demo organization ar', 'Address 1', 'Location 1', 'email1@example.com', 'Industry 1', '123-456-7890', 'linkedin1', 'X 1', 'facebook1', 'instagram1', 'youtube1', 'About EN 1', 'About AR 1', 'Vision EN 1', 'Vision AR 1', 'Mission EN 1', 'Mission AR 1', GETDATE(), 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE strategic_profiles;
-- +goose StatementEnd
