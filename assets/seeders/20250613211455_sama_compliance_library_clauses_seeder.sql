-- +goose Up
-- +goose StatementBegin

DECLARE @version_id INT;
DECLARE @CybersecurityLeadershipGovernanceID INT;
DECLARE @CyberSecurityRiskManagementComplianceID INT;
DECLARE @CyberSecurityOperationsTechnologyID INT;
DECLARE @ThirdPartyCybersecurityID INT;

-- Get version ID 
SELECT @version_id = id FROM compliance_library_versions WHERE version = 'version 1.0' AND year = 2024;

-- DOMAIN 1:
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '3.1', N'Cyber Security Leadership and Governance', N'القيادة والحوكمة في مجال الأمن السيبراني', N'The ultimate responsibility for cybersecurity rests with the board of the Member Organization. The board of the Member Organization can delegate its cybersecurity responsibilities to a cybersecurity committee (or a senior manager from a control function). The cybersecurity committee could be responsible for defining the cybersecurity governance and setting the Member Organization’s cybersecurity strategy. The cybersecurity committee can also be responsible for defining a cybersecurity policy and ensuring the operational effectiveness of this cybersecurity policy.To develop and maintain the cybersecurity policy and to execute the cybersecurity activities across the Member Organization, an independent cybersecurity function should be established.', N'تقع المسؤولية النهائية عن الأمن السيبراني على عاتق مجلس إدارة المنظمة العضو. ويمكن لمجلس إدارة المنظمة العضو تفويض مسؤولياته المتعلقة بالأمن السيبراني إلى لجنة للأمن السيبراني (أو إلى مدير رفيع المستوى من إدارة الرقابة). ويمكن للجنة الأمن السيبراني أن تتولى مسؤولية تحديد حوكمة الأمن السيبراني ووضع استراتيجية المنظمة العضو في هذا المجال. كما يمكن للجنة الأمن السيبراني أن تتولى مسؤولية تحديد سياسة الأمن السيبراني وضمان فعاليتها التشغيلية. ولتطوير سياسة الأمن السيبراني وصيانتها، وتنفيذ أنشطة الأمن السيبراني في جميع أنحاء المنظمة العضو، ينبغي إنشاء إدارة مستقلة للأمن السيبراني.', NULL, GETDATE(), GETDATE()
);
SELECT @CybersecurityLeadershipGovernanceID = id FROM compliance_library_clauses
WHERE clause_number = '3.1' AND compliance_library_version_id = @version_id;

-- DOMAIN 2: 
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '3.2', N'Cyber Security Risk Management and Compliance', N'إدارة مخاطر الأمن السيبراني والامتثال', N'Risk management is the ongoing process of identifying, analyzing, responding and monitoring and reviewing risks. The cyber security risk management process focusses specifically on managing risks related to cyber security. In order to manage cyber security risks, Member Organizations should: 
a. identify their cyber security risks – cyber security risk identification;
b. determine the likelihood that cyber security risks will occur and the resulting impact – cyber security risk analysis;
c. determine the appropriate response to cyber security risks and select relevant controls – cyber security risk response;
d. monitor the cyber security risk treatment and review control effectiveness – cyber security risk monitoring and review.
The compliance with the cyber security controls should be subject to periodic review and audit.', N'إدارة المخاطر هي عملية مستمرة لتحديد المخاطر وتحليلها والاستجابة لها ومراقبتها ومراجعتها. تركز عملية إدارة مخاطر الأمن السيبراني بشكل خاص على إدارة المخاطر المتعلقة بالأمن السيبراني. ولإدارة مخاطر الأمن السيبراني، ينبغي على المنظمات الأعضاء:
أ. تحديد مخاطر الأمن السيبراني لديها - تحديد مخاطر الأمن السيبراني؛
ب. تحديد احتمالية حدوث مخاطر الأمن السيبراني والأثر الناتج عنها - تحليل مخاطر الأمن السيبراني؛
ج. تحديد الاستجابة المناسبة لمخاطر الأمن السيبراني واختيار الضوابط ذات الصلة - الاستجابة لمخاطر الأمن السيبراني؛
د. مراقبة معالجة مخاطر الأمن السيبراني ومراجعة فعالية الضوابط - مراقبة مخاطر الأمن السيبراني ومراجعتها.
يجب أن يخضع الامتثال لضوابط الأمن السيبراني لمراجعة وتدقيق دوريين.', NULL, GETDATE(), GETDATE()
);
SELECT @CyberSecurityRiskManagementComplianceID = id FROM compliance_library_clauses
WHERE clause_number = '3.2' AND compliance_library_version_id = @version_id;

