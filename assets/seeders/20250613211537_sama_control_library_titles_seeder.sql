-- +goose Up
-- +goose StatementBegin
DECLARE @compliance_library_id INT;
SELECT @compliance_library_id = id FROM compliance_libraries
WHERE title_en = 'Cyber Security Framework | Saudi Arabian Monetary Authority';

DECLARE @cl_1 INT;
SELECT @cl_1 = id FROM control_library WHERE title_en = N'Define and Approve Cybersecurity Strategy';

DECLARE @cl_2 INT;
SELECT @cl_2 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Strategy Roadmap';

DECLARE @cl_3 INT;
SELECT @cl_3 = id FROM control_library WHERE title_en = N'Review Cybersecurity Strategy Periodically';

DECLARE @cl_4 INT;
SELECT @cl_4 = id FROM control_library WHERE title_en = N'Establish Independent Cybersecurity Function';

DECLARE @cl_5 INT;
SELECT @cl_5 = id FROM control_library WHERE title_en = N'Employ Qualified Cybersecurity Professionals';

DECLARE @cl_6 INT;
SELECT @cl_6 = id FROM control_library WHERE title_en = N'Establish Cybersecurity Steering Committee';

DECLARE @cl_7 INT;
SELECT @cl_7 = id FROM control_library WHERE title_en = N'Define and Document Cybersecurity Policies';

DECLARE @cl_8 INT;
SELECT @cl_8 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Policies';

DECLARE @cl_9 INT;
SELECT @cl_9 = id FROM control_library WHERE title_en = N'Support Policies with Technical Standards';

DECLARE @cl_10 INT;
SELECT @cl_10 = id FROM control_library WHERE title_en = N'Review Cybersecurity Policies Periodically';

DECLARE @cl_11 INT;
SELECT @cl_11 = id FROM control_library WHERE title_en = N'Define Cybersecurity Organizational Structure';

DECLARE @cl_12 INT;
SELECT @cl_12 = id FROM control_library WHERE title_en = N'Review Cybersecurity Roles Periodically';

DECLARE @cl_13 INT;
SELECT @cl_13 = id FROM control_library WHERE title_en = N'Define Cybersecurity Risk Management Procedures';

DECLARE @cl_14 INT;
SELECT @cl_14 = id FROM control_library WHERE title_en = N'Implement Risk Management Procedures';

DECLARE @cl_15 INT;
SELECT @cl_15 = id FROM control_library WHERE title_en = N'Conduct Risk Assessments for Key Scenarios';

DECLARE @cl_16 INT;
SELECT @cl_16 = id FROM control_library WHERE title_en = N'Assess Risks in Technology Projects';

DECLARE @cl_17 INT;
SELECT @cl_17 = id FROM control_library WHERE title_en = N'Assess Risks Before Infrastructure Changes';

DECLARE @cl_18 INT;
SELECT @cl_18 = id FROM control_library WHERE title_en = N'Assess Risks for Third-Party Services';

DECLARE @cl_19 INT;
SELECT @cl_19 = id FROM control_library WHERE title_en = N'Assess Risks for New Services';

DECLARE @cl_20 INT;
SELECT @cl_20 = id FROM control_library WHERE title_en = N'Review Risk Management Procedures Periodically';

DECLARE @cl_21 INT;
SELECT @cl_21 = id FROM control_library WHERE title_en = N'Include Cybersecurity in Project Management';

DECLARE @cl_22 INT;
SELECT @cl_22 = id FROM control_library WHERE title_en = N'Implement Cybersecurity in Change Management';

DECLARE @cl_23 INT;
SELECT @cl_23 = id FROM control_library WHERE title_en = N'Conduct Vulnerability Assessments';

DECLARE @cl_24 INT;
SELECT @cl_24 = id FROM control_library WHERE title_en = N'Ensure Secure Configuration and Patching';

DECLARE @cl_25 INT;
SELECT @cl_25 = id FROM control_library WHERE title_en = N'Secure Software Development Processes';

DECLARE @cl_26 INT;
SELECT @cl_26 = id FROM control_library WHERE title_en = N'Use Secure Coding Standards';

DECLARE @cl_27 INT;
SELECT @cl_27 = id FROM control_library WHERE title_en = N'Use Trusted Software Development Tools';

DECLARE @cl_28 INT;
SELECT @cl_28 = id FROM control_library WHERE title_en = N'Test Software Compliance';

DECLARE @cl_29 INT;
SELECT @cl_29 = id FROM control_library WHERE title_en = N'Ensure Secure Software Integration';

DECLARE @cl_30 INT;
SELECT @cl_30 = id FROM control_library WHERE title_en = N'Secure Software Configuration Before Deployment';

DECLARE @cl_31 INT;
SELECT @cl_31 = id FROM control_library WHERE title_en = N'Review Project Management Cybersecurity Requirements';

DECLARE @cl_32 INT;
SELECT @cl_32 = id FROM control_library WHERE title_en = N'Comply with Cybersecurity Laws';

DECLARE @cl_33 INT;
SELECT @cl_33 = id FROM control_library WHERE title_en = N'Adhere to International Cybersecurity Agreements';

DECLARE @cl_34 INT;
SELECT @cl_34 = id FROM control_library WHERE title_en = N'Conduct Periodic Cybersecurity Reviews';

DECLARE @cl_35 INT;
SELECT @cl_35 = id FROM control_library WHERE title_en = N'Perform Independent Cybersecurity Audits';

DECLARE @cl_36 INT;
SELECT @cl_36 = id FROM control_library WHERE title_en = N'Document and Report Audit Results';

DECLARE @cl_37 INT;
SELECT @cl_37 = id FROM control_library WHERE title_en = N'Define Personnel Cybersecurity Requirements';

DECLARE @cl_38 INT;
SELECT @cl_38 = id FROM control_library WHERE title_en = N'Implement Personnel Cybersecurity Requirements';

DECLARE @cl_39 INT;
SELECT @cl_39 = id FROM control_library WHERE title_en = N'Establish Pre-Employment Cybersecurity Requirements';

DECLARE @cl_40 INT;
SELECT @cl_40 = id FROM control_library WHERE title_en = N'Include Cybersecurity in Employment Contracts';

DECLARE @cl_41 INT;
SELECT @cl_41 = id FROM control_library WHERE title_en = N'Screen Candidates for Cybersecurity Roles';

DECLARE @cl_42 INT;
SELECT @cl_42 = id FROM control_library WHERE title_en = N'Enforce Cybersecurity During Employment';

DECLARE @cl_43 INT;
SELECT @cl_43 = id FROM control_library WHERE title_en = N'Provide Cybersecurity Awareness Training';

DECLARE @cl_44 INT;
SELECT @cl_44 = id FROM control_library WHERE title_en = N'Ensure Compliance with Cybersecurity Policies';

DECLARE @cl_45 INT;
SELECT @cl_45 = id FROM control_library WHERE title_en = N'Revoke Access Upon Termination';

DECLARE @cl_46 INT;
SELECT @cl_46 = id FROM control_library WHERE title_en = N'Review Personnel Cybersecurity Requirements Periodically';

DECLARE @cl_47 INT;
SELECT @cl_47 = id FROM control_library WHERE title_en = N'Develop Cybersecurity Awareness Program';

DECLARE @cl_48 INT;
SELECT @cl_48 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Awareness Program';

DECLARE @cl_49 INT;
SELECT @cl_49 = id FROM control_library WHERE title_en = N'Cover Key Cybersecurity Threats in Awareness Program';

DECLARE @cl_50 INT;
SELECT @cl_50 = id FROM control_library WHERE title_en = N'Train on Secure Email Handling';

DECLARE @cl_51 INT;
SELECT @cl_51 = id FROM control_library WHERE title_en = N'Train on Secure Mobile Device Use';

DECLARE @cl_52 INT;
SELECT @cl_52 = id FROM control_library WHERE title_en = N'Train on Secure Internet Browsing';

DECLARE @cl_53 INT;
SELECT @cl_53 = id FROM control_library WHERE title_en = N'Train on Secure Social Media Use';

DECLARE @cl_54 INT;
SELECT @cl_54 = id FROM control_library WHERE title_en = N'Provide Specialized Cybersecurity Training';

DECLARE @cl_55 INT;
SELECT @cl_55 = id FROM control_library WHERE title_en = N'Train Cybersecurity Personnel';

DECLARE @cl_56 INT;
SELECT @cl_56 = id FROM control_library WHERE title_en = N'Train Software Development Personnel';

DECLARE @cl_57 INT;
SELECT @cl_57 = id FROM control_library WHERE title_en = N'Train Executive Personnel';

DECLARE @cl_58 INT;
SELECT @cl_58 = id FROM control_library WHERE title_en = N'Review Cybersecurity Awareness Program Periodically';

DECLARE @cl_59 INT;
SELECT @cl_59 = id FROM control_library WHERE title_en = N'Define Asset Management Cybersecurity Requirements';

DECLARE @cl_60 INT;
SELECT @cl_60 = id FROM control_library WHERE title_en = N'Implement Asset Management Cybersecurity Requirements';

DECLARE @cl_61 INT;
SELECT @cl_61 = id FROM control_library WHERE title_en = N'Define Acceptable Use Policy for Assets';

DECLARE @cl_62 INT;
SELECT @cl_62 = id FROM control_library WHERE title_en = N'Implement Acceptable Use Policy for Assets';

DECLARE @cl_63 INT;
SELECT @cl_63 = id FROM control_library WHERE title_en = N'Classify and Handle Assets per Regulations';

DECLARE @cl_64 INT;
SELECT @cl_64 = id FROM control_library WHERE title_en = N'Review Asset Management Requirements Periodically';

DECLARE @cl_65 INT;
SELECT @cl_65 = id FROM control_library WHERE title_en = N'Define Identity and Access Management Requirements';

DECLARE @cl_66 INT;
SELECT @cl_66 = id FROM control_library WHERE title_en = N'Implement Identity and Access Management Requirements';

DECLARE @cl_67 INT;
SELECT @cl_67 = id FROM control_library WHERE title_en = N'Establish Identity and Access Management Controls';

DECLARE @cl_68 INT;
SELECT @cl_68 = id FROM control_library WHERE title_en = N'Implement Single-Factor Authentication';

DECLARE @cl_69 INT;
SELECT @cl_69 = id FROM control_library WHERE title_en = N'Implement Multi-Factor Authentication';

DECLARE @cl_70 INT;
SELECT @cl_70 = id FROM control_library WHERE title_en = N'Enforce Access Control Principles';

DECLARE @cl_71 INT;
SELECT @cl_71 = id FROM control_library WHERE title_en = N'Manage Privileged Access';

DECLARE @cl_72 INT;
SELECT @cl_72 = id FROM control_library WHERE title_en = N'Review User Access Rights Periodically';

DECLARE @cl_73 INT;
SELECT @cl_73 = id FROM control_library WHERE title_en = N'Review Identity and Access Management Periodically';

DECLARE @cl_74 INT;
SELECT @cl_74 = id FROM control_library WHERE title_en = N'Define Information Systems Protection Requirements';

DECLARE @cl_75 INT;
SELECT @cl_75 = id FROM control_library WHERE title_en = N'Implement Information Systems Protection Requirements';

DECLARE @cl_76 INT;
SELECT @cl_76 = id FROM control_library WHERE title_en = N'Establish Information Systems Protection Controls';

DECLARE @cl_77 INT;
SELECT @cl_77 = id FROM control_library WHERE title_en = N'Manage Malware Protection';

DECLARE @cl_78 INT;
SELECT @cl_78 = id FROM control_library WHERE title_en = N'Restrict External Storage Media';

DECLARE @cl_79 INT;
SELECT @cl_79 = id FROM control_library WHERE title_en = N'Manage Software Patch Deployment';

DECLARE @cl_80 INT;
SELECT @cl_80 = id FROM control_library WHERE title_en = N'Synchronize System Clocks';

DECLARE @cl_81 INT;
SELECT @cl_81 = id FROM control_library WHERE title_en = N'Review Information Systems Protection Periodically';

DECLARE @cl_82 INT;
SELECT @cl_82 = id FROM control_library WHERE title_en = N'Define Email Security Requirements';

DECLARE @cl_83 INT;
SELECT @cl_83 = id FROM control_library WHERE title_en = N'Implement Email Security Requirements';

DECLARE @cl_84 INT;
SELECT @cl_84 = id FROM control_library WHERE title_en = N'Establish Email Security Controls';

DECLARE @cl_85 INT;
SELECT @cl_85 = id FROM control_library WHERE title_en = N'Filter Phishing and Spam Emails';

DECLARE @cl_86 INT;
SELECT @cl_86 = id FROM control_library WHERE title_en = N'Secure Email with Multi-Factor Authentication';

DECLARE @cl_87 INT;
SELECT @cl_87 = id FROM control_library WHERE title_en = N'Archive and Backup Emails';

DECLARE @cl_88 INT;
SELECT @cl_88 = id FROM control_library WHERE title_en = N'Protect Against Advanced Threats';

DECLARE @cl_89 INT;
SELECT @cl_89 = id FROM control_library WHERE title_en = N'Secure Email Domains';

DECLARE @cl_90 INT;
SELECT @cl_90 = id FROM control_library WHERE title_en = N'Review Email Security Requirements Periodically';

DECLARE @cl_91 INT;
SELECT @cl_91 = id FROM control_library WHERE title_en = N'Define Network Security Requirements';

DECLARE @cl_92 INT;
SELECT @cl_92 = id FROM control_library WHERE title_en = N'Implement Network Security Requirements';

DECLARE @cl_93 INT;
SELECT @cl_93 = id FROM control_library WHERE title_en = N'Establish Network Security Controls';

DECLARE @cl_94 INT;
SELECT @cl_94 = id FROM control_library WHERE title_en = N'Segment Network Infrastructure';

DECLARE @cl_95 INT;
SELECT @cl_95 = id FROM control_library WHERE title_en = N'Segregate Network Environments';

DECLARE @cl_96 INT;
SELECT @cl_96 = id FROM control_library WHERE title_en = N'Secure Internet Connectivity';

DECLARE @cl_97 INT;
SELECT @cl_97 = id FROM control_library WHERE title_en = N'Protect Wireless Networks';

DECLARE @cl_98 INT;
SELECT @cl_98 = id FROM control_library WHERE title_en = N'Manage Network Services and Ports';

DECLARE @cl_99 INT;
SELECT @cl_99 = id FROM control_library WHERE title_en = N'Deploy Intrusion Prevention Systems';

DECLARE @cl_100 INT;
SELECT @cl_100 = id FROM control_library WHERE title_en = N'Secure Domain Name Service (DNS)';

