-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

INSERT INTO risk_matrix_likelihoods (risk_matrix_id, title_en, title_ar, score, description_en, description_ar, probability, status)
VALUES
    (@RiskMatrixId, 'Negligible', N'تافه', 1, 
     'Expected to occur more than once per week.',
     N'من المتوقع أن يحدث أكثر من مرة في الأسبوع.',
     '< 90%',
     'active'),
     
    (@RiskMatrixId, 'Minor', N'صغير', 2,
     'Expected to occur about once every month.',
     N'من المتوقع أن يحدث حوالي مرة كل شهر.',
     '< 90%',
     'active'),
     
    (@RiskMatrixId, 'Moderate', N'متوسط', 3,
     'Expected to occur about once every year.',
     N'من المتوقع أن يحدث حوالي مرة كل سنة.',
     '>30% and < 60%',
     'active'),
     
    (@RiskMatrixId, 'Critical', N'حرج', 4,
     'Expected to occur between 1-3 years.',
     N'من المتوقع أن يحدث بين 1-3 سنوات.',
     '>10% and < 30%',
     'active'),
     
    (@RiskMatrixId, 'Catastrophic', N'كارثي', 5,
     'Expected to occur after a period greater than 3 years.',
     N'من المتوقع أن يحدث بعد فترة تزيد عن 3 سنوات.',
     '< 10%',
     'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_matrix_likelihoods;
-- +goose StatementEnd
