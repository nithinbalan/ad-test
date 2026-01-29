-- +goose Up
-- +goose StatementBegin
/* ================================
   legislative_bodies
================================ */
EXEC sp_rename 
    'dbo.legislative_bodies', 
    'compliance_legislative_bodies';

EXEC sp_rename 
    'dbo.CHK_legislative_bodies_status',
    'CHK_compliance_legislative_bodies_status',
    'OBJECT';


/* ================================
   legislation_types
================================ */
EXEC sp_rename 
    'dbo.legislation_types',
    'compliance_legislation_types';

EXEC sp_rename 
    'dbo.CHK_legislation_types_status',
    'CHK_compliance_legislation_types_status',
    'OBJECT';


/* ================================
   issuing_authorities
================================ */
EXEC sp_rename 
    'dbo.issuing_authorities',
    'compliance_issuing_authorities';

EXEC sp_rename 
    'dbo.CHK_issuing_authorities_status',
    'CHK_compliance_issuing_authorities_status',
    'OBJECT';


/* ================================
   importance_classifications
================================ */
EXEC sp_rename 
    'dbo.importance_classifications',
    'compliance_importance_classifications';

-- UNIQUE constraint
EXEC sp_rename 
    'dbo.CHK_importance_classifications_type',
    'UQ_compliance_importance_classifications_type',
    'OBJECT';

-- CHECK constraint
EXEC sp_rename 
    'dbo.CHK_importance_classifications_status',
    'CHK_compliance_importance_classifications_status',
    'OBJECT';


/* ================================
   secrecy_classifications
================================ */
EXEC sp_rename 
    'dbo.secrecy_classifications',
    'compliance_secrecy_classifications';

EXEC sp_rename 
    'dbo.CHK_secrecy_classifications_status',
    'CHK_compliance_secrecy_classifications_status',
    'OBJECT';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
/* ================================
   secrecy_classifications
================================ */
EXEC sp_rename 
    'dbo.CHK_compliance_secrecy_classifications_status',
    'CHK_secrecy_classifications_status',
    'OBJECT';

EXEC sp_rename 
    'dbo.compliance_secrecy_classifications',
    'secrecy_classifications';


/* ================================
   importance_classifications
================================ */
EXEC sp_rename 
    'dbo.CHK_compliance_importance_classifications_status',
    'CHK_importance_classifications_status',
    'OBJECT';

EXEC sp_rename 
    'dbo.UQ_compliance_importance_classifications_type',
    'CHK_importance_classifications_type',
    'OBJECT';

EXEC sp_rename 
    'dbo.compliance_importance_classifications',
    'importance_classifications';


/* ================================
   issuing_authorities
================================ */
EXEC sp_rename 
    'dbo.CHK_compliance_issuing_authorities_status',
    'CHK_issuing_authorities_status',
    'OBJECT';

EXEC sp_rename 
    'dbo.compliance_issuing_authorities',
    'issuing_authorities';


/* ================================
   legislation_types
================================ */
EXEC sp_rename 
    'dbo.CHK_compliance_legislation_types_status',
    'CHK_legislation_types_status',
    'OBJECT';

EXEC sp_rename 
    'dbo.compliance_legislation_types',
    'legislation_types';


/* ================================
   legislative_bodies
================================ */
EXEC sp_rename 
    'dbo.CHK_compliance_legislative_bodies_status',
    'CHK_legislative_bodies_status',
    'OBJECT';

EXEC sp_rename 
    'dbo.compliance_legislative_bodies',
    'legislative_bodies';
-- +goose StatementEnd