DECLARE @cl_101 INT;
SELECT @cl_101 = id FROM control_library WHERE title_en = N'Protect Browsing Against Advanced Threats';

DECLARE @cl_102 INT;
SELECT @cl_102 = id FROM control_library WHERE title_en = N'Protect Against DDoS Attacks';

DECLARE @cl_103 INT;
SELECT @cl_103 = id FROM control_library WHERE title_en = N'Review Network Security Requirements Periodically';

DECLARE @cl_104 INT;
SELECT @cl_104 = id FROM control_library WHERE title_en = N'Define Mobile Device Security Requirements';

DECLARE @cl_105 INT;
SELECT @cl_105 = id FROM control_library WHERE title_en = N'Implement Mobile Device Security Requirements';

DECLARE @cl_106 INT;
SELECT @cl_106 = id FROM control_library WHERE title_en = N'Establish Mobile Device Security Controls';

DECLARE @cl_107 INT;
SELECT @cl_107 = id FROM control_library WHERE title_en = N'Secure Mobile Device Data';

DECLARE @cl_108 INT;
SELECT @cl_108 = id FROM control_library WHERE title_en = N'Restrict Mobile Device Use';

DECLARE @cl_109 INT;
SELECT @cl_109 = id FROM control_library WHERE title_en = N'Secure Data Wiping on Mobile Devices';

DECLARE @cl_110 INT;
SELECT @cl_110 = id FROM control_library WHERE title_en = N'Promote Mobile Device Security Awareness';

DECLARE @cl_111 INT;
SELECT @cl_111 = id FROM control_library WHERE title_en = N'Review Mobile Device Security Requirements Periodically';

DECLARE @cl_112 INT;
SELECT @cl_112 = id FROM control_library WHERE title_en = N'Define Data Protection Requirements';

DECLARE @cl_113 INT;
SELECT @cl_113 = id FROM control_library WHERE title_en = N'Implement Data Protection Requirements';

DECLARE @cl_114 INT;
SELECT @cl_114 = id FROM control_library WHERE title_en = N'Comply with Data Protection Standards';

DECLARE @cl_115 INT;
SELECT @cl_115 = id FROM control_library WHERE title_en = N'Review Data Protection Requirements Periodically';

DECLARE @cl_116 INT;
SELECT @cl_116 = id FROM control_library WHERE title_en = N'Define Cryptography Requirements';

DECLARE @cl_117 INT;
SELECT @cl_117 = id FROM control_library WHERE title_en = N'Implement Cryptography Requirements';

DECLARE @cl_118 INT;
SELECT @cl_118 = id FROM control_library WHERE title_en = N'Establish Cryptography Controls';

DECLARE @cl_119 INT;
SELECT @cl_119 = id FROM control_library WHERE title_en = N'Use Approved Cryptographic Systems';

DECLARE @cl_120 INT;
SELECT @cl_120 = id FROM control_library WHERE title_en = N'Manage Cryptographic Keys';

DECLARE @cl_121 INT;
SELECT @cl_121 = id FROM control_library WHERE title_en = N'Encrypt Data Based on Classification';

DECLARE @cl_122 INT;
SELECT @cl_122 = id FROM control_library WHERE title_en = N'Review Cryptography Requirements Periodically';

DECLARE @cl_123 INT;
SELECT @cl_123 = id FROM control_library WHERE title_en = N'Define Backup and Recovery Requirements';

DECLARE @cl_124 INT;
SELECT @cl_124 = id FROM control_library WHERE title_en = N'Implement Backup and Recovery Requirements';

DECLARE @cl_125 INT;
SELECT @cl_125 = id FROM control_library WHERE title_en = N'Establish Backup and Recovery Controls';

DECLARE @cl_126 INT;
SELECT @cl_126 = id FROM control_library WHERE title_en = N'Define Backup Scope for Critical Assets';

DECLARE @cl_127 INT;
SELECT @cl_127 = id FROM control_library WHERE title_en = N'Enable Quick System Recovery';

DECLARE @cl_128 INT;
SELECT @cl_128 = id FROM control_library WHERE title_en = N'Test Backups Periodically';

DECLARE @cl_129 INT;
SELECT @cl_129 = id FROM control_library WHERE title_en = N'Review Backup and Recovery Requirements Periodically';

DECLARE @cl_130 INT;
SELECT @cl_130 = id FROM control_library WHERE title_en = N'Define Vulnerability Management Requirements';

DECLARE @cl_131 INT;
SELECT @cl_131 = id FROM control_library WHERE title_en = N'Implement Vulnerability Management Requirements';

DECLARE @cl_132 INT;
SELECT @cl_132 = id FROM control_library WHERE title_en = N'Establish Vulnerability Management Controls';

DECLARE @cl_133 INT;
SELECT @cl_133 = id FROM control_library WHERE title_en = N'Conduct Periodic Vulnerability Assessments';

DECLARE @cl_134 INT;
SELECT @cl_134 = id FROM control_library WHERE title_en = N'Classify Vulnerabilities by Criticality';

DECLARE @cl_135 INT;
SELECT @cl_135 = id FROM control_library WHERE title_en = N'Remediate Vulnerabilities Based on Risk';

DECLARE @cl_136 INT;
SELECT @cl_136 = id FROM control_library WHERE title_en = N'Manage Security Patches';

DECLARE @cl_137 INT;
SELECT @cl_137 = id FROM control_library WHERE title_en = N'Subscribe to Vulnerability Information Sources';

DECLARE @cl_138 INT;
SELECT @cl_138 = id FROM control_library WHERE title_en = N'Review Vulnerability Management Requirements Periodically';

DECLARE @cl_139 INT;
SELECT @cl_139 = id FROM control_library WHERE title_en = N'Define Penetration Testing Requirements';

DECLARE @cl_140 INT;
SELECT @cl_140 = id FROM control_library WHERE title_en = N'Implement Penetration Testing Requirements';

DECLARE @cl_141 INT;
SELECT @cl_141 = id FROM control_library WHERE title_en = N'Establish Penetration Testing Controls';

DECLARE @cl_142 INT;
SELECT @cl_142 = id FROM control_library WHERE title_en = N'Define Penetration Testing Scope';

DECLARE @cl_143 INT;
SELECT @cl_143 = id FROM control_library WHERE title_en = N'Conduct Periodic Penetration Tests';

DECLARE @cl_144 INT;
SELECT @cl_144 = id FROM control_library WHERE title_en = N'Define Event Logging and Monitoring Requirements';

DECLARE @cl_145 INT;
SELECT @cl_145 = id FROM control_library WHERE title_en = N'Implement Event Logging and Monitoring Requirements';

DECLARE @cl_146 INT;
SELECT @cl_146 = id FROM control_library WHERE title_en = N'Establish Event Logging and Monitoring Controls';

DECLARE @cl_147 INT;
SELECT @cl_147 = id FROM control_library WHERE title_en = N'Enable Event Logging for Critical Assets';

DECLARE @cl_148 INT;
SELECT @cl_148 = id FROM control_library WHERE title_en = N'Log Remote and Privileged Access'; 



INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_1,
  2,
  N'3.1.1.1',
  N'A cyber security committee should be established and be mandated by the board.',
  N'ينبغي إنشاء لجنة للأمن السيبراني وتفويضها من قبل مجلس الإدارة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_2,
  2,
  N'3.1.1.2',
  N'The cyber security committee should be headed by an independent senior manager from a control function.',
  N'ينبغي أن يرأس لجنة الأمن السيبراني مدير كبير مستقل من وظيفة الرقابة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_3,
  2,
  N'3.1.1.3',
  N'The following positions should be represented in the cyber security committee:
a. senior managers from all relevant departments (e.g., COO, CIO, compliance officer, heads of
relevant business departments);
b. Chief information security officer (CISO);
c. Internal audit may attend as an observer. ',
  N'ينبغي تمثيل المناصب التالية في لجنة الأمن السيبراني:
أ. كبار المديرين من جميع الإدارات ذات الصلة (مثل: مدير العمليات، مدير المعلومات، مسؤول الامتثال، رؤساء أقسام الأعمال ذات الصلة)؛
ب. كبير مسؤولي أمن المعلومات؛
ج. يجوز لقسم التدقيق الداخلي الحضور بصفة مراقب.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_4,
  2,
  N'3.1.1.4',
  N'A cyber security committee charter should be developed, approved and reflect:
a. committee objectives;
b. roles and responsibilities;
c. minimum number of meeting participants;
d. meeting frequency (minimum on quarterly basis).',
  N'ينبغي وضع ميثاق لجنة الأمن السيبراني والموافقة عليه، وأن يعكس ما يلي:
أ. أهداف اللجنة؛
ب. الأدوار والمسؤوليات؛
ج. الحد الأدنى لعدد المشاركين في الاجتماعات؛
د. وتيرة الاجتماعات (على الأقل ربع سنوي).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_5,
  2,
  N'3.1.1.5',
  N' A cyber security function should be established. ',
  N'ينبغي إنشاء وظيفة للأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_6,
  2,
  N'3.1.1.6',
  N'The cyber security function should be independent from the information technology function. To
avoid any conflict of interest, the cyber security function and information technology function should
have separate reporting lines, budgets and staff evaluations. ',
  N'ينبغي أن تكون وظيفة الأمن السيبراني مستقلة عن وظيفة تكنولوجيا المعلومات. لتجنب أي تضارب في المصالح، ينبغي أن يكون لكلٍّ من وظيفة الأمن السيبراني ووظيفة تكنولوجيا المعلومات خطوط إبلاغ وميزانيات وتقييمات منفصلة للموظفين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_7,
  2,
  N'3.1.1.7',
  N'The cyber security function should report directly to the CEO/managing director of the Member
Organization or general manager of a control function.',
  N'يجب أن ترفع وظيفة الأمن السيبراني تقاريرها مباشرة إلى الرئيس التنفيذي/المدير الإداري للمنظمة العضو أو المدير العام لوظيفة الرقابة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_8,
  2,
  N'3.1.1.8',
  N'A full-time senior manager for the cyber security function, referred to as CISO, should be appointed
at senior management level. ',
  N'يجب تعيين مدير أول متفرغ لوظيفة الأمن السيبراني، يشار إليه باسم CISO، على مستوى الإدارة العليا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_9,
  2,
  N'3.1.1.9',
  N'The Member Organization should :
a. ensure the CISO has a Saudi nationality;
b. ensure the CISO is sufficiently qualified;
c. obtain no objection from SAMA to assign the CISO.',
  N'ينبغي على المنظمة العضوة:
أ. التأكد من أن مسؤول أمن المعلومات يحمل الجنسية السعودية؛
ب. التأكد من أن مسؤول أمن المعلومات مؤهل بشكل كافٍ؛
ج. الحصول على عدم ممانعة من مؤسسة النقد العربي السعودي لتعيين مسؤول أمن المعلومات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_10,
  2,
  N'3.1.1.10',
  N'The board of the Member Organization should allocate sufficient budget to execute the required
cyber security activities. ',
  N'ينبغي لمجلس إدارة المنظمة العضوة تخصيص ميزانية كافية لتنفيذ أنشطة الأمن السيبراني المطلوبة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_11,
  2,
  N'3.1.2 .1',
  N'The cyber security strategy should be defined, approved, maintained and executed. ',
  N'ينبغي تحديد استراتيجية الأمن السيبراني والموافقة عليها وصيانتها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_12,
  2,
  N'3.1.2 .2',
  N'The cyber security strategy should be aligned with:
a. the Member Organization’s overall objectives;
b. the legal and regulatory compliance requirements of the Member Organization;
c. the Banking Sector’s cyber security strategy.',
  N'ينبغي أن تتوافق استراتيجية الأمن السيبراني مع:
أ. الأهداف العامة للمنظمة العضو؛
ب. متطلبات الامتثال القانوني والتنظيمي للمنظمة العضو؛
ج. استراتيجية الأمن السيبراني للقطاع المصرفي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_13,
  2,
  N'3.1.2 .3',
  N'The cyber security strategy should address:
a. the importance and benefits of cyber security for the Member Organization;
b. the anticipated future state of cyber security for the Member Organization to become and remain
resilient to (emerging) cyber security threats;
c. which and when cyber security initiatives and projects should be executed to achieve the
anticipated future state.',
  N'ينبغي أن تتناول استراتيجية الأمن السيبراني ما يلي:
أ. أهمية وفوائد الأمن السيبراني للمنظمة العضو؛
ب. الحالة المستقبلية المتوقعة للأمن السيبراني للمنظمة العضو لتصبح وتظل قادرة على الصمود في وجه التهديدات الأمنية السيبرانية (الناشئة)؛
ج. أي مبادرات ومشاريع الأمن السيبراني ينبغي تنفيذها ومتى ينبغي تنفيذها لتحقيق الحالة المستقبلية المتوقعة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_14,
  2,
  N'3.1.3.1',
  N'The cyber security policy should be defined, approved and communicated.',
  N'ينبغي تحديد سياسة الأمن السيبراني والموافقة عليها وتوصيلها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_15,
  2,
  N'3.1.3.2',
  N'The cyber security policy should be reviewed periodically according to a predefined and structured
review process. ',
  N'ينبغي مراجعة سياسة الأمن السيبراني بشكل دوري وفقًا لعملية مراجعة محددة ومنظمة مسبقًا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_16,
  2,
  N'3.1.3.3',
  N'The cyber security policy should be:
a. considered as input for other corporate policies of the Member Organization (e.g., HR policy,
finance policy and IT policy);
b. supported by detailed security standards (e.g., password standard, firewall standard) and
procedures;
c. based on best practices and (inter)national standards;
d. communicated to relevant stakeholders. ',
  N'ينبغي أن تكون سياسة الأمن السيبراني:
أ. مُدخلات للسياسات المؤسسية الأخرى للمنظمة العضو (مثل سياسة الموارد البشرية، وسياسة المالية، وسياسة تكنولوجيا المعلومات)؛
ب. مدعومة بمعايير وإجراءات أمنية مفصلة (مثل معيار كلمة المرور، ومعيار جدار الحماية)؛
ج. مبنية على أفضل الممارسات والمعايير (الدولية)؛
د. مُبلغة إلى أصحاب المصلحة المعنيين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_17,
  2,
  N'3.1.3.4',
  N'The cyber security policy should include:
a. a definition of cyber security;
b. the Member Organization’s overall cyber security objectives and scope;
c. a statement of the board’s intent, supporting the cyber security objectives;
d. a definition of general and specific responsibilities for cyber security;
e. the reference to supporting cyber security standards and procedures;
f. cyber security requirements that ensure:
1. information is classified in a way that indicates its importance to the Member Organization;
2. information is protected in terms of cyber security requirements, in line with the risk appetite;
3. owners are appointed for all information assets;
4. cyber security risk assessments are conducted for information assets;
5. relevant stakeholders are made aware of cyber security and their expected behavior (cyber
security awareness program);
6. compliance with regulatory and contractual obligations are being met;
7. cyber security breaches and suspected cyber security weaknesses are reported;
8. cyber security is reflected in business continuity management. ',
  N'ينبغي أن تتضمن سياسة الأمن السيبراني ما يلي:
