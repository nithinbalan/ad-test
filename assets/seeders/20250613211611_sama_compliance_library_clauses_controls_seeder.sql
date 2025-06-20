-- +goose Up
-- +goose StatementBegin
DECLARE @clause_3_1_1 INT;
SELECT @clause_3_1_1 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.1';
DECLARE @ctrl_1 INT;
SELECT @ctrl_1 = id FROM control_library WHERE title_en = N'Define and Approve Cybersecurity Strategy';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_1, GETDATE(), GETDATE()
);
DECLARE @ctrl_2 INT;
SELECT @ctrl_2 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Strategy Roadmap';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_2, GETDATE(), GETDATE()
);
DECLARE @ctrl_3 INT;
SELECT @ctrl_3 = id FROM control_library WHERE title_en = N'Review Cybersecurity Strategy Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_3, GETDATE(), GETDATE()
);
DECLARE @ctrl_4 INT;
SELECT @ctrl_4 = id FROM control_library WHERE title_en = N'Establish Independent Cybersecurity Function';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_4, GETDATE(), GETDATE()
);
DECLARE @ctrl_5 INT;
SELECT @ctrl_5 = id FROM control_library WHERE title_en = N'Employ Qualified Cybersecurity Professionals';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_5, GETDATE(), GETDATE()
);
DECLARE @ctrl_6 INT;
SELECT @ctrl_6 = id FROM control_library WHERE title_en = N'Establish Cybersecurity Steering Committee';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_6, GETDATE(), GETDATE()
);
DECLARE @ctrl_7 INT;
SELECT @ctrl_7 = id FROM control_library WHERE title_en = N'Define and Document Cybersecurity Policies';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_7, GETDATE(), GETDATE()
);
DECLARE @ctrl_8 INT;
SELECT @ctrl_8 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Policies';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_8, GETDATE(), GETDATE()
);
DECLARE @ctrl_9 INT;
SELECT @ctrl_9 = id FROM control_library WHERE title_en = N'Support Policies with Technical Standards';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_9, GETDATE(), GETDATE()
);
DECLARE @ctrl_10 INT;
SELECT @ctrl_10 = id FROM control_library WHERE title_en = N'Review Cybersecurity Policies Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_1, @ctrl_10, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_2  INT;
SELECT @clause_3_1_2  = id FROM compliance_library_clauses WHERE clause_number = N'3.1.2';
DECLARE @ctrl_11 INT;
SELECT @ctrl_11 = id FROM control_library WHERE title_en = N'Define Cybersecurity Organizational Structure';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_2 , @ctrl_11, GETDATE(), GETDATE()
);
DECLARE @ctrl_12 INT;
SELECT @ctrl_12 = id FROM control_library WHERE title_en = N'Review Cybersecurity Roles Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_2 , @ctrl_12, GETDATE(), GETDATE()
);
DECLARE @ctrl_13 INT;
SELECT @ctrl_13 = id FROM control_library WHERE title_en = N'Define Cybersecurity Risk Management Procedures';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_2 , @ctrl_13, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_3 INT;
SELECT @clause_3_1_3 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.3';
DECLARE @ctrl_14 INT;
SELECT @ctrl_14 = id FROM control_library WHERE title_en = N'Implement Risk Management Procedures';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_3, @ctrl_14, GETDATE(), GETDATE()
);
DECLARE @ctrl_15 INT;
SELECT @ctrl_15 = id FROM control_library WHERE title_en = N'Conduct Risk Assessments for Key Scenarios';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_3, @ctrl_15, GETDATE(), GETDATE()
);
DECLARE @ctrl_16 INT;
SELECT @ctrl_16 = id FROM control_library WHERE title_en = N'Assess Risks in Technology Projects';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_3, @ctrl_16, GETDATE(), GETDATE()
);
DECLARE @ctrl_17 INT;
SELECT @ctrl_17 = id FROM control_library WHERE title_en = N'Assess Risks Before Infrastructure Changes';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_3, @ctrl_17, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_4 INT;
SELECT @clause_3_1_4 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.4';
DECLARE @ctrl_18 INT;
SELECT @ctrl_18 = id FROM control_library WHERE title_en = N'Assess Risks for Third-Party Services';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_18, GETDATE(), GETDATE()
);
DECLARE @ctrl_19 INT;
SELECT @ctrl_19 = id FROM control_library WHERE title_en = N'Assess Risks for New Services';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_19, GETDATE(), GETDATE()
);
DECLARE @ctrl_20 INT;
SELECT @ctrl_20 = id FROM control_library WHERE title_en = N'Review Risk Management Procedures Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_20, GETDATE(), GETDATE()
);
DECLARE @ctrl_21 INT;
SELECT @ctrl_21 = id FROM control_library WHERE title_en = N'Include Cybersecurity in Project Management';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_21, GETDATE(), GETDATE()
);
DECLARE @ctrl_22 INT;
SELECT @ctrl_22 = id FROM control_library WHERE title_en = N'Implement Cybersecurity in Change Management';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_22, GETDATE(), GETDATE()
);
DECLARE @ctrl_23 INT;
SELECT @ctrl_23 = id FROM control_library WHERE title_en = N'Conduct Vulnerability Assessments';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_4, @ctrl_23, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_5 INT;
SELECT @clause_3_1_5 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.5';
DECLARE @ctrl_24 INT;
SELECT @ctrl_24 = id FROM control_library WHERE title_en = N'Ensure Secure Configuration and Patching';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_5, @ctrl_24, GETDATE(), GETDATE()
);
DECLARE @ctrl_25 INT;
SELECT @ctrl_25 = id FROM control_library WHERE title_en = N'Secure Software Development Processes';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_5, @ctrl_25, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_6 INT;
SELECT @clause_3_1_6 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.6';
DECLARE @ctrl_26 INT;
SELECT @ctrl_26 = id FROM control_library WHERE title_en = N'Use Secure Coding Standards';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_26, GETDATE(), GETDATE()
);
DECLARE @ctrl_27 INT;
SELECT @ctrl_27 = id FROM control_library WHERE title_en = N'Use Trusted Software Development Tools';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_27, GETDATE(), GETDATE()
);
DECLARE @ctrl_28 INT;
SELECT @ctrl_28 = id FROM control_library WHERE title_en = N'Test Software Compliance';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_28, GETDATE(), GETDATE()
);
DECLARE @ctrl_29 INT;
SELECT @ctrl_29 = id FROM control_library WHERE title_en = N'Ensure Secure Software Integration';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_29, GETDATE(), GETDATE()
);
DECLARE @ctrl_30 INT;
SELECT @ctrl_30 = id FROM control_library WHERE title_en = N'Secure Software Configuration Before Deployment';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_30, GETDATE(), GETDATE()
);
DECLARE @ctrl_31 INT;
SELECT @ctrl_31 = id FROM control_library WHERE title_en = N'Review Project Management Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_31, GETDATE(), GETDATE()
);
DECLARE @ctrl_32 INT;
SELECT @ctrl_32 = id FROM control_library WHERE title_en = N'Comply with Cybersecurity Laws';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_6, @ctrl_32, GETDATE(), GETDATE()
);
DECLARE @clause_3_1_7 INT;
SELECT @clause_3_1_7 = id FROM compliance_library_clauses WHERE clause_number = N'3.1.7';
DECLARE @ctrl_33 INT;
SELECT @ctrl_33 = id FROM control_library WHERE title_en = N'Adhere to International Cybersecurity Agreements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_7, @ctrl_33, GETDATE(), GETDATE()
);
DECLARE @ctrl_34 INT;
SELECT @ctrl_34 = id FROM control_library WHERE title_en = N'Conduct Periodic Cybersecurity Reviews';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_1_7, @ctrl_34, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_1 INT;
SELECT @clause_3_2_1 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.1';
DECLARE @ctrl_35 INT;
SELECT @ctrl_35 = id FROM control_library WHERE title_en = N'Perform Independent Cybersecurity Audits';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_35, GETDATE(), GETDATE()
);
DECLARE @ctrl_36 INT;
SELECT @ctrl_36 = id FROM control_library WHERE title_en = N'Document and Report Audit Results';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_36, GETDATE(), GETDATE()
);
DECLARE @ctrl_37 INT;
SELECT @ctrl_37 = id FROM control_library WHERE title_en = N'Define Personnel Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_37, GETDATE(), GETDATE()
);
DECLARE @ctrl_38 INT;
SELECT @ctrl_38 = id FROM control_library WHERE title_en = N'Implement Personnel Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_38, GETDATE(), GETDATE()
);
DECLARE @ctrl_39 INT;
SELECT @ctrl_39 = id FROM control_library WHERE title_en = N'Establish Pre-Employment Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_39, GETDATE(), GETDATE()
);
DECLARE @ctrl_40 INT;
SELECT @ctrl_40 = id FROM control_library WHERE title_en = N'Include Cybersecurity in Employment Contracts';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_40, GETDATE(), GETDATE()
);
DECLARE @ctrl_41 INT;
SELECT @ctrl_41 = id FROM control_library WHERE title_en = N'Screen Candidates for Cybersecurity Roles';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_41, GETDATE(), GETDATE()
);
DECLARE @ctrl_42 INT;
SELECT @ctrl_42 = id FROM control_library WHERE title_en = N'Enforce Cybersecurity During Employment';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_42, GETDATE(), GETDATE()
);
DECLARE @ctrl_43 INT;
SELECT @ctrl_43 = id FROM control_library WHERE title_en = N'Provide Cybersecurity Awareness Training';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_43, GETDATE(), GETDATE()
);
DECLARE @ctrl_44 INT;
SELECT @ctrl_44 = id FROM control_library WHERE title_en = N'Ensure Compliance with Cybersecurity Policies';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_44, GETDATE(), GETDATE()
);
DECLARE @ctrl_45 INT;
SELECT @ctrl_45 = id FROM control_library WHERE title_en = N'Revoke Access Upon Termination';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1, @ctrl_45, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_1_1 INT;
SELECT @clause_3_2_1_1 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.1.1';
DECLARE @ctrl_46 INT;
SELECT @ctrl_46 = id FROM control_library WHERE title_en = N'Review Personnel Cybersecurity Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_1, @ctrl_46, GETDATE(), GETDATE()
);
DECLARE @ctrl_47 INT;
SELECT @ctrl_47 = id FROM control_library WHERE title_en = N'Develop Cybersecurity Awareness Program';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_1, @ctrl_47, GETDATE(), GETDATE()
);
DECLARE @ctrl_48 INT;
SELECT @ctrl_48 = id FROM control_library WHERE title_en = N'Implement Cybersecurity Awareness Program';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_1, @ctrl_48, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_1_2 INT;
SELECT @clause_3_2_1_2 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.1.2';
DECLARE @ctrl_49 INT;
SELECT @ctrl_49 = id FROM control_library WHERE title_en = N'Cover Key Cybersecurity Threats in Awareness Program';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_2, @ctrl_49, GETDATE(), GETDATE()
);
DECLARE @ctrl_50 INT;
SELECT @ctrl_50 = id FROM control_library WHERE title_en = N'Train on Secure Email Handling';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_2, @ctrl_50, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_1_3 INT;
SELECT @clause_3_2_1_3 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.1.3';
DECLARE @ctrl_51 INT;
SELECT @ctrl_51 = id FROM control_library WHERE title_en = N'Train on Secure Mobile Device Use';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_51, GETDATE(), GETDATE()
);
DECLARE @ctrl_52 INT;
SELECT @ctrl_52 = id FROM control_library WHERE title_en = N'Train on Secure Internet Browsing';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_52, GETDATE(), GETDATE()
);
DECLARE @ctrl_53 INT;
SELECT @ctrl_53 = id FROM control_library WHERE title_en = N'Train on Secure Social Media Use';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_53, GETDATE(), GETDATE()
);
DECLARE @ctrl_54 INT;
SELECT @ctrl_54 = id FROM control_library WHERE title_en = N'Provide Specialized Cybersecurity Training';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_54, GETDATE(), GETDATE()
);
DECLARE @ctrl_55 INT;
SELECT @ctrl_55 = id FROM control_library WHERE title_en = N'Train Cybersecurity Personnel';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_55, GETDATE(), GETDATE()
);
DECLARE @ctrl_56 INT;
SELECT @ctrl_56 = id FROM control_library WHERE title_en = N'Train Software Development Personnel';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_56, GETDATE(), GETDATE()
);
DECLARE @ctrl_57 INT;
SELECT @ctrl_57 = id FROM control_library WHERE title_en = N'Train Executive Personnel';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_3, @ctrl_57, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_1_4 INT;
SELECT @clause_3_2_1_4 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.1.4';
DECLARE @ctrl_58 INT;
SELECT @ctrl_58 = id FROM control_library WHERE title_en = N'Review Cybersecurity Awareness Program Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_4, @ctrl_58, GETDATE(), GETDATE()
);
DECLARE @ctrl_59 INT;
SELECT @ctrl_59 = id FROM control_library WHERE title_en = N'Define Asset Management Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_1_4, @ctrl_59, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_2 INT;
SELECT @clause_3_2_2 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.2';
DECLARE @ctrl_60 INT;
SELECT @ctrl_60 = id FROM control_library WHERE title_en = N'Implement Asset Management Cybersecurity Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_2, @ctrl_60, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_3 INT;
SELECT @clause_3_2_3 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.3';
DECLARE @ctrl_61 INT;
SELECT @ctrl_61 = id FROM control_library WHERE title_en = N'Define Acceptable Use Policy for Assets';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_3, @ctrl_61, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_4 INT;
SELECT @clause_3_2_4 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.4';
DECLARE @ctrl_62 INT;
SELECT @ctrl_62 = id FROM control_library WHERE title_en = N'Implement Acceptable Use Policy for Assets';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_4, @ctrl_62, GETDATE(), GETDATE()
);
DECLARE @ctrl_63 INT;
SELECT @ctrl_63 = id FROM control_library WHERE title_en = N'Classify and Handle Assets per Regulations';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_4, @ctrl_63, GETDATE(), GETDATE()
);
DECLARE @ctrl_64 INT;
SELECT @ctrl_64 = id FROM control_library WHERE title_en = N'Review Asset Management Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_4, @ctrl_64, GETDATE(), GETDATE()
);
DECLARE @ctrl_65 INT;
SELECT @ctrl_65 = id FROM control_library WHERE title_en = N'Define Identity and Access Management Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_4, @ctrl_65, GETDATE(), GETDATE()
);
DECLARE @ctrl_66 INT;
SELECT @ctrl_66 = id FROM control_library WHERE title_en = N'Implement Identity and Access Management Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_4, @ctrl_66, GETDATE(), GETDATE()
);
DECLARE @clause_3_2_5 INT;
SELECT @clause_3_2_5 = id FROM compliance_library_clauses WHERE clause_number = N'3.2.5';
DECLARE @ctrl_67 INT;
SELECT @ctrl_67 = id FROM control_library WHERE title_en = N'Establish Identity and Access Management Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_5, @ctrl_67, GETDATE(), GETDATE()
);
DECLARE @ctrl_68 INT;
SELECT @ctrl_68 = id FROM control_library WHERE title_en = N'Implement Single-Factor Authentication';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_2_5, @ctrl_68, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_1 INT;
SELECT @clause_3_3_1 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.1';
DECLARE @ctrl_69 INT;
SELECT @ctrl_69 = id FROM control_library WHERE title_en = N'Implement Multi-Factor Authentication';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_1, @ctrl_69, GETDATE(), GETDATE()
);
DECLARE @ctrl_70 INT;
SELECT @ctrl_70 = id FROM control_library WHERE title_en = N'Enforce Access Control Principles';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_1, @ctrl_70, GETDATE(), GETDATE()
);
DECLARE @ctrl_71 INT;
SELECT @ctrl_71 = id FROM control_library WHERE title_en = N'Manage Privileged Access';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_1, @ctrl_71, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_2 INT;
SELECT @clause_3_3_2 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.2';
DECLARE @ctrl_72 INT;
SELECT @ctrl_72 = id FROM control_library WHERE title_en = N'Review User Access Rights Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_2, @ctrl_72, GETDATE(), GETDATE()
);
DECLARE @ctrl_73 INT;
SELECT @ctrl_73 = id FROM control_library WHERE title_en = N'Review Identity and Access Management Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_2, @ctrl_73, GETDATE(), GETDATE()
);
DECLARE @ctrl_74 INT;
SELECT @ctrl_74 = id FROM control_library WHERE title_en = N'Define Information Systems Protection Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_2, @ctrl_74, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_3 INT;
SELECT @clause_3_3_3 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.3';
DECLARE @ctrl_75 INT;
SELECT @ctrl_75 = id FROM control_library WHERE title_en = N'Implement Information Systems Protection Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_3, @ctrl_75, GETDATE(), GETDATE()
);
DECLARE @ctrl_76 INT;
SELECT @ctrl_76 = id FROM control_library WHERE title_en = N'Establish Information Systems Protection Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_3, @ctrl_76, GETDATE(), GETDATE()
);
DECLARE @ctrl_77 INT;
SELECT @ctrl_77 = id FROM control_library WHERE title_en = N'Manage Malware Protection';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_3, @ctrl_77, GETDATE(), GETDATE()
);

