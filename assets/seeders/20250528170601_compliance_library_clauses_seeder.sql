-- +goose Up
-- +goose StatementBegin

DECLARE @version_id INT;
DECLARE @CybersecurityGovernanceID INT;
DECLARE @CybersecurityDefenseID INT;
DECLARE @CybersecurityResilienceID INT;
DECLARE @ThirdPartyandCloudComputingCybersecurityID INT;

-- Get version ID for ECC 1
SELECT @version_id = id FROM compliance_library_versions WHERE version = 'ECC 1' AND year = 2024;

-- DOMAIN 1: GOVERNANCE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '1', 'Cybersecurity Governance', 'حوكمة الأمن السيبراني',
    NULL, NULL, NULL, GETDATE(), GETDATE()
);
SELECT @CybersecurityGovernanceID = id FROM compliance_library_clauses
WHERE clause_number = '1' AND compliance_library_version_id = @version_id;

-- DOMAIN 2: CYBER DEFENSE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '2', 'Cybersecurity Defense', 'الدفاع عن الأمن السيبراني',
    NULL, NULL, NULL, GETDATE(), GETDATE()
);
SELECT @CybersecurityDefenseID = id FROM compliance_library_clauses
WHERE clause_number = '2' AND compliance_library_version_id = @version_id;

-- DOMAIN 3: CYBER RESILIENCE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '3', 'Cybersecurity Resilience', 'مرونة الأمن السيبراني',
    NULL, NULL, NULL, GETDATE(), GETDATE()
);
SELECT @CybersecurityResilienceID = id FROM compliance_library_clauses
WHERE clause_number = '3' AND compliance_library_version_id = @version_id;

-- DOMAIN 4: THIRD-PARTY & CLOUD COMPUTING
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '4', 'Third-Party and Cloud Computing Cybersecurity', 'الأمن السيبراني للحوسبة السحابية والجهات الخارجية',
    NULL, NULL, NULL, GETDATE(), GETDATE()
);
SELECT @ThirdPartyandCloudComputingCybersecurityID = id FROM compliance_library_clauses
WHERE clause_number = '4' AND compliance_library_version_id = @version_id;

-- INSERT ALL SUB-CLAUSES FOR DOMAIN 1: GOVERNANCE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES
(@version_id, '1.1', 'Cybersecurity Strategy', 'استراتيجية الأمن السيبراني', 'To ensure that cybersecurity plans, goals, initiatives, and projects are contributing to compliance with related laws and regulations.', 'لضمان أن خطط وأهداف ومبادرات ومشاريع الأمن السيبراني تساهم في الامتثال للقوانين واللوائح ذات الصلة.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.2', 'Cybersecurity Management', 'إدارة الأمن السيبراني', 'To ensure Authorizing Official''s support in implementing and managing cybersecurity programs within the organization as per related laws and regulations.', 'لضمان دعم المسؤول المخول في تنفيذ وإدارة برامج الأمن السيبراني داخل المنظمة وفقًا للقوانين واللوائح ذات الصلة.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.3', 'Cybersecurity Policies and Procedures', 'سياسات وإجراءات الأمن السيبراني', 'To ensure that cybersecurity requirements are documented, communicated, and complied with by the organization.', 'لضمان توثيق متطلبات الأمن السيبراني وتوصيلها والامتثال لها من قبل المنظمة.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.4', 'Cybersecurity Roles and Responsibilities', 'أدوار ومسؤوليات الأمن السيبراني', 'To ensure that roles and responsibilities are defined for all parties.', 'لضمان تحديد الأدوار والمسؤوليات لجميع الأطراف.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.5', 'Cybersecurity Risk Management', 'إدارة مخاطر الأمن السيبراني', 'To ensure managing cybersecurity risks in a methodological approach.', 'لضمان إدارة مخاطر الأمن السيبراني بطريقة منهجية.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.6', 'Cybersecurity in Information and Technology Project Management', 'الأمن السيبراني في إدارة مشاريع المعلومات والتكنولوجيا', 'To ensure that cybersecurity requirements are included in project management methodology.', 'لضمان إدراج متطلبات الأمن السيبراني في منهجية إدارة المشاريع.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.7', 'Compliance with Cybersecurity Standards, Laws and Regulations', 'الامتثال لمعايير وقوانين ولوائح الأمن السيبراني', 'To ensure that the organization''s cybersecurity program is in compliance with related laws and regulations.', 'لضمان أن برنامج الأمن السيبراني الخاص بالمنظمة يتوافق مع القوانين واللوائح ذات الصلة.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.8', 'Periodical Cybersecurity Review and Audit', 'المراجعة والتدقيق الدوري للأمن السيبراني', 'To ensure that cybersecurity controls are implemented and in compliance with organizational policies.', 'لضمان تنفيذ ضوابط الأمن السيبراني والامتثال لسياسات المنظمة.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.9', 'Cybersecurity in Human Resources', 'الأمن السيبراني في الموارد البشرية', 'To ensure that cybersecurity risks and requirements related to personnel.', 'لضمان مواجهة مخاطر الأمن السيبراني ومتطلباته المتعلقة بالموظفين.', @CybersecurityGovernanceID, GETDATE(), GETDATE()),
(@version_id, '1.10', 'Cybersecurity Awareness and Training Program', 'برنامج التوعية والتدريب في مجال الأمن السيبراني', 'To ensure that personnel are aware of their cybersecurity responsibilities.', 'لضمان أن يكون الموظفون على دراية بمسؤولياتهم المتعلقة بالأمن السيبراني.', @CybersecurityGovernanceID, GETDATE(), GETDATE());