أ. تعريف للأمن السيبراني؛
ب. الأهداف العامة للمنظمة العضو المتعلقة بالأمن السيبراني ونطاقها؛
ج. بيان نوايا مجلس الإدارة، يدعم أهداف الأمن السيبراني؛
د. تعريف للمسؤوليات العامة والخاصة المتعلقة بالأمن السيبراني؛
هـ. الإشارة إلى دعم معايير وإجراءات الأمن السيبراني؛
و. متطلبات الأمن السيبراني التي تضمن:
1. تصنيف المعلومات بطريقة تشير إلى أهميتها للمنظمة العضو؛
2. حماية المعلومات من حيث متطلبات الأمن السيبراني، بما يتماشى مع مستوى المخاطرة؛
3. تعيين مالكين لجميع أصول المعلومات؛
4. إجراء تقييمات لمخاطر الأمن السيبراني لأصول المعلومات؛
5. توعية أصحاب المصلحة المعنيين بالأمن السيبراني وسلوكهم المتوقع (برنامج التوعية بالأمن السيبراني)؛
6. الالتزام بالالتزامات التنظيمية والتعاقدية؛
7. الإبلاغ عن خروقات الأمن السيبراني ونقاط الضعف المشتبه بها في الأمن السيبراني؛
8. انعكاس الأمن السيبراني في إدارة استمرارية الأعمال.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_18,
  2,
  N'3.1.4.1',
  N'The Board of Directors has the ultimate responsibility for cyber security, including:
a. ensuring that sufficient budget for cyber security is allocated;
b. approving the cyber security committee charter;
c. endorsing (after being approved by the cyber security committee):
1. the cyber security governance;
2. the cyber security strategy;
3. the cyber security policy.',
  N'يتحمل مجلس الإدارة المسؤولية النهائية عن الأمن السيبراني، بما في ذلك:
أ. ضمان تخصيص ميزانية كافية للأمن السيبراني؛
ب. الموافقة على ميثاق لجنة الأمن السيبراني؛
ج. التصديق (بعد موافقة لجنة الأمن السيبراني):
1. حوكمة الأمن السيبراني؛
2. استراتيجية الأمن السيبراني؛
3. سياسة الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_19,
  2,
  N'3.1.4.2',
  N'The cyber security committee should be responsible for:
a. monitoring, reviewing and communicating the Member Organization’s cyber security risk appetite
periodically or upon a material change in the risk appetite;
b. reviewing the cyber security strategy to ensure that it supports the Member Organization
objectives;
c. approving, communicating, supporting and monitoring: 
1. the cyber security governance;
2. the cyber security strategy;
3. the cyber security policy;
4. cyber security programs (e.g., awareness program, data classification program, data privacy,
data leakage prevention, key cyber security improvements);
5. cyber security risk management process;
6. the key risk indicators (KRIs) and key performance indicators (KPIs) for cyber security. ',
  N'تتولى لجنة الأمن السيبراني مسؤولية:
أ. رصد ومراجعة وإبلاغ مستوى تقبل المنظمة العضو للمخاطر في مجال الأمن السيبراني بشكل دوري أو عند حدوث تغيير جوهري فيه؛
ب. مراجعة استراتيجية الأمن السيبراني لضمان دعمها لأهداف المنظمة العضو؛
ج. الموافقة على ما يلي، وإبلاغه، ودعمه، ورصده:
1. حوكمة الأمن السيبراني؛
2. استراتيجية الأمن السيبراني؛
3. سياسة الأمن السيبراني؛
4. برامج الأمن السيبراني (مثل برنامج التوعية، وبرنامج تصنيف البيانات، وخصوصية البيانات، ومنع تسرب البيانات، والتحسينات الرئيسية للأمن السيبراني)؛
5. عملية إدارة مخاطر الأمن السيبراني؛
6. مؤشرات المخاطر الرئيسية (KRIs) ومؤشرات الأداء الرئيسية (KPIs) للأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_20,
  2,
  N'3.1.4.3',
  N'The senior management should be responsible for:
ensuring that standards, processes and procedures reflect security requirements (if applicable);
b. ensuring that individuals accept and comply with the cyber security policy, supporting standards
and procedures when they are issued and updated;
c. ensuring that cyber security responsibilities are incorporated in the job descriptions of key
positions and cyber security staff. ',
  N'ينبغي أن تكون الإدارة العليا مسؤولة عن:
ضمان أن تعكس المعايير والعمليات والإجراءات متطلبات الأمن (إن وجدت)؛
ب. ضمان قبول الأفراد لسياسة الأمن السيبراني والامتثال لها، ودعم المعايير والإجراءات عند إصدارها وتحديثها؛
ج. ضمان دمج مسؤوليات الأمن السيبراني في الوصف الوظيفي للمناصب الرئيسية وموظفي الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_21,
  2,
  N'3.1.4.4',
  N'The CISO should be responsible for:
a. developing and maintaining:
1. cyber security strategy;
2. cyber security policy;
3. cyber security architecture;
4. cyber security risk management process;
b.ensuring that detailed security standards and procedures are established, approved and
implemented;
c.delivering risk-based cyber security solutions that address people, process and technology;
d.developing the cyber security staff to deliver cyber security solutions in a business context;
e.the cyber security activities across the Member Organization, including:
1. monitoring of the cyber security activities (SOC monitoring);
2. monitoring of compliance with cyber security regulations, policies, standards and procedures;
3. overseeing the investigation of cyber security incidents; 
4. gathering and analyzing threat intelligence from internal and external sources;
5. performing cyber security reviews;
f. conducting cyber security risk assessments on the Members Organization’s information assets;
g. proactively supporting other functions on cyber security, including:
1. performing information and system classifications;
2. determining cyber security requirements for important projects;
3. performing cyber security reviews.
h. defining and conducting the cyber security awareness programs;
i. measuring and reporting the KRIs and KPIs on:
1. cyber security strategy;
2. cyber security policy compliance;
3. cyber security standards and procedures;
4. Cyber security programs (e.g., awareness, data classification, key cyber security improvements)',
  N'ينبغي أن يكون مسؤول أمن المعلومات مسؤولاً عن:
أ. تطوير وصيانة:
1. استراتيجية الأمن السيبراني؛
2. سياسة الأمن السيبراني؛
3. بنية الأمن السيبراني؛
4. عملية إدارة مخاطر الأمن السيبراني؛
ب. ضمان وضع معايير وإجراءات أمنية مفصلة والموافقة عليها وتنفيذها؛
ج. تقديم حلول أمن سيبراني قائمة على المخاطر تُعنى بالأفراد والعمليات والتكنولوجيا؛
د. تطوير موظفي الأمن السيبراني لتقديم حلول أمن سيبراني في سياق الأعمال؛
هـ. أنشطة الأمن السيبراني في جميع أنحاء المنظمة العضو، بما في ذلك:
1. مراقبة أنشطة الأمن السيبراني (مراقبة مركز العمليات الأمنية)؛
2. مراقبة الامتثال للوائح وسياسات ومعايير وإجراءات الأمن السيبراني؛
3. الإشراف على التحقيق في حوادث الأمن السيبراني؛
4. جمع وتحليل معلومات استخبارات التهديدات من مصادر داخلية وخارجية؛
5. إجراء مراجعات الأمن السيبراني؛
و. إجراء تقييمات مخاطر الأمن السيبراني على أصول المعلومات الخاصة بالمنظمة العضو؛
ز. دعم وظائف أخرى في مجال الأمن السيبراني بشكل استباقي، بما في ذلك:
1. تصنيف المعلومات والأنظمة؛
2. تحديد متطلبات الأمن السيبراني للمشاريع المهمة؛
3. إجراء مراجعات للأمن السيبراني.

ح. تحديد وتنفيذ برامج التوعية بالأمن السيبراني؛

ط. قياس مؤشرات الأداء الرئيسية (KRIs) ومؤشرات الأداء الرئيسية (KPIs) وإعداد التقارير عنها فيما يتعلق بـ:
1. استراتيجية الأمن السيبراني؛
2. الامتثال لسياسة الأمن السيبراني؛
3. معايير وإجراءات الأمن السيبراني؛
4. برامج الأمن السيبراني (مثل التوعية، وتصنيف البيانات، والتحسينات الرئيسية في الأمن السيبراني)',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_22,
  2,
  N'3.1.4.5',
  N'The internal audit function should be responsible for:
a. Performing cyber security audits',
  N'ينبغي أن تكون وظيفة التدقيق الداخلي مسؤولة عن: أ. إجراء عمليات تدقيق الأمن السيبراني',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_23,
  2,
  N'3.1.4.6',
  N'All Member Organization’s staff should be responsible for:
a. Complying with cyber security policy, standards, and procedures',
  N'يجب على جميع موظفي المنظمة العضوة أن يكونوا مسؤولين عن: أ. الالتزام بسياسة ومعايير وإجراءات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_24,
  2,
  N'3.1.5.1',
  N'Cyber security should be integrated into the Member Organization''s project management methodology to ensure risks are identified and addressed during the project lifecycle.',
  N'ينبغي دمج الأمن السيبراني في منهجية إدارة المشاريع الخاصة بالمنظمة العضو لضمان تحديد المخاطر ومعالجتها أثناء دورة حياة المشروع.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_25,
  2,
  N'3.1.5.2',
  N'The methodology should ensure:
a. Cyber security objectives are included in project objectives.
b. The cyber security function is part of all project phases.
c. A risk assessment is performed at project start to identify cyber security risks and ensure they are addressed using existing or new controls.
d. Cyber security risks are registered and tracked in the project risk register.
e. Responsibilities for cyber security are clearly defined and assigned.
f. A cyber security review is conducted by an independent internal or external party.',
  N'يجب أن تضمن المنهجية ما يلي:
أ. تضمين أهداف الأمن السيبراني ضمن أهداف المشروع.
ب. أن تكون وظيفة الأمن السيبراني جزءًا لا يتجزأ من جميع مراحل المشروع.
ج. إجراء تقييم للمخاطر عند بدء المشروع لتحديد مخاطر الأمن السيبراني وضمان معالجتها باستخدام الضوابط الحالية أو الجديدة.
د. تسجيل مخاطر الأمن السيبراني وتتبعها في سجل مخاطر المشروع.
هـ. تحديد مسؤوليات الأمن السيبراني وتوزيعها بوضوح.
و. إجراء مراجعة للأمن السيبراني من قبل جهة مستقلة، داخلية أو خارجية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_26,
  2,
  N'3.1.6.1',
  N'Cyber security awareness programs should be defined, approved, and conducted to promote awareness and foster a positive security culture.',
  N'ينبغي تحديد برامج التوعية بالأمن السيبراني والموافقة عليها وتنفيذها لتعزيز الوعي وتعزيز ثقافة أمنية إيجابية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_27,
  2,
  N'3.1.6.2',
  N'Programs should target the following groups:
a. Staff of the Member Organization
b. Third parties of the Member Organization
c. Customers of the Member Organization',
  N'ينبغي أن تستهدف البرامج الفئات التالية:
أ. موظفو المنظمة العضوة
ب. أطراف ثالثة تابعة للمنظمة العضوة
ج. عملاء المنظمة العضوة',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_28,
  2,
  N'3.1.6.3',
  N'Tailor programs to specific audiences and deliver through multiple channels to influence behavior.',
  N'قم بتصميم برامج مخصصة لجمهور محدد وتقديمها عبر قنوات متعددة للتأثير على السلوك.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_29,
  2,
  N'3.1.6.4',
  N' Conduct awareness activities periodically and throughout the year.',
  N'إجراء أنشطة توعوية بشكل دوري وعلى مدار العام.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_30,
  2,
  N'3.1.6.5',
  N'At a minimum, the program should include:
a. Explanation of security measures
b. Roles and responsibilities
c. Information on emerging threats (e.g., spear-phishing, whaling)',
  N'على الأقل، ينبغي أن يتضمن البرنامج ما يلي:
أ. شرح للإجراءات الأمنية
ب. الأدوار والمسؤوليات
ج. معلومات عن التهديدات الناشئة (مثل التصيد الاحتيالي، وصيد الحيتان)',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_31,
  2,
  N'3.1.6.6',
  N'Evaluate the effectiveness of the awareness program.
a. measure the effectiveness of the awareness activities;
b. formulate recommendations to improve the cyber security awareness program.',
  N'تقييم فعالية برنامج التوعية.
أ. قياس فعالية أنشطة التوعية؛
ب. صياغة توصيات لتحسين برنامج التوعية بالأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_32,
  2,
  N'3.1.6.7',
  N'Customer awareness should address for both retail and commercial customers and, at a minimum,
include a listing of suggested cyber security mechanisms which customers may consider implementing
to mitigate their own risk(s). ',
  N'يجب أن يتناول برنامج توعية العملاء كلاً من عملاء التجزئة والعملاء التجاريين، ويجب أن يتضمن على الأقل قائمة بآليات الأمن السيبراني المقترحة التي قد يفكر العملاء في تنفيذها للتخفيف من مخاطرهم الخاصة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_33,
  2,
  N'3.1.7.1',
  N'Specialist or security-related skills training should be provided to staff in the Member Organization’s
relevant functional area categories in line with their job descriptions, including:
a. key roles within the organization;
b. staff of the cyber security function;
c. staff involved in developing and (technically) maintaining information assets;
d. staff involved in risk assessments.',
  N'ينبغي توفير تدريب على المهارات المتخصصة أو المتعلقة بالأمن للموظفين في فئات المجالات الوظيفية ذات الصلة بالمنظمة العضوة بما يتماشى مع أوصاف وظائفهم، بما في ذلك:
أ. الأدوار الرئيسية داخل المنظمة؛
ب. موظفو وظيفة الأمن السيبراني؛
ج. الموظفون المشاركون في تطوير أصول المعلومات وصيانتها (فنيًا)؛
د. الموظفون المشاركون في تقييمات المخاطر.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_34,
  2,
  N'3.1.7.2',
  N'Education should be provided in order to equip staff with the skills and required knowledge to
securely operate the Member Organization’s information assets.',
  N'ينبغي توفير التعليم لتزويد الموظفين بالمهارات والمعرفة اللازمة لتشغيل أصول المعلومات الخاصة بالمنظمة العضوة بشكل آمن.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_35,
  2,
  N'3.2.1.1',
  N'The cyber security risk management process should be defined, approved and implemented. ',
  N'ينبغي تحديد عملية إدارة مخاطر الأمن السيبراني والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_36,
  2,
  N'3.2.1.2',
  N'The cyber security risk management process should focus on safeguarding the confidentiality,
