-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixID INT;

SELECT @RiskMatrixID = id
FROM risk_matrix
WHERE type = 'isms';

INSERT INTO risk_matrix_impacts
(
    risk_matrix_id,
    title_en,
    title_ar,
    score,
    description_en,
    description_ar,
    status
)
VALUES
(
@RiskMatrixID,
'Very Low',
N'منخفض جدًا',
1,
'Causes negligible or insignificant impact on operations, finances, compliance, or reputation. Little or no corrective action is required.',
N'يسبب تأثيرًا ضئيلًا أو غير ملحوظ على العمليات أو الجوانب المالية أو الامتثال أو السمعة. لا يتطلب سوى القليل أو لا يتطلب أي إجراءات تصحيحية.',
'active'
),
(
@RiskMatrixID,
'Low',
N'منخفض',
2,
'Causes minor disruption with limited financial or operational impact. Recovery is straightforward with minimal business consequences.',
N'يسبب اضطرابًا بسيطًا مع تأثير محدود على الجوانب التشغيلية أو المالية. يمكن التعافي بسهولة مع آثار محدودة على الأعمال.',
'active'
),
(
@RiskMatrixID,
'Moderate',
N'متوسط',
3,
'Causes noticeable operational disruption, moderate financial loss, limited legal or regulatory impact, and recoverable reputational damage.',
N'يسبب اضطرابًا ملحوظًا في العمليات، وخسائر مالية متوسطة، وتأثيرًا قانونيًا أو تنظيميًا محدودًا، مع أضرار يمكن استعادة آثارها على السمعة.',
'active'
),
(
@RiskMatrixID,
'High',
N'عالٍ',
4,
'Causes major operational disruption, substantial financial loss, serious regulatory consequences, or significant reputational impact requiring executive attention.',
N'يسبب اضطرابًا كبيرًا في العمليات، وخسائر مالية كبيرة، وعواقب تنظيمية جسيمة، أو تأثيرًا كبيرًا على السمعة يتطلب اهتمام الإدارة العليا.',
'active'
),
(
@RiskMatrixID,
'Catastrophic',
N'كارثي',
5,
'Causes severe damage to the organization, such as complete service disruption, significant financial loss, major regulatory penalties, or long-term reputational damage.',
N'يسبب أضرارًا جسيمة للمنظمة، مثل توقف كامل للخدمات، أو خسائر مالية كبيرة، أو عقوبات تنظيمية كبيرة، أو أضرار طويلة الأمد بالسمعة.',
'active'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM risk_matrix_impacts
WHERE risk_matrix_id = (
SELECT id FROM risk_matrix WHERE type='isms'
);

-- +goose StatementEnd