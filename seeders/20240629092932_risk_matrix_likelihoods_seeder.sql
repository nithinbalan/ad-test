-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

INSERT INTO risk_matrix_likelihoods (risk_matrix_id, title_en, title_ar, score, description_en, description_ar, status)
VALUES
    (@RiskMatrixId, 'Very Low', N'منخفض جداً', 1, 
     'There is a very low likelihood of this risk event occurring. It is highly unlikely and may not significantly impact operations if it does.',
     N'هناك احتمالية منخفضة جداً لحدوث هذا الحدث الخطر. من غير المرجح بشكل كبير وقد لا يؤثر بشكل كبير على العمليات إذا حدث.',
     'active'),
     
    (@RiskMatrixId, 'Low', N'منخفض', 2,
     'There is a low likelihood of this risk event occurring. It may occur occasionally but is not expected to cause significant disruptions.',
     N'هناك احتمالية منخفضة لحدوث هذا الحدث الخطر. قد يحدث بين الحين والآخر ولكن لا يُتوقع أن يسبب اضطرابات كبيرة.',
     'active'),
     
    (@RiskMatrixId, 'Medium', N'متوسط', 3,
     'There is a moderate likelihood of this risk event occurring. It could occur under certain circumstances and may cause noticeable disruptions.',
     N'هناك احتمالية متوسطة لحدوث هذا الحدث الخطر. قد يحدث تحت ظروف معينة وقد يسبب اضطرابات ملحوظة.',
     'active'),
     
    (@RiskMatrixId, 'High', N'عالي', 4,
     'There is a high likelihood of this risk event occurring. It is expected to occur frequently and could cause significant disruptions.',
     N'هناك احتمالية عالية لحدوث هذا الحدث الخطر. من المتوقع أن يحدث بشكل متكرر وقد يسبب اضطرابات كبيرة.',
     'active'),
     
    (@RiskMatrixId, 'Very High', N'عالية جداً', 5,
     'There is a very high likelihood of this risk event occurring. It is expected to occur frequently and could have severe impacts on operations.',
     N'هناك احتمالية عالية جداً لحدوث هذا الحدث الخطر. من المتوقع أن يحدث بشكل متكرر وقد يكون له تأثيرات خطيرة على العمليات.',
     'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_matrix_likelihoods;
-- +goose StatementEnd