DECLARE @clause_3_3_4 INT;
SELECT @clause_3_3_4 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.4';
DECLARE @ctrl_78 INT;
SELECT @ctrl_78 = id FROM control_library WHERE title_en = N'Restrict External Storage Media';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_4, @ctrl_78, GETDATE(), GETDATE()
);
DECLARE @ctrl_79 INT;
SELECT @ctrl_79 = id FROM control_library WHERE title_en = N'Manage Software Patch Deployment';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_4, @ctrl_79, GETDATE(), GETDATE()
);
DECLARE @ctrl_80 INT;
SELECT @ctrl_80 = id FROM control_library WHERE title_en = N'Synchronize System Clocks';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_4, @ctrl_80, GETDATE(), GETDATE()
);

DECLARE @clause_3_3_5 INT;
SELECT @clause_3_3_5 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.5';
DECLARE @ctrl_81 INT;
SELECT @ctrl_81 = id FROM control_library WHERE title_en = N'Review Information Systems Protection Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_5, @ctrl_81, GETDATE(), GETDATE()
);
DECLARE @ctrl_82 INT;
SELECT @ctrl_82 = id FROM control_library WHERE title_en = N'Define Email Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_5, @ctrl_82, GETDATE(), GETDATE()
);
DECLARE @ctrl_83 INT;
SELECT @ctrl_83 = id FROM control_library WHERE title_en = N'Implement Email Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_5, @ctrl_83, GETDATE(), GETDATE()
);
DECLARE @ctrl_84 INT;
SELECT @ctrl_84 = id FROM control_library WHERE title_en = N'Establish Email Security Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_5, @ctrl_84, GETDATE(), GETDATE()
);