integrity and availability of information assets. ',
  N'ينبغي أن تركز عملية إدارة مخاطر الأمن السيبراني على حماية سرية وسلامة وتوافر أصول المعلومات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_37,
  2,
  N'3.2.1.3',
  N'The cyber security risk management process should be aligned with the existing enterprise risk
management process. ',
  N'ينبغي أن تتوافق عملية إدارة مخاطر الأمن السيبراني مع عملية إدارة المخاطر المؤسسية الحالية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_38,
  2,
  N'3.2.1.4',
  N'The cyber security risk management process should be documented and address:
a. risk identification;
b. risk analysis;
c. risk response;
d. risk monitoring and review. ',
  N'ينبغي توثيق عملية إدارة مخاطر الأمن السيبراني ومعالجة ما يلي:
أ. تحديد المخاطر؛
ب. تحليل المخاطر؛
ج. الاستجابة للمخاطر؛
د. مراقبة المخاطر ومراجعتها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_39,
  2,
  N'3.2.1.5',
  N'The cyber security risk management process should address the Member Organization’s information
assets, including (but not limited to):
a. business processes;
b. business applications;
c. infrastructure components. ',
  N'ينبغي أن تتناول عملية إدارة مخاطر الأمن السيبراني أصول المعلومات الخاصة بالمنظمة العضو، بما في ذلك (على سبيل المثال لا الحصر):
أ. العمليات التجارية؛
ب. تطبيقات الأعمال؛
ج. مكونات البنية التحتية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_40,
  2,
  N'3.2.1.6',
  N'The cyber security risk management process should be initiated:
a. at an early stage of the project;
b. prior to critical change;
c. when outsourcing is being considered;
d. when launching new products and technologies.',
  N'ينبغي البدء في عملية إدارة مخاطر الأمن السيبراني:
أ. في مرحلة مبكرة من المشروع؛
ب. قبل التغيير الحرج؛
ج. عند النظر في الاستعانة بمصادر خارجية؛
د. عند إطلاق منتجات وتقنيات جديدة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_41,
  2,
  N'3.2.1.7',
  N'Existing information assets should be periodically subject to cyber security risk assessment based on
their classification or risk profile. ',
  N'ينبغي أن تخضع أصول المعلومات الموجودة بشكل دوري لتقييم مخاطر الأمن السيبراني بناءً على تصنيفها أو ملف تعريف المخاطر الخاص بها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_42,
  2,
  N'3.2.1.8',
  N'The cyber security risk management activities should involve:
a. business owners;
b. IT specialists;
c. cyber security specialists;
d. key user representatives. ',
  N'ينبغي أن تشمل أنشطة إدارة مخاطر الأمن السيبراني ما يلي:
أ. أصحاب الأعمال؛
ب. متخصصو تكنولوجيا المعلومات؛
ج. متخصصو الأمن السيبراني؛
د. ممثلو المستخدمين الرئيسيون.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_43,
  2,
  N'3.2.1.9',
  N'The result of the risk assessment should be reported to the relevant business owner (i.e., risk owner)
within the Member Organization; ',
  N'يجب الإبلاغ عن نتيجة تقييم المخاطر إلى صاحب العمل المعني (أي مالك المخاطر) داخل المنظمة العضو؛',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_44,
  2,
  N'3.2.1.10',
  N'The relevant business owner (i.e., risk owner) within the Member Organization should accept and
endorse the risk assessment results. ',
  N'يجب على صاحب العمل المعني (أي صاحب المخاطر) داخل المنظمة العضو قبول نتائج تقييم المخاطر والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_45,
  2,
  N'3.2.1.11',
  N'The Member Organization’s cyber security risk appetite and risk tolerance should be clearly defined
and formally approved.',
  N'ينبغي تحديد مدى تقبل المنظمة العضوة للمخاطر المتعلقة بالأمن السيبراني وقدرتها على تحمل المخاطر بشكل واضح، والموافقة عليها رسميًا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_46,
  2,
  N'3.2.1.1.1',
  N'Cyber security risk identification should be performed.',
  N'ينبغي إجراء تحديد مخاطر الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_47,
  2,
  N'3.2.1.1.2',
  N'Identified cyber security risks should be documented (in a central register).',
  N'ينبغي توثيق مخاطر الأمن السيبراني التي تم تحديدها (في سجل مركزي).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_48,
  2,
  N'3.2.1.1.3',
  N'Cyber security risk identification should address relevant information assets, threats, vulnerabilities
and the key existing cyber security controls. ',
  N'ينبغي أن يتناول تحديد مخاطر الأمن السيبراني الأصول المعلوماتية ذات الصلة والتهديدات والثغرات الأمنية وضوابط الأمن السيبراني الرئيسية الموجودة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_49,
  2,
  N'3.2.1.2.1',
  N'A cyber security risk analysis should be performed. ',
  N'ينبغي إجراء تحليل لمخاطر الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_50,
  2,
  N'3.2.1.2.2',
  N'The cyber security risk analysis should address the level of potential business impact and likelihood of
cyber security threat events materializing.',
  N'ينبغي أن يتناول تحليل مخاطر الأمن السيبراني مستوى التأثير التجاري المحتمل واحتمالية وقوع أحداث تهديد الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_51,
  2,
  N'3.2.1.3.1',
  N'The relevant determined cyber security risks should be treated according to the Member
Organization’s risk appetite and cyber security requirements. ',
  N'ينبغي التعامل مع مخاطر الأمن السيبراني المحددة ذات الصلة وفقًا لرغبة المنظمة العضو في تحمل المخاطر ومتطلبات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_52,
  2,
  N'3.2.1.3.2',
  N'Cyber security risk response should ensure that the list of risk treatment options are documented (i.e.,
accepting, avoiding, transferring or mitigating risks by applying cyber security controls).',
  N'ينبغي أن تضمن الاستجابة لمخاطر الأمن السيبراني توثيق قائمة خيارات معالجة المخاطر (أي قبول المخاطر أو تجنبها أو نقلها أو التخفيف منها من خلال تطبيق ضوابط الأمن السيبراني).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_53,
  2,
  N'3.2.1.3.3',
  N'Accepting cyber security risks should include:
a. the consideration of predefined limits for levels of cyber security risk;
b. the approval and sign-off by the business owner, ensuring that:
1. the accepted cyber security risk is within the risk appetite and is reported to the cyber security
committee;
2. the accepted cyber security risk does not contradict SAMA regulations. ',
  N'يجب أن يشمل قبول مخاطر الأمن السيبراني ما يلي:
أ. مراعاة حدود محددة مسبقًا لمستويات مخاطر الأمن السيبراني؛
ب. موافقة وتوقيع صاحب العمل، مع ضمان ما يلي:
1. أن تكون مخاطر الأمن السيبراني المقبولة ضمن حدود المخاطر المقبولة، وأن يتم الإبلاغ عنها إلى لجنة الأمن السيبراني؛
2. أن لا تتعارض مخاطر الأمن السيبراني المقبولة مع لوائح مؤسسة النقد العربي السعودي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_54,
  2,
  N'3.2.1.3.4',
  N'Avoiding cyber security risks should involve a decision by a business owner to cancel or postpone a
particular activity or project that introduces an unacceptable cyber security risk. ',
  N'يجب أن يتضمن تجنب مخاطر الأمن السيبراني اتخاذ صاحب العمل قرارًا بإلغاء أو تأجيل نشاط أو مشروع معين يمثل خطرًا غير مقبول على الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_55,
  2,
  N'3.2.1.3.5',
  N'Transferring or sharing the cyber security risks should:
a. involve sharing the cyber security risks with relevant (internal or external) providers;
b. be accepted by the receiving (internal or external) provider(s);
c. eventually lead to the actual transferring or sharing of the cyber security risk. ',
  N'يجب أن يتضمن نقل أو مشاركة مخاطر الأمن السيبراني ما يلي:
أ. مشاركة مخاطر الأمن السيبراني مع مقدمي الخدمات ذوي الصلة (الداخليين أو الخارجيين)؛
ب. أن يتم قبولها من قبل مقدمي الخدمات المستلمين (الداخليين أو الخارجيين)؛
ج. أن تؤدي في النهاية إلى النقل الفعلي أو مشاركة مخاطر الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_56,
  2,
  N'3.2.1.3.6',
  N'Applying cyber security controls to mitigate cyber security risks should include:
a. identifying appropriate cyber security controls;
b. evaluating the strengths and weaknesses of the cyber security controls;
1. assessing the cost of implementing the cyber security controls;
2. assessing the feasibility of implementing the cyber security controls;
3. reviewing relevant compliance requirements for the cyber security controls;
c. selecting cyber security controls;
d. identifying, documenting and obtaining sign-off for any residual risk by the business owner. ',
  N'يجب أن يشمل تطبيق ضوابط الأمن السيبراني للتخفيف من مخاطر الأمن السيبراني ما يلي:
أ. تحديد ضوابط الأمن السيبراني المناسبة؛
ب. تقييم نقاط القوة والضعف في ضوابط الأمن السيبراني؛
1. تقييم تكلفة تطبيق ضوابط الأمن السيبراني؛
2. تقييم جدوى تطبيق ضوابط الأمن السيبراني؛
3. مراجعة متطلبات الامتثال ذات الصلة لضوابط الأمن السيبراني؛
ج. اختيار ضوابط الأمن السيبراني؛
د. تحديد وتوثيق والحصول على موافقة صاحب العمل على أي مخاطر متبقية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_57,
  2,
  N'3.2.1.3.7',
  N'Cyber security risk treatment actions should be documented in a risk treatment plan. ',
  N'ينبغي توثيق إجراءات معالجة مخاطر الأمن السيبراني في خطة معالجة المخاطر.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_58,
  2,
  N'3.2.1.4.1',
  N'The cyber security treatment should be monitored, including:
a. tracking progress in accordance to treatment plan;
b. the selected and agreed cyber security controls are being implemented.',
  N'يجب مراقبة معالجة الأمن السيبراني، بما في ذلك:
أ. تتبع التقدم وفقًا لخطة المعالجة؛
ب. تنفيذ ضوابط الأمن السيبراني المختارة والمتفق عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_59,
  2,
  N'3.2.1.4.2',
  N'The design and effectiveness of the revised or newly implemented cyber security controls should be
reviewed. ',
  N'ينبغي مراجعة تصميم وفعالية ضوابط الأمن السيبراني المنقحة أو التي تم تنفيذها حديثًا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_60,
  2,
  N'3.2.2.1',
  N'1. A process should be established for ensuring compliance with relevant regulatory requirements
affecting cyber security across the Member Organization. The process of ensuring compliance should:
a. be performed periodically or when new regulatory requirements become effective;
b. involve representatives from key areas of the Member Organization;
c. result in the update of cyber security policy, standards and procedures to accommodate any
necessary changes (if applicable). ',
  N'١. ينبغي وضع آلية لضمان الامتثال للمتطلبات التنظيمية ذات الصلة التي تؤثر على الأمن السيبراني في جميع أنحاء المنظمة العضو. ينبغي أن تتضمن عملية ضمان الامتثال ما يلي:
أ. أن تُنفذ بشكل دوري أو عند سريان المتطلبات التنظيمية الجديدة؛
ب. أن تشمل ممثلين من القطاعات الرئيسية في المنظمة العضو؛
ج. أن تؤدي إلى تحديث سياسة ومعايير وإجراءات الأمن السيبراني لاستيعاب أي تغييرات ضرورية (إن وجدت).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_61,
  2,
  N'3.2.3.1',
  N'Control considerations
1. The Member Organization should comply with:
a. Payment Card Industry Data Security Standard (PCI-DSS);
b. EMV (Europay, MasterCard and Visa) technical standard;
c. SWIFT Customer Security Controls Framework – March 2017.',
  N'اعتبارات الرقابة
1. ينبغي على المنظمة العضوة الامتثال لما يلي:
أ. معيار أمن بيانات صناعة بطاقات الدفع (PCI-DSS)؛
ب. المعيار الفني EMV (يوروباي، ماستركارد، وفيزا)؛
ج. إطار عمل ضوابط أمن عملاء سويفت - مارس 2017.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_62,
  2,
  N'3.2.4.1',
  N'Cyber security reviews should be periodically performed for critical information assets. ',
  N'ينبغي إجراء مراجعات الأمن السيبراني بشكل دوري للأصول المعلوماتية الهامة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_63,
  2,
  N'3.2.4.2',
  N'Customer and internet facing services should be subject to annual review and penetration tests. ',
  N'يجب أن تخضع الخدمات المقدمة للعملاء والإنترنت للمراجعة السنوية واختبارات الاختراق.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_64,
  2,
  N'3.2.4.3',
  N'Details of cyber security review performed should be recorded, including the results of review, issues
identified and recommended actions. ',
  N'ينبغي تسجيل تفاصيل مراجعة الأمن السيبراني التي تم إجراؤها، بما في ذلك نتائج المراجعة والمشكلات التي تم تحديدها والإجراءات الموصى بها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_65,
  2,
  N'3.2.4.4',
  N'The results of cyber security review should be reported to business owner.',
  N'ينبغي الإبلاغ عن نتائج مراجعة الأمن السيبراني إلى صاحب العمل.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_66,
  2,
  N'3.2.4.5',
  N'Cyber security review should be subject to follow-up reviews to check that:
a. all identified issues have been addressed; 
b. critical risks have been treated effectively; 
c. all agreed actions are being managed on an ongoing basis. ',
  N'يجب أن تخضع مراجعة الأمن السيبراني لمراجعات متابعة للتأكد من: أ. معالجة جميع القضايا التي تم تحديدها؛ ب. معالجة المخاطر الحرجة بشكل فعال؛ ج. إدارة جميع الإجراءات المتفق عليها بشكل مستمر.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_67,
  2,
  N'3.2.5.1',
  N'Cyber security audits should be performed independently and according to generally accepted
auditing standards and SAMA cyber security framework. ',
  N'يجب إجراء عمليات تدقيق الأمن السيبراني بشكل مستقل ووفقًا لمعايير التدقيق المقبولة عمومًا وإطار عمل الأمن السيبراني لمؤسسة النقد العربي السعودي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_68,
  2,
  N'3.2.5.2',
  N'Cyber security audits should be performed according to the Member Organization’s audit manual and
