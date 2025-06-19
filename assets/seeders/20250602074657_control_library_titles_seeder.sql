-- +goose Up
-- +goose StatementBegin
DECLARE @compliance_library_id INT;
SELECT @compliance_library_id = id FROM compliance_libraries
WHERE title_en = 'Essential Cybersecurity Controls (ECC – 2: 2024)';

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
DECLARE @cl_149 INT;
SELECT @cl_149 = id FROM control_library WHERE title_en = N'Identify Monitoring Technologies';
DECLARE @cl_150 INT;
SELECT @cl_150 = id FROM control_library WHERE title_en = N'Monitor Cybersecurity Events Continuously';
DECLARE @cl_151 INT;
SELECT @cl_151 = id FROM control_library WHERE title_en = N'Retain Event Logs for Minimum Period';
DECLARE @cl_152 INT;
SELECT @cl_152 = id FROM control_library WHERE title_en = N'Review Event Logging and Monitoring Requirements';
DECLARE @cl_153 INT;
SELECT @cl_153 = id FROM control_library WHERE title_en = N'Define Cybersecurity Incident Management Requirements';
DECLARE @cl_154 INT;
SELECT @cl_154 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Incident Management Requirements';
DECLARE @cl_155 INT;
SELECT @cl_155 = id FROM control_library WHERE title_en = N'Establish Cybersecurity Incident Management Controls';
DECLARE @cl_156 INT;
SELECT @cl_156 = id FROM control_library WHERE title_en = N'Develop Escalation-Based Response Plans';
DECLARE @cl_157 INT;
SELECT @cl_157 = id FROM control_library WHERE title_en = N'Classify Cybersecurity Incidents';
DECLARE @cl_158 INT;
SELECT @cl_158 = id FROM control_library WHERE title_en = N'Report Cybersecurity Incidents';
DECLARE @cl_159 INT;
SELECT @cl_159 = id FROM control_library WHERE title_en = N'Share Threat Intelligence';
DECLARE @cl_160 INT;
SELECT @cl_160 = id FROM control_library WHERE title_en = N'Collect Threat Intelligence Feeds';
DECLARE @cl_161 INT;
SELECT @cl_161 = id FROM control_library WHERE title_en = N'Review Incident Management Requirements Periodically';
DECLARE @cl_162 INT;
SELECT @cl_162 = id FROM control_library WHERE title_en = N'Define Physical Security Requirements';
DECLARE @cl_163 INT;
SELECT @cl_163 = id FROM control_library WHERE title_en = N'Implement Physical Security Requirements';
DECLARE @cl_164 INT;
SELECT @cl_164 = id FROM control_library WHERE title_en = N'Establish Physical Security Controls';
DECLARE @cl_165 INT;
SELECT @cl_165 = id FROM control_library WHERE title_en = N'Control Access to Sensitive Areas';
DECLARE @cl_166 INT;
SELECT @cl_166 = id FROM control_library WHERE title_en = N'Monitor Facility Access';
DECLARE @cl_167 INT;
SELECT @cl_167 = id FROM control_library WHERE title_en = N'Protect Facility Access Records';

DECLARE @cl_168 INT;
SELECT @cl_168 = id FROM control_library WHERE title_en = N'Secure Destruction of Physical Assets';

DECLARE @cl_169 INT;
SELECT @cl_169 = id FROM control_library WHERE title_en = N'Secure Devices Inside and Outside Facilities';

DECLARE @cl_170 INT;
SELECT @cl_170 = id FROM control_library WHERE title_en = N'Review Physical Security Requirements Periodically';

DECLARE @cl_171 INT;
SELECT @cl_171 = id FROM control_library WHERE title_en = N'Define Web Application Security Requirements';

DECLARE @cl_172 INT;
SELECT @cl_172 = id FROM control_library WHERE title_en = N'Implement Web Application Security Requirements';

DECLARE @cl_173 INT;
SELECT @cl_173 = id FROM control_library WHERE title_en = N'Establish Web Application Security Controls';

DECLARE @cl_174 INT;
SELECT @cl_174 = id FROM control_library WHERE title_en = N'Use Web Application Firewall';

DECLARE @cl_175 INT;
SELECT @cl_175 = id FROM control_library WHERE title_en = N'Adopt Multi-Tier Architecture';

DECLARE @cl_176 INT;
SELECT @cl_176 = id FROM control_library WHERE title_en = N'Use Secure Protocols';

DECLARE @cl_177 INT;
SELECT @cl_177 = id FROM control_library WHERE title_en = N'Define Secure Usage Policies';

DECLARE @cl_178 INT;
SELECT @cl_178 = id FROM control_library WHERE title_en = N'Implement Risk-Based Authentication';

DECLARE @cl_179 INT;
SELECT @cl_179 = id FROM control_library WHERE title_en = N'Review Web Application Security Requirements Periodically';

DECLARE @cl_180 INT;
SELECT @cl_180 = id FROM control_library WHERE title_en = N'Define Business Continuity Requirements';

DECLARE @cl_181 INT;
SELECT @cl_181 = id FROM control_library WHERE title_en = N'Implement Business Continuity Requirements';

DECLARE @cl_182 INT;
SELECT @cl_182 = id FROM control_library WHERE title_en = N'Establish Business Continuity Controls';

DECLARE @cl_183 INT;
SELECT @cl_183 = id FROM control_library WHERE title_en = N'Ensure Continuity of Cybersecurity Systems';

DECLARE @cl_184 INT;
SELECT @cl_184 = id FROM control_library WHERE title_en = N'Develop Incident Response Plans';

DECLARE @cl_185 INT;
SELECT @cl_185 = id FROM control_library WHERE title_en = N'Develop Disaster Recovery Plans';

DECLARE @cl_186 INT;
SELECT @cl_186 = id FROM control_library WHERE title_en = N'Review Business Continuity Requirements Periodically';

DECLARE @cl_187 INT;
SELECT @cl_187 = id FROM control_library WHERE title_en = N'Define Third-Party Cybersecurity Requirements';

DECLARE @cl_188 INT;
SELECT @cl_188 = id FROM control_library WHERE title_en = N'Establish Third-Party Cybersecurity Controls';

DECLARE @cl_189 INT;
SELECT @cl_189 = id FROM control_library WHERE title_en = N'Include Non-Disclosure Clauses';

DECLARE @cl_190 INT;
SELECT @cl_190 = id FROM control_library WHERE title_en = N'Define Incident Communication Procedures';

DECLARE @cl_191 INT;
SELECT @cl_191 = id FROM control_library WHERE title_en = N'Enforce Compliance with Cybersecurity Policies';

DECLARE @cl_192 INT;
SELECT @cl_192 = id FROM control_library WHERE title_en = N'Establish IT Outsourcing Cybersecurity Controls';

DECLARE @cl_193 INT;
SELECT @cl_193 = id FROM control_library WHERE title_en = N'Assess Risks for Outsourcing Services';

DECLARE @cl_194 INT;
SELECT @cl_194 = id FROM control_library WHERE title_en = N'Ensure Local Cybersecurity Service Centers';

DECLARE @cl_195 INT;
SELECT @cl_195 = id FROM control_library WHERE title_en = N'Review Third-Party Cybersecurity Requirements Periodically';

DECLARE @cl_196 INT;
SELECT @cl_196 = id FROM control_library WHERE title_en = N'Define Cloud Computing Security Requirements';

DECLARE @cl_197 INT;
SELECT @cl_197 = id FROM control_library WHERE title_en = N'Implement Cloud Computing Security Requirements';

DECLARE @cl_198 INT;
SELECT @cl_198 = id FROM control_library WHERE title_en = N'Establish Cloud Computing Security Controls';

