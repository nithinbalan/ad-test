-- +goose Up
-- +goose StatementBegin

-- Declare variable for Risk Matrix Risk Rating IDs
DECLARE @VeryLowRatingID INT, @LowRatingID INT, @MediumRatingID INT, @HighRatingID INT, @VeryHighRatingID INT;

-- Fetching Risk Matrix Risk Rating IDs
SELECT @VeryLowRatingID = id FROM risk_matrix_risk_ratings WHERE type = 'very_low';
SELECT @LowRatingID = id FROM risk_matrix_risk_ratings WHERE type = 'low';
SELECT @MediumRatingID = id FROM risk_matrix_risk_ratings WHERE type = 'medium';
SELECT @HighRatingID = id FROM risk_matrix_risk_ratings WHERE type = 'high';
SELECT @VeryHighRatingID = id FROM risk_matrix_risk_ratings WHERE type = 'very_high';

-- Inserting Risk Treatment Guidelines
INSERT INTO risk_treatment_guidelines (risk_matrix_risk_rating_id, description_en, description_ar, status)
VALUES
    -- Very Low Risk Guidelines
    (@VeryLowRatingID, 
     'No treatment is needed for very low risks, only monitoring.',
     'لا يلزم علاج للمخاطر المنخفضة جدًا، فقط المراقبة.',
     'active'),

    -- Low Risk Guidelines
    (@LowRatingID, 
     'Low risks require minimal action and periodic review.',
     'تتطلب المخاطر المنخفضة إجراءً بسيطًا ومراجعة دورية.',
     'active'),

    -- Medium Risk Guidelines
    (@MediumRatingID, 
     'Medium risks need mitigation plans and monitoring.',
     'تتطلب المخاطر المتوسطة خطط تخفيف ومراقبة.',
     'active'),

    -- High Risk Guidelines
    (@HighRatingID, 
     'High risks need urgent attention and risk reduction plans.',
     'تتطلب المخاطر العالية اهتمامًا عاجلاً وخطط تقليل المخاطر.',
     'active'),

    -- Very High Risk Guidelines
    (@VeryHighRatingID, 
     'Very high risks must be mitigated immediately.',
     'يجب تخفيف المخاطر العالية جدًا على الفور.',
     'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Truncate the `risk_treatment_guidelines` table
TRUNCATE TABLE risk_treatment_guidelines;

-- +goose StatementEnd