DECLARE @clause_3_3_6 INT;
SELECT @clause_3_3_6 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.6';
DECLARE @ctrl_85 INT;
SELECT @ctrl_85 = id FROM control_library WHERE title_en = N'Filter Phishing and Spam Emails';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_6, @ctrl_85, GETDATE(), GETDATE()
);
DECLARE @ctrl_86 INT;
SELECT @ctrl_86 = id FROM control_library WHERE title_en = N'Secure Email with Multi-Factor Authentication';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_6, @ctrl_86, GETDATE(), GETDATE()
);
DECLARE @ctrl_87 INT;
SELECT @ctrl_87 = id FROM control_library WHERE title_en = N'Archive and Backup Emails';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_6, @ctrl_87, GETDATE(), GETDATE()
);
DECLARE @ctrl_88 INT;
SELECT @ctrl_88 = id FROM control_library WHERE title_en = N'Protect Against Advanced Threats';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_6, @ctrl_88, GETDATE(), GETDATE()
);
DECLARE @ctrl_89 INT;
SELECT @ctrl_89 = id FROM control_library WHERE title_en = N'Secure Email Domains';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_6, @ctrl_89, GETDATE(), GETDATE()
);

DECLARE @clause_3_3_7 INT;
SELECT @clause_3_3_7 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.7';
DECLARE @ctrl_90 INT;
SELECT @ctrl_90 = id FROM control_library WHERE title_en = N'Review Email Security Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_7, @ctrl_90, GETDATE(), GETDATE()
);
DECLARE @ctrl_91 INT;
SELECT @ctrl_91 = id FROM control_library WHERE title_en = N'Define Network Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_7, @ctrl_91, GETDATE(), GETDATE()
);
DECLARE @ctrl_92 INT;
SELECT @ctrl_92 = id FROM control_library WHERE title_en = N'Implement Network Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_7, @ctrl_92, GETDATE(), GETDATE()
);
DECLARE @ctrl_93 INT;
SELECT @ctrl_93 = id FROM control_library WHERE title_en = N'Establish Network Security Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_7, @ctrl_93, GETDATE(), GETDATE()
);