audit plan. ',
  N'يجب إجراء عمليات تدقيق الأمن السيبراني وفقًا لدليل التدقيق وخطة التدقيق الخاصة بالمنظمة العضو.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_69,
  2,
  N'3.3.1.1',
  N'The human resources process should define, approve and implement cyber security requirements.',
  N'ينبغي لعملية الموارد البشرية أن تحدد وتوافق على وتنفذ متطلبات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_70,
  2,
  N'3.3.1.2',
  N'The effectiveness of the human resources process should be monitored, measured and periodically evaluated.',
  N'ينبغي مراقبة فعالية عملية الموارد البشرية وقياسها وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_71,
  2,
  N'3.3.1.3',
  N'The human resource process should include:
a. cyber security responsibilities and non-disclosure clauses within staff agreements (during and
after the employment);
b. staff should receive cyber security awareness at the start and during their employment;
c. when disciplinary actions will be applicable;
d. screening and background check; 
e. post-employment cyber security activities, such as:
1. revoking access rights;
2. returning information assets assigned (e.g., access badge, tokens, mobile devices, all electronic
and physical information). ',
  N'ينبغي أن تشمل عملية الموارد البشرية ما يلي:
أ. مسؤوليات الأمن السيبراني وبنود عدم الإفصاح ضمن اتفاقيات الموظفين (أثناء وبعد التوظيف)؛
ب. ينبغي أن يتلقى الموظفون وعيًا بالأمن السيبراني في بداية عملهم وأثناءه؛
ج. متى ستكون الإجراءات التأديبية قابلة للتطبيق؛
د. الفحص والتحقق من الخلفية؛
هـ. أنشطة الأمن السيبراني بعد التوظيف، مثل:
1. إلغاء حقوق الوصول؛
2. إعادة أصول المعلومات المخصصة (مثل شارة الوصول، والرموز، والأجهزة المحمولة، وجميع المعلومات الإلكترونية والمادية).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_72,
  2,
  N'3.3.2.1',
  N'The physical security process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية الأمن المادي والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_73,
  2,
  N'3.3.2.2',
  N'The effectiveness of the physical security process should be monitored, measured and periodically evaluated.',
  N'ينبغي مراقبة فعالية عملية الأمن المادي وقياسها وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_74,
  2,
  N'3.3.2.3',
  N'The physical security process should include (but not limited to):
