-- +goose Up
-- +goose StatementBegin

DECLARE @library_id INT;
SELECT @library_id = id FROM compliance_libraries
WHERE title_en = 'Essential Cybersecurity Controls (ECC – 2: 2024)';

INSERT INTO compliance_library_versions (
    compliance_library_id,
    version,
    year,
    overview,
    introduction,
    domains,
    created_at,
    updated_at
)
VALUES (
    @library_id,
    'ECC 1',
    2024,
    N'{
        "executive_summary": "The Kingdom of Saudi Arabia’s Vision 2030 aims at a comprehensive improvement of the nation and its security, economy and citizens’ well-being. One of the essential goals is continued transformation into the digital world and improvement of the digital infrastructure in order to keep up with the accelerated global progress in digital services, renewable global networks and IT/OT systems in line with improved computer processing, massive data storage and exchange capabilities for data, in order to be prepared for handling artificial intelligence and 4th industrial revolution transformations.",
        "tlp": {
            "description": "A TLP was created for sharing the maximum amount of sensitive information, and it is widely used world wide. There are Four colours.",
            "levels": [
                {
                    "label": "Red",
                    "title": "Personal, Confidential, and for the Intended Recipient only",
                    "description": "The recipient has no rights to share information classified in red with any person outside the defined range of recipients, either inside or outside the organization, beyond the scope specified for receipt."
                },
                {
                    "label": "Amber",
                    "title": "Restricted Sharing",
                    "description": "The recipient may share information classified in amber only with intended recipients inside the organization and with recipients who are required to take action related to the shared information."
                },
                {
                    "label": "Green",
                    "title": "Sharing within the Same Community",
                    "description": "The recipient may share information classified in green with other recipients inside the organization or outside it within the same sector or related to the organization. However, it is not allowed to exchange or publish this information on public channels."
                },
                {
                    "label": "White",
                    "title": "No Restrictions",
                    "description": "The recipient may share information classified in white with other recipients inside the organization or outside it within the same sector or related to the organization. However, it is not allowed to exchange or publish this information on public channels."
                }
            ]
        }
    }',
    N'{
        "introduction": {
            "description": "The National Cybersecurity Authority (referred to in this document as ''The Authority'' or ''NCA'') updated the essential cybersecurity controls (ECC – 2: 2024) after conducting a comprehensive study of multiple national and international cybersecurity frameworks and standards, studying related national decisions, law and regulatory requirements, reviewing and leveraging cybersecurity best practices, analyzing previous cybersecurity incidents and attacks on government and other critical organizations, and surveying and considering opinions of multiple national organizations.",
            "controls_summary": "The Essential Cybersecurity Controls consist of the following:",
            "control_breakdown": [
                "4 Cybersecurity Main Domains.",
                "28 Cybersecurity Subdomains.",
                "110 Cybersecurity Controls.",
                "90 Cybersecurity Subcontrols."
            ],
            "note": "These cybersecurity controls are linked to related national and international law and regulatory requirements."
        },
        "objectives": {
            "description": "The main objective of these controls is to set the minimum cybersecurity requirements for information and technology assets in organizations. These requirements are based on industry leading practices which will help organizations minimize the cybersecurity risks that originate from internal and external threats. The following key objectives must be focused on in order to protect the organization''s information and technology assets:",
            "key_objectives": [
                "Confidentiality",
                "Integrity",
                "Availability"
            ],
            "pillars_note": "These controls take into consideration the following four main cybersecurity pillars:",
            "pillars": [
                "Strategy",
                "People",
                "Processes",
                "Technology"
            ]
        },
        "scope_of_work_and_applicability": {
            "scope_of_work": "These controls are applicable to government organizations in the Kingdom of Saudi Arabia (including ministries, authorities, establishments and others) and its companies and entities (inside and outside the kingdom), as well as private sector organizations owning, operating or hosting Critical National Infrastructures (CNIs), which are all referred to herein as ''The Organization''. The NCA strongly encourages all other organizations in the Kingdom to leverage these controls to implement best practices to improve and enhance their cybersecurity.",
            "statement_of_applicability": "These controls have been updated after taking into consideration the cybersecurity needs of all organizations and sectors in the Kingdom of Saudi Arabia. Every organization must comply with all applicable controls in this document.",
            "note": "Applicability to implement these cybersecurity controls depends on the organization''s business and its use of certain technologies. For example: Controls in subdomain 4.2 (Cloud Computing and Hosting Cybersecurity) are applicable and must be implemented by organizations currently using or planning to use cloud computing and hosting services."
        },
        "implementation_and_compliance": {
            "description": "To comply with item 3 of article 10 of NCA’s mandate and as per the Royal Decree number 57231 dated 10/11/1439H, all organizations within the scope of these controls must implement whatever necessary to ensure continuous compliance with the controls.",
            "methods": "NCA evaluates organizations’ compliance with the ECC through multiple means such as self-assessments by the organizations, periodic reports of the compliance tool or on-site audits.",
            "tool": "NCA will issue a tool (ECC-2: 2024 Assessment and Compliance Tool) to organize the process of evaluation and compliance measurement against the ECC."
        },
        "update_and_review": {
            "description": "NCA will periodically review and update the ECC as per the cybersecurity requirements and related industry updates. NCA will communicate and publish the updated version of ECC for implementation and compliance. Appendix (C) illustrates the updates on the previous version (i.e., ECC-1:2018)."
        }
    }',
    N'{
        "domains": [
            {
                "title": "Cybersecurity Governance",
                "subdomains": [
                    { "code": "1-1", "name": "Cybersecurity Strategy" },
                    { "code": "1-2", "name": "Cybersecurity Management" },
                    { "code": "1-3", "name": "Cybersecurity Policies and Procedures" },
                    { "code": "1-4", "name": "Cybersecurity Roles and Responsibilities" },
                    { "code": "1-5", "name": "Cybersecurity Risk Management" },
                    { "code": "1-6", "name": "Cybersecurity in Information and Technology Project Management" },
                    { "code": "1-7", "name": "Compliance with Cybersecurity Standards, Laws and Regulations" },
                    { "code": "1-8", "name": "Periodical Cybersecurity Review and Audit" },
                    { "code": "1-9", "name": "Cybersecurity in Human Resources" },
                    { "code": "1-10", "name": "Cybersecurity Awareness and Training Program" }
                ]
            },
            {
                "title": "Code Review & Quality Assurance",
                "subdomains": [
                    { "code": "2-1", "name": "Asset Management" },
                    { "code": "2-2", "name": "Identity and Access Management" },
                    { "code": "2-3", "name": "Information Systems and Information Processing Facilities Protection" },
                    { "code": "2-4", "name": "Email Protection" },
                    { "code": "2-5", "name": "Network Security Management" },
                    { "code": "2-6", "name": "Mobile Devices Security" },
                    { "code": "2-7", "name": "Data and Information Protection" },
                    { "code": "2-8", "name": "Cryptography" },
                    { "code": "2-9", "name": "Backup and Recovery Management" },
                    { "code": "2-10", "name": "Vulnerability Management" },
                    { "code": "2-11", "name": "Penetration Testing" },
                    { "code": "2-12", "name": "Cybersecurity Event Logs and Monitoring Management" },
                    { "code": "2-13", "name": "Cybersecurity Incident and Threat Management" },
                    { "code": "2-14", "name": "Physical Security" },
                    { "code": "2-15", "name": "Web Application Security" }
                ]
            },
            {
                "title": "Software Deployment Process",
                "subdomains": [
                    { "code": "3-1", "name": "Cybersecurity Resilience Aspects of Business Continuity Management (BCM)" }
                ]
            },
            {
                "title": "Third-Party and Cloud Computing Cybersecurity",
                "subdomains": [
                    { "code": "4-1", "name": "Third-Party Cybersecurity" },
                    { "code": "4-2", "name": "Cloud Computing and Hosting Cybersecurity" }
                ]
            }
        ]
    }',
    GETDATE(),
    GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_library_versions;
-- +goose StatementEnd
