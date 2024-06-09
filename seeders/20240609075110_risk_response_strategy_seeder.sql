-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_response_strategies (title_en, title_ar, description_en, description_ar, status) VALUES 
('Accept', 'قبول', 'Risk event that is acknowledged and accepted without any further action to mitigate, transfer or avoid.', 'حدث المخاطر الذي يتم الاعتراف به وقبوله دون اتخاذ أي إجراء آخر للتخفيف أو النقل أو التجنب.', 'active'),
('Avoid', 'تجنب', 'Risk event that would affect the business continuity of Almarai and is fundamental to achieving goals and objectives or will have critical impact on achieving desired outcomes, to the extent that one or more critical objectives will not be achieved.', 'حدث مخاطر من شأنه أن يؤثر على استمرارية أعمال المراعي وهو أساسي لتحقيق الأهداف أو سيكون له تأثير حاسم على تحقيق النتائج المرجوة، إلى الحد الذي لن يتم فيه تحقيق هدف أو أكثر من الأهداف الحرجة.', 'active'),
('Transfer', 'نقل', 'Risk event that is transferred to another party to manage, often through insurance or outsourcing.', 'حدث المخاطر الذي يتم نقله إلى طرف آخر لإدارته، غالبًا من خلال التأمين أو الاستعانة بمصادر خارجية.', 'active'),
('Mitigate', 'تخفيف', 'Risk event that is reduced in likelihood or impact through specific actions.', 'حدث المخاطر الذي يتم تقليله في الاحتمالية أو التأثير من خلال إجراءات محددة.', 'active'),
('Tolerate', 'تحمل', 'Risk event that is accepted but monitored for changes in conditions.', 'حدث المخاطر الذي يتم قبوله ولكن يتم مراقبته للتغيرات في الظروف.', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM risk_response_strategies WHERE title_en IN ('Accept', 'Avoid', 'Transfer', 'Mitigate', 'Tolerate');
-- +goose StatementEnd
