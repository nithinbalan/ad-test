-- +goose Up
-- +goose StatementBegin
UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur after a period greater than 3 years.', description_ar = 'ومن المتوقع أن يحدث بعد فترة أطول من 3 سنوات.'
WHERE title_en = 'Negligible';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur between 1-3 years.', description_ar = 'ومن المتوقع أن يحدث خلال فترة تتراوح بين 1-3 سنوات.'
WHERE title_en = 'Minor';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur about once every month.', description_ar = 'ومن المتوقع أن يحدث ذلك مرة واحدة كل شهر تقريبًا.'
WHERE title_en = 'Critical';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur more than once per week.', description_ar = 'ومن المتوقع أن يحدث أكثر من مرة في الأسبوع.'
WHERE title_en = 'Catastrophic';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur more than once per week.', description_ar = 'ومن المتوقع أن يحدث أكثر من مرة في الأسبوع.'
WHERE title_en = 'Negligible';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur about once every month.', description_ar = 'ومن المتوقع أن يحدث ذلك مرة واحدة كل شهر تقريبًا.'
WHERE title_en = 'Minor';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur between 1-3 years.', description_ar = 'ومن المتوقع أن يحدث خلال فترة تتراوح بين 1-3 سنوات.'
WHERE title_en = 'Critical';

UPDATE risk_matrix_likelihoods
SET description_en = 'Expected to occur after a period greater than 3 years.', description_ar = 'ومن المتوقع أن يحدث بعد فترة أطول من 3 سنوات.'
WHERE title_en = 'Catastrophic';
-- +goose StatementEnd
