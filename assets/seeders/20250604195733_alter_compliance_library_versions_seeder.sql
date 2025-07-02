-- +goose Up
-- +goose StatementBegin

DECLARE @sama_id INT;
SELECT @sama_id = id FROM compliance_libraries
WHERE title_en = 'Cyber Security Framework | Saudi Arabian Monetary Authority';

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
    @sama_id,
    'version 1.0',
    2024,
    N'{"foreword": "In view of the ever-growing seriousness of cyber-attacks, we are conscious of the need to stay one-step ahead. The issuance of a Cyber Security Framework (“Framework”) seeks to support our regulated entities in their efforts to have an appropriate cyber security governance and to build a robust infrastructure along with the necessary detective and preventive controls. The Framework articulates appropriate controls and provides guidance on how to assess maturity level. The adoption and implementation of the Framework is a vital step for ensuring that Saudi Arabian Banking, Insurance and Financing Companies sectors can manage and withstand cyber security threats. In designing the Framework, we have considered the ways that our regulated entities are leveraging technology and felt that each entity will be able to adopt a common approach for addressing cyber security. This will ensure cyber security risks are properly managed throughout the sectors. To achieve the above, the full support and oversight from the Board of Directors and Senior Management are required for its implementation. The Information Technology Risk team within the Deputiship of Supervision is at your disposal for any clarifications and we remain committed to guiding our regulated entities in creating a safer cyber environment."}',
   N'{"introduction": {"title": "1.1 Introduction To The Framework", "description": ["The current digital society has high expectations of flawless customer experience, continuous availability of services and effective protection of sensitive data. Information assets and online services are now strategically important to all public and private organizations, as well as to broader society. These services are vital to the creation of a vibrant digital economy. They are also becoming systemically important to the economy and to broader national security. All of which underlines the need to safeguard sensitive data and transactions, and thereby ensure confidence in the overall Saudi Financial Sector.", "The stakes are high when it comes to the confidentiality, integrity and availability of information assets, and applying new online services and new developments (e.g. Fintech, block chain); while improving resilience against cyber threats. Not only is the dependency on these services growing, but the threat landscape is rapidly changing. The Financial Sector recognizes the rate at which the cyber threats and risks are evolving, as well as the changing technology and business landscape.", "SAMA established a Cyber Security Framework (“the Framework”) to enable Financial Institutions regulated by SAMA (“the Member Organizations”) to effectively identify and address risks related to cyber security. To maintain the protection of information assets and online services, the Member Organizations must adopt the Framework."], "objectives_title": "The objective of the Framework is as follows:", "objectives": ["1. To create a common approach for addressing cyber security within the Member Organizations.", "2. To achieve an appropriate maturity level of cyber security controls within the Member Organizations.", "3. To ensure cyber security risks are properly managed throughout the Member Organizations."], "assessment": "The Framework will be used to periodically assess the maturity level and evaluate the effectiveness of the cyber security controls at Member Organizations, and to compare these with other Member Organizations.", "basis": "The Framework is based on the SAMA requirements and industry cyber security standards, such as NIST, ISF, ISO, BASEL and PCI.", "supersedes": "The Framework supersedes all previous issued SAMA circulars with regard to cyber security. Please refer to ‘Appendix A – Overview previous issued SAMA circulars’ for more details."}, "definition": {"title": "1.2 Definition Of Cyber Security", "description": ["Cyber security is defined as the collection of tools, policies, security concepts, security safeguards, guidelines, risk management approaches, actions, training, best practices, assurance, and technologies that can be used to protect the member organization''s information assets against internal and external threats."], "objectives": ["• Confidentiality – Information assets are accessible only to those authorized to have access (i.e., protected from unauthorized disclosure or (un)intended leakage of sensitive data).", "• Integrity – Information assets are accurate, complete and processed correctly (i.e., protected from unauthorized modification, which may include authenticity and non-repudiation).", "• Availability – Information assets are resilient and accessible when required (i.e., protected from unauthorized disruption)."]}, "scope": {"title": "1.3 Scope", "description": ["The Framework defines principles and objectives for initiating, implementing, maintaining, monitoring and improving cyber security controls in Member Organizations.", "The Framework provides cyber security controls which are applicable to the information assets of the Member Organization, including:"], "applicability": ["• Electronic information.", "• Physical information (hardcopy).", "• Applications, software, electronic services and databases.", "• Computers and electronic machines (e.g., ATM).", "• Information storage devices (e.g., hard disk, USB stick).", "• Premises, equipment and communication networks (technical infrastructure)."], "additional": ["The Framework provides direction for cyber security requirements for Member Organizations and its subsidiaries, staff, third parties and customers.", "For business continuity related requirements please refer to the SAMA Business Continuity Minimum Requirements.", "The Framework has an interrelationship with other corporate policies for related areas, such as physical security and fraud management. This framework does not address the non-cyber security requirements for those areas."]}, "applicability": {"title": "1.4 Applicability", "description": ["The Framework is applicable to all Member Organizations regulated by SAMA, which include the following:"], "entities": ["• All Banks operating in Saudi Arabia;", "• All Insurance and/or Reinsurance Companies operating in Saudi Arabia;", "• All Financing Companies operating in Saudi Arabia;", "• All Credit Bureaus operating In Saudi Arabia;", "• The Financial Market Infrastructure"], "exceptions": ["All domains are applicable for the banking sector. However, for other financial institutions the following exceptions apply:", "• Sub-domain (3.1.2): the alignment with cyber security strategy of banking sector is mandatory when applicable.", "• Exclude sub-domain (3.2.3). However, if the organization stores, processes or transmits cardholder data or deals with SWIFT services, then PCI standard and/or SWIFT Customer Security Controls Framework should be implemented.", "• Exclude sub-domain (3.3.12).", "• Exclude sub-domain (3.3.13). However, if the organization provides online services for customers, a Multi Factor Authentication capability should be implemented."]}, "responsibilities": {"title": "1.5 Responsibilities", "description": ["The framework is mandated by SAMA. SAMA is the owner and is responsible for periodically updating the Framework.", "The Member Organizations are responsible for adopting and implementing the Framework."]}, "interpretation": {"title": "1.6 Interpretation", "description": ["SAMA, as the owner of the Framework, is solely responsible for providing interpretations of the principles, objectives and control considerations, if required."]}, "target_audience": {"title": "1.7 Target Audience", "description": ["The Framework is intended for senior and executive management, business owners, owners of information assets, CISOs and those who are responsible for and involved in defining, implementing and reviewing cyber security controls within the Member Organizations."]}, "review_updates_maintenance": {"title": "1.8 Review, Updates and Maintenance", "description": ["The Framework will be reviewed and maintained by SAMA.", "SAMA will review the Framework periodically to determine the Framework’s effectiveness, including the effectiveness of the Framework to address emerging cyber security threats and risks. If applicable, SAMA will update the Framework based on the outcome of the review.", "If a Member Organization considers that an update to the Framework is required, the Member Organization should formally submit the requested update to SAMA. SAMA will review the requested update, and when approved, the Framework will be adjusted.", "The Member Organization will remain responsible to be compliant with the Framework pending the requested update.", "Please refer to ‘Appendix B – How to request an Update to the Framework’ for the process of requesting an update to the Framework.", "Version control will be implemented for maintaining the Framework. Whenever any changes are made, the preceding version shall be retired and the new version shall be published and communicated to all Member Organizations. For the convenience of the Member Organizations, changes to the Framework shall be clearly indicated."]}, "reading_guide": {"title": "1.9 Reading Guide", "description": ["The Framework is structured as follows. Chapter 2 elaborates on the structure of the Framework, and provides instructions on how to apply the Framework. Chapter 3 presents the actual Framework, including the cyber security domains and subdomains, principles, objectives and control considerations."]}}',
     N'{
    "structure": {
      "title": "2.1 Structure",
      "description": "The Framework is structured around four main domains, namely:",
      "domains": [
        "- Cyber Security Leadership and Governance.",
        "- Cyber Security Risk Management and Compliance.",
        "- Cyber Security Operations and Technology.",
        "- Third Party Cyber Security."
      ],
      "subdomain_context": "For each domain, several subdomains are defined. A subdomain focusses on a specific cyber security topic. Per subdomain, the Framework states a principle, objective and control considerations.",
      "subdomain_explanations": [
        "- A principle summarizes the main set of required cyber security controls related to the subdomain.",
        "- The objective describes the purpose of the principle and what the set of required cyber security controls are expected to achieve.",
        "- The control considerations reflect the mandated cyber security controls that should be considered."
      ],
      "numbering_system_note": "Control considerations have been uniquely numbered throughout the Framework. Where applicable, a control consideration can consist of up to 4 levels.",
      "numbering_system_structure": {
        "note": "The control considerations are numbered according to the following numbering system:",
        "example_format": {
          "domain_reference": "3.3",
          "subdomain_reference": "12",
          "separator": "-",
          "control_reference": "4.a.6.a",
          "explanation": "Reference to Control Consideration (Level 1, 2, 3 and 4)"
        }
      },
      "figure_caption": "The figure below illustrates the overall structure of the Framework and indicates the cyber security domains and subdomains, including a reference to the applicable section of the Framework.",
      "domain_breakdown": [
        {
          "heading": "Cyber Security Framework",
          "domain": "3.1 Cyber Security Leadership & Governance",
          "subdomains": [
            "Cyber Security Governance",
            "Cyber Security Strategy",
            "Cyber Security Policy",
            "Cyber Security Roles & Responsibilities",
            "Cyber Security in Project Management",
            "Cyber Security Awareness",
            "Cyber Security Training"
          ]
        },
        {
          "domain": "3.2 Cyber Security Risk Management & Compliance",
          "subdomains": [
            "Cyber Security Risk Management",
            "Regulatory Compliance",
            "Compliance with (inter)national industry standards",
            "Cyber Security Review",
            "Cyber Security Audit"
          ]
        },
        {
          "domain": "3.3 Cyber Security Operations & Technology",
          "subdomains": [
            "Human Resources",
            "Physical Security",
            "Asset Management",
            "Cyber Security Architecture",
            "Identity & Access Management",
            "Application Security",
            "Change Management",
            "Infrastructure Security",
            "Cryptography",
            "Bring Your Own Device (BYOD)",
            "Secure Disposal of Information Assets",
            "Payment Systems",
            "Electronic Banking Services",
            "Cyber Security Event Management",
            "Cyber Security Incident Management",
            "Threat Management",
            "Vulnerability Management"
          ]
        },
        {
          "domain": "3.4 Third Party Cyber Security",
          "subdomains": [
            "Contract & Vendor Management",
            "Outsourcing",
            "Cloud Computing"
          ]
        }
      ]
    },
    "principle_based": {
      "title": "2.2 Principle-based",
      "description": "The Framework is principle based, also referred to as risk based. This means that it prescribes key cyber security principles and objectives to be embedded and achieved by the Member Organization. The list of mandated control considerations provides additional direction and should be considered by the Member Organization in achieving the objectives. When a certain control consideration cannot be tailored or implemented, the Member Organization should consider applying compensating controls, pursuing an internal risk acceptance and requesting a formal waiver from SAMA.",
      "note": "Please refer to Appendix D for details for the – How to request a Waiver from the Framework – process."
    },
    "self_assessment_review_audit": {
      "title": "2.3 Self-Assessment, Review and Audit",
      "description": "The implementation of the Framework at the Member Organization will be subject to a periodic self-assessment. The self-assessment will be performed by the Member Organization based on a questionnaire. The self-assessments will be reviewed and audited by SAMA to determine the level of compliance with the Framework and the cyber security maturity level of the Member Organization.",
      "note": "Please refer to ‘2.4 Cyber Security Maturity Model’ for more details about the cyber security maturity model."
    },
    "cyber_security_maturity_model": {
      "title": "2.4 Cyber Security Maturity Model",
      "description": "The cyber security maturity level will be measured with the help of a predefined cyber security maturity model. The cyber security maturity model distinguishes 6 maturity levels (0, 1, 2, 3, 4 and 5), which are summarized in the table below. In order to achieve levels 3, 4 or 5, a Member Organization must first meet all criteria of the preceding maturity levels.",
      "maturity_levels": [
        {
          "level": 0,
          "name": "Non-existent",
          "criteria": [
            "No documentation.",
            "There is no awareness or attention for certain cyber security control."
          ],
          "explanation": "Cyber security controls are not in place. There may be no awareness of the particular risk area or no current plans to implement such cyber security controls."
        },
        {
          "level": 1,
          "name": "Ad-hoc",
          "criteria": [
            "Cyber security controls is not or partially defined.",
            "Cyber security controls are performed in an inconsistent way.",
            "Cyber security controls are not fully defined."
          ],
          "explanation": "Cyber security control design and execution varies by department or owner. Cyber security control design may only partially mitigate the identified risk and execution may be inconsistent."
        },
        {
          "level": 2,
          "name": "Repeatable but informal",
          "criteria": [
            "The execution of the cyber security control is based on an informal and unwritten, though standardized, practice."
          ],
          "explanation": "Repeatable cyber security controls are in place. However, the control objectives and design are not formally defined or approved. There is limited consideration for a structured review or testing of a control."
        },
        {
          "level": 3,
          "name": "Structured and formalized",
          "criteria": [
            "Cyber security controls are defined, approved and implemented in a structured and formalized way.",
            "The implementation of cyber security controls can be demonstrated."
          ],
          "explanation": "Cyber security policies, standards and procedures are established. Compliance with cyber security documentation i.e., policies, standards and procedures is monitored, preferably using a governance, risk and compliance tool (GRC). Key performance indicators are defined, monitored and reported to evaluate the implementation."
        },
        {
          "level": 4,
          "name": "Managed and measurable",
          "criteria": [
            "The effectiveness of the cyber security controls are periodically assessed and improved when necessary.",
            "This periodic measurement, evaluations and opportunities for improvement are documented."
          ],
          "explanation": "Effectiveness of cyber security controls are measured and periodically evaluated. Key risk indicators and trend reporting are used to determine the effectiveness of the cyber security controls. Results of measurement and evaluation are used to identify opportunities for improvement of the cyber security controls."
        },
        {
          "level": 5,
          "name": "Adaptive",
          "criteria": [
            "Cyber security controls are subject to a continuous improvement plan."
          ],
          "explanation": "The enterprise-wide cyber security program focuses on continuous compliance, effectiveness and improvement of the cyber security controls. Cyber security controls are integrated with enterprise risk management framework and practices. Performance of cyber security controls are evaluated using peer and sector data."
        }
      ],
      "maturity_level_3": {
        "title": "2.4.1 Maturity Level 3",
        "description": "To achieve level 3 maturity, a Member Organization should define, approve and implement cyber security controls. In addition, it should monitor compliance with the cyber security documentation.",
        "documentation_details": {
          "note": "The cyber security documentation should clearly indicate “why”, “what” and “how” cyber security controls should be implemented.",
          "components": [
            "Cyber security policies",
            "Cyber security standards",
            "Cyber security procedures"
          ]
        }
      },
      "documentation_hierarchy": {
        "policy": "Cyber Security Policies - why/what",
        "standard": "Cyber Security Standard - what",
        "procedures": "Cyber Security Procedures - how",
        "note": "Process refers to the collective implementation of policies, standards, and procedures."
      },
      "maturity_level_4": {
        "title": "2.4.2 Maturity Level 4",
        "description": "To achieve maturity level 4, the Member Organization should periodically measure and evaluate the effectiveness of implemented cyber security controls.",
        "measurement_method": "Key risk indicators (KRIs) should be defined. A KRI indicates the norm for effectiveness measurement and should define thresholds to determine whether the actual result of measurement is below, on, or above the targeted norm. KRIs are used for trend reporting and identification of potential improvements."
      },
      "maturity_level_5": {
        "title": "2.4.3 Maturity Level 5",
        "description": "Maturity level 5 focuses on the continuous improvement of cyber security controls. Continuous improvement is achieved through continuously analyzing the goals and achievements of cyber security and identifying structural improvements.",
        "integrations": [
          "Cyber security controls should be integrated with enterprise risk management practices and supported with automated real-time monitoring.",
          "Business process owners should be accountable for monitoring the compliance of the cyber security controls, measuring their effectiveness and incorporating them into the enterprise risk management framework.",
          "Performance of cyber security controls should be evaluated using peer and sector data."
        ]
      }
    }
  }',
    GETDATE(),
    GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM compliance_library_versions
WHERE compliance_library_id = @sama_id
  AND version = 'version 1.0';

-- +goose StatementEnd