DECLARE @clause_3_3_8 INT;
SELECT @clause_3_3_8 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.8';
DECLARE @ctrl_94 INT;
SELECT @ctrl_94 = id FROM control_library WHERE title_en = N'Segment Network Infrastructure';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_94, GETDATE(), GETDATE()
);
DECLARE @ctrl_95 INT;
SELECT @ctrl_95 = id FROM control_library WHERE title_en = N'Segregate Network Environments';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_95, GETDATE(), GETDATE()
);
DECLARE @ctrl_96 INT;
SELECT @ctrl_96 = id FROM control_library WHERE title_en = N'Secure Internet Connectivity';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_96, GETDATE(), GETDATE()
);
DECLARE @ctrl_97 INT;
SELECT @ctrl_97 = id FROM control_library WHERE title_en = N'Protect Wireless Networks';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_97, GETDATE(), GETDATE()
);
DECLARE @ctrl_98 INT;
SELECT @ctrl_98 = id FROM control_library WHERE title_en = N'Manage Network Services and Ports';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_98, GETDATE(), GETDATE()
);
DECLARE @ctrl_99 INT;
SELECT @ctrl_99 = id FROM control_library WHERE title_en = N'Deploy Intrusion Prevention Systems';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_8, @ctrl_99, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_9 INT;
SELECT @clause_3_3_9 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.9';
DECLARE @ctrl_100 INT;
SELECT @ctrl_100 = id FROM control_library WHERE title_en = N'Secure Domain Name Service (DNS)';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_9, @ctrl_100, GETDATE(), GETDATE()
);
DECLARE @ctrl_101 INT;
SELECT @ctrl_101 = id FROM control_library WHERE title_en = N'Protect Browsing Against Advanced Threats';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_9, @ctrl_101, GETDATE(), GETDATE()
);
DECLARE @ctrl_102 INT;
SELECT @ctrl_102 = id FROM control_library WHERE title_en = N'Protect Against DDoS Attacks';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_9, @ctrl_102, GETDATE(), GETDATE()
);
DECLARE @ctrl_103 INT;
SELECT @ctrl_103 = id FROM control_library WHERE title_en = N'Review Network Security Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_9, @ctrl_103, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_10 INT;
SELECT @clause_3_3_10 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.10';
DECLARE @ctrl_104 INT;
SELECT @ctrl_104 = id FROM control_library WHERE title_en = N'Define Mobile Device Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_10, @ctrl_104, GETDATE(), GETDATE()
);
DECLARE @ctrl_105 INT;
SELECT @ctrl_105 = id FROM control_library WHERE title_en = N'Implement Mobile Device Security Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_10, @ctrl_105, GETDATE(), GETDATE()
);
DECLARE @ctrl_106 INT;
SELECT @ctrl_106 = id FROM control_library WHERE title_en = N'Establish Mobile Device Security Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_10, @ctrl_106, GETDATE(), GETDATE()
);
DECLARE @ctrl_107 INT;
SELECT @ctrl_107 = id FROM control_library WHERE title_en = N'Secure Mobile Device Data';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_10, @ctrl_107, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_11 INT;
SELECT @clause_3_3_11 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.11';
DECLARE @ctrl_108 INT;
SELECT @ctrl_108 = id FROM control_library WHERE title_en = N'Restrict Mobile Device Use';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_108, GETDATE(), GETDATE()
);
DECLARE @ctrl_109 INT;
SELECT @ctrl_109 = id FROM control_library WHERE title_en = N'Secure Data Wiping on Mobile Devices';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_109, GETDATE(), GETDATE()
);
DECLARE @ctrl_110 INT;
SELECT @ctrl_110 = id FROM control_library WHERE title_en = N'Promote Mobile Device Security Awareness';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_110, GETDATE(), GETDATE()
);
DECLARE @ctrl_111 INT;
SELECT @ctrl_111 = id FROM control_library WHERE title_en = N'Review Mobile Device Security Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_111, GETDATE(), GETDATE()
);
DECLARE @ctrl_112 INT;
SELECT @ctrl_112 = id FROM control_library WHERE title_en = N'Define Data Protection Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_112, GETDATE(), GETDATE()
);
DECLARE @ctrl_113 INT;
SELECT @ctrl_113 = id FROM control_library WHERE title_en = N'Implement Data Protection Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_11, @ctrl_113, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_12 INT;
SELECT @clause_3_3_12 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.12';
DECLARE @ctrl_114 INT;
SELECT @ctrl_114 = id FROM control_library WHERE title_en = N'Comply with Data Protection Standards';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_12, @ctrl_114, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_13 INT;
SELECT @clause_3_3_13 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.13';
DECLARE @ctrl_115 INT;
SELECT @ctrl_115 = id FROM control_library WHERE title_en = N'Review Data Protection Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_13, @ctrl_115, GETDATE(), GETDATE()
);
DECLARE @ctrl_116 INT;
SELECT @ctrl_116 = id FROM control_library WHERE title_en = N'Define Cryptography Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_13, @ctrl_116, GETDATE(), GETDATE()
);
DECLARE @ctrl_117 INT;
SELECT @ctrl_117 = id FROM control_library WHERE title_en = N'Implement Cryptography Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_13, @ctrl_117, GETDATE(), GETDATE()
);
DECLARE @ctrl_118 INT;
SELECT @ctrl_118 = id FROM control_library WHERE title_en = N'Establish Cryptography Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_13, @ctrl_118, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_14 INT;
SELECT @clause_3_3_14 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.14';
DECLARE @ctrl_119 INT;
SELECT @ctrl_119 = id FROM control_library WHERE title_en = N'Use Approved Cryptographic Systems';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_14, @ctrl_119, GETDATE(), GETDATE()
);
DECLARE @ctrl_120 INT;
SELECT @ctrl_120 = id FROM control_library WHERE title_en = N'Manage Cryptographic Keys';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_14, @ctrl_120, GETDATE(), GETDATE()
);
DECLARE @ctrl_121 INT;
SELECT @ctrl_121 = id FROM control_library WHERE title_en = N'Encrypt Data Based on Classification';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_14, @ctrl_121, GETDATE(), GETDATE()
);
DECLARE @ctrl_122 INT;
SELECT @ctrl_122 = id FROM control_library WHERE title_en = N'Review Cryptography Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_14, @ctrl_122, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_15 INT;
SELECT @clause_3_3_15 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.15';
DECLARE @ctrl_123 INT;
SELECT @ctrl_123 = id FROM control_library WHERE title_en = N'Define Backup and Recovery Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_123, GETDATE(), GETDATE()
);
DECLARE @ctrl_124 INT;
SELECT @ctrl_124 = id FROM control_library WHERE title_en = N'Implement Backup and Recovery Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_124, GETDATE(), GETDATE()
);
DECLARE @ctrl_125 INT;
SELECT @ctrl_125 = id FROM control_library WHERE title_en = N'Establish Backup and Recovery Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_125, GETDATE(), GETDATE()
);
DECLARE @ctrl_126 INT;
SELECT @ctrl_126 = id FROM control_library WHERE title_en = N'Define Backup Scope for Critical Assets';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_126, GETDATE(), GETDATE()
);
DECLARE @ctrl_127 INT;
SELECT @ctrl_127 = id FROM control_library WHERE title_en = N'Enable Quick System Recovery';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_127, GETDATE(), GETDATE()
);
DECLARE @ctrl_128 INT;
SELECT @ctrl_128 = id FROM control_library WHERE title_en = N'Test Backups Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_128, GETDATE(), GETDATE()
);
DECLARE @ctrl_129 INT;
SELECT @ctrl_129 = id FROM control_library WHERE title_en = N'Review Backup and Recovery Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_15, @ctrl_129, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_16 INT;
SELECT @clause_3_3_16 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.16';
DECLARE @ctrl_130 INT;
SELECT @ctrl_130 = id FROM control_library WHERE title_en = N'Define Vulnerability Management Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_16, @ctrl_130, GETDATE(), GETDATE()
);
DECLARE @ctrl_131 INT;
SELECT @ctrl_131 = id FROM control_library WHERE title_en = N'Implement Vulnerability Management Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_16, @ctrl_131, GETDATE(), GETDATE()
);
DECLARE @ctrl_132 INT;
SELECT @ctrl_132 = id FROM control_library WHERE title_en = N'Establish Vulnerability Management Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_16, @ctrl_132, GETDATE(), GETDATE()
);
DECLARE @clause_3_3_17 INT;
SELECT @clause_3_3_17 = id FROM compliance_library_clauses WHERE clause_number = N'3.3.17';
DECLARE @ctrl_133 INT;
SELECT @ctrl_133 = id FROM control_library WHERE title_en = N'Conduct Periodic Vulnerability Assessments';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_17, @ctrl_133, GETDATE(), GETDATE()
);
DECLARE @ctrl_134 INT;
SELECT @ctrl_134 = id FROM control_library WHERE title_en = N'Classify Vulnerabilities by Criticality';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_17, @ctrl_134, GETDATE(), GETDATE()
);
DECLARE @ctrl_135 INT;
SELECT @ctrl_135 = id FROM control_library WHERE title_en = N'Remediate Vulnerabilities Based on Risk';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_3_17, @ctrl_135, GETDATE(), GETDATE()
);
DECLARE @clause_3_4_1 INT;
SELECT @clause_3_4_1 = id FROM compliance_library_clauses WHERE clause_number = N'3.4.1';
DECLARE @ctrl_136 INT;
SELECT @ctrl_136 = id FROM control_library WHERE title_en = N'Manage Security Patches';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_136, GETDATE(), GETDATE()
);
DECLARE @ctrl_137 INT;
SELECT @ctrl_137 = id FROM control_library WHERE title_en = N'Subscribe to Vulnerability Information Sources';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_137, GETDATE(), GETDATE()
);
DECLARE @ctrl_138 INT;
SELECT @ctrl_138 = id FROM control_library WHERE title_en = N'Review Vulnerability Management Requirements Periodically';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_138, GETDATE(), GETDATE()
);
DECLARE @ctrl_139 INT;
SELECT @ctrl_139 = id FROM control_library WHERE title_en = N'Define Penetration Testing Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_139, GETDATE(), GETDATE()
);
DECLARE @ctrl_140 INT;
SELECT @ctrl_140 = id FROM control_library WHERE title_en = N'Implement Penetration Testing Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_140, GETDATE(), GETDATE()
);
DECLARE @ctrl_141 INT;
SELECT @ctrl_141 = id FROM control_library WHERE title_en = N'Establish Penetration Testing Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_1, @ctrl_141, GETDATE(), GETDATE()
);
DECLARE @clause_3_4_2 INT;
SELECT @clause_3_4_2 = id FROM compliance_library_clauses WHERE clause_number = N'3.4.2';
DECLARE @ctrl_142 INT;
SELECT @ctrl_142 = id FROM control_library WHERE title_en = N'Define Penetration Testing Scope';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_2, @ctrl_142, GETDATE(), GETDATE()
);
DECLARE @ctrl_143 INT;
SELECT @ctrl_143 = id FROM control_library WHERE title_en = N'Conduct Periodic Penetration Tests';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_2, @ctrl_143, GETDATE(), GETDATE()
);
DECLARE @ctrl_144 INT;
SELECT @ctrl_144 = id FROM control_library WHERE title_en = N'Define Event Logging and Monitoring Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_2, @ctrl_144, GETDATE(), GETDATE()
);
DECLARE @clause_3_4_3 INT;
SELECT @clause_3_4_3 = id FROM compliance_library_clauses WHERE clause_number = N'3.4.3';
DECLARE @ctrl_145 INT;
SELECT @ctrl_145 = id FROM control_library WHERE title_en = N'Implement Event Logging and Monitoring Requirements';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_3, @ctrl_145, GETDATE(), GETDATE()
);
DECLARE @ctrl_146 INT;
SELECT @ctrl_146 = id FROM control_library WHERE title_en = N'Establish Event Logging and Monitoring Controls';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_3, @ctrl_146, GETDATE(), GETDATE()
);
DECLARE @ctrl_147 INT;
SELECT @ctrl_147 = id FROM control_library WHERE title_en = N'Enable Event Logging for Critical Assets';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_3, @ctrl_147, GETDATE(), GETDATE()
);
DECLARE @ctrl_148 INT;
SELECT @ctrl_148 = id FROM control_library WHERE title_en = N'Log Remote and Privileged Access';
INSERT INTO compliance_library_clauses_controls (
    compliance_library_clause_id,
    control_library_id,
    created_at,
    updated_at
) VALUES (
    @clause_3_4_3, @ctrl_148, GETDATE(), GETDATE()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE compliance_library_clauses_controls;
-- +goose StatementEnd