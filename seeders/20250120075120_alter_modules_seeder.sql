-- +goose Up
-- +goose StatementBegin
-- update test and exercise route
UPDATE modules
SET route = '/business-continuity/test-and-exercise',
title_en = 'Test and Exercise',
title_ar = 'الاختبار والتمرين',
identifier = 'test_and_exercise'
WHERE identifier = 'test_and_excercise';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET route = '/business-continuity/test-and-excercise'
title_en = 'Test and Excercise',
title_ar = 'الاختبار والتمرين',
identifier = 'test_and_excercise'
WHERE identifier = 'test_and_excercise';
-- +goose StatementEnd