-- INSERT ALL SUB-CLAUSES FOR DOMAIN 2: CYBERSECURITY DEFENSE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES
(@version_id, '2.1', 'Asset Management', 'إدارة الأصول', 'To ensure that the organization has an accurate and detailed inventory.', 'لضمان أن لدى المؤسسة جردًا دقيقًا ومفصلًا للأصول.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.2', 'Identity and Access Management', 'إدارة الهوية والوصول', 'To ensure secure and restricted logical access to information.', 'لضمان الوصول المنطقي الآمن والمقيد إلى المعلومات.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.3', 'Information System and Information Processing Facilities Protection', 'حماية نظم المعلومات ومرافق معالجة المعلومات', 'To ensure the protection of information systems.', 'لضمان حماية نظم المعلومات.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.4', 'Email Protection', 'حماية البريد الإلكتروني', 'To ensure the protection of email service from cyber risks.', 'لضمان حماية خدمة البريد الإلكتروني من المخاطر السيبرانية.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.5', 'Networks Security Management', 'إدارة أمن الشبكات', 'To ensure the protection of organization''s network from cyber risks.', 'لضمان حماية شبكة المؤسسة من المخاطر السيبرانية.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.6', 'Mobile Devices Security', 'أمن الأجهزة المحمولة', 'To ensure the protection of mobile devices (including laptops).', 'لضمان حماية الأجهزة المحمولة (بما في ذلك الحواسيب المحمولة).', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.7', 'Data and Information Protection', 'حماية البيانات والمعلومات', 'To ensure the confidentiality, integrity, and availability of data.', 'لضمان سرية وسلامة وتوفر البيانات.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.8', 'Cryptography', 'التشفير', 'To ensure the proper and efficient use of cryptography to protect.', 'لضمان الاستخدام السليم والفعال للتشفير للحماية.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.9', 'Backup and Recovery Management', 'إدارة النسخ الاحتياطي والاستعادة', 'To ensure the protection of organization''s data and software.', 'لضمان حماية بيانات وبرمجيات المؤسسة.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.10', 'Vulnerabilities Management', 'إدارة الثغرات', 'To ensure timely detection and effective remediation of vulnerabilities.', 'لضمان الكشف في الوقت المناسب والمعالجة الفعالة للثغرات.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.11', 'Penetration Testing', 'اختبار الاختراق', 'To assess and evaluate the efficiency of the organization''s cybersecurity defense.', 'لتقييم فعالية الدفاع السيبراني للمؤسسة.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.12', 'Cybersecurity Event Logs and Monitoring Management', 'إدارة سجلات الأحداث والمراقبة السيبرانية', 'To ensure timely collection, analysis, and monitoring of events.', 'لضمان جمع وتحليل ومراقبة الأحداث في الوقت المناسب.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.13', 'Cybersecurity Incident and Threat Management', 'إدارة الحوادث والتهديدات السيبرانية', 'To ensure timely identification, detection, and effective management of incidents.', 'لضمان التعرف والكشف والإدارة الفعالة للحوادث في الوقت المناسب.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.14', 'Physical Security', 'الأمن المادي', 'To ensure the protection of assets from unauthorized physical access.', 'لضمان حماية الأصول من الوصول المادي غير المصرح به.', @CybersecurityDefenseID, GETDATE(), GETDATE()),
(@version_id, '2.15', 'Web Application Security', 'أمن تطبيقات الويب', 'To ensure the protection of external web applications against cyber risks.', 'لضمان حماية تطبيقات الويب الخارجية من المخاطر السيبرانية.', @CybersecurityDefenseID, GETDATE(), GETDATE());

-- INSERT SUB-CLAUSE FOR DOMAIN 3: CYBERSECURITY RESILIENCE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES
(@version_id, '3.1', 'Cybersecurity Resilience Aspects of Business Continuity Management (BCM)', 'جوانب مرونة الأمن السيبراني في إدارة استمرارية الأعمال (BCM)', 'To ensure inclusion of cybersecurity resiliency in BCM to minimize impacts from cyber incidents.', 'لضمان إدراج مرونة الأمن السيبراني في إدارة استمرارية الأعمال لتقليل التأثيرات الناجمة عن الحوادث السيبرانية.', @CybersecurityResilienceID, GETDATE(), GETDATE());

-- INSERT SUB-CLAUSES FOR DOMAIN 4: THIRD-PARTY AND CLOUD COMPUTING
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES
(@version_id, '4.1', 'Third-Party Cybersecurity', 'الأمن السيبراني للطرف الثالث', 'To ensure the protection of assets against cybersecurity risks related to third-parties.', 'لضمان حماية الأصول ضد مخاطر الأمن السيبراني المتعلقة بأطراف ثالثة.', @ThirdPartyandCloudComputingCybersecurityID, GETDATE(), GETDATE()),
(@version_id, '4.2', 'Cloud Computing and Hosting Cybersecurity', 'الحوسبة السحابية والاستضافة والأمن السيبراني', 'To ensure remediation of cyber risks and implementation of cybersecurity requirements related to cloud and hosting.', 'لضمان معالجة المخاطر السيبرانية وتنفيذ متطلبات الأمن السيبراني المتعلقة بالسحابة والاستضافة.', @ThirdPartyandCloudComputingCybersecurityID, GETDATE(), GETDATE());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_library_clauses;
-- +goose StatementEnd