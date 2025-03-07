-- +goose Up
-- +goose StatementBegin
UPDATE risk_statuses
SET color_code = 
    CASE type
        WHEN 'assessed' THEN '#FFC107'  -- Changed from #FFFF00 (pure yellow) to #FFC107 (yellow from supported list)
        WHEN 'accepted' THEN '#28A745'  -- Changed from #00FF00 (pure green) to #28A745 (green from supported list)
        WHEN 'closed' THEN '#6C757D'    -- Changed from #808080 (neutral gray) to #6C757D (gray from supported list)
        WHEN 'pending' THEN '#6F42C1'   -- Changed from #800080 (pure purple) to #6F42C1 (purple from supported list)
        WHEN 'identified' THEN '#007BFF' -- Changed from #0000FF (pure blue) to #007BFF (blue from supported list)
        WHEN 'submitted' THEN '#17A2B8' -- Changed from #1E90FF (lighter blue) to #17A2B8 (light blue from supported list)
        WHEN 'reverted' THEN '#FD7E14'  -- Changed from #FF4500 (reddish-orange) to #FD7E14 (orange from supported list)
        WHEN 'rejected' THEN '#DC3545'  -- Changed from #DC143C (crimson red) to #DC3545 (red from supported list)
        WHEN 'approved' THEN '#28A745'  -- Changed from #008000 (dark green) to #28A745 (green from supported list)
        WHEN 'resubmitted' THEN '#FFC107' -- Changed from #FFD700 (gold) to #FFC107 (yellow from supported list)
    END
WHERE type IN ('assessed', 'accepted', 'closed', 'pending', 'identified', 'submitted', 'reverted', 'rejected', 'approved', 'resubmitted');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE risk_statuses
SET color_code = 
    CASE type
        WHEN 'assessed' THEN '#FFFF00'
        WHEN 'accepted' THEN '#00FF00'
        WHEN 'closed' THEN '#808080'
        WHEN 'pending' THEN '#800080'
        WHEN 'identified' THEN '#0000FF'
        WHEN 'submitted' THEN '#1E90FF'
        WHEN 'reverted' THEN '#FF4500'
        WHEN 'rejected' THEN '#DC143C'
        WHEN 'approved' THEN '#008000'
        WHEN 'resubmitted' THEN '#FFD700'
    END
WHERE type IN ('assessed', 'accepted', 'closed', 'pending', 'identified', 'submitted', 'reverted', 'rejected', 'approved', 'resubmitted');
-- +goose StatementEnd