-- DOMAIN 3:
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '3.3', N'Cyber Security Operations and Technology', N'عمليات وتكنولوجيا الأمن السيبراني', N'In order to safeguard the protection of the operations and technology of the Member Organization''s information assets and its staff, third parties and customers, the Member Organizations have to ensure that security requirements for their information assets and the supporting processes are defined, approved and implemented. The compliance with these cyber security requirements should be monitored and the effectiveness of the cyber security controls should be periodically measured and evaluated in order to identify potential revisions of the controls or measurements.', N'لضمان حماية عمليات وتقنيات أصول المعلومات الخاصة بالمنظمة العضوة وموظفيها والجهات الخارجية وعملائها، يتعين على المنظمات الأعضاء ضمان تحديد واعتماد وتنفيذ متطلبات أمن أصول المعلومات الخاصة بها والعمليات الداعمة لها. وينبغي مراقبة الامتثال لمتطلبات الأمن السيبراني هذه، وقياس وتقييم فعالية ضوابط الأمن السيبراني بشكل دوري لتحديد أي تعديلات محتملة عليها.', NULL, GETDATE(), GETDATE()
);
SELECT @CyberSecurityOperationsTechnologyID = id FROM compliance_library_clauses
WHERE clause_number = '3.3' AND compliance_library_version_id = @version_id;

-- DOMAIN 4:
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en,
    clause_title_ar, clause_description_en, clause_description_ar,
    compliance_library_clause_id, created_at, updated_at
) VALUES (
    @version_id, '3.4', N'Third Party Cyber Security', N'الأمن السيبراني للجهات الخارجية', N'When Member Organizations do rely on, or have to deal with third party services, it is key to ensure the same level of cyber security protection is implemented at the third party, as within the Member Organization. This paragraph describes how the cyber security requirements between the Member Organization and Third Parties should be organized, implemented and monitored. Third Parties in this Framework are defined as, information services providers, outsourcing providers, cloud computing providers, vendors, suppliers, governmental agencies, etc.', N'عندما تعتمد المنظمات الأعضاء على خدمات جهات خارجية، أو تضطر للتعامل معها، من الضروري ضمان تطبيق نفس مستوى حماية الأمن السيبراني لدى هذه الجهات، كما هو الحال داخل المنظمة العضو. توضح هذه الفقرة كيفية تنظيم متطلبات الأمن السيبراني بين المنظمة العضو والجهات الخارجية، وتطبيقها، ومراقبتها. تُعرّف الجهات الخارجية في هذا الإطار بأنها: مقدمو خدمات المعلومات، ومقدمو خدمات التعهيد الخارجي، ومقدمو خدمات الحوسبة السحابية، والبائعون، والموردون، والهيئات الحكومية، إلخ.', NULL, GETDATE(), GETDATE()
);
SELECT @ThirdPartyCybersecurityID = id FROM compliance_library_clauses
WHERE clause_number = '3.4' AND compliance_library_version_id = @version_id;

