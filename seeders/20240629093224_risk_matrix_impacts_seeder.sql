-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

INSERT INTO risk_matrix_impacts (risk_matrix_id, title_en, title_ar, score, description_en, description_ar, status)
VALUES
    (@RiskMatrixId, 'Minor', N'طفيف', 1, 
     'Minor impact. This could result in slight inconvenience or minor delays but is easily manageable.',
     N'تأثير طفيف. يمكن أن يتسبب في إزعاج طفيف أو تأخيرات بسيطة ولكن يمكن إدارته بسهولة.',
     'active'),
     
    (@RiskMatrixId, 'Moderate', N'متوسط', 2,
     'Moderate impact. This could lead to noticeable disruptions requiring some resources to mitigate.',
     N'تأثير متوسط. يمكن أن يؤدي إلى اضطرابات ملحوظة تتطلب بعض الموارد للتخفيف منها.',
     'active'),
     
    (@RiskMatrixId, 'Significant', N'كبير', 3,
     'Significant impact. This would cause substantial disruptions and may require significant resources to recover from.',
     N'تأثير كبير. قد يؤدي إلى اضطرابات كبيرة وقد يتطلب موارد كبيرة للتعافي منها.',
     'active'),
     
    (@RiskMatrixId, 'Major', N'كبير جداً', 4,
     'Major impact. This could have severe consequences and may significantly impact operations or services.',
     N'تأثير كبير جداً. يمكن أن يكون له عواقب خطيرة وقد يؤثر بشكل كبير على العمليات أو الخدمات.',
     'active'),
     
    (@RiskMatrixId, 'Catastrophic', N'كارثي', 5,
     'Catastrophic impact. This would have devastating consequences and could severely impact business continuity.',
     N'تأثير كارثي. قد يكون له عواقب مدمرة وقد يؤثر بشكل كبير على استمرارية الأعمال.',
     'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
