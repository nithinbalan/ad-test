-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixID INT;

SELECT @RiskMatrixID = id
FROM risk_matrix
WHERE type='isms';
INSERT INTO risk_matrix_likelihoods
(
    risk_matrix_id,
    title_en,
    title_ar,
    score,
    description_en,
    description_ar,
    probability,
    status
)
VALUES
(
@RiskMatrixID,
'Very Rare',
N'نادر جدًا',
1,
'The event is highly unlikely to occur and would only happen under exceptional circumstances.',
N'من غير المرجح للغاية أن يحدث هذا الحدث، ولن يقع إلا في ظروف استثنائية.',
'Very Rare',
'active'
),
(
@RiskMatrixID,
'Unlikely',
N'غير محتمل',
2,
'The event is not expected to occur under normal circumstances but remains possible.',
N'لا يُتوقع حدوث هذا الحدث في الظروف العادية، لكنه يظل ممكنًا.',
'Unlikely',
'active'
),
(
@RiskMatrixID,
'Possible',
N'محتمل',
3,
'The event could occur at some time. There is a moderate probability of occurrence.',
N'قد يحدث هذا الحدث في وقت ما، مع وجود احتمال متوسط لحدوثه.',
'Possible',
'active'
),
(
@RiskMatrixID,
'Likely',
N'مرجح',
4,
'The event will probably occur in many circumstances. There is a high chance that it will happen.',
N'من المرجح أن يحدث هذا الحدث في العديد من الظروف، مع وجود احتمال مرتفع لحدوثه.',
'Likely',
'active'
),
(
@RiskMatrixID,
'Almost Certain',
N'شبه مؤكد',
5,
'The event is expected to occur in most circumstances. It has a very high probability of occurring and may happen frequently.',
N'من المتوقع أن يحدث هذا الحدث في معظم الظروف، مع احتمال مرتفع جدًا لحدوثه وقد يتكرر بشكل متكرر.',
'Almost Certain',
'active'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM risk_matrix_likelihoods
WHERE risk_matrix_id = (
SELECT id FROM risk_matrix WHERE type='isms'
);

-- +goose StatementEnd