-- INSERT ALL SUB-CLAUSES FOR DOMAIN 1: GOVERNANCE
INSERT INTO compliance_library_clauses (
    compliance_library_version_id, clause_number, clause_title_en, clause_title_ar,
    clause_description_en, clause_description_ar, compliance_library_clause_id, created_at, updated_at
) VALUES
(@version_id, '3.1.1', N'Cyber Security Governance', N'حوكمة الأمن السيبراني', N'Principle
A cyber security governance structure should be defined and implemented, and should be endorsed by the board.
Objective
To direct and control the overall approach to cyber security within the Member Organization.', N'المبدأ
يجب تحديد هيكل حوكمة الأمن السيبراني وتنفيذه، ويجب أن يحظى بموافقة مجلس الإدارة.
الهدف
توجيه ومراقبة النهج العام للأمن السيبراني داخل المنظمة العضو.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()),
(@version_id, '3.1.2', N'Cyber Security Strategy', N'استراتيجية الأمن السيبراني', N'Principle
A cyber security strategy should be defined and aligned with the Member Organization’s strategic objectives, as well as with the Banking Sector’s cyber security strategy.
Objective To ensure that cyber security initiatives and projects within the Member Organization contribute to the
Member Organization’s strategic objectives and are aligned with the Banking Sector’s cyber security strategy.', N'المبدأ
يجب تحديد استراتيجية الأمن السيبراني ومواءمتها مع الأهداف الاستراتيجية للمنظمة العضو، وكذلك مع استراتيجية الأمن السيبراني للقطاع المصرفي.
الهدف: ضمان مساهمة مبادرات ومشاريع الأمن السيبراني داخل المنظمة العضو في تحقيق الأهداف الاستراتيجية للمنظمة العضو، ومواءمتها مع استراتيجية الأمن السيبراني للقطاع المصرفي.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()
),
(@version_id, '3.1.3', N'Cyber Security Policy', N'سياسة الأمن السيبراني', N'Principle
A cyber security policy should be defined, approved and communicated.
Objective
To document the Member Organization’s commitment and objectives of cyber security, and to communicate this to the relevant stakeholders.', N'المبدأ
يجب تحديد سياسة الأمن السيبراني، والموافقة عليها، وتعميمها.
الهدف
توثيق التزام المنظمة العضوة وأهدافها المتعلقة بالأمن السيبراني، وإيصال ذلك إلى الجهات المعنية.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()
),
(@version_id, '3.1.4', N'Cyber Security Roles and Responsibilities', N'أدوار ومسؤوليات الأمن السيبراني', N'Principle
Responsibilities to implement, maintain, support and promote cyber security should be defined throughout the Member Organization. Additionally, all parties involved in cyber security should understand and take their role and responsibilities.
Objective
To ensure that relevant stakeholders are aware of the responsibilities with regard to cyber security and apply cyber security controls throughout the Member Organization.', N'المبدأ
يجب تحديد مسؤوليات تطبيق الأمن السيبراني وصيانته ودعمه وتعزيزه في جميع أنحاء المنظمة العضو. بالإضافة إلى ذلك، ينبغي على جميع الأطراف المعنية بالأمن السيبراني فهم أدوارها ومسؤولياتها والاضطلاع بها.

الهدف
ضمان وعي الجهات المعنية بمسؤولياتها المتعلقة بالأمن السيبراني وتطبيق ضوابط الأمن السيبراني في جميع أنحاء المنظمة العضو.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()
),
(@version_id, '3.1.5', N'Cyber Security in Project Management', N'الأمن السيبراني في إدارة المشاريع', N'Principle:
Cyber security should be addressed in project management and project governance.
Objective:
To ensure that all the Member Organization’s projects meet cyber security requirements.', N'المبدأ: ينبغي مراعاة الأمن السيبراني في إدارة المشاريع وحوكمتها. الهدف: ضمان استيفاء جميع مشاريع المنظمة العضوة لمتطلبات الأمن السيبراني.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()
),
(@version_id, '3.1.6', N'Cyber Security Awareness', N'التوعية بالأمن السيبراني', N'Principle:
A cyber security awareness program should be defined and conducted for staff, third parties, and customers of the Member Organization.
Objective:
To create a cyber security risk-aware culture where the Member Organization’s staff, third parties, and customers make effective risk-based decisions that protect the Member Organization’s information.', N'المبدأ:
يجب وضع برنامج توعية بالأمن السيبراني وتنفيذه لموظفي المنظمة العضوة، والجهات الخارجية، وعملائها.
الهدف:
بناء ثقافة واعية بمخاطر الأمن السيبراني، حيث يتخذ موظفو المنظمة العضوة، والجهات الخارجية، وعملاؤها قرارات فعّالة قائمة على المخاطر، تحمي معلومات المنظمة العضوة.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()
),
(@version_id, '3.1.7', N'Cyber Security Training', N'تدريب الأمن السيبراني', N'Principle
Staff of the Member Organization should be provided with training regarding how to operate the Member Organization’s systems securely and to address and apply cyber security controls.
Objective
To ensure that staff of the Member Organization are equipped with the skills and required knowledge to protect the Member Organization’s information assets and to fulfil their cyber security responsibilities.', N'المبدأ
يجب تدريب موظفي المنظمة العضوة على كيفية تشغيل أنظمة المنظمة العضوة بأمان، وكيفية التعامل مع ضوابط الأمن السيبراني وتطبيقها.
الهدف
ضمان تزويد موظفي المنظمة العضوة بالمهارات والمعرفة اللازمة لحماية أصول معلومات المنظمة العضوة، والوفاء بمسؤولياتهم في مجال الأمن السيبراني.', @CybersecurityLeadershipGovernanceID, GETDATE(), GETDATE()),
(@version_id, '3.2.1', N'Cyber Security Risk Management', N'إدارة مخاطر الأمن السيبراني', N'Principle
A cyber security risk management process should be defined, approved and implemented, and should be aligned with the Member Organization’s enterprise risk management process.
Objective
To ensure cyber security risks are properly managed to protect the confidentiality, integrity and availability of the Member Organization’s information assets, and to ensure the cyber security risk management process is aligned with the Member Organization’s enterprise risk management process.', N'المبدأ
يجب تحديد عملية إدارة مخاطر الأمن السيبراني، والموافقة عليها، وتنفيذها، وأن تكون متوافقة مع عملية إدارة المخاطر المؤسسية للمنظمة العضو.
الهدف
ضمان إدارة مخاطر الأمن السيبراني بشكل سليم لحماية سرية وسلامة وتوافر أصول المعلومات للمنظمة العضو، وضمان توافق عملية إدارة مخاطر الأمن السيبراني مع عملية إدارة المخاطر المؤسسية للمنظمة العضو.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.1.1', N'Cyber Security Risk Identification', N'تحديد مخاطر الأمن السيبراني', N'Principle
Cyber security risk identification should be performed and should include the Member Organization’s relevant assets, threats, existing controls and vulnerabilities.
Objective
To find, recognize and describe the Member Organization’s cyber security risks.', N'المبدأ
يجب تحديد مخاطر الأمن السيبراني، ويجب أن يشمل ذلك الأصول والتهديدات والضوابط ونقاط الضعف ذات الصلة بالمنظمة العضو.
الهدف
تحديد مخاطر الأمن السيبراني التي تواجهها المنظمة العضو، والتعرف عليها، ووصفها.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.1.2', N'Cyber Security Risk Analysis', N'تحليل مخاطر الأمن السيبراني', N'Principle
A cyber security risk analysis should be conducted based on the likelihood that the identified cyber security risks will occur and their resulting impact.
Objective
To analyze and determine the nature and the level of the identified cyber security risks.', N'المبدأ
يجب إجراء تحليل مخاطر الأمن السيبراني بناءً على احتمالية حدوث مخاطر الأمن السيبراني المُحددة وأثرها الناتج.
الهدف
تحليل وتحديد طبيعة ومستوى مخاطر الأمن السيبراني المُحددة.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.1.3', N'Cyber Security Risk Response', N'الاستجابة لمخاطر الأمن السيبراني', N'Principle
The cyber security risks of a Member Organization should be treated.
Objective
To ensure cyber security risks are treated (i.e., accepted, avoided, transferred or mitigated).', N'المبدأ
يجب معالجة مخاطر الأمن السيبراني التي تواجهها أي منظمة عضو.
الهدف
ضمان معالجة مخاطر الأمن السيبراني (أي قبولها، تجنبها، نقلها، أو تخفيف آثارها).', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.1.4', N'Cyber Risk Monitoring and Review', N'مراقبة ومراجعة المخاطر السيبرانية', N'Principle
The progress cyber security risk treatment should be monitored and the effectiveness of revised or newly implemented cyber security controls should be reviewed.
Objective
To ensure that the cyber security risk treatment is performed according to the treatment plans. To ensure that the revised or newly implemented cyber security controls are effective.', N'المبدأ
يجب مراقبة تقدم معالجة مخاطر الأمن السيبراني، ومراجعة فعالية ضوابط الأمن السيبراني المُعدّلة أو المُطبّقة حديثًا.
الهدف
ضمان تنفيذ معالجة مخاطر الأمن السيبراني وفقًا لخطط المعالجة، وضمان فعالية ضوابط الأمن السيبراني المُعدّلة أو المُطبّقة حديثًا.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.2', N'Regulatory Compliance', N'الامتثال التنظيمي', N'Principle
A process should be established by the Member Organization to identify, communicate and comply with the cyber security implications of relevant regulations. 
Objective
To comply with regulations affecting cyber security of the Member Organization.', N'المبدأ

ينبغي على المنظمة العضو وضع آلية لتحديد آثار اللوائح ذات الصلة على الأمن السيبراني، والتواصل بشأنها، والامتثال لها.

الهدف

الامتثال للوائح التي تؤثر على الأمن السيبراني للمنظمة العضو.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.3', N'Compliance with (inter)national industry standards', N'الامتثال للمعايير الصناعية (الدولية)', N'Principle
The Member Organization should comply with mandatory (inter)national industry standards. 
Objective
To comply with mandatory (inter)national industry standards.', N'المبدأ
يجب على المنظمة العضوة الالتزام بمعايير الصناعة الدولية (الدولية) الإلزامية.
الهدف
الالتزام بمعايير الصناعة الدولية (الدولية) الإلزامية.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.4', N'Cyber Security Review', N'مراجعة الأمن السيبراني', N'Principle
The cyber security status of the Member Organization’s information assets should be subject to periodic cyber security review. 
Objective
To ascertain whether the cyber security controls are securely designed and implemented, and the effectiveness of these controls is being monitored.', N'المبدأ
يجب أن تخضع حالة الأمن السيبراني لأصول معلومات المنظمة العضو لمراجعة دورية.
الهدف
التأكد من تصميم وتنفيذ ضوابط الأمن السيبراني بشكل آمن، ومراقبة فعالية هذه الضوابط.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
), (
    @version_id, '3.2.5', N'Cyber Security Audits', N'عمليات تدقيق الأمن السيبراني', N'Principle
The cyber security status of the Member Organization’s information assets should be subject to thorough, independent and regular cyber security audits performed in accordance with generally accepted auditing standards and SAMA cyber security framework. 
Objective
To ascertain with reasonable assurance whether the cyber security controls are securely designed and implemented, and whether the effectiveness of these controls is being monitored.', N'المبدأ
يجب أن تخضع حالة الأمن السيبراني لأصول معلومات المنظمة العضو لعمليات تدقيق شاملة ومستقلة ومنتظمة للأمن السيبراني، تُجرى وفقًا لمعايير التدقيق المتعارف عليها وإطار عمل مؤسسة النقد العربي السعودي للأمن السيبراني.

الهدف
التأكد بدرجة معقولة من أن ضوابط الأمن السيبراني مصممة ومطبقة بشكل آمن، وما إذا كانت فعاليتها تخضع للمراقبة.', @CyberSecurityRiskManagementComplianceID, GETDATE(), GETDATE()
),
(@version_id, '3.3.1', N'Human Resources', N'الموارد البشرية', N'Principle
The Member Organization should incorporate cyber security requirements into human resources processes.
Objective
To ensure that Member Organization staff’s cyber security responsibilities are embedded in staff agreements and staff are being screened before and during their employment lifecycle.', N'المبدأ

ينبغي على المنظمة العضو دمج متطلبات الأمن السيبراني في عمليات الموارد البشرية.

الهدف

ضمان تضمين مسؤوليات الأمن السيبراني لموظفي المنظمة العضو في اتفاقيات الموظفين، وفحصهم قبل وأثناء دورة حياتهم الوظيفية.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.2', N'Physical Security', N'الأمن المادي', N'Principle
The Member Organization should ensure all facilities which host information assets are physically protected against intentional and unintentional security events.
Objective
To prevent unauthorized physical access to the Member Organization information assets and to ensure its protection.', N'المبدأ

يتعين على المنظمة العضو ضمان حماية جميع المرافق التي تستضيف أصول المعلومات ماديًا من الحوادث الأمنية المتعمدة وغير المتعمدة.

الهدف

منع الوصول المادي غير المصرح به إلى أصول معلومات المنظمة العضو وضمان حمايتها.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.3', N'Asset Management', N'إدارة الأصول', N'Principle
The Member Organization should define, approve, implement, communicate and monitor an asset management process, which supports an accurate, up-to-date and unified asset register. 
Objective
To support the Member Organization in having an accurate and up-to-date inventory and central insight in the physical / logical location and relevant details of all available information assets, in order to support its processes, such as financial, procurement, IT and cyber security processes.', N'المبدأ

ينبغي على المنظمة العضو تحديد عملية إدارة أصول، والموافقة عليها، وتنفيذها، وتعميمها، ومراقبتها، بما يدعم سجل أصول دقيقًا وحديثًا وموحدًا.

الهدف

دعم المنظمة العضو في الحصول على جرد دقيق وحديث، ورؤية مركزية للموقع الفعلي/المكاني لجميع أصول المعلومات المتاحة، والتفاصيل ذات الصلة بها، لدعم عملياتها، مثل العمليات المالية، وعمليات المشتريات، وتكنولوجيا المعلومات، والأمن السيبراني.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.4', N'Cyber Security Architecture', N'هندسة الأمن السيبراني', N'Principle
The Member Organization should define, follow and review the cyber security architecture, which outlines the cyber security requirements in the enterprise architecture and addresses the design principles for developing cyber security capabilities.
Objective
To support the Member Organization in achieving a strategic, consistent, cost effective and end-to-end cyber security architecture.', N'المبدأ

يتعين على المنظمة العضوة تحديد ومتابعة ومراجعة بنية الأمن السيبراني، التي تحدد متطلبات الأمن السيبراني في بنية المؤسسة، وتتناول مبادئ التصميم لتطوير قدرات الأمن السيبراني.

الهدف

دعم المنظمة العضوة في تحقيق بنية أمن سيبراني استراتيجية، ومتسقة، وفعالة من حيث التكلفة، وشاملة.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.5', N'Identity and Access Management', N'إدارة الهوية والوصول', N'Principle
The Member Organization should restrict access to its information assets in line with their business requirements based on the need-to-have or need-to-know principles.
Objective
To ensure that the Member Organization only provides authorized and sufficient access privileges to approved users.', N'المبدأ

ينبغي على المنظمة العضو تقييد الوصول إلى أصول معلوماتها بما يتماشى مع متطلبات أعمالها، استنادًا إلى مبدأي "الحاجة إلى المعرفة".

الهدف

ضمان أن المنظمة العضو تُقدم امتيازات الوصول المصرح بها والكافية فقط للمستخدمين المعتمدين.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.6', N'Application Security', N'أمان التطبيق', N'Principle
The Member Organization should define, approve and implement cyber security standards for application systems. The compliance with these standards should be monitored and the effectiveness of these controls should be measured and periodically evaluated.
Objective
To ensure that sufficient cyber security controls are formally documented and implemented for all applications, and that the compliance is monitored and its effectiveness is evaluated periodically within the Member Organization.', N'المبدأ

يتعين على المنظمة العضو تحديد معايير الأمن السيبراني لأنظمة التطبيقات واعتمادها وتطبيقها. ويجب مراقبة الامتثال لهذه المعايير وقياس فعالية هذه الضوابط وتقييمها دوريًا.

الهدف

ضمان توثيق ضوابط الأمن السيبراني الكافية رسميًا وتطبيقها على جميع التطبيقات، ومراقبة الامتثال وتقييم فعاليته دوريًا داخل المنظمة العضو.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.7', N'Change Management', N'إدارة التغيير', N'Principle
The Member Organization should define, approve and implement a change management process that controls all changes to information assets. The compliance with the process should be monitored and the effectiveness should be measured and periodically evaluated.
Objective
To ensure that all change in the information assets within the Member Organization follow a strict change control process.', N'المبدأ

ينبغي على المنظمة العضو تحديد واعتماد وتنفيذ عملية إدارة التغيير التي تتحكم في جميع التغييرات التي تطرأ على أصول المعلومات. ويجب مراقبة الامتثال لهذه العملية وقياس فعاليتها وتقييمها دوريًا.

الهدف

ضمان اتباع جميع التغييرات في أصول المعلومات داخل المنظمة العضو لعملية صارمة لمراقبة التغيير.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.8', N'Infrastructure Security', N'أمن البنية التحتية', N'Principle
The Member Organization should define, approve and implement cyber security standards for their infrastructure components. The compliance with these standards should be monitored and the effectiveness should be measured and periodically evaluated.
Objective
To support that all cyber security controls within the infrastructure are formally documented and the compliance is monitored and its effectiveness is evaluated periodically within the Member Organization.', N'المبدأ

يتعين على المنظمة العضوة تحديد معايير الأمن السيبراني لمكونات بنيتها التحتية واعتمادها وتطبيقها. وينبغي مراقبة الامتثال لهذه المعايير وقياس فعاليتها وتقييمها دوريًا.

الهدف

ضمان توثيق جميع ضوابط الأمن السيبراني داخل البنية التحتية رسميًا، ومراقبة الامتثال وتقييم فعاليتها دوريًا داخل المنظمة العضوة.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.9', N'Cryptography', N'التشفير', N'Principle
The use of cryptographic solutions within the Member Organizations should be defined, approved and implemented.
Objective
To ensure that access to and integrity of sensitive information is protected and the originator of communication or transactions can be confirmed.', N'المبدأ
يجب تحديد استخدام حلول التشفير داخل المنظمات الأعضاء، والموافقة عليه، وتطبيقه.
الهدف
ضمان حماية الوصول إلى المعلومات الحساسة وسلامتها، والتأكد من هوية مُنشئ الاتصالات أو المعاملات.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.10', N'Bring Your Own Device (BYOD)', N'أحضر جهازك الخاص (BYOD)', N'Principle
When the Member Organization allows the use of personal devices (e.g., smartphones, tablets, laptops) for business purposes, the use should be supported by a defined, approved and implemented cyber security standard, additional staff agreements and a cyber security awareness training. 
Objective
To ensure that business and sensitive information of the Member Organization is securely handled by staff and protected during transmission and storage, when using personal devices.', N'المبدأ
عندما تسمح المنظمة العضو باستخدام الأجهزة الشخصية (مثل الهواتف الذكية، والأجهزة اللوحية، وأجهزة الكمبيوتر المحمولة) لأغراض العمل، يجب أن يكون الاستخدام مدعومًا بمعيار أمن سيبراني محدد ومعتمد ومطبق، واتفاقيات إضافية للموظفين، وتدريب على التوعية بالأمن السيبراني.

الهدف

ضمان تعامل موظفي المنظمة العضو مع معلوماتها التجارية والحساسة بشكل آمن، وحمايتها أثناء النقل والتخزين، عند استخدام الأجهزة الشخصية.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.11', N'Secure Disposal of Information Assets', N'التخلص الآمن من أصول المعلومات', N'Principle
The information assets of the Member Organization should be securely disposed when the information assets are no longer required.
Objective
To ensure that the Member Organization’s business, customer and other sensitive information are protected from leakage or unauthorized disclosure when disposed.', N'المبدأ
يجب التخلص من أصول المعلومات الخاصة بالمنظمة العضوة بشكل آمن عند عدم الحاجة إليها.
الهدف
ضمان حماية معلومات أعمال المنظمة العضوة وعملائها وغيرها من المعلومات الحساسة من التسرب أو الإفصاح غير المصرح به عند التخلص منها.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.12', N'Payment Systems', N'أنظمة الدفع', N'Principle
The Member Organization should define, approve, implement and monitor a cyber security standard for payment systems. The effectiveness of this process should be measured and periodically evaluated.
Objective
To ensure the Member Organization safeguards the confidentiality and integrity of shared banking systems.', N'المبدأ

يتعين على المنظمة العضو تحديد معيار للأمن السيبراني لأنظمة الدفع، واعتماده، وتطبيقه، ومراقبته. وينبغي قياس فعالية هذه العملية وتقييمها دوريًا.

الهدف

ضمان حماية المنظمة العضو لسرية وسلامة الأنظمة المصرفية المشتركة.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.13', N'Electronic Banking Services', N'الخدمات المصرفية الإلكترونية', N'Principle
The Member Organization should define, approve, implement and monitor a cyber security standard for electronic banking services. The effectiveness of this standard should be measured and periodically evaluated.
Objective
To ensure the Member Organization safeguards the confidentiality and integrity of the customer information and transactions.', N'المبدأ
على المؤسسة العضو تحديد معيار للأمن السيبراني للخدمات المصرفية الإلكترونية، واعتماده، وتطبيقه، ومراقبته. ويجب قياس فعالية هذا المعيار وتقييمه دوريًا.

الهدف
ضمان حماية المؤسسة العضو لسرية وسلامة معلومات ومعاملات العملاء.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.14', N'Cyber Security Event Management', N'إدارة أحداث الأمن السيبراني', N'Principle
The Member Organization should define, approve and implement a security event management process
to analyze operational and security loggings and respond to security events. The effectiveness of this
process should be measured and periodically evaluated.
Objective
To ensure timely identification and response to anomalies or suspicious events within regard to
information assets.', N'المبدأ

ينبغي على المنظمة العضو تحديد عملية لإدارة الأحداث الأمنية، واعتمادها، وتطبيقها، وذلك لتحليل سجلات التشغيل والأمن، والاستجابة للأحداث الأمنية. وينبغي قياس فعالية هذه العملية وتقييمها دوريًا.

الهدف

ضمان تحديد الشذوذ أو الأحداث المشبوهة المتعلقة بأصول المعلومات والاستجابة لها في الوقت المناسب.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.15', N'Cyber Security Incident Management', N'إدارة حوادث الأمن السيبراني', N'Principle
The Member Organization should define, approve and implement a cyber security incident management
that is aligned with the enterprise incident management process, to identify, respond to and recover from
cyber security incidents. The effectiveness of this process should be measured and periodically evaluated.
Objective
To ensure timely identification and handling of cyber security incidents in order to reduce the (potential)
business impact for the Member Organization.', N'المبدأ

ينبغي على المنظمة العضو تحديد واعتماد وتنفيذ نظام لإدارة حوادث الأمن السيبراني، بما يتماشى مع عملية إدارة حوادث المؤسسة، وذلك لتحديد حوادث الأمن السيبراني والاستجابة لها والتعافي منها. وينبغي قياس فعالية هذه العملية وتقييمها دوريًا.

الهدف

ضمان تحديد حوادث الأمن السيبراني ومعالجتها في الوقت المناسب للحد من التأثير (المحتمل) على أعمال المنظمة العضو.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.16', N'Threat Management', N'إدارة التهديدات', N'Principle
The Member Organization should define, approve and implement a threat intelligence management
process to identify, assess and understand threats to the Member Organization information assets, using multiple reliable sources. The effectiveness of this process should be measured and periodically
evaluated.
Objective
To obtain an adequate understanding of the Member Organization’s emerging threat posture.', N'المبدأ

ينبغي على المنظمة العضو تحديد واعتماد وتنفيذ عملية لإدارة معلومات التهديدات لتحديد وتقييم وفهم التهديدات التي تواجه أصول معلومات المنظمة العضو، وذلك باستخدام مصادر موثوقة متعددة. وينبغي قياس فعالية هذه العملية وتقييمها دوريًا.

الهدف

الحصول على فهم كافٍ لوضع التهديدات الناشئة التي تواجهها المنظمة العضو.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
), (
    @version_id, '3.3.17', N'Vulnerability Management', N'إدارة الثغرات الأمنية', N'Principle
The Member Organization should define, approve and implement a vulnerability management process for the identification and mitigation of application and infrastructural vulnerabilities. The effectiveness of this process should be measured and the effectiveness should be periodically evaluated.
Objective
To ensure timely identification and effective mitigation of application and infrastructure vulnerabilities in order to reduce the likelihood and business impact for the Member Organization.', N'المبدأ

ينبغي على المنظمة العضو تحديد واعتماد وتنفيذ عملية لإدارة الثغرات الأمنية لتحديد ثغرات التطبيقات والبنية التحتية والتخفيف من حدتها. وينبغي قياس فعالية هذه العملية وتقييمها دوريًا.

الهدف

ضمان التحديد الفوري لثغرات التطبيقات والبنية التحتية والتخفيف الفعال من حدتها، وذلك لتقليل احتمالية حدوثها وتأثيرها على أعمال المنظمة العضو.', @CyberSecurityOperationsTechnologyID, GETDATE(), GETDATE()
),
(@version_id, '3.4.1', N'Contract and Vendor Management', N'إدارة العقود والموردين', N'Principle
The Member Organization should define, approve, implement and monitor the required cyber security controls within the contract and vendor management processes.
Objective
To ensure that the Member Organization’s approved cyber security requirements are appropriately addressed before signing the contract, and the compliance with the cyber security requirements is being monitored and evaluated during the contract life-cycle.', N'المبدأ

يتعين على المنظمة العضو تحديد ضوابط الأمن السيبراني المطلوبة، والموافقة عليها، وتنفيذها، ومراقبتها ضمن عمليات إدارة العقود والموردين.

الهدف

ضمان تلبية متطلبات الأمن السيبراني المعتمدة من المنظمة العضو بشكل مناسب قبل توقيع العقد، ومراقبة وتقييم الامتثال لمتطلبات الأمن السيبراني طوال مدة العقد.', @ThirdPartyCybersecurityID, GETDATE(), GETDATE()
), (
    @version_id, '3.4.2', N'Outsourcing', N'الاستعانة بمصادر خارجية', N'Principle
The Member Organization should define, implement and monitor the required cyber security controls within outsourcing policy and outsourcing process. The effectiveness of the defined cyber security controls should periodically be measured and evaluated.
Objective
To ensure that the Member Organization’s cyber security requirements are appropriately addressed before, during and while exiting outsourcing contracts.', N'المبدأ

يتعين على المنظمة العضو تحديد ضوابط الأمن السيبراني اللازمة وتنفيذها ومراقبتها ضمن سياسة التعهيد الخارجي وعملية التعهيد. وينبغي قياس فعالية ضوابط الأمن السيبراني المحددة وتقييمها دوريًا.

الهدف

ضمان تلبية متطلبات الأمن السيبراني للمنظمة العضو بشكل مناسب قبل وأثناء وبعد إنهاء عقود التعهيد الخارجي.', @ThirdPartyCybersecurityID, GETDATE(), GETDATE()
), (
    @version_id, '3.4.3', N'Cloud Computing', N'الحوسبة السحابية', N'Principle
The Member Organization should define, implement and monitor the required cyber security controls within the cloud computing policy and process for hybrid and public cloud services. The effectiveness of the defined cyber security controls should periodically be measured and evaluated. Please note that this requirement is not applicable to private cloud services (= internal cloud).
Objective
To ensure that all functions and staff within the Member Organization are aware of the agreed direction and position on hybrid and public cloud services, the required process to apply for hybrid and public cloud services, the risk appetite on hybrid and public cloud services and the specific cyber security requirements for hybrid and public cloud services.', N'المبدأ

يتعين على المنظمة العضو تحديد ضوابط الأمن السيبراني المطلوبة ضمن سياسة وإجراءات الحوسبة السحابية لخدمات السحابة الهجينة والعامة، وتطبيقها، ومراقبتها. وينبغي قياس وتقييم فعالية ضوابط الأمن السيبراني المحددة بشكل دوري. يُرجى العلم أن هذا الشرط لا ينطبق على خدمات السحابة الخاصة (= السحابة الداخلية).

الهدف

ضمان إلمام جميع الوظائف والموظفين داخل المنظمة العضو بالتوجه والموقف المتفق عليه بشأن خدمات السحابة الهجينة والعامة، والإجراءات اللازمة للتقدم بطلبات الحصول على خدمات السحابة الهجينة والعامة، ومستوى المخاطرة المطلوب لخدمات السحابة الهجينة والعامة، ومتطلبات الأمن السيبراني المحددة لها.', @ThirdPartyCybersecurityID, GETDATE(), GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_library_clauses;
-- +goose StatementEnd