DECLARE @cl_199 INT;
SELECT @cl_199 = id FROM control_library WHERE title_en = N'Classify Data for Cloud Hosting';

DECLARE @cl_200 INT;
SELECT @cl_200 = id FROM control_library WHERE title_en = N'Separate Cloud Environments';

DECLARE @cl_201 INT;
SELECT @cl_201 = id FROM control_library WHERE title_en = N'Review Cloud Computing Security Requirements Periodically';


INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_1,
  @compliance_library_id,
  N'1.1.1',
  N'A cybersecurity strategy must be defined, documented, and approved. It must be supported by the head of the organization or his/her delegate (referred to in this document as Authorizing Official). The strategy goals must be in line with related laws and regulations.',
  N'يجب تحديد استراتيجية الأمن السيبراني وتوثيقها واعتمادها. ويجب أن يدعمها رئيس المؤسسة أو من ينوب عنه (المشار إليه في هذه الوثيقة بالمسؤول المُصرّح له). ويجب أن تتوافق أهداف الاستراتيجية مع القوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_2,
  @compliance_library_id,
  N'1.1.2',
  N'A roadmap must be executed to implement the cybersecurity strategy.',
  N'يجب تنفيذ خارطة طريق لتطبيق استراتيجية الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_3,
  @compliance_library_id,
  N'1.1.3',
  N'The cybersecurity strategy must be reviewed periodically according to planned intervals or upon changes to related laws and regulations.',
  N'يجب مراجعة استراتيجية الأمن السيبراني بشكل دوري وفقًا لفترات مخططة أو عند حدوث تغييرات في القوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_4,
  @compliance_library_id,
  N'1.2.1',
  N'A dedicated cybersecurity function (e.g., division, department) must be established within the organization. This function must be independent from the Information Technology/Information Communication and Technology (IT/ICT) functions (as per the Royal Decree number 37140 dated 14/8/1438H). It is highly recommended that this cybersecurity function reports directly to the head of the organization or his/her delegate while ensuring that this does not result in a conflict of interest.',
  N'يجب إنشاء إدارة متخصصة للأمن السيبراني (مثل قسم أو إدارة) داخل المؤسسة. يجب أن تكون هذه الإدارة مستقلة عن إدارة تقنية المعلومات/اتصالات وتكنولوجيا المعلومات (وفقًا للمرسوم الملكي رقم 37140 بتاريخ 14/8/1438هـ). يُوصى بشدة بأن ترفع إدارة الأمن السيبراني تقاريرها مباشرةً إلى رئيس المؤسسة أو من ينوب عنه، مع ضمان عدم نشوء أي تضارب في المصالح.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_5,
  @compliance_library_id,
  N'1.2.2',
  N'All cybersecurity positions must be filled with full-time and qualified Saudi cybersecurity professionals.',
  N'يجب أن يتم شغل جميع وظائف الأمن السيبراني من قبل متخصصين سعوديين مؤهلين بدوام كامل في مجال الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_6,
  @compliance_library_id,
  N'1.2.3',
  N'A cybersecurity steering committee must be established by the Authorizing Official to ensure the support and implementation of the cybersecurity.',
  N'يجب على المسؤول المخول إنشاء لجنة توجيهية للأمن السيبراني لضمان دعم وتنفيذ الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_7,
  @compliance_library_id,
  N'1.3.1',
  N'Cybersecurity policies and procedures must be defined and documented by the cybersecurity function, approved by the Authorizing Official, and disseminated to relevant parties inside and outside the organization.',
  N'يجب أن يتم تحديد سياسات وإجراءات الأمن السيبراني وتوثيقها من قبل وظيفة الأمن السيبراني، والموافقة عليها من قبل المسؤول المخول، ونشرها على الأطراف ذات الصلة داخل وخارج المنظمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_8,
  @compliance_library_id,
  N'1.3.2',
  N'The cybersecurity function must ensure that the cybersecurity policies and procedures are implemented.',
  N'يجب على وظيفة الأمن السيبراني ضمان تنفيذ سياسات وإجراءات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_9,
  @compliance_library_id,
  N'1.3.3',
  N'The cybersecurity policies and procedures must be supported by technical security standards (e.g., operating systems, databases, and firewall technical security standards).',
  N'يجب أن تكون سياسات وإجراءات الأمن السيبراني مدعومة بمعايير الأمان الفنية (على سبيل المثال، أنظمة التشغيل وقواعد البيانات ومعايير الأمان الفنية لجدران الحماية).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_10,
  @compliance_library_id,
  N'1.3.4',
  N'The cybersecurity policies and procedures must be reviewed periodically according to planned intervals or upon changes to related laws and regulations. Changes and reviews must be approved and documented.',
  N'يجب مراجعة سياسات وإجراءات الأمن السيبراني دوريًا وفقًا للفترات الزمنية المخطط لها أو عند حدوث أي تغييرات في القوانين واللوائح ذات الصلة. ويجب اعتماد هذه التغييرات والمراجعات وتوثيقها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_11,
  @compliance_library_id,
  N'1.4.1',
  N'Cybersecurity organizational structure and related roles and responsibilities must be defined, documented, approved, supported, and assigned by the Authorizing Official, while ensuring that this does not result in a conflict of interest.',
  N'يجب تحديد الهيكل التنظيمي للأمن السيبراني والأدوار والمسؤوليات ذات الصلة وتوثيقها والموافقة عليها ودعمها وتعيينها من قبل المسؤول المخول، مع ضمان عدم تسبب ذلك في تضارب في المصالح.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_12,
  @compliance_library_id,
  N'1.4.2',
  N'The cybersecurity roles and responsibilities must be reviewed periodically according to planned intervals or upon changes to related laws and regulations.',
  N'يجب مراجعة أدوار ومسؤوليات الأمن السيبراني بشكل دوري وفقًا لفترات مخططة أو عند حدوث تغييرات في القوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_13,
  @compliance_library_id,
  N'1.5.1',
  N'Cybersecurity risk management methodology and procedures must be defined, documented, and approved as per confidentiality, integrity, and availability considerations of information and technology assets.',
  N'يجب تحديد منهجية وإجراءات إدارة مخاطر الأمن السيبراني وتوثيقها والموافقة عليها وفقًا لاعتبارات السرية والنزاهة وتوافر أصول المعلومات والتكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_14,
  @compliance_library_id,
  N'1.5.2',
  N'The cybersecurity risk management methodology and procedures must be implemented by the cybersecurity function.',
  N'يجب على وظيفة الأمن السيبراني تنفيذ منهجية وإجراءات إدارة مخاطر الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_15,
  @compliance_library_id,
  N'1.5.3',
  N'The cybersecurity risk assessment procedures must be implemented at least in the following cases',
  N'يجب تنفيذ إجراءات تقييم مخاطر الأمن السيبراني على الأقل في الحالات التالية',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_16,
  @compliance_library_id,
  N'1.5.3.1',
  N'Early stages of technology projects',
  N'المراحل المبكرة من مشاريع التكنولوجيا',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_17,
  @compliance_library_id,
  N'1.5.3.2',
  N'Before making major changes to technology infrastructure',
  N'قبل إجراء تغييرات كبيرة على البنية التحتية للتكنولوجيا',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_18,
  @compliance_library_id,
  N'1.5.3.3',
  N'During the planning phase of obtaining third-party services',
  N'أثناء مرحلة التخطيط للحصول على خدمات الطرف الثالث',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_19,
  @compliance_library_id,
  N'1.5.3.4',
  N'During the planning phase and before going live for new technology services and products.',
  N'أثناء مرحلة التخطيط وقبل إطلاق الخدمات والمنتجات التكنولوجية الجديدة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_20,
  @compliance_library_id,
  N'1.5.4',
  N'The cybersecurity risk management methodology and procedures must be reviewed periodically according to planned intervals or upon changes to related laws and regulations. Changes and reviews must be approved and documented.',
  N'يجب مراجعة منهجية وإجراءات إدارة مخاطر الأمن السيبراني دوريًا وفقًا للفترات الزمنية المخطط لها أو عند حدوث أي تغييرات في القوانين واللوائح ذات الصلة. ويجب اعتماد هذه التغييرات والمراجعات وتوثيقها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_21,
  @compliance_library_id,
  N'1.6.1',
  N'Cybersecurity requirements must be included in project and asset (information/technology) change management methodology and procedures to identify and manage cybersecurity risks as part of the project management lifecycle. The cybersecurity requirements must be a key part of the overall requirements of technology projects',
  N'يجب تضمين متطلبات الأمن السيبراني في منهجية وإجراءات إدارة التغيير في المشاريع والأصول (المعلومات/التكنولوجيا) لتحديد مخاطر الأمن السيبراني وإدارتها كجزء من دورة إدارة المشروع. ويجب أن تكون متطلبات الأمن السيبراني جزءًا أساسيًا من المتطلبات العامة لمشاريع التكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_22,
  @compliance_library_id,
  N'1.6.2',
  N'The cybersecurity requirements in project and asset (information/technology) change management must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني في إدارة تغيير المشروع والأصول (المعلومات/التكنولوجيا) ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_23,
  @compliance_library_id,
  N'1.6.2.1',
  N'Vulnerability assessment and remediation.',
  N'تقييم نقاط الضعف ومعالجتها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_24,
  @compliance_library_id,
  N'1.6.2.2',
  N'Conducting a configurations’ review, secure configuration, and hardening and patching before changes or going live for technology projects.',
  N'إجراء مراجعة التكوينات وتأمين التكوين وتقويته وتصحيحه قبل إجراء التغييرات أو إطلاق المشاريع التكنولوجية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_25,
  @compliance_library_id,
  N'1.6.3',
  N'The cybersecurity requirements related to software and application development projects must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني المتعلقة بمشاريع تطوير البرامج والتطبيقات ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_26,
  @compliance_library_id,
  N'1.6.3.1',
  N'Using secure coding standards.',
  N'استخدام معايير الترميز الآمنة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_27,
  @compliance_library_id,
  N'1.6.3.2',
  N'Using trusted and licensed sources for software development tools and libraries.',
  N'استخدام مصادر موثوقة ومرخصة لأدوات تطوير البرامج والمكتبات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_28,
  @compliance_library_id,
  N'1.6.3.3',
  N'Conducting compliance tests for software against the defined organizational cybersecurity requirements.',
  N'إجراء اختبارات التوافق للبرامج مع متطلبات الأمن السيبراني المنظمة المحددة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_29,
  @compliance_library_id,
  N'1.6.3.4',
  N'Secure integration between software components.',
  N'التكامل الآمن بين مكونات البرنامج.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_30,
  @compliance_library_id,
  N'1.6.3.5',
  N'Conducting a configurations’ review, secure configuration, and hardening and patching before going live for software products.',
  N'إجراء مراجعة التكوينات وتأمين التكوين والتحصين والتصحيح قبل إطلاق منتجات البرامج.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_31,
  @compliance_library_id,
  N'1.6.4',
  N'The cybersecurity requirements in project management must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني في إدارة المشاريع بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_32,
  @compliance_library_id,
  N'1.7.1',
  N'The organization must comply with related national cybersecurity laws and regulations.',
  N'يجب على المنظمة الامتثال لقوانين وأنظمة الأمن السيبراني الوطنية ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_33,
  @compliance_library_id,
  N'1.7.2',
  N'The organization must comply with any nationally-approved international agreements and commitments related to cybersecurity.',
  N'يجب على المنظمة الالتزام بأي اتفاقيات والتزامات دولية معتمدة على المستوى الوطني تتعلق بالأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_34,
  @compliance_library_id,
  N'1.8.1',
  N'Cybersecurity reviews must be conducted periodically by the cybersecurity function in the organization to assess compliance with the cybersecurity controls in the organization.',
  N'يجب إجراء مراجعات الأمن السيبراني بشكل دوري من قبل وظيفة الأمن السيبراني في المنظمة لتقييم الامتثال لضوابط الأمن السيبراني في المنظمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_35,
  @compliance_library_id,
  N'1.8.2',
  N'Cybersecurity audits and reviews must be conducted by independent parties outside the cybersecurity function (e.g., Internal Audit function) to assess compliance with the cybersecurity controls in the organization. Audits and reviews must be conducted independently, while ensuring that this does not result in a conflict of interest, as per the Generally Accepted Auditing Standards (GAAS), and related laws and regulations.',
  N'يجب أن تُجري جهات مستقلة خارج نطاق إدارة الأمن السيبراني (مثل إدارة التدقيق الداخلي) عمليات تدقيق ومراجعة الأمن السيبراني لتقييم الامتثال لضوابط الأمن السيبراني في المؤسسة. ويجب إجراء عمليات التدقيق والمراجعة بشكل مستقل، مع ضمان عدم نشوء أي تضارب في المصالح، وفقًا لمعايير التدقيق المقبولة عمومًا (GAAS) والقوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_36,
  @compliance_library_id,
  N'1.8.3',
  N'Results from the cybersecurity audits and reviews must be documented and presented to the cybersecurity steering committee and Authorizing Official. Results must include the audit/review scope, observations, recommendations, and remediation plans.',
  N'يجب توثيق نتائج عمليات التدقيق والمراجعة المتعلقة بالأمن السيبراني وتقديمها إلى اللجنة التوجيهية للأمن السيبراني والمسؤول المُخوّل. يجب أن تتضمن النتائج نطاق التدقيق/المراجعة، والملاحظات، والتوصيات، وخطط المعالجة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_37,
  @compliance_library_id,
  N'1.9.1',
  N'Personnel cybersecurity requirements (prior to employment, during employment, and after termination/separation) must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني للموظفين (قبل التوظيف، وأثناء التوظيف، وبعد انتهاء الخدمة/الفصل) وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_38,
  @compliance_library_id,
  N'1.9.2',
  N'The personnel cybersecurity requirements must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني للموظفين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_39,
  @compliance_library_id,
  N'1.9.3',
  N'The personnel cybersecurity requirements prior to employment must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني للموظفين قبل التوظيف ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_40,
  @compliance_library_id,
  N'1.9.3.1',
  N'Inclusion of personnel cybersecurity responsibilities and non-disclosure clauses (covering the cybersecurity requirements during employment and after termination/separation) in employment contracts.',
  N'إدراج مسؤوليات الأمن السيبراني للموظفين وبنود عدم الإفصاح (التي تغطي متطلبات الأمن السيبراني أثناء العمل وبعد الإنهاء/الفصل) في عقود العمل.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_41,
  @compliance_library_id,
  N'1.9.3.2',
  N'Screening or vetting candidates for cybersecurity and critical/privileged positions.',
  N'فحص أو التحقق من المرشحين لوظائف الأمن السيبراني والمناصب الحرجة/المتميزة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_42,
  @compliance_library_id,
  N'1.9.4',
  N'The personnel cybersecurity requirements during employment must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني للموظفين أثناء العمل ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_43,
  @compliance_library_id,
  N'1.9.4.1',
  N'Cybersecurity awareness (during onboarding and during employment).',
  N'التوعية بالأمن السيبراني (أثناء التوجيه وأثناء العمل).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_44,
  @compliance_library_id,
  N'1.9.4.2',
  N'Implementation of and compliance with the cybersecurity requirements as per the organizational cybersecurity policies and procedures.',
  N'تنفيذ متطلبات الأمن السيبراني والامتثال لها وفقًا لسياسات وإجراءات الأمن السيبراني التنظيمية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_45,
  @compliance_library_id,
  N'1.9.5',
  N'Personnel access to information and technology assets must be reviewed and removed immediately upon termination/separation.',
  N'يجب مراجعة وصول الموظفين إلى المعلومات والأصول التكنولوجية وإزالتها فورًا عند انتهاء الخدمة/الفصل.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_46,
  @compliance_library_id,
  N'1.9.6',
  N'Personnel cybersecurity requirements must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني للموظفين بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_47,
  @compliance_library_id,
  N'1.10.1',
  N'A cybersecurity awareness program must be developed and approved. The program must be conducted periodically through multiple channels to strengthen awareness about cybersecurity, cyber threats, and risks, and to build a positive cybersecurity awareness culture.',
  N'يجب تطوير برنامج توعية بالأمن السيبراني واعتماده. يجب تنفيذ البرنامج دوريًا عبر قنوات متعددة لتعزيز الوعي بالأمن السيبراني والتهديدات والمخاطر السيبرانية، وبناء ثقافة إيجابية للتوعية بالأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_48,
  @compliance_library_id,
  N'1.10.2',
  N'The cybersecurity awareness program must be implemented.',
  N'يجب تنفيذ برنامج التوعية بالأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_49,
  @compliance_library_id,
  N'1.10.3',
  N'The cybersecurity awareness program must cover the latest cyber threats and how to protect against them and must include at least the following subjects',
  N'يجب أن يغطي برنامج التوعية بالأمن السيبراني أحدث التهديدات السيبرانية وكيفية الحماية منها ويجب أن يشمل على الأقل الموضوعات التالية',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_50,
  @compliance_library_id,
  N'1.10.3.1',
  N'Secure handling of email services, especially phishing emails.',
  N'التعامل الآمن مع خدمات البريد الإلكتروني، وخاصة رسائل التصيد الاحتيالي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_51,
  @compliance_library_id,
  N'1.10.3.2',
  N'Secure handling of mobile devices and storage media.',
  N'التعامل الآمن مع الأجهزة المحمولة ووسائط التخزين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_52,
  @compliance_library_id,
  N'1.10.3.3',
  N'Secure Internet browsing.',
  N'تصفح الإنترنت بشكل آمن.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_53,
  @compliance_library_id,
  N'1.10.3.4',
  N'Secure use of social media.',
  N'الاستخدام الآمن لوسائل التواصل الاجتماعي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_54,
  @compliance_library_id,
  N'1.10.4',
  N'Essential and customized (i.e., tailored to job functions as it relates to cybersecurity) training and access to professional skillsets must be made available to personnel working directly on tasks related to cybersecurity, including',
  N'يجب توفير التدريب الأساسي والمخصص (أي المصمم خصيصًا لوظائف العمل فيما يتعلق بالأمن السيبراني) والوصول إلى مجموعات المهارات المهنية للموظفين الذين يعملون مباشرة في المهام المتعلقة بالأمن السيبراني، بما في ذلك',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_55,
  @compliance_library_id,
  N'1.10.4.1',
  N'Cybersecurity function’s personnel.',
  N'موظفي قسم الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_56,
  @compliance_library_id,
  N'1.10.4.2',
  N'Personnel working on software/application development and information and technology assets operations.',
  N'الموظفين العاملين في تطوير البرمجيات/التطبيقات وعمليات أصول المعلومات والتكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_57,
  @compliance_library_id,
  N'1.10.4.3',
  N'Executive and supervisory positions.',
  N'المناصب التنفيذية والإشرافية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_58,
  @compliance_library_id,
  N'1.10.5',
  N'The implementation of the cybersecurity awareness program must be reviewed periodically.',
  N'يجب مراجعة تنفيذ برنامج التوعية بالأمن السيبراني بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_59,
  @compliance_library_id,
  N'2.1.1',
  N' Cybersecurity requirements for managing information and technology assets must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة أصول المعلومات والتكنولوجيا وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_60,
  @compliance_library_id,
  N'2.1.2',
  N' The cybersecurity requirements for managing information and technology assets must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة أصول المعلومات والتكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_61,
  @compliance_library_id,
  N'2.1.3',
  N' Acceptable use policy of information and technology assets must be defined, documented, and approved.',
  N'يجب تحديد سياسة الاستخدام المقبولة لأصول المعلومات والتكنولوجيا وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_62,
  @compliance_library_id,
  N'2.1.4',
  N' Acceptable use policy of information and technology assets must be implemented.',
  N'يجب تنفيذ سياسة الاستخدام المقبولة لأصول المعلومات والتكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_63,
  @compliance_library_id,
  N'2.1.5',
  N' Information and technology assets must be classified, labeled, and handled as per related law and regulatory requirements.',
  N'يجب تصنيف أصول المعلومات والتكنولوجيا ووضع العلامات عليها ومعالجتها وفقًا للقانون والمتطلبات التنظيمية ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_64,
  @compliance_library_id,
  N'2.1.6',
  N' The cybersecurity requirements for managing information and technology assets must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لإدارة أصول المعلومات والتكنولوجيا بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_65,
  @compliance_library_id,
  N'2.2.1',
  N' Cybersecurity requirements for identity and access management must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة الهوية والوصول وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_66,
  @compliance_library_id,
  N'2.2.2',
  N' The cybersecurity requirements for identity and access management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة الهوية والوصول.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_67,
  @compliance_library_id,
  N'2.2.3',
  N'The cybersecurity requirements for identity and access management must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لإدارة الهوية والوصول ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_68,
  @compliance_library_id,
  N'2.2.3.1',
  N'Single-factor authentication based on username and password.',
  N'المصادقة أحادية العامل بناءً على اسم المستخدم وكلمة المرور.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_69,
  @compliance_library_id,
  N'2.2.3.2',
  N'Multi-factor authentication for remote access, defining suitable authentication factors, number of factors, and suitable technique based on the result of the impact assessment of authentication failure and bypass for remote access.',
  N'المصادقة متعددة العوامل للوصول عن بعد، وتحديد عوامل المصادقة المناسبة، وعدد العوامل، والتقنية المناسبة بناءً على نتيجة تقييم تأثير فشل المصادقة وتجاوز الوصول عن بعد.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_70,
  @compliance_library_id,
  N'2.2.3.3',
  N'User authorization based on identity and access control principles Need-to-Know and Need-to-Use, Least Privilege, and Segregation of Duties.',
  N'تفويض المستخدم على أساس مبادئ التحكم في الهوية والوصول والحاجة إلى المعرفة والحاجة إلى الاستخدام، والحد الأدنى من الامتيازات، وفصل الواجبات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_71,
  @compliance_library_id,
  N'2.2.3.4',
  N'Privileged access management.',
  N'إدارة الوصول المتميز.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_72,
  @compliance_library_id,
  N'2.2.3.5',
  N'Periodic review of users’ identities and access rights.',
  N'المراجعة الدورية لهويات المستخدمين وحقوق الوصول.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_73,
  @compliance_library_id,
  N'2.2.4',
  N'The implementation of the cybersecurity requirements for identity and access management must be reviewed periodically.',
  N'يجب مراجعة تنفيذ متطلبات الأمن السيبراني لإدارة الهوية والوصول بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_74,
  @compliance_library_id,
  N'2.3.1',
  N'Cybersecurity requirements for protecting information systems and information processing facilities must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لحماية أنظمة المعلومات ومرافق معالجة المعلومات وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_75,
  @compliance_library_id,
  N'2.3.2',
  N'The cybersecurity requirements for protecting information systems and information processing facilities must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لحماية أنظمة المعلومات ومرافق معالجة المعلومات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_76,
  @compliance_library_id,
  N'2.3.3',
  N'The cybersecurity requirements for protecting information systems and information processing facilities must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لحماية أنظمة المعلومات ومرافق معالجة المعلومات ما يلي على الأقل:',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_77,
  @compliance_library_id,
  N'2.3.3.1',
  N'Advanced, up-to-date, and secure management of malware and virus protection on servers and workstations.',
  N'إدارة متقدمة وحديثة وآمنة لحماية البرامج الضارة والفيروسات على الخوادم ومحطات العمل.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_78,
  @compliance_library_id,
  N'2.3.3.2',
  N'Restricted use and secure handling of external storage media.',
  N'الاستخدام المحدود والتعامل الآمن مع وسائط التخزين الخارجية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_79,
  @compliance_library_id,
  N'2.3.3.3',
  N'Patch management for information systems, software, and devices.',
  N'إدارة التصحيحات لأنظمة المعلومات والبرامج والأجهزة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_80,
  @compliance_library_id,
  N'2.3.3.4',
  N'Centralized clock synchronization with an accurate and trusted source (e.g., Saudi Standards, Metrology, and Quality Organization (SASO)).',
  N'مزامنة الساعة المركزية مع مصدر دقيق وموثوق (على سبيل المثال، الهيئة السعودية للمواصفات والمقاييس والجودة (SASO)).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_81,
  @compliance_library_id,
  N'2.3.4',
  N'The cybersecurity requirements for protecting information systems and information processing facilities must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لحماية أنظمة المعلومات ومرافق معالجة المعلومات بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_82,
  @compliance_library_id,
  N'2.4.1',
  N'Cybersecurity requirements for protecting email service must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لحماية خدمة البريد الإلكتروني وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_83,
  @compliance_library_id,
  N'2.4.2',
  N'The cybersecurity requirements for email service must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لخدمة البريد الإلكتروني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_84,
  @compliance_library_id,
  N'2.4.3',
  N'The cybersecurity requirements for protecting the email service must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لحماية خدمة البريد الإلكتروني ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_85,
  @compliance_library_id,
  N'2.4.3.1',
  N'Analyzing and filtering email messages (specifically phishing emails and spam) using advanced and up-to-date email protection techniques.',
  N'تحليل وتصفية رسائل البريد الإلكتروني (وخاصة رسائل التصيد والبريد العشوائي) باستخدام تقنيات حماية البريد الإلكتروني المتقدمة والحديثة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_86,
  @compliance_library_id,
  N'2.4.3.2',
  N'Multi-factor authentication for remote and webmail access to email service, defining authentication factors, number of factors, and suitable techniques based on the result of the impact assessment of authentication failure and bypass.',
  N'المصادقة متعددة العوامل للوصول عن بعد وعبر البريد الإلكتروني إلى خدمة البريد الإلكتروني، وتحديد عوامل المصادقة وعدد العوامل والتقنيات المناسبة بناءً على نتيجة تقييم تأثير فشل المصادقة وتجاوزها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_87,
  @compliance_library_id,
  N'2.4.3.3',
  N'Email archiving and backup.',
  N'أرشفة البريد الإلكتروني والنسخ الاحتياطي.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_88,
  @compliance_library_id,
  N'2.4.3.4',
  N'Secure management and protection against Advanced Persistent Threats (APT), which normally utilize zero-day viruses and malware.',
  N'إدارة وحماية آمنة ضد التهديدات المتقدمة المستمرة (APT)، والتي تستخدم عادةً فيروسات اليوم صفر والبرامج الضارة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_89,
  @compliance_library_id,
  N'2.4.3.5',
  N'Validation of the organization’s email service domains through the Haseen platform by using Sender Policy Framework (SPF), Domain Keys Identified Mail (DKIM), and Domain Message Authentication Reporting and Conformance (DMARC).',
  N'التحقق من صحة نطاقات خدمة البريد الإلكتروني للمنظمة من خلال منصة Haseen باستخدام إطار سياسة المرسل (SPF)، والبريد المحدد بمفاتيح المجال (DKIM)، وتقارير مصادقة رسائل المجال والتوافق (DMARC).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_90,
  @compliance_library_id,
  N'2.4.4',
  N'The cybersecurity requirements for email service must be reviewed periodically',
  N'يجب مراجعة متطلبات الأمن السيبراني لخدمة البريد الإلكتروني بشكل دوري',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_91,
  @compliance_library_id,
  N'2.5.1',
  N' Cybersecurity requirements for network security management must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة أمن الشبكات وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_92,
  @compliance_library_id,
  N'2.5.2',
  N' The cybersecurity requirements for network security management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة أمن الشبكات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_93,
  @compliance_library_id,
  N'2.5.3',
  N'The cybersecurity requirements for network security management must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لإدارة أمن الشبكة ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_94,
  @compliance_library_id,
  N'2.5.3.1',
  N'Logical or physical segregation and segmentation of network segments using firewalls and defense-in-depth principles.',
  N'الفصل والتجزئة المنطقية أو المادية لأجزاء الشبكة باستخدام جدران الحماية ومبادئ الدفاع المتعمق.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_95,
  @compliance_library_id,
  N'2.5.3.2',
  N'Network segregation between production, test, and development environments.',
  N'الفصل الشبكي بين بيئات الإنتاج والاختبار والتطوير.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_96,
  @compliance_library_id,
  N'2.5.3.3',
  N'Secure browsing and Internet connectivity including restrictions on the use of file storage/sharing and remote access websites, and protection against suspicious websites.',
  N'التصفح الآمن والاتصال بالإنترنت بما في ذلك القيود المفروضة على استخدام مواقع تخزين الملفات/مشاركتها ومواقع الوصول عن بعد، والحماية من المواقع المشبوهة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_97,
  @compliance_library_id,
  N'2.5.3.4',
  N'Wireless network protection using strong authentication and encryption techniques. A comprehensive risk assessment and management exercise must be conducted to assess and manage the cyber risks prior to connecting any wireless networks to the organization’s internal network.',
  N'حماية الشبكات اللاسلكية باستخدام تقنيات مصادقة وتشفير قوية. يجب إجراء تقييم شامل للمخاطر وإدارتها قبل ربط أي شبكات لاسلكية بالشبكة الداخلية للمؤسسة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_98,
  @compliance_library_id,
  N'2.5.3.5',
  N'Management and restrictions on network services, protocols, and ports.',
  N'إدارة وقيود خدمات الشبكة والبروتوكولات والمنافذ.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_99,
  @compliance_library_id,
  N'2.5.3.6',
  N'Intrusion Prevention Systems (IPS).',
  N'أنظمة منع التطفل (IPS).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_100,
  @compliance_library_id,
  N'2.5.3.7',
  N'Security of Domain Name Service (DNS) through Haseen platform.',
  N'أمن خدمة اسم النطاق (DNS) عبر منصة Haseen.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_101,
  @compliance_library_id,
  N'2.5.3.8',
  N'Secure management and protection of Internet browsing channel against Advanced Persistent Threats (APT), which normally utilize zero-day viruses and malware.',
  N'إدارة وحماية آمنة لقناة تصفح الإنترنت ضد التهديدات المستمرة المتقدمة (APT)، والتي تستخدم عادةً فيروسات اليوم صفر والبرامج الضارة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_102,
  @compliance_library_id,
  N'2.5.3.9',
  N'Protecting against Distributed Denial of Service (DDoS) attacks to limit risks arising from these attacks.',
  N'الحماية ضد هجمات رفض الخدمة الموزعة (DDoS) للحد من المخاطر الناجمة عن هذه الهجمات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_103,
  @compliance_library_id,
  N'2.5.4',
  N' The cybersecurity requirements for network security management must be reviewed periodically',
  N'يجب مراجعة متطلبات الأمن السيبراني لإدارة أمن الشبكات بشكل دوري',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_104,
  @compliance_library_id,
  N'2.6.1',
  N'Cybersecurity requirements for mobile devices security and BYOD must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لأمن الأجهزة المحمولة و BYOD وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_105,
  @compliance_library_id,
  N'2.6.2',
  N'The cybersecurity requirements for mobile devices security and BYOD must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لأمن الأجهزة المحمولة و BYOD.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_106,
  @compliance_library_id,
  N'2.6.3',
  N' The cybersecurity requirements for mobile devices security and BYOD must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لأمن الأجهزة المحمولة و BYOD ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_107,
  @compliance_library_id,
  N'2.6.3.1',
  N' Separation and encryption of organization’s data and information stored on mobile devices and BYODs.',
  N'فصل وتشفير بيانات المنظمة ومعلوماتها المخزنة على الأجهزة المحمولة وأجهزة BYOD.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_108,
  @compliance_library_id,
  N'2.6.3.2',
  N' Controlled and restricted use based on job requirements.',
  N'الاستخدام الخاضع للرقابة والمقيد بناءً على متطلبات الوظيفة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_109,
  @compliance_library_id,
  N'2.6.3.3',
  N' Secure wiping of organization’s data and information stored on mobile devices and BYOD in cases of device loss, theft or after termination/separation from the organization.',
  N'المسح الآمن لبيانات المنظمة ومعلوماتها المخزنة على الأجهزة المحمولة و BYOD في حالات فقدان الجهاز أو سرقته أو بعد إنهاء الخدمة/الانفصال عن المنظمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_110,
  @compliance_library_id,
  N'2.6.3.4',
  N'Security awareness for mobile devices users.',
  N'الوعي الأمني لمستخدمي الأجهزة المحمولة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_111,
  @compliance_library_id,
  N'2.6.4',
  N' The cybersecurity requirements for mobile devices security and BYOD must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لأمن الأجهزة المحمولة و BYOD بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_112,
  @compliance_library_id,
  N'2.7.1',
  N' Cybersecurity requirements for protecting and handling data and information must be defined, documented, and approved as per the related laws and regulations.',
  N'يجب تحديد متطلبات الأمن السيبراني لحماية البيانات والمعلومات ومعالجتها وتوثيقها والموافقة عليها وفقًا للقوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_113,
  @compliance_library_id,
  N'2.7.2',
  N' The cybersecurity requirements for protecting and handling data and information must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لحماية البيانات والمعلومات ومعالجتها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_114,
  @compliance_library_id,
  N'2.7.3',
  N' The cybersecurity requirements for protecting and handling data and information must include at least the applicable requirements in Data Cybersecurity Controls published by NCA.',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لحماية البيانات والمعلومات ومعالجتها على الأقل المتطلبات المعمول بها في ضوابط الأمن السيبراني للبيانات التي نشرتها NCA.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_115,
  @compliance_library_id,
  N'2.7.4',
  N' The cybersecurity requirements for protecting and handling data and information must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لحماية البيانات والمعلومات ومعالجتها بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_116,
  @compliance_library_id,
  N'2.8.1',
  N' Cybersecurity requirements for cryptography must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني للتشفير وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_117,
  @compliance_library_id,
  N'2.8.2',
  N' The cybersecurity requirements for cryptography must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني للتشفير.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_118,
  @compliance_library_id,
  N'2.8.3',
  N' The cybersecurity requirements for cryptography must include at least the requirements in the National Cryptographic Standards; published by NCA and each national entity is required to choose and implement the appropriate cryptographic standard level based on the nature and sensitivity of the data, systems, and networks to be protected, and based on the risk assessment by the entity; and as per related laws and regulations; according to the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني للتشفير على الأقل المتطلبات الواردة في معايير التشفير الوطنية؛ التي نشرتها الهيئة الوطنية للتشفير، ويُطلب من كل كيان وطني اختيار وتنفيذ مستوى معيار التشفير المناسب بناءً على طبيعة وحساسية البيانات والأنظمة والشبكات التي يجب حمايتها، وبناءً على تقييم المخاطر من قبل الكيان؛ ووفقًا للقوانين واللوائح ذات الصلة؛ وفقًا لما يلي',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_119,
  @compliance_library_id,
  N'2.8.3.1',
  N' Approved cryptographic systems and solutions standards and its technical and regulatory limitations.',
  N'معايير أنظمة وحلول التشفير المعتمدة وحدودها الفنية والتنظيمية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_120,
  @compliance_library_id,
  N'2.8.3.2',
  N' Secure management of cryptographic keys during their lifecycles.',
  N'إدارة آمنة للمفاتيح التشفيرية أثناء دورات حياتها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_121,
  @compliance_library_id,
  N'2.8.3.3',
  N' Encryption of data in-transit, at-rest, and while processing as per classification and related laws and regulations.',
  N'تشفير البيانات أثناء النقل، وفي حالة السكون، وأثناء المعالجة وفقًا للتصنيف والقوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_122,
  @compliance_library_id,
  N'2.8.4',
  N' The cybersecurity requirements for cryptography must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني للتشفير بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_123,
  @compliance_library_id,
  N'2.9.1',
  N' Cybersecurity requirements for backup and recovery management must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة النسخ الاحتياطي والاسترداد وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_124,
  @compliance_library_id,
  N'2.9.2',
  N' The cybersecurity requirements for backup and recovery management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة النسخ الاحتياطي والاسترداد.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_125,
  @compliance_library_id,
  N'2.9.3',
  N' The cybersecurity requirements for backup and recovery management must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لإدارة النسخ الاحتياطي والاسترداد ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_126,
  @compliance_library_id,
  N'2.9.3.1',
  N' Scope and coverage of backups to cover critical technology and information assets.',
  N'نطاق وتغطية النسخ الاحتياطية لتغطية الأصول التكنولوجية والمعلوماتية الهامة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_127,
  @compliance_library_id,
  N'2.9.3.2',
  N' Ability to perform quick recovery of data and systems after cybersecurity incidents.',
  N'القدرة على القيام باسترداد سريع للبيانات والأنظمة بعد حوادث الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_128,
  @compliance_library_id,
  N'2.9.3.3',
  N' Periodic tests of backups.',
  N'اختبارات دورية للنسخ الاحتياطية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_129,
  @compliance_library_id,
  N'2.9.4',
  N' The cybersecurity requirements for backup and recovery management must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لإدارة النسخ الاحتياطي والاسترداد بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_130,
  @compliance_library_id,
  N'2.10.1',
  N'Cybersecurity requirements for technical vulnerabilities management must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة الثغرات التقنية وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_131,
  @compliance_library_id,
  N'2.10.2',
  N' The cybersecurity requirements for technical vulnerabilities management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة الثغرات التقنية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_132,
  @compliance_library_id,
  N'2.10.3',
  N' The cybersecurity requirements for technical vulnerabilities management must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لإدارة الثغرات التقنية ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_133,
  @compliance_library_id,
  N'2.10.3.1',
  N' Periodic vulnerabilities assessments.',
  N'تقييمات دورية للثغرات الأمنية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_134,
  @compliance_library_id,
  N'2.10.3.2',
  N' Vulnerabilities classification based on criticality level.',
  N'تصنيف الثغرات الأمنية على أساس مستوى الخطورة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_135,
  @compliance_library_id,
  N'2.10.3.3',
  N' Vulnerabilities remediation based on classification and associated risk levels.',
  N'معالجة الثغرات الأمنية بناءً على التصنيف ومستويات المخاطر المرتبطة بها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_136,
  @compliance_library_id,
  N'2.10.3.4',
  N' Security patch management.',
  N'إدارة تصحيحات الأمان.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_137,
  @compliance_library_id,
  N'2.10.3.5',
  N' Subscription with authorized and trusted cybersecurity resources for up-to-date information and notifications on technical vulnerabilities.',
  N'الاشتراك في موارد الأمن السيبراني المعتمدة والموثوقة للحصول على معلومات وإشعارات محدثة حول الثغرات التقنية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_138,
  @compliance_library_id,
  N'2.10.4',
  N' The cybersecurity requirements for technical vulnerabilities management must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لإدارة الثغرات التقنية بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_139,
  @compliance_library_id,
  N'2.11.1',
  N' Cybersecurity requirements for penetration testing exercises must be defined, documented, and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لممارسات اختبار الاختراق وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_140,
  @compliance_library_id,
  N'2.11.2',
  N' The cybersecurity requirements for penetration testing processes must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لعمليات اختبار الاختراق.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_141,
  @compliance_library_id,
  N'2.11.3',
  N' The cybersecurity requirements for penetration testing processes must include at least the following',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لعمليات اختبار الاختراق ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_142,
  @compliance_library_id,
  N'2.11.3.1',
  N' Scope of penetration tests which must cover Internet-facing services and its technical components including infrastructure, websites, web applications, mobile apps, email, and remote access.',
  N'نطاق اختبارات الاختراق التي يجب أن تغطي الخدمات التي تواجه الإنترنت ومكوناتها التقنية بما في ذلك البنية التحتية ومواقع الويب وتطبيقات الويب وتطبيقات الهاتف المحمول والبريد الإلكتروني والوصول عن بعد.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_143,
  @compliance_library_id,
  N'2.11.3.2',
  N' Conducting penetration tests periodically.',
  N'إجراء اختبارات الاختراق بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_144,
  @compliance_library_id,
  N'2.12.1',
  N' Cybersecurity requirements for event logs and monitoring management must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لسجلات الأحداث وإدارة المراقبة وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_145,
  @compliance_library_id,
  N'2.12.2',
  N' The cybersecurity requirements for event logs and monitoring management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لسجلات الأحداث وإدارة المراقبة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_146,
  @compliance_library_id,
  N'2.12.3',
  N' The cybersecurity requirements must include',
  N'يجب أن تشمل متطلبات الأمن السيبراني ما يلي:',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_147,
  @compliance_library_id,
  N'2.12.3.1',
  N' Activation of cybersecurity event logs on critical information assets.',
  N'تفعيل سجلات أحداث الأمن السيبراني على الأصول المعلوماتية الهامة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_148,
  @compliance_library_id,
  N'2.12.3.2',
  N' Activation on remote access and privileged user accounts.',
  N'التنشيط على الوصول عن بعد وحسابات المستخدمين المميزين.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_149,
  @compliance_library_id,
  N'2.12.3.3',
  N' Identification of required technologies (e.g., SIEM).',
  N'تحديد التقنيات المطلوبة (على سبيل المثال، SIEM).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_150,
  @compliance_library_id,
  N'2.12.3.4',
  N' Continuous monitoring of cybersecurity events.',
  N'المراقبة المستمرة لأحداث الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_151,
  @compliance_library_id,
  N'2.12.3.5',
  N' Retention period for cybersecurity event logs (minimum 12 months).',
  N'فترة الاحتفاظ بسجلات أحداث الأمن السيبراني (12 شهرًا كحد أدنى).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_152,
  @compliance_library_id,
  N'2.12.4',
  N' The cybersecurity requirements must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_153,
  @compliance_library_id,
  N'2.13.1',
  N' Requirements for cybersecurity incidents and threat management must be defined, documented and approved.',
  N'يجب تحديد متطلبات حوادث الأمن السيبراني وإدارة التهديدات وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_154,
  @compliance_library_id,
  N'2.13.2',
  N' The requirements must be implemented.',
  N'يجب تنفيذ المتطلبات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_155,
  @compliance_library_id,
  N'2.13.3',
  N' The requirements for cybersecurity incidents and threat management must include at least the following',
  N'يجب أن تتضمن متطلبات حوادث الأمن السيبراني وإدارة التهديدات ما يلي على الأقل',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_156,
  @compliance_library_id,
  N'2.13.3.1',
  N' Cybersecurity incident response plans and escalation procedures.',
  N'خطط الاستجابة لحوادث الأمن السيبراني وإجراءات التصعيد.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_157,
  @compliance_library_id,
  N'2.13.3.2',
  N' Cybersecurity incidents classification.',
  N'تصنيف حوادث الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_158,
  @compliance_library_id,
  N'2.13.3.3',
  N' Cybersecurity incidents reporting to NCA.',
  N'الإبلاغ عن حوادث الأمن السيبراني إلى الوكالة الوطنية لمكافحة الجريمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_158,
  @compliance_library_id,
  N'2.13.3.4',
  N' Sharing incident notifications, threat intelligence, breach indicators and reports with NCA.',
  N'مشاركة إشعارات الحوادث، ومعلومات التهديد، ومؤشرات الاختراق والتقارير مع NCA.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_160,
  @compliance_library_id,
  N'2.13.3.5',
  N' Collecting and handling threat intelligence feeds.',
  N'جمع ومعالجة معلومات استخبارات التهديدات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_161,
  @compliance_library_id,
  N'2.13.4',
  N' The requirements for cybersecurity incidents and threat management must be reviewed periodically.',
  N'يجب مراجعة متطلبات حوادث الأمن السيبراني وإدارة التهديدات بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_162,
  @compliance_library_id,
  N'2.14.1',
  N' Cybersecurity requirements for physical protection of information and technology assets must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني للحماية المادية لأصول المعلومات والتكنولوجيا وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_163,
  @compliance_library_id,
  N'2.14.2',
  N' The cybersecurity requirements for physical protection of information and technology assets must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني للحماية المادية لأصول المعلومات والتكنولوجيا.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_164,
  @compliance_library_id,
  N'2.14.3',
  N' The cybersecurity requirements must include at least',
  N'يجب أن تتضمن متطلبات الأمن السيبراني ما لا يقل عن',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_165,
  @compliance_library_id,
  N'2.14.3.1',
  N' Authorized access to sensitive areas (e.g., data centers, disaster recovery, surveillance centers, etc.).',
  N'الوصول المسموح به إلى المناطق الحساسة (على سبيل المثال، مراكز البيانات، واستعادة الكوارث، ومراكز المراقبة، وما إلى ذلك).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_166,
  @compliance_library_id,
  N'2.14.3.2',
  N' Facility entry/exit records and CCTV monitoring.',
  N'سجلات الدخول والخروج من المنشأة ومراقبة كاميرات المراقبة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_167,
  @compliance_library_id,
  N'2.14.3.3',
  N' Protection of facility entry/exit and surveillance records.',
  N'حماية سجلات الدخول والخروج والمراقبة الخاصة بالمنشأة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_168,
  @compliance_library_id,
  N'2.14.3.4',
  N' Secure destruction and re-use of physical assets with classified information.',
  N'التدمير الآمن وإعادة استخدام الأصول المادية التي تحتوي على معلومات سرية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_169,
  @compliance_library_id,
  N'2.14.3.5',
  N' Security of devices and equipment inside/outside facilities.',
  N'أمن الأجهزة والمعدات داخل وخارج المنشأة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_170,
  @compliance_library_id,
  N'2.14.4',
  N' The cybersecurity requirements must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_171,
  @compliance_library_id,
  N'2.15.1',
  N' Cybersecurity requirements for external web applications must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لتطبيقات الويب الخارجية وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_172,
  @compliance_library_id,
  N'2.15.2',
  N' The cybersecurity requirements must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_173,
  @compliance_library_id,
  N'2.15.3',
  N' The cybersecurity requirements must include at least',
  N'يجب أن تتضمن متطلبات الأمن السيبراني ما لا يقل عن',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_174,
  @compliance_library_id,
  N'2.15.3.1',
  N' Use of web application firewall.',
  N'استخدام جدار حماية تطبيقات الويب.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_175,
  @compliance_library_id,
  N'2.15.3.2',
  N' Adoption of multi-tier architecture.',
  N'اعتماد الهندسة المعمارية متعددة الطبقات.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_176,
  @compliance_library_id,
  N'2.15.3.3',
  N' Use of secure protocols (e.g., HTTPS).',
  N'استخدام بروتوكولات آمنة (على سبيل المثال، HTTPS).',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_177,
  @compliance_library_id,
  N'2.15.3.4',
  N' Clarification of secure usage policies.',
  N'توضيح سياسات الاستخدام الآمن.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_178,
  @compliance_library_id,
  N'2.15.3.5',
  N' Authentication based on factors and risk-based impact.',
  N'المصادقة على أساس العوامل والتأثير المبني على المخاطر.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_179,
  @compliance_library_id,
  N'2.15.4',
  N' The cybersecurity requirements must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_180,
  @compliance_library_id,
  N'3.1.1',
  N' Cybersecurity requirements for business continuity management must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني لإدارة استمرارية الأعمال وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_181,
  @compliance_library_id,
  N'3.1.2',
  N' The cybersecurity requirements for business continuity management must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني لإدارة استمرارية الأعمال.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_182,
  @compliance_library_id,
  N'3.1.3',
  N' The cybersecurity requirements must include at least',
  N'يجب أن تتضمن متطلبات الأمن السيبراني ما لا يقل عن',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_183,
  @compliance_library_id,
  N'3.1.3.1',
  N' Ensuring the continuity of cybersecurity systems and procedures.',
  N'ضمان استمرارية أنظمة وإجراءات الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_184,
  @compliance_library_id,
  N'3.1.3.2',
  N' Developing response plans for cybersecurity incidents that may affect the business continuity.',
  N'إعداد خطط الاستجابة لحوادث الأمن السيبراني التي قد تؤثر على استمرارية الأعمال.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_185,
  @compliance_library_id,
  N'3.1.3.3',
  N' Developing disaster recovery plans.',
  N'تطوير خطط التعافي من الكوارث.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_186,
  @compliance_library_id,
  N'3.1.4',
  N' The cybersecurity requirements for business continuity management must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني لإدارة استمرارية الأعمال بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_187,
  @compliance_library_id,
  N'4.1.1',
  N' Cybersecurity requirements for contracts and agreements with third-parties must be identified, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني للعقود والاتفاقيات مع أطراف ثالثة وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_188,
  @compliance_library_id,
  N'4.1.2',
  N' The cybersecurity requirements for contracts and agreements with third-parties must include at least',
  N'يجب أن تتضمن متطلبات الأمن السيبراني للعقود والاتفاقيات مع أطراف ثالثة ما لا يقل عن',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_189,
  @compliance_library_id,
  N'4.1.2.1',
  N' Non-disclosure clauses and secure removal of organization’s data by third-parties upon end of service.',
  N'بنود عدم الإفصاح والإزالة الآمنة لبيانات المنظمة من قبل أطراف ثالثة عند انتهاء الخدمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_190,
  @compliance_library_id,
  N'4.1.2.2',
  N' Communication procedures in case of cybersecurity incidents.',
  N'إجراءات التواصل في حالة وقوع حوادث الأمن السيبراني.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_191,
  @compliance_library_id,
  N'4.1.2.3',
  N' Requirements for third-parties to comply with related organizational policies and procedures, laws and regulations.',
  N'متطلبات امتثال الأطراف الثالثة للسياسات والإجراءات التنظيمية والقوانين واللوائح ذات الصلة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_192,
  @compliance_library_id,
  N'4.1.3',
  N' The cybersecurity requirements for IT outsourcing and managed services third-parties must include',
  N'يجب أن تتضمن متطلبات الأمن السيبراني لخدمات الاستعانة بمصادر خارجية في مجال تكنولوجيا المعلومات والخدمات المُدارة من قِبل جهات خارجية ما يلي:',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_193,
  @compliance_library_id,
  N'4.1.3.1',
  N' Conducting a cybersecurity risk assessment to ensure mitigation controls before signing or when laws change.',
  N'إجراء تقييم لمخاطر الأمن السيبراني لضمان ضوابط التخفيف قبل التوقيع على القوانين أو عند تغييرها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_194,
  @compliance_library_id,
  N'4.1.3.2',
  N' Cybersecurity managed services centers must be fully located within the Kingdom of Saudi Arabia.',
  N'يجب أن تتواجد مراكز خدمات الأمن السيبراني المُدارة بالكامل داخل المملكة العربية السعودية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_195,
  @compliance_library_id,
  N'4.1.4',
  N' The cybersecurity requirements for contracts and agreements with third-parties must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني للعقود والاتفاقيات مع أطراف ثالثة بشكل دوري.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_196,
  @compliance_library_id,
  N'4.2.1',
  N' Cybersecurity requirements related to the use of hosting and cloud computing services must be defined, documented and approved.',
  N'يجب تحديد متطلبات الأمن السيبراني المتعلقة باستخدام خدمات الاستضافة والحوسبة السحابية وتوثيقها والموافقة عليها.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_197,
  @compliance_library_id,
  N'4.2.2',
  N' The cybersecurity requirements related to the use of hosting and cloud computing services must be implemented.',
  N'يجب تنفيذ متطلبات الأمن السيبراني المتعلقة باستخدام خدمات الاستضافة والحوسبة السحابية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_198,
  @compliance_library_id,
  N'4.2.3',
  N' In line with applicable laws, and in addition to domain controls, the cybersecurity requirements must include',
  N'وفقًا للقوانين المعمول بها، وبالإضافة إلى ضوابط النطاق، يجب أن تتضمن متطلبات الأمن السيبراني ما يلي:',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_199,
  @compliance_library_id,
  N'4.2.3.1',
  N' Classification of data prior to hosting and returning data in a usable format upon service completion.',
  N'تصنيف البيانات قبل استضافتها وإرجاع البيانات بتنسيق قابل للاستخدام عند اكتمال الخدمة.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_200,
  @compliance_library_id,
  N'4.2.3.2',
  N' Separation of the organization’s environments (specifically virtual servers) from other environments hosted at the cloud provider.',
  N'فصل بيئات المنظمة (الخوادم الافتراضية على وجه التحديد) عن البيئات الأخرى المستضافة لدى مزود الخدمة السحابية.',
  GETDATE(), GETDATE()
);

INSERT INTO control_library_titles (
  control_library_id, compliance_library_id, control_code, control_title_en, control_title_ar, created_at, updated_at
)
VALUES (
  @cl_201,
  @compliance_library_id,
  N'4.2.4',
  N' The cybersecurity requirements related to the use of hosting and cloud computing services must be reviewed periodically.',
  N'يجب مراجعة متطلبات الأمن السيبراني المتعلقة باستخدام خدمات الاستضافة والحوسبة السحابية بشكل دوري.',
  GETDATE(), GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE control_library_titles;
-- +goose StatementEnd