a. physical entry controls (including visitor security);
b. monitoring and surveillance (e.g., CCTV, ATMs GPS tracking, sensitivity sensors);
c. protection of data centers and data rooms;
d. environmental protection;
e. protection of information assets during lifecycle (including transport and secure disposal, avoiding
unauthorized access and (un)intended data leakage. ',
  N'يجب أن تشمل عملية الأمن المادي (على سبيل المثال لا الحصر):
أ. ضوابط الدخول المادي (بما في ذلك أمن الزوار)؛
ب. المراقبة والرصد (مثل كاميرات المراقبة التلفزيونية المغلقة، وتتبع أجهزة الصراف الآلي بنظام تحديد المواقع العالمي (GPS)، وأجهزة استشعار الحساسية)؛
ج. حماية مراكز البيانات وغرف البيانات؛
د. حماية البيئة؛
هـ. حماية أصول المعلومات أثناء دورة حياتها (بما في ذلك النقل والتخلص الآمن منها، وتجنب الوصول غير المصرح به وتسرب البيانات (غير المقصود).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_75,
  2,
  N'3.3.3.1',
  N'The asset management process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية إدارة الأصول والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_76,
  2,
  N'3.3.3.2',
  N'The effectiveness of the asset management process should be monitored, measured and periodically
evaluated.',
  N'ينبغي مراقبة فعالية عملية إدارة الأصول وقياسها وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_77,
  2,
  N'3.3.3.3',
  N'The asset management process should include:
a. a unified register;
b. ownership and custodianship of information assets;
c. the reference to relevant other processes, depending on asset management;
d. information asset classification, labeling and handling;
e. the discovery of new information assets.',
  N'ينبغي أن تشمل عملية إدارة الأصول ما يلي:
أ. سجل موحد؛
ب. ملكية أصول المعلومات وحراستها؛
ج. الإشارة إلى العمليات الأخرى ذات الصلة، حسب إدارة الأصول؛
د. تصنيف أصول المعلومات، ووضع العلامات عليها، ومعالجتها؛
هـ. اكتشاف أصول معلومات جديدة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_78,
  2,
  N'3.3.4.1',
  N'The cyber security architecture should be defined, approved and implemented.',
  N'ينبغي تحديد بنية الأمن السيبراني والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_79,
  2,
  N'3.3.4.2',
  N'The compliance with the cyber security architecture should be monitored.',
  N'ينبغي مراقبة الامتثال لهندسة الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_80,
  2,
  N'3.3.4.3',
  N'The cyber security architecture should include:
a. a strategic outline of cyber security capabilities and controls based on the business requirements;
b. approval of the defined cyber security architecture;
c. the requirement of having qualified cyber security architects;
d. design principles for developing cyber security controls and applying cyber security requirements
(i.e., the security-by-design principle);
e. periodic review of the cyber security architecture.',
  N'يجب أن تتضمن بنية الأمن السيبراني ما يلي:
أ. مخطط استراتيجي لقدرات وضوابط الأمن السيبراني بناءً على متطلبات العمل؛
ب. الموافقة على بنية الأمن السيبراني المحددة؛
ج. شرط وجود مهندسي أمن سيبراني مؤهلين؛
د. مبادئ التصميم لتطوير ضوابط الأمن السيبراني وتطبيق متطلباته (أي مبدأ الأمن بالتصميم)؛
هـ. مراجعة دورية لبنية الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_81,
  2,
  N'3.3.5.1',
  N'The identity and access management policy, including the responsibilities and accountabilities, should
be defined, approved and implemented.  ',
  N'ينبغي تحديد سياسة إدارة الهوية والوصول، بما في ذلك المسؤوليات والالتزامات، والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_82,
  2,
  N'3.3.5.2',
  N'The compliance with the identity and access policy should be monitored.',
  N'ينبغي مراقبة الامتثال لسياسة الهوية والوصول.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_83,
  2,
  N'3.3.5.3',
  N'The effectiveness of the cyber security controls within the identity and access management policy should be measured and periodically evaluated.',
  N'ينبغي قياس وتقييم فعالية ضوابط الأمن السيبراني ضمن سياسة إدارة الهوية والوصول بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_84,
  2,
  N'3.3.5.4',
  N'The identity and access management policy should include:
a. business requirements for access control (i.e., need-to-have and need-to-know);
b. user access management (e.g., joiners, movers, leavers):
1. all identified user types should be covered (i.e., internal staff, third parties);
2. changes of job status or job positions for internal staff (e.g. joiner, mover and leaver) should
be instigated by the human resources department;
3. changes for external staff or third parties should be instigated by the appointed accountable
party;
4. user access requests are formally approved in accordance with business and compliance
requirements (i.e., need-to-have and need-to-know to avoid unauthorized access and
(un)intended data leakage));
5. changes in access rights should be processed in a timely manner;
6. periodically user access rights and profiles should be reviewed;
7. an audit trail of submitted, approved and processed user access requests and revocation
requests should be established;
c. user access management should be supported by automation;
d. centralization of the identity and access management function;
e. multi-factor authentication for sensitive and critical systems and profiles;
f.
privileged and remote access management, which should address:
1. the allocation and restricted use of privileged and remote access, specifying:
a. multi-factor authentication should be used for all remote access;
b. multi-factor authentication should be used for privilege access on critical systems based
on a risk assessment;
2. the periodic review of users with privileged and remote accounts;
3. individual accountability;
4. the use of non-personal privileged accounts, including:
a. limitation and monitoring;
b. confidentiality of passwords;
c. changing passwords frequently and at the end of each session. ',
  N'يجب أن تتضمن سياسة إدارة الهوية والوصول ما يلي:
أ. متطلبات العمل للتحكم في الوصول (أي، ما يلزم امتلاكه وما يلزم معرفته)؛
ب. إدارة وصول المستخدمين (مثل: المنضمون، والمنتقلون، والمغادرون):
1. يجب تغطية جميع أنواع المستخدمين المحددة (أي: الموظفون الداخليون، والأطراف الثالثة)؛
2. يجب أن تتم الموافقة على تغييرات الحالة الوظيفية أو المناصب الوظيفية للموظفين الداخليين (مثل: المنضمون، والمنتقلون، والمغادرون) من قِبل إدارة الموارد البشرية؛
3. يجب أن تتم الموافقة على التغييرات المتعلقة بالموظفين الخارجيين أو الأطراف الثالثة من قِبل الجهة المسؤولة المُعيّنة؛
4. تتم الموافقة رسميًا على طلبات وصول المستخدمين وفقًا لمتطلبات العمل والامتثال (أي، ما يلزم امتلاكه وما يلزم معرفته لتجنب الوصول غير المصرح به وتسريب البيانات (غير المقصود))؛
5. يجب معالجة التغييرات في حقوق الوصول في الوقت المناسب؛
6. يجب مراجعة حقوق وصول المستخدمين وملفاتهم الشخصية بشكل دوري؛
٧. ينبغي إنشاء سجل تدقيق لطلبات وصول المستخدمين المُقدمة والمُعتمدة والمُعالجة وطلبات الإلغاء؛

ج. ينبغي دعم إدارة وصول المستخدمين بالأتمتة؛

د. مركزية وظيفة إدارة الهوية والوصول؛

هـ. المصادقة متعددة العوامل للأنظمة والملفات الشخصية الحساسة والحرجة؛

و. إدارة الوصول المتميز والبعيد، والتي ينبغي أن تتناول:

١. تخصيص وتقييد استخدام الوصول المتميز والبعيد، مع تحديد:

أ. ينبغي استخدام المصادقة متعددة العوامل لجميع عمليات الوصول عن بُعد؛

ب. ينبغي استخدام المصادقة متعددة العوامل للوصول المتميز على الأنظمة الحرجة بناءً على تقييم المخاطر؛

٢. المراجعة الدورية للمستخدمين ذوي الحسابات المتميزة والبعيدة؛

٣. المساءلة الفردية؛

٤. استخدام الحسابات المتميزة غير الشخصية، بما في ذلك:

أ. التقييد والمراقبة؛

ب. سرية كلمات المرور؛

ج. تغيير كلمات المرور بشكل متكرر وفي نهاية كل جلسة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_85,
  2,
  N'3.3.6.1',
  N'The application cyber security standards should be defined, approved and implemented. ',
  N'ينبغي تحديد معايير الأمن السيبراني للتطبيق والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_86,
  2,
  N'3.3.6.2',
  N'The compliance with the application security standards should be monitored.',
  N'ينبغي مراقبة الامتثال لمعايير أمان التطبيق.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_87,
  2,
  N'3.3.6.3',
  N'The effectiveness of the application cyber security controls should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية ضوابط الأمن السيبراني للتطبيق وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_88,
  2,
  N'3.3.6.4',
  N'Application development should follow the approved secure system development life cycle methodology (SDLC).',
  N'يجب أن يتبع تطوير التطبيقات منهجية دورة حياة تطوير النظام الآمنة المعتمدة (SDLC).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_89,
  2,
  N'3.3.6.5',
  N'The application security standard should include:
a. secure coding standards;
b. the cyber security controls implemented (e.g., configuration parameters, events to monitor and
retain [including system access and data], identity and access management);
c. the segregation of duties within the application (supported with a documented authorization
matrix);
d. the protection of data aligned with the (agreed) classification scheme (including privacy of
customer data and, avoiding unauthorized access and (un)intended data leakage);
e. vulnerability and patch management; 
f.
back-up and recovery procedures;
g. periodic cyber security compliance review.',
  N'يجب أن يتضمن معيار أمان التطبيق ما يلي:
أ. معايير الترميز الآمن؛
ب. ضوابط الأمن السيبراني المُطبقة (مثل معلمات التكوين، والأحداث التي يجب مراقبتها والاحتفاظ بها [بما في ذلك الوصول إلى النظام والبيانات]، وإدارة الهوية والوصول)؛
ج. فصل المهام داخل التطبيق (مدعومًا بمصفوفة تفويض موثقة)؛
د. حماية البيانات بما يتماشى مع نظام التصنيف (المتفق عليه) (بما في ذلك خصوصية بيانات العميل، وتجنب الوصول غير المصرح به وتسريب البيانات (غير المقصود)؛
هـ. إدارة الثغرات الأمنية والتصحيحات؛
و.
إجراءات النسخ الاحتياطي والاسترداد؛
ز. مراجعة دورية للامتثال للأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_90,
  2,
  N'3.3.7.1',
  N' The change management process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية إدارة التغيير والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_91,
  2,
  N'3.3.7.2',
  N'The compliance with the change management process should be monitored. ',
  N'ينبغي مراقبة الامتثال لعملية إدارة التغيير.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_92,
  2,
  N'3.3.7.3',
  N'The effectiveness of the cyber security controls within the change management process should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية ضوابط الأمن السيبراني ضمن عملية إدارة التغيير وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_93,
  2,
  N'3.3.7.4',
  N'The change management process should include:
a. cyber security requirements for controlling changes to information assets, such as assessing the
impact of requested changes, classification of changes and the review of changes;
b. security testing, which should (if applicable) include:
1. penetration testing;
2. code review if applications are developed internally;
3. code review of externally developed applications and if the source code is available
4. a code review report (or equivalent, such as an independent assurance statement) in case the source code cannot be provided;
c. approval of changes by the business owner;
d. approval from the cyber security function before submitting to Change Advisory Board (CAB); 
e. approval by CAB;
f. post-implementation review of the related cyber security controls;
g. development, testing and implementation are segregated for both the (technical) environment and involved individuals;
h. the procedure for emergency changes and fixes;
i. fall-back and roll-back procedures. ',
  N'ينبغي أن تشمل عملية إدارة التغيير ما يلي:
أ. متطلبات الأمن السيبراني للتحكم في التغييرات على أصول المعلومات، مثل تقييم أثر التغييرات المطلوبة، وتصنيف التغييرات، ومراجعتها؛
ب. اختبار الأمان، والذي ينبغي أن يشمل (إن وجد) ما يلي:
1. اختبار الاختراق؛
2. مراجعة الكود إذا طُوّرت التطبيقات داخليًا؛
3. مراجعة كود التطبيقات المُطوّرة خارجيًا، وما إذا كان الكود المصدر متاحًا.
4. تقرير مراجعة الكود (أو ما يعادله، مثل بيان ضمان مستقل) في حال تعذر توفير الكود المصدر؛
ج. موافقة صاحب العمل على التغييرات؛
د. موافقة إدارة الأمن السيبراني قبل تقديمها إلى مجلس استشارات التغيير؛
هـ. موافقة مجلس استشارات التغيير؛
و. مراجعة ما بعد التنفيذ لضوابط الأمن السيبراني ذات الصلة؛
ز. يتم فصل التطوير والاختبار والتنفيذ لكل من البيئة (الفنية) والأفراد المعنيين؛
ح. إجراءات التغييرات والإصلاحات الطارئة؛
ط. إجراءات الرجوع إلى الوضع السابق والتراجع عنه.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_94,
  2,
  N'3.3.8.1',
  N'The infrastructure security standards should be defined, approved and implemented. ',
  N'ينبغي تحديد معايير أمن البنية التحتية والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_95,
  2,
  N'3.3.8.2',
  N'The compliance with the infrastructure security standards should be monitored.',
  N'ينبغي مراقبة الامتثال لمعايير أمن البنية التحتية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_96,
  2,
  N'3.3.8.3',
  N'The effectiveness of the infrastructure cyber security controls should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية ضوابط الأمن السيبراني للبنية التحتية وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_97,
  2,
  N'3.3.8.4',
  N'The infrastructure security standards should cover all instances of infrastructure available in the main datacenter(s), the disaster recovery data site(s) and office spaces.',
  N'ينبغي أن تغطي معايير أمن البنية التحتية جميع حالات البنية التحتية المتوفرة في مراكز البيانات الرئيسية، وموقع (مواقع) بيانات الاسترداد من الكوارث، ومساحات المكاتب.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_98,
  2,
  N'3.3.8.5',
  N'The infrastructure security standards should cover all instances of infrastructure (e.g., operating
systems, servers, virtual machines, firewalls, network devices, IDS, IPS, wireless network, gateway
servers, proxy servers, email gateways, external connections, databases, file-shares, workstations,
laptops, tablets, mobile devices, PBX).',
  N'يجب أن تغطي معايير أمن البنية التحتية جميع حالات البنية التحتية (على سبيل المثال، أنظمة التشغيل، والخوادم، والآلات الافتراضية، وجدران الحماية، وأجهزة الشبكة، وأنظمة كشف التسلل، وأنظمة منع التسلل، والشبكات اللاسلكية، وخوادم البوابة، وخوادم الوكيل، وبوابات البريد الإلكتروني، والاتصالات الخارجية، وقواعد البيانات، ومشاركات الملفات، ومحطات العمل، وأجهزة الكمبيوتر المحمولة، والأجهزة اللوحية، والأجهزة المحمولة، ونظام PBX).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_99,
  2,
  N'3.3.8.6',
  N'The infrastructure security standard should include:
a. the cyber security controls implemented (e.g., configuration parameters, events to monitor and retain [including system access and data], data-leakage prevention [DLP], identity and access management, remote maintenance);
b. the segregation of duties within the infrastructure component (supported with a documented authorization matrix);
c. the protection of data aligned with the (agreed) classification scheme (including privacy of
customer data and, avoiding unauthorized access and (un)intended data leakage);
d. the use of approved software and secure protocols;
e. segmentation of networks;
f. malicious code/software and virus protection (and applying application whitelisting and APT protection);
g. vulnerability and patch management;
h. DDOS protection (where applicable); this should include:
1. the use of scrubbing services;
2. specification of the bandwidth agreed;
3. 24x7 monitoring by Security Operating Center (SOC), Service Provider (SP) and scrubbing provider;
4. testing of DDOS scrubbing (minimum twice a year);
5. DDOS services should be implemented for the main datacenter(s) as well as the disaster
recovery site(s);
i. back-up and recovery procedures;
j. periodic cyber security compliance review.',
  N'يجب أن يتضمن معيار أمن البنية التحتية ما يلي:
أ. ضوابط الأمن السيبراني المُطبقة (مثل معلمات التكوين، وأحداث المراقبة والاحتفاظ [بما في ذلك الوصول إلى النظام والبيانات]، ومنع تسرب البيانات [DLP]، وإدارة الهوية والوصول، والصيانة عن بُعد)؛
ب. فصل المهام داخل مكون البنية التحتية (مع دعم مصفوفة تفويض موثقة)؛
ج. حماية البيانات بما يتماشى مع نظام التصنيف (المتفق عليه) (بما في ذلك خصوصية بيانات العميل، وتجنب الوصول غير المصرح به وتسريب البيانات (غير المقصود))؛
د. استخدام برامج معتمدة وبروتوكولات آمنة؛
هـ. تجزئة الشبكات؛
و. الحماية من البرامج/الأكواد الضارة والفيروسات (وتطبيق القائمة البيضاء للتطبيقات وحماية التهديدات المتقدمة المستمرة)؛
ز. إدارة الثغرات الأمنية والتصحيحات؛
ح. الحماية من هجمات الحرمان من الخدمة الموزعة (DDOS) (إن وجدت)؛ يجب أن يشمل ذلك:
1. استخدام خدمات التنظيف؛
2. تحديد النطاق الترددي المتفق عليه؛
٣. مراقبة على مدار الساعة طوال أيام الأسبوع من قِبل مركز عمليات الأمن (SOC)، ومُقدِّم الخدمة (SP)، ومُقدِّم خدمات التنظيف؛
٤. اختبار تنظيف هجمات الحرمان من الخدمة الموزعة (DDOS) (مرتين سنويًا على الأقل)؛
٥. ينبغي تطبيق خدمات الحرمان من الخدمة الموزعة (DDOS) في مراكز البيانات الرئيسية، بالإضافة إلى مواقع التعافي من الكوارث؛
١. إجراءات النسخ الاحتياطي والاسترداد؛
ي. مراجعة دورية للامتثال للأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_100,
  2,
  N'3.3.9.1',
  N'A cryptographic security standard should be defined, approved and implemented.',
  N'ينبغي تحديد معيار أمان التشفير والموافقة عليه وتنفيذه.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_101,
  2,
  N'3.3.9.2',
  N'The compliance with the cryptographic security standard should be monitored.',
  N'ينبغي مراقبة الامتثال لمعايير الأمان التشفيرية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_102,
  2,
  N'3.3.9.3',
  N'The effectiveness of the cryptographic security controls should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية ضوابط الأمن التشفيرية وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_103,
  2,
  N'3.3.9.4',
  N'The cryptographic security standard should include:
a. an overview of the approved cryptographic solutions and relevant restrictions (e.g., technically,
legally);
b. the circumstances when the approved cryptographic solutions should be applied;
c. the management of encryption keys, including lifecycle management, archiving and recovery.',
  N'يجب أن يتضمن معيار أمان التشفير ما يلي:
أ. لمحة عامة عن حلول التشفير المعتمدة والقيود ذات الصلة (على سبيل المثال، من الناحية الفنية والقانونية)؛
ب. الظروف التي ينبغي فيها تطبيق حلول التشفير المعتمدة؛
ج. إدارة مفاتيح التشفير، بما في ذلك إدارة دورة الحياة والأرشفة والاسترداد.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_104,
  2,
  N'3.3.10.1',
  N'The BYOD cyber security standard should be defined, approved and implemented. ',
  N'ينبغي تحديد معيار الأمن السيبراني BYOD والموافقة عليه وتنفيذه.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_105,
  2,
  N'3.3.10.2',
  N'The compliance with the BYOD cyber security standard should be monitored.  ',
  N'ينبغي مراقبة الامتثال لمعيار الأمن السيبراني BYOD.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_106,
  2,
  N'3.3.10.3',
  N'Effectiveness of the BYOD cyber security controls should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية ضوابط الأمن السيبراني الخاصة بـ BYOD وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_107,
  2,
  N'3.3.10.4',
  N'The BYOD standard should include:
a. responsibilities of the user (including awareness training);
b. information regarding the restrictions and consequences for staff when the Member Organization
implements cyber security controls on their personal devices; for example when using modified
devices (jailbreaking), terminating the employment or in case of loss or theft of the personal
device;
c. the isolation of business information from personal information (e.g., containerization);
d. the regulation of corporate mobile applications or approved “public” mobile applications;
e. the use of mobile device management (MDM); applying access controls to the device and business
container and encryption mechanisms on the personal device (to ensure secure transmission and
storage). ',
  N'يجب أن يتضمن معيار إحضار أجهزتك الخاصة (BYOD) ما يلي:
أ. مسؤوليات المستخدم (بما في ذلك التدريب على التوعية)؛
ب. معلومات بشأن القيود والعواقب التي قد يتعرض لها الموظفون عند تطبيق المنظمة العضو ضوابط الأمن السيبراني على أجهزتهم الشخصية؛ على سبيل المثال عند استخدام أجهزة معدلة (كسر الحماية)، أو إنهاء الخدمة، أو في حالة فقدان أو سرقة الجهاز الشخصي؛
ج. عزل معلومات العمل عن المعلومات الشخصية (مثل: النقل بالحاويات)؛
د. تنظيم تطبيقات الهاتف المحمول للشركات أو تطبيقات الهاتف المحمول "العامة" المعتمدة؛
هـ. استخدام إدارة الأجهزة المحمولة (MDM)؛ وتطبيق ضوابط الوصول إلى الجهاز وحاوية العمل وآليات التشفير على الجهاز الشخصي (لضمان النقل والتخزين الآمنين).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_108,
  2,
  N'3.3.11.1',
  N'The secure disposal standard and procedure should be defined, approved and implemented. ',
  N'ينبغي تحديد معيار وإجراءات التخلص الآمن من النفايات والموافقة عليهما وتنفيذهما.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_109,
  2,
  N'3.3.11.2',
  N'The compliance with the secure disposal standard and procedure should be monitored.  ',
  N'ينبغي مراقبة الامتثال لمعايير وإجراءات التخلص الآمن.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_110,
  2,
  N'3.3.11.3',
  N'The effectiveness of the secure disposal cyber security controls should be measured and periodically
evaluated. ',
  N'ينبغي قياس فعالية ضوابط الأمن السيبراني للتخلص الآمن من النفايات وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_111,
  2,
  N'3.3.11.4',
  N'Information assets should be disposed in accordance with legal and regulatory requirements, when
no longer required (i.e. meeting data privacy regulations to avoid unauthorized access and avoid 
(un)intended data leakage). ',
  N'يجب التخلص من أصول المعلومات وفقًا للمتطلبات القانونية والتنظيمية، عندما لا تكون هناك حاجة إليها بعد الآن (أي تلبية لوائح خصوصية البيانات لتجنب الوصول غير المصرح به وتجنب تسرب البيانات (غير المقصود)).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_112,
  2,
  N'3.3.11.5',
  N'Sensitive information should be destroyed using techniques to make the information non-retrievable
(e.g., secure erase, secure wiping, incineration, double crosscut, shredding).',
  N'يجب تدمير المعلومات الحساسة باستخدام تقنيات تجعلها غير قابلة للاسترجاع (على سبيل المثال، المسح الآمن، والحرق، والقطع المزدوج، والتقطيع).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_113,
  2,
  N'3.3.11.6',
  N'The Member Organization should ensure that third party service providers used for secure disposal,
transport and storage comply with the secure disposal standard and procedure and the effectiveness
is periodically measured and evaluated.',
  N'يجب على المنظمة العضو التأكد من أن مقدمي الخدمات من الجهات الخارجية المستخدمة للتخلص الآمن من النفايات ونقلها وتخزينها يمتثلون لمعايير وإجراءات التخلص الآمن، وأن الفعالية يتم قياسها وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_114,
  2,
  N'3.3.12.1',
  N'For Saudi Arabian Riyal Interbank Express (SARIE) information, please refer to the SARIE Information Security Policy, Version Issue 1.0 - June 2016.
2. For mada information, please refer to the following sections in the mada Rules and Standards Technical Book (see appendix A):
3. Part IIIa - Security Framework, Version Issue 6.0.0 - May 2016
4. Part IIIb - HSM Requirements, Version Issue 6.0.0 - May 2016
5. SAMA CA IPK Certificate Procedures, Version Issue 6.0.1 – October 2016',
  N'للحصول على معلومات حول نظام سريع للتحويلات المالية السريعة (SARIE)، يُرجى مراجعة سياسة أمن المعلومات الخاصة بنظام سريع، الإصدار 1.0 - يونيو 2016.

2. للحصول على معلومات حول مدى، يُرجى مراجعة الأقسام التالية في الدليل الفني لقواعد ومعايير مدى (انظر الملحق أ):

3. الجزء الثالث (أ) - إطار العمل الأمني، الإصدار 6.0.0 - مايو 2016

4. الجزء الثالث (ب) - متطلبات وحدة الأمن (HSM)، الإصدار 6.0.0 - مايو 2016

5. إجراءات شهادة IPK الصادرة عن مؤسسة النقد العربي السعودي (SAMA CA)، الإصدار 6.0.1 - أكتوبر 2016',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_115,
  2,
  N'3.3.13.1',
  N'The cyber security standards for electronic banking services should be defined, approved and implemented.',
  N'ينبغي تحديد معايير الأمن السيبراني للخدمات المصرفية الإلكترونية والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_116,
  2,
  N'3.3.13.2',
  N'The compliance with cyber security standards for electronic banking services should be monitored.',
  N'ينبغي مراقبة الامتثال لمعايير الأمن السيبراني للخدمات المصرفية الإلكترونية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_117,
  2,
  N'3.3.13.3',
  N'The effectiveness of the cyber security standard for electronic banking services should be measured and periodically evaluated.',
  N'ينبغي قياس وتقييم فعالية معيار الأمن السيبراني للخدمات المصرفية الإلكترونية بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_118,
  2,
  N'3.3.13.4',
  N'Electronic banking services security standard should cover:
a. use of brand protection measures to protect online services including social media.
b. online, mobile and phone banking:
1. use of official application stores and websites (applicable for online and mobile banking);
2. use of detection measures and take-down of malicious apps and websites (applicable for online and mobile banking);
3. use of sandboxing (applicable for online and mobile banking);
4. use of non-caching techniques (applicable for online and mobile banking);
5. use of communication techniques to avoid ‘man-in-the-middle’-attacks (applicable for online and mobile banking);
6. use of multi-factor authentication mechanisms:
a. multi-factor authentication should be used during the registration process for the customer in order to use of electronic banking services;
b. multi-factor authentication should be implemented for all electronic banking services available to customers;
c. the use of hard and soft tokens should be password protected;
d. revoking the access of customers after 3 successive incorrect passwords or invalid PINs;
e. the process for changing the customer mobile number should only be done from either a branch or ATM;
f. the processes for requesting and activating of the multi-factor authentication should be done through different delivery channels;
g. multi-factor authentication should be implemented for the following processes:
1. sign-on;
2. adding or modifying beneficiaries;
3. adding utility and government payment services;
4. high-risk transactions (when it exceeds predefined limits);
5. password reset;
7. the processes for adding and activating beneficiaries should be done through different delivery channels (applicable for mobile and online banking);
8. high availability of the electronic banking services should be ensured;
9. scheduled downtime of the electronic banking services should be timely communicated to SAMA and customers;
10. contractual agreements between the Member Organization and the customer addressing the roles, responsibilities and liabilities for both the Member Organization and the customers; 
11. obtaining approval of SAMA before launching a new electronic banking service.
c. ATMs and POSs:
1. prevention and detection of exploiting the ATM/POS application and infrastructure vulnerabilities (e.g., cables, (USB)-ports, rebooting);
2. cyber security measures, such as hardening of operating systems, malware protection, privacy
screens, masking of passwords or account numbers (e.g., screen and receipt), geo-blocking (e.g., disable cards per default for outside GCC countries, disable magnetic strip transactions), video monitoring (CCTV), revoking cards after 3 successive invalid PINs, anti-skimming
solutions (hardware/software), and PIN-pad protection;
3. remote stopping of ATMs in case of malicious activities.
d. SMS instant notification services:
1. SMS messages should not contain sensitive data (e.g., account balance - except for credit cards);
2. SMS alert should be sent to both mobile numbers (old and new) when the customer’s mobile number has been changed;
3. SMS notification should be sent to the customer’s mobile number when requesting a new multi-factor authentication mechanism.
4. SMS notification should be sent to the customer’s mobile number for all retail and personal financial transactions. 
5. SMS notification should be sent to the customer’s mobile number when beneficiaries are added, modified and activated. ',
  N'يجب أن يشمل معيار أمن الخدمات المصرفية الإلكترونية ما يلي:
أ. استخدام تدابير حماية العلامة التجارية لحماية الخدمات الإلكترونية، بما في ذلك وسائل التواصل الاجتماعي.
ب. الخدمات المصرفية الإلكترونية والهواتف المحمولة:
1. استخدام متاجر التطبيقات والمواقع الإلكترونية الرسمية (ينطبق على الخدمات المصرفية الإلكترونية والهواتف المحمولة)؛
2. استخدام تدابير الكشف عن التطبيقات والمواقع الإلكترونية الضارة وحذفها (ينطبق على الخدمات المصرفية الإلكترونية والهواتف المحمولة)؛
3. استخدام الحماية (ينطبق على الخدمات المصرفية الإلكترونية والهواتف المحمولة)؛
4. استخدام تقنيات عدم التخزين المؤقت (ينطبق على الخدمات المصرفية الإلكترونية والهواتف المحمولة)؛
5. استخدام تقنيات الاتصال لتجنب هجمات "الرجل في المنتصف" (ينطبق على الخدمات المصرفية الإلكترونية والهواتف المحمولة)؛
6. استخدام آليات المصادقة متعددة العوامل:
أ. يجب استخدام المصادقة متعددة العوامل أثناء عملية تسجيل العميل لاستخدام الخدمات المصرفية الإلكترونية؛
ب. يجب تطبيق المصادقة متعددة العوامل لجميع الخدمات المصرفية الإلكترونية المتاحة للعملاء؛
ج. يجب حماية استخدام الرموز الصلبة واللينة بكلمة مرور؛
د. إلغاء وصول العملاء بعد ثلاث كلمات مرور خاطئة متتالية أو أرقام تعريف شخصية غير صالحة؛
هـ. يجب أن تتم عملية تغيير رقم جوال العميل من الفرع أو جهاز الصراف الآلي فقط؛
و. يجب أن تتم عمليات طلب وتفعيل المصادقة متعددة العوامل من خلال قنوات توصيل مختلفة؛
ز. يجب تطبيق المصادقة متعددة العوامل للعمليات التالية:
1. تسجيل الدخول؛
2. إضافة أو تعديل المستفيدين؛
3. إضافة خدمات المرافق والمدفوعات الحكومية؛
4. المعاملات عالية المخاطر (عندما تتجاوز الحدود المحددة مسبقًا)؛
5. إعادة تعيين كلمة المرور؛
7. يجب أن تتم عمليات إضافة وتفعيل المستفيدين من خلال قنوات توصيل مختلفة (تنطبق على الخدمات المصرفية عبر الهاتف المحمول والإنترنت)؛
8. يجب ضمان توافر الخدمات المصرفية الإلكترونية بشكل كبير؛
9. يجب إبلاغ مؤسسة النقد العربي السعودي والعملاء في الوقت المناسب بتوقف الخدمات المصرفية الإلكترونية المجدولة؛
١٠. اتفاقيات تعاقدية بين المؤسسة العضوة والعميل تتناول الأدوار والمسؤوليات والالتزامات لكل من المؤسسة العضوة والعملاء؛
١١. الحصول على موافقة مؤسسة النقد العربي السعودي (ساما) قبل إطلاق خدمة مصرفية إلكترونية جديدة.
ج. أجهزة الصراف الآلي ونقاط البيع:
١. منع وكشف استغلال ثغرات تطبيقات أجهزة الصراف الآلي/نقاط البيع والبنية التحتية (مثل الكابلات، ومنافذ USB، وإعادة التشغيل)؛
٢. تدابير الأمن السيبراني، مثل تقوية أنظمة التشغيل، والحماية من البرامج الضارة، وشاشات الخصوصية، وإخفاء كلمات المرور أو أرقام الحسابات (مثل الشاشة والإيصال)، والحظر الجغرافي (مثل تعطيل البطاقات افتراضيًا خارج دول مجلس التعاون الخليجي، وتعطيل معاملات الشريط المغناطيسي)، ومراقبة الفيديو (كاميرات المراقبة)، وإلغاء البطاقات بعد ٣ أرقام تعريف شخصية غير صالحة متتالية، وحلول مكافحة التزييف (الأجهزة/البرمجيات)، وحماية لوحة أرقام التعريف الشخصية؛
٣. إيقاف أجهزة الصراف الآلي عن بُعد في حالة حدوث أنشطة ضارة.
د. خدمات الإشعارات الفورية عبر الرسائل النصية القصيرة:
1. يجب ألا تحتوي الرسائل النصية القصيرة على بيانات حساسة (مثل رصيد الحساب - باستثناء بطاقات الائتمان)؛
2. يجب إرسال تنبيهات الرسائل النصية القصيرة إلى رقمي الجوال (القديم والجديد) عند تغيير رقم جوال العميل؛
3. يجب إرسال إشعارات الرسائل النصية القصيرة إلى رقم جوال العميل عند طلب آلية مصادقة متعددة العوامل جديدة.
4. يجب إرسال إشعارات الرسائل النصية القصيرة إلى رقم جوال العميل لجميع المعاملات المالية الشخصية والتجارية.
5. يجب إرسال إشعارات الرسائل النصية القصيرة إلى رقم جوال العميل عند إضافة مستفيدين أو تعديلهم أو تفعيلهم.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_119,
  2,
  N'3.3.14.1',
  N'The security event management process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية إدارة الأحداث الأمنية والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_120,
  2,
  N'3.3.14.2',
  N'The effectiveness of the cyber security controls within the security event management process should be measured and periodically evaluated.',
  N'ينبغي قياس وتقييم فعالية ضوابط الأمن السيبراني ضمن عملية إدارة الأحداث الأمنية بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_121,
  2,
  N'3.3.14.3',
  N'To support this process a security event monitoring standard should be defined, approved and
implemented. 
a. the standard should address for all information assets the mandatory events which should be
monitored, based on the classification or risk profile of the information asset. ',
  N'لدعم هذه العملية، يجب تحديد معيار لمراقبة الأحداث الأمنية والموافقة عليه وتنفيذه. أ. يجب أن يتناول المعيار لجميع أصول المعلومات الأحداث الإلزامية التي يجب مراقبتها، بناءً على تصنيف أو ملف تعريف مخاطر أصل المعلومات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_122,
  2,
  N'3.3.14.4',
  N'The security event management process should include requirements for:
a. the establishment of a designated team responsible for security monitoring (i.e., Security Operations Center (SOC));
b. skilled and (continuously) trained staff;
c. a restricted area to facilitate SOC activities and workspaces;
d. resources required continuous security event monitoring activities (24x7);
e. detection and handling of malicious code and software;
f. detection and handling of security or suspicious events and anomalies;
g. deployment of security network packet analysis solution;
h. adequately protected logs;
i. periodic compliance monitoring of applications and infrastructure cyber security standards
j. automated and centralized analysis of security loggings and correlation of event or patterns (i.e., Security Information and Event Management (SIEM));
k. reporting of cyber security incidents;
l. independent periodic testing of the effectiveness of the security operations center (e.g., red teaming). ',
  N'ينبغي أن تتضمن عملية إدارة الأحداث الأمنية متطلباتٍ لما يلي:
أ. إنشاء فريق مُعيَّن مسؤول عن مراقبة الأمن (أي مركز عمليات الأمن (SOC))؛
ب. موظفون ماهرون ومدربون (بشكلٍ مستمر)؛
ج. منطقة مُقيَّدة لتسهيل أنشطة مركز عمليات الأمن ومساحات العمل؛
د. الموارد اللازمة لأنشطة مراقبة الأحداث الأمنية المُستمرة (على مدار الساعة طوال أيام الأسبوع)؛
هـ. الكشف عن الأكواد البرمجية والبرمجيات الضارة ومعالجتها؛
و. الكشف عن الأحداث والشذوذ الأمني ​​أو المُريب ومعالجته؛
ز. نشر حلول تحليل حزم شبكة الأمان؛
ح. سجلات محمية بشكلٍ كافٍ؛
ط. مراقبة الامتثال الدوري لمعايير الأمن السيبراني للتطبيقات والبنية التحتية
ي. تحليل آلي ومركزي لسجلات الأمن وربط الأحداث أو الأنماط (أي إدارة معلومات الأمن والأحداث (SIEM))؛
ك. الإبلاغ عن حوادث الأمن السيبراني؛
ل. اختبار دوري مُستقل لفعالية مركز عمليات الأمن (مثل الفريق الأحمر).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_123,
  2,
  N'3.3.15.1',
  N'The cyber security incident management process should be defined, approved, implemented and aligned with the enterprise incident management process.',
  N'ينبغي تحديد عملية إدارة حوادث الأمن السيبراني والموافقة عليها وتنفيذها ومواءمتها مع عملية إدارة الحوادث في المؤسسة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_124,
  2,
  N'3.3.15.2',
  N'The effectiveness of the cyber security controls within the cyber security incident management process should be measured and periodically evaluated.',
  N'ينبغي قياس وتقييم فعالية ضوابط الأمن السيبراني ضمن عملية إدارة حوادث الأمن السيبراني بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_125,
  2,
  N'3.3.15.3',
  N'The standard should address the mandatory and suspicious security events which should be responded to.',
  N'ينبغي أن يتناول المعيار الأحداث الأمنية الإلزامية والمريبة التي يجب الاستجابة لها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_126,
  2,
  N'3.3.15.4',
  N'The security incident management process should include requirements for:
a. the establishment of a designated team responsible for security incident management;
b. skilled and (continuously) trained staff;
c. sufficient capacity available of certified forensic staff for handling major incidents (e.g., internal staff or contracting an external forensic team);
d. a restricted area to facilitate the computer emergency response team (CERT) workspaces;
e. the classification of cyber security incidents;
f. the timely handling of cyber security incidents, recording and monitoring progress;
g. the protection of relevant evidence and loggings;
h. post-incident activities, such as forensics, root-cause analysis of the incidents; 
i. reporting of suggested improvements to the CISO and the Committee;
j. establish a cyber security incident repository. ',
  N'ينبغي أن تتضمن عملية إدارة الحوادث الأمنية متطلباتٍ لما يلي:
أ. إنشاء فريق مُعيَّن مسؤول عن إدارة الحوادث الأمنية؛
ب. موظفون ماهرون ومدربون (بشكلٍ مستمر)؛
ج. توافر قدرات كافية من موظفي الطب الشرعي المعتمدين للتعامل مع الحوادث الكبرى (مثل الموظفين الداخليين أو التعاقد مع فريق طب شرعي خارجي)؛
د. منطقة مُقيَّدة لتسهيل عمل فريق الاستجابة للطوارئ الحاسوبية (CERT)؛
هـ. تصنيف حوادث الأمن السيبراني؛
و. التعامل مع حوادث الأمن السيبراني في الوقت المناسب، وتسجيل التقدم ومراقبته؛
ز. حماية الأدلة والسجلات ذات الصلة؛
ح. أنشطة ما بعد الحادث، مثل الطب الشرعي، وتحليل السبب الجذري للحوادث؛
ط. الإبلاغ عن التحسينات المُقترحة إلى مسؤول أمن المعلومات واللجنة؛
ي. إنشاء مستودع لحوادث الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_127,
  2,
  N'3.3.15.5',
  N'The Member Organization should inform ‘SAMA IT Risk Supervision’ immediately when a medium or high classified security incident has occurred and identified.',
  N'يجب على المنظمة العضو إبلاغ "إدارة مخاطر تكنولوجيا المعلومات في مؤسسة النقد العربي السعودي" على الفور عند وقوع حادثة أمنية متوسطة أو عالية التصنيف وتحديدها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_128,
  2,
  N'3.3.15.6',
  N'The Member Organization should obtain ‘no objection’ from ‘SAMA IT Risk Supervision’ before any media interaction related to the incident.',
  N'يجب على المنظمة العضو الحصول على "عدم اعتراض" من "إدارة الإشراف على مخاطر تكنولوجيا المعلومات في مؤسسة النقد العربي السعودي" قبل أي تفاعل إعلامي يتعلق بالحادث.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_129,
  2,
  N'3.3.15.7',
  N'The Member Organization should submit a formal incident report ‘SAMA IT Risk Supervision’ after
resuming operations, including the following incident details:
a. title of incident;
b. classification of the incident (medium or high);
c. date and time of incident occurred;
d. date and time of incident detected;
e. information assets involved;
f. (technical) details of the incident; 
g. root-cause analysis;
h. corrective activities performed and planned;
i. description of impact (e.g., loss of data, disruption of services, unauthorized modification of data, (un)intended data leakage, number of customers impacted);
j. total estimated cost of incident;
k. estimated cost of corrective actions. ',
  N'يجب على المنظمة العضو تقديم تقرير رسمي عن الحادث "إشراف ساما على مخاطر تكنولوجيا المعلومات" بعد استئناف العمليات، متضمنًا تفاصيل الحادث التالية:
أ. عنوان الحادث؛
ب. تصنيف الحادث (متوسط ​​أو مرتفع)؛
ج. تاريخ ووقت وقوع الحادث؛
د. تاريخ ووقت اكتشاف الحادث؛
هـ. أصول المعلومات المعنية؛
و. التفاصيل (الفنية) للحادث؛
ز. تحليل السبب الجذري؛
ح. الأنشطة التصحيحية التي تم تنفيذها والمخطط لها؛
ط. وصف التأثير (مثل: فقدان البيانات، انقطاع الخدمات، التعديل غير المصرح به للبيانات، تسرب البيانات (غير المقصود)، عدد العملاء المتأثرين)؛
ي. التكلفة الإجمالية المقدرة للحادث؛
ك. التكلفة المقدرة للإجراءات التصحيحية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_130,
  2,
  N'3.3.16.1',
  N'The threat intelligence management process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية إدارة معلومات التهديد والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_131,
  2,
  N'3.3.16.2',
  N'The effectiveness of the threat intelligence management process should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية عملية إدارة معلومات التهديدات وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_132,
  2,
  N'3.3.16.3',
  N'The threat intelligence management process should include:
a. the use of internal sources, such as access control, application and infrastructure logs, IDS, IPS,
security tooling, Security Information and Event Monitoring (SIEM), support functions (e.g., Legal,
Audit, IT Helpdesk, Forensics, Fraud Management, Risk Management, Compliance);
b. the use of reliable and relevant external sources, such as SAMA, government agencies, security forums, (security) vendors, security organizations and specialist notification services;
c. a defined methodology to analyze the threat information periodically;
d. the relevant details on identified or collected threats, such as modus operandi, actors, motivation and type of threats;
e. the relevance of the derived intelligence and the action-ability for follow-up (for e.g., SOC, Risk Management);
f. sharing the relevant intelligence with the relevant stakeholders (e.g., SAMA, BCIS members).',
  N'ينبغي أن تشمل عملية إدارة معلومات التهديدات ما يلي:
أ. استخدام المصادر الداخلية، مثل التحكم في الوصول، وسجلات التطبيقات والبنية التحتية، وأنظمة كشف التسلل (IDS)، وأنظمة منع التطفل (IPS)، وأدوات الأمن، ومراقبة معلومات الأمن والأحداث (SIEM)، ووظائف الدعم (مثل: الشؤون القانونية، والتدقيق، ومركز دعم تكنولوجيا المعلومات، والتحليل الجنائي، وإدارة الاحتيال، وإدارة المخاطر، والامتثال)؛
ب. استخدام مصادر خارجية موثوقة وذات صلة، مثل هيئة النقد السعودية (SAMA)، والهيئات الحكومية، ومنتديات الأمن، وموردي (الأمن)، والمنظمات الأمنية، وخدمات الإخطار المتخصصة؛
ج. منهجية محددة لتحليل معلومات التهديدات بشكل دوري؛
د. التفاصيل ذات الصلة بالتهديدات المحددة أو المجمعة، مثل أسلوب العمل، والجهات الفاعلة، والدوافع، ونوع التهديدات؛
هـ. أهمية المعلومات الاستخباراتية المُستقاة وإمكانية اتخاذ إجراءات بشأنها (مثل: مركز العمليات الأمنية (SOC)، وإدارة المخاطر)؛
و. مشاركة المعلومات الاستخباراتية ذات الصلة مع الجهات المعنية (مثل: هيئة النقد السعودية (SAMA)، وأعضاء نظام معلومات الأمن (BCIS)).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_133,
  2,
  N'3.3.17.1',
  N'The vulnerability management process should be defined, approved and implemented.',
  N'ينبغي تحديد عملية إدارة الثغرات الأمنية والموافقة عليها وتنفيذها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_134,
  2,
  N'3.3.17.2',
  N'The effectiveness of the vulnerability management process should be measured and periodically evaluated.',
  N'ينبغي قياس فعالية عملية إدارة الثغرات وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_135,
  2,
  N'3.3.17.3',
  N'The vulnerability management process should include:
a. all information assets;
b. frequency of performing the vulnerability scan (risk-based);
c. classification of vulnerabilities;
d. defined timelines to mitigate (per classification);
e. prioritization for classified information assets;
f. patch management and method of deployment. ',
  N'ينبغي أن تشمل عملية إدارة الثغرات الأمنية ما يلي:
أ. جميع أصول المعلومات؛
ب. وتيرة إجراء فحص الثغرات الأمنية (بناءً على المخاطر)؛
ج. تصنيف الثغرات الأمنية؛
د. جداول زمنية محددة للتخفيف من حدتها (لكل تصنيف)؛
هـ. تحديد أولويات أصول المعلومات السرية؛
و. إدارة التصحيحات وطريقة النشر.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_136,
  2,
  N'3.4.1.1',
  N'The cyber security requirements should be defined, approved, implemented and communicated within the contract and vendor management processes. ',
  N'ينبغي تحديد متطلبات الأمن السيبراني والموافقة عليها وتنفيذها وتوصيلها ضمن عمليات إدارة العقود والبائعين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_137,
  2,
  N'3.4.1.2',
  N'The compliance with contract and vendor management process should be monitored.',
  N'ينبغي مراقبة الامتثال لعملية إدارة العقود والبائعين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_138,
  2,
  N'3.4.1.3',
  N'The effectiveness of the cyber security controls within the contract and vendor management process should be measured and periodically evaluated.',
  N'ينبغي قياس وتقييم فعالية ضوابط الأمن السيبراني ضمن عملية إدارة العقود والموردين بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_139,
  2,
  N'3.4.1.4',
  N'These contract and vendor management processes should cover:
a. whether the involvement of the cyber security function is actively required (e.g., in case of due diligence);
b. the baseline cyber security requirements which should be applied in all cases;
c. the right to periodically perform cyber security reviews and audits.',
  N'ينبغي أن تشمل عمليات إدارة العقود والموردين ما يلي:
أ. ما إذا كانت مشاركة وظيفة الأمن السيبراني مطلوبة بشكل نشط (على سبيل المثال، في حالة العناية الواجبة)؛
ب. متطلبات الأمن السيبراني الأساسية التي ينبغي تطبيقها في جميع الحالات؛
ج. الحق في إجراء مراجعات وتدقيقات دورية للأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_140,
  2,
  N'3.4.1.5',
  N'The contract management process should cover requirements for:
a. executing a cyber security risk assessment as part of the procurement process;
b. defining the specific cyber security requirements as part of the tender process;
c. evaluating the replies of potential vendors on the defined cyber security requirements;
d. testing of the agreed cyber security requirements (risk-based);
e. defining the communication or escalation process in case of cyber security incidents;
f. ensuring cyber security requirements are defined for exiting, terminating or renewing the contract (including escrow agreements if applicable);
g. defining a mutual confidentiality agreement. ',
  N'ينبغي أن تشمل عملية إدارة العقود متطلبات ما يلي:
أ. إجراء تقييم لمخاطر الأمن السيبراني كجزء من عملية الشراء؛
ب. تحديد متطلبات الأمن السيبراني المحددة كجزء من عملية تقديم العطاءات؛
ج. تقييم ردود البائعين المحتملين على متطلبات الأمن السيبراني المحددة؛
د. اختبار متطلبات الأمن السيبراني المتفق عليها (على أساس المخاطر)؛
هـ. تحديد عملية الاتصال أو التصعيد في حالة وقوع حوادث أمن سيبراني؛
و. ضمان تحديد متطلبات الأمن السيبراني لإنهاء العقد أو إنهائه أو تجديده (بما في ذلك اتفاقيات الضمان إن وجدت)؛
ز. تحديد اتفاقية سرية متبادلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_141,
  2,
  N'3.4.1.6',
  N'The vendor management process (i.e., service level management) should cover requirements for:
a. periodic reporting, reviewing and evaluating the contractually agreed cyber security requirements (in SLAs). ',
  N'يجب أن تغطي عملية إدارة البائعين (أي إدارة مستوى الخدمة) متطلبات ما يلي: أ. إعداد التقارير الدورية ومراجعة وتقييم متطلبات الأمن السيبراني المتفق عليها تعاقديًا (في اتفاقيات مستوى الخدمة).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_142,
  2,
  N'3.4.2.1',
  N'The cyber security requirements within the outsourcing policy and process should be defined, approved, implemented and communicated within Member Organization. ',
  N'يجب تحديد متطلبات الأمن السيبراني ضمن سياسة وعملية الاستعانة بمصادر خارجية، والموافقة عليها وتنفيذها وتوصيلها داخل المنظمة العضو.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_143,
  2,
  N'3.4.2.2',
  N'The cyber security requirements regarding the outsourcing policy and process should be measured
and periodically evaluated. ',
  N'ينبغي قياس متطلبات الأمن السيبراني المتعلقة بسياسة وعملية الاستعانة بمصادر خارجية وتقييمها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_144,
  2,
  N'3.4.2.3',
  N'The outsourcing process should include:
a. the approval from SAMA prior to material outsourcing; 
b. the involvement of the cyber security function;
c. compliance with the SAMA circular on outsourcing. ',
  N'يجب أن تشمل عملية الاستعانة بمصادر خارجية ما يلي:
أ. موافقة مؤسسة النقد العربي السعودي قبل الاستعانة بمصادر خارجية للمواد؛
ب. إشراك قسم الأمن السيبراني؛
ج. الامتثال لتعميم مؤسسة النقد العربي السعودي بشأن الاستعانة بمصادر خارجية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_145,
  2,
  N'3.4.3.1',
  N'The cyber security controls within the cloud computing policy for hybrid and public cloud services
should be defined, approved and implemented and communicated within Member Organization. ',
  N'ينبغي تحديد ضوابط الأمن السيبراني ضمن سياسة الحوسبة السحابية للخدمات السحابية الهجينة والعامة، والموافقة عليها وتنفيذها وتوصيلها داخل المنظمة العضو.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_146,
  2,
  N'3.4.3.2',
  N'The compliance with the cloud computing policy should be monitored.',
  N'ينبغي مراقبة الامتثال لسياسة الحوسبة السحابية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_147,
  2,
  N'3.4.3.3',
  N'The cyber security controls regarding the cloud computing policy and process for hybrid and public cloud services should be periodically measured and evaluated. ',
  N'ينبغي قياس وتقييم ضوابط الأمن السيبراني المتعلقة بسياسة الحوسبة السحابية وعملية الخدمات السحابية الهجينة والعامة بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
) VALUES (
  @cl_148,
  2,
  N'3.4.3.4',
  N'The cloud computing policy for hybrid and public cloud services should address requirements for:
a. the process for adopting cloud services, including that:
1. a cyber security risk assessment and due diligence on the cloud service provider and its cloud services should be performed;
2. the Member Organization should obtain SAMA approval prior to using cloud services or signing the contract with the cloud provider;
3. a contract should be in place, including the cyber security requirements, before using cloud services;
b. data location, including that:
1. in principle only cloud services should be used that are located in Saudi Arabia, or when cloud services are to be used outside Saudi Arabia that the Member Organization should obtain explicit approval from SAMA;
c. data use limitations, including that:
1. the cloud service provider should not use the Member Organization’s data for secondary purposes;
d. security, including that:
1. the cloud service provider should implement and monitor the cyber security controls as determined in the risk assessment for protecting the confidentiality, integrity and availability of the Member Organization’s data;
e. data segregation, including that:
1. the Member Organization’s data is logically segregated from other data held by the cloud service provider, including that the cloud service provider should be able to identify the Member Organization’s data and at all times should be able to distinguish it from other data.
f. business continuity, including that:
1. business continuity requirements are met in accordance with the Member Organization’s business continuity policy;
g. audit, review and monitoring, including that:
1. the Member Organization has the right to perform a cyber security review at the cloud service provider;
2. the Member Organization has the right to perform a cyber security audit at the cloud service provider;
3. the Member Organization has the right to perform a cyber security examination at the cloud service provider;
h. exit, including that:
1. the Member Organization has termination rights;
2. the cloud service provider has to return the Member Organization’s data on termination;
3. the cloud service provider has to irreversibly delete the Member Organization’s data on termination.',
  N'ينبغي أن تتناول سياسة الحوسبة السحابية للخدمات السحابية الهجينة والعامة متطلبات ما يلي:
أ. عملية اعتماد الخدمات السحابية، بما في ذلك:
1. إجراء تقييم لمخاطر الأمن السيبراني والعناية الواجبة على مزود الخدمة السحابية وخدماته السحابية؛
2. حصول المنظمة العضوة على موافقة مؤسسة النقد العربي السعودي (SAMA) قبل استخدام الخدمات السحابية أو توقيع العقد مع مزود الخدمة السحابية؛
3. إبرام عقد، بما في ذلك متطلبات الأمن السيبراني، قبل استخدام الخدمات السحابية؛
ب. تحديد موقع البيانات، بما في ذلك:
1. من حيث المبدأ، ينبغي استخدام الخدمات السحابية الموجودة في المملكة العربية السعودية فقط، أو عند استخدام الخدمات السحابية خارج المملكة العربية السعودية، ينبغي للمنظمة العضوة الحصول على موافقة صريحة من مؤسسة النقد العربي السعودي؛
ج. قيود استخدام البيانات، بما في ذلك:
1. عدم استخدام مزود الخدمة السحابية لبيانات المنظمة العضوة لأغراض ثانوية؛
د. الأمن، بما في ذلك:
1. تنفيذ ومراقبة ضوابط الأمن السيبراني كما هو محدد في تقييم المخاطر لحماية سرية بيانات المنظمة العضوة وسلامتها وتوافرها؛
هـ. فصل البيانات، بما في ذلك:
1. فصل بيانات المنظمة العضو منطقيًا عن البيانات الأخرى التي يحتفظ بها مزود خدمة السحابة، بما في ذلك أن يكون مزود خدمة السحابة قادرًا على تحديد بيانات المنظمة العضو وأن يكون قادرًا في جميع الأوقات على تمييزها عن البيانات الأخرى.
و. استمرارية الأعمال، بما في ذلك:
1. استيفاء متطلبات استمرارية الأعمال وفقًا لسياسة استمرارية الأعمال الخاصة بالمنظمة العضو؛
ز. التدقيق والمراجعة والمراقبة، بما في ذلك:
1. يحق للمنظمة العضو إجراء مراجعة للأمن السيبراني لدى مزود خدمة السحابة؛
2. يحق للمنظمة العضو إجراء تدقيق للأمن السيبراني لدى مزود خدمة السحابة؛
3. يحق للمنظمة العضو إجراء فحص للأمن السيبراني لدى مزود خدمة السحابة؛
ح. الخروج، بما في ذلك:
1. تتمتع المنظمة العضو بحقوق الإنهاء؛
2. يتعين على مزود خدمة السحابة إعادة بيانات المنظمة العضو عند الإنهاء؛
3. يتعين على مزود الخدمة السحابية حذف بيانات المنظمة العضو بشكل نهائي عند إنهاء الخدمة.',
  GETDATE(), GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_library_titles;
-- +goose StatementEnd
