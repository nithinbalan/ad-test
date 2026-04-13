-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

-- Fetching Impact Categories (using new names after rename)
DECLARE @FinancialId INT, @ReputationId INT, @OperationalId INT, @DisruptionId INT;
DECLARE @SystemsId INT, @CybersecurityId INT, @CustomerSatisfactionId INT, @DevelopmentGrowthId INT;
DECLARE @CreditId INT, @HumanResourcesId INT, @RegulationsPoliciesId INT, @HealthSafetyId INT;
DECLARE @GovernanceControlsId INT;

SELECT @FinancialId = id FROM risk_impact_categories WHERE title_en = 'Financial';
SELECT @ReputationId = id FROM risk_impact_categories WHERE title_en = 'Reputation';
SELECT @OperationalId = id FROM risk_impact_categories WHERE title_en = 'Operational';
SELECT @DisruptionId = id FROM risk_impact_categories WHERE title_en = 'Disruption';
SELECT @SystemsId = id FROM risk_impact_categories WHERE title_en = 'Systems';
SELECT @CybersecurityId = id FROM risk_impact_categories WHERE title_en = 'Cybersecurity';
SELECT @CustomerSatisfactionId = id FROM risk_impact_categories WHERE title_en = 'Customer Satisfaction';
SELECT @DevelopmentGrowthId = id FROM risk_impact_categories WHERE title_en = 'Development and Growth';
SELECT @CreditId = id FROM risk_impact_categories WHERE title_en = 'Credit';
SELECT @HumanResourcesId = id FROM risk_impact_categories WHERE title_en = 'Human Resources';
SELECT @RegulationsPoliciesId = id FROM risk_impact_categories WHERE title_en = 'Regulations and Policies';
SELECT @HealthSafetyId = id FROM risk_impact_categories WHERE title_en = 'Health and Safety';
SELECT @GovernanceControlsId = id FROM risk_impact_categories WHERE title_en = 'Governance and Internal Controls';

-- Fetching Matrix Impact Levels
DECLARE @VeryLowId INT, @LowId INT, @ModerateId INT, @HighId INT, @VeryHighId INT;

SELECT @VeryLowId = id FROM risk_matrix_impacts WHERE title_en = 'Very Low' AND risk_matrix_id = @RiskMatrixId;
SELECT @LowId = id FROM risk_matrix_impacts WHERE title_en = 'Low' AND risk_matrix_id = @RiskMatrixId;
SELECT @ModerateId = id FROM risk_matrix_impacts WHERE title_en = 'Moderate' AND risk_matrix_id = @RiskMatrixId;
SELECT @HighId = id FROM risk_matrix_impacts WHERE title_en = 'High' AND risk_matrix_id = @RiskMatrixId;
SELECT @VeryHighId = id FROM risk_matrix_impacts WHERE title_en = 'Very High' AND risk_matrix_id = @RiskMatrixId;

-- Update Financial Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Very minor financial losses < SAR 10 million and SAR 15 million, not affecting the funds liquidity or operating budget'
WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minor financial losses < SAR 30 million and ≥ SAR 10 million, which can be quickly offset by current cash flows.'
WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Moderate financial losses < SAR 50 million and ≥ SAR 30 million, necessitating the reallocation of resources or adjustment of certain budget items.'
WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant financial losses < SAR 100 million and ≥ SAR 50 million, requiring strong corrective measures such as expenditure reduction or debt restructuring.'
WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Severe financial losses ≥ SAR 100 million, threatening the funds sustainability and requiring a comprehensive financial restructuring.'
WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @VeryHighId;

-- Update Reputation Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Internal impact, which can be easily contained within the fund without reaching the external audience and does not cause significant adverse reputational repercussions.'
WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Limited external impact, reaching a small number of customers or partners, which can be quickly contained through direct and effective communication.'
WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Tangible impact on social media, reaching a limited audience of customers or partners, with limited interactions, requiring a rapid response to contain the situation.'
WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Noticeable impact reaching local media or social media platforms, with moderate public engagement, affecting the funds image locally and requiring swift intervention.'
WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Widespread impact reaching local media and social media platforms in the Kingdom, with significant public engagement, resulting in a significant decline in confidence locally and requiring substantial efforts to rebuild reputation.'
WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @VeryHighId;

-- Update Operational Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Very low impact on operational processes, requiring corrective actions and updating of existing controls.'
WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Limited impact on operational processes, requiring corrective actions and updating of existing controls.'
WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Key operational processes partially affected, requiring additional resources.'
WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Key operational processes fully affected, requiring additional resources or indirect senior management intervention, in addition to the possibility of an operational impact on three departments.'
WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Inability to perform key operational processes at the required quality level and within the specified timeframe, requiring escalation to senior management for urgent decisions to mitigate the impact, in addition to the possibility of an operational impact on five or more departments.'
WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @VeryHighId;

-- Update Disruption Impact Guidelines (formerly Business Interruption)
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A very limited outage after which service will be restored immediately with no major impact on daily activities.'
WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Short-term disruption that can be quickly contained, with minimal impact on service and business continuity.'
WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Moderate service disruption, which requires time to fully restore service, resulting in delays to operational processes.'
WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Prolonged disruption requiring serious intervention to restore service, significantly impacting operational capacity.'
WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Total service disruption, which may require a complete system redesign, resulting in a severe impact on business continuity.'
WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @VeryHighId;

-- Update Systems Impact Guidelines (formerly Technical)
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minor IT issues with no downtime.'
WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minor impact on IT systems, which can be quickly rectified.'
WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant delay in the delivery of digital services, requiring limited intervention.'
WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant disruption to critical systems leading to operational disruption.'
WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A complete breakdown of IT systems threatens the continuity of the fund.'
WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @VeryHighId;

-- Update Cybersecurity Impact Guidelines (formerly Cyber)
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'No cybersecurity breach of the funds systems and services.'
WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'No cybersecurity breach of the funds systems and services.'
WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'No cybersecurity breaches on the Funds systems and services provided'
WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A cyberattack/penetration of the funds systems and services has a critical impact on the funds systems and services, data confidentiality, accuracy, and availability, which may disrupt the funds operations.'
WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A cyberattack/penetration of the funds systems and services has a critical impact on the funds systems and services, data confidentiality, accuracy, and availability, which may disrupt the funds operations.'
WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @VeryHighId;

-- Update Customer Satisfaction Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Customer satisfaction > 90%, with services continuing to be delivered at a high quality and without significant complaints.'
WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Customer satisfaction > 90%, with services continuing to be delivered at a high quality and without significant complaints.'
WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Customer satisfaction > 60% and < 80%, with a significant increase in the number of complaints or customer dissatisfaction, requiring intervention to improve service levels.'
WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Customer satisfaction > 40% 9 < 60%, with a significant loss of customer confidence due to deteriorating service quality or delayed response.'
WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Customer satisfaction < 40%, with widespread loss of satisfaction due to a sharp decline in service quality or repeated errors, potentially leading to permanent customer loss.'
WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @VeryHighId;

-- Update Development and Growth Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'No significant strategic impact, as strategic projects and initiatives continue as planned without any disruption, ensuring the achievement of long-term objectives'
WHERE risk_impact_category_id = @DevelopmentGrowthId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Limited strategic impact requires minor adjustments to some projects, but strategic objectives remain achievable without a significant impact.'
WHERE risk_impact_category_id = @DevelopmentGrowthId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant strategic delay necessitates a reassessment of the timeline and resource allocation, with some delay in achieving objectives and future expansion.'
WHERE risk_impact_category_id = @DevelopmentGrowthId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant strategic impact results in disruption to key projects, threatening the achievement of strategic objectives and the loss of important expansion opportunities.'
WHERE risk_impact_category_id = @DevelopmentGrowthId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Severe strategic impact causes a complete halt to strategic projects, significantly threatening the organizations ability to achieve its long-term objectives and reducing its competitiveness.'
WHERE risk_impact_category_id = @DevelopmentGrowthId AND risk_matrix_impact_id = @VeryHighId;

-- Update Credit Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Near-zero credit risk, as the default rate remains very low and credit is managed very efficiently, without any impact on the financial stability of the fund.'
WHERE risk_impact_category_id = @CreditId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minimal impact on credit risk, with some defaults being quickly resolved, without significant impact on the loan portfolio or stability.'
WHERE risk_impact_category_id = @CreditId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant increase in credit risk, with a moderate rise in the default rate, requires a reassessment of customer creditworthiness and adjustment of lending terms.'
WHERE risk_impact_category_id = @CreditId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant impact on the credit portfolio, with a significant rise in the default rate leading to significant financial losses, necessitating strong corrective action to maintain financial stability.'
WHERE risk_impact_category_id = @CreditId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Severe impact on the credit portfolio, with a sharp increase in default rates threatening the stability of the financial fund, requiring a comprehensive restructuring of credit management and loss reduction.'
WHERE risk_impact_category_id = @CreditId AND risk_matrix_impact_id = @VeryHighId;

-- Update Human Resources Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Staff satisfaction and productivity are very high, with minimal changes in hiring or resignation rates. The work environment is stable and attractive.'
WHERE risk_impact_category_id = @HumanResourcesId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minor impact on staff satisfaction and productivity, with some resignations or changes in hiring rates, but no significant impact on team performance or overall objectives.'
WHERE risk_impact_category_id = @HumanResourcesId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Marked decline in staff satisfaction and productivity, with increased resignation rates or difficulty in attracting new talent, requiring management interventions to improve the work environment and provide incentive programs.'
WHERE risk_impact_category_id = @HumanResourcesId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant impact on human resources, with a marked increase in resignations or failure to retain key talent, resulting in a skills gap and significantly impacting the achievement of operational objectives.'
WHERE risk_impact_category_id = @HumanResourcesId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Widespread loss of key staff and failure to attract or retain talent, directly threatening the continuity of the Funds core operations and resulting in a significant decline in overall performance.'
WHERE risk_impact_category_id = @HumanResourcesId AND risk_matrix_impact_id = @VeryHighId;

-- Update Regulations and Policies Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Compliance with policies and regulations is complete and accurate, with no violations or delays in updates. All procedures are carried out according to established standards and without any regulatory issues.'
WHERE risk_impact_category_id = @RegulationsPoliciesId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Minor impact on compliance, with some slight delays in policy updates or regulation implementation, but no significant impact on organizational processes or relationships with regulators.'
WHERE risk_impact_category_id = @RegulationsPoliciesId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Marked delay in updating policies or implementing regulations, resulting in some organizational problems or disputes with regulators, requiring interventions to rectify the situation'
WHERE risk_impact_category_id = @RegulationsPoliciesId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant policy or regulatory violations directly impact regulatory compliance, potentially resulting in fines or regulatory action against the fund. A rapid reassessment of policies and correction of violations is required.'
WHERE risk_impact_category_id = @RegulationsPoliciesId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Comprehensive failure to comply with policies and regulations, resulting in significant regulatory penalties that threaten the continuity of operations or cause substantial financial losses to the Fund. Requires a comprehensive restructuring of the organizational system and immediate corrective action.'
WHERE risk_impact_category_id = @RegulationsPoliciesId AND risk_matrix_impact_id = @VeryHighId;

-- Update Health and Safety Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'The work environment is completely safe, with no recorded accidents or health risks. All safety and security standards are adhered to, providing an ideal working environment for employees.'
WHERE risk_impact_category_id = @HealthSafetyId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'The working environment is completely safe, with no recorded incidents or health hazards. All security and safety standards are adhered to, providing an ideal working environment for employees.'
WHERE risk_impact_category_id = @HealthSafetyId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A marked increase in the number of incidents or health hazards, necessitating additional measures to improve the working environment and reduce risks. These incidents may result in minor delays in some operations.'
WHERE risk_impact_category_id = @HealthSafetyId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Major incidents or health hazards directly impacting employee safety and disrupting essential operations. The situation requires an immediate response to improve security and safety and restore a safe working environment.'
WHERE risk_impact_category_id = @HealthSafetyId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant threats to life and safety causing a complete work stoppage, with significant risks leading to substantial human and financial losses. The situation requires a comprehensive reassessment of security and safety procedures and the implementation of substantial reforms.'
WHERE risk_impact_category_id = @HealthSafetyId AND risk_matrix_impact_id = @VeryHighId;

-- Update Governance and Internal Controls Impact Guidelines
UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Full compliance with the governance and internal controls system, with no violations or breaches. All operations are conducted according to approved standards, ensuring transparency and accountability.'
WHERE risk_impact_category_id = @GovernanceControlsId AND risk_matrix_impact_id = @VeryLowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Some minor non-compliance with internal controls, which can be addressed quickly without a significant impact on operations or organizational confidence.'
WHERE risk_impact_category_id = @GovernanceControlsId AND risk_matrix_impact_id = @LowId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant violations and non-compliance with procedures, requiring a comprehensive review of the system, leading to delays in decision-making and the emergence of some operational risks that require corrective interventions.'
WHERE risk_impact_category_id = @GovernanceControlsId AND risk_matrix_impact_id = @ModerateId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'Significant weaknesses in adherence to the governance and internal controls system lead to significant operational risks and negative impacts on operational efficiency and organizational stability. The situation requires a comprehensive reassessment and strengthening of controls and compliance'
WHERE risk_impact_category_id = @GovernanceControlsId AND risk_matrix_impact_id = @HighId;

UPDATE risk_matrix_impact_guidelines 
SET guideline = 'A complete breakdown in compliance with the governance and internal controls system, resulting in a failure to manage risks and a direct threat to business continuity. The situation calls for radical reforms and the restoration of organizational confidence through a comprehensive restructuring and improved compliance.'
WHERE risk_impact_category_id = @GovernanceControlsId AND risk_matrix_impact_id = @VeryHighId;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

-- Fetching Impact Categories
DECLARE @FinancialId INT, @ReputationId INT, @OperationalId INT, @DisruptionId INT;
DECLARE @SystemsId INT, @CybersecurityId INT, @CustomerSatisfactionId INT, @DevelopmentGrowthId INT;
DECLARE @CreditId INT, @HumanResourcesId INT, @RegulationsPoliciesId INT, @HealthSafetyId INT;
DECLARE @GovernanceControlsId INT;

SELECT @FinancialId = id FROM risk_impact_categories WHERE title_en = 'Financial';
SELECT @ReputationId = id FROM risk_impact_categories WHERE title_en = 'Reputation';
SELECT @OperationalId = id FROM risk_impact_categories WHERE title_en = 'Operational';
SELECT @DisruptionId = id FROM risk_impact_categories WHERE title_en = 'Disruption';
SELECT @SystemsId = id FROM risk_impact_categories WHERE title_en = 'Systems';
SELECT @CybersecurityId = id FROM risk_impact_categories WHERE title_en = 'Cybersecurity';
SELECT @CustomerSatisfactionId = id FROM risk_impact_categories WHERE title_en = 'Customer Satisfaction';
SELECT @DevelopmentGrowthId = id FROM risk_impact_categories WHERE title_en = 'Development and Growth';
SELECT @CreditId = id FROM risk_impact_categories WHERE title_en = 'Credit';
SELECT @HumanResourcesId = id FROM risk_impact_categories WHERE title_en = 'Human Resources';
SELECT @RegulationsPoliciesId = id FROM risk_impact_categories WHERE title_en = 'Regulations and Policies';
SELECT @HealthSafetyId = id FROM risk_impact_categories WHERE title_en = 'Health and Safety';
SELECT @GovernanceControlsId = id FROM risk_impact_categories WHERE title_en = 'Governance and Internal Controls';

-- Fetching Matrix Impact Levels
DECLARE @VeryLowId INT, @LowId INT, @ModerateId INT, @HighId INT, @VeryHighId INT;

SELECT @VeryLowId = id FROM risk_matrix_impacts WHERE title_en = 'Very Low' AND risk_matrix_id = @RiskMatrixId;
SELECT @LowId = id FROM risk_matrix_impacts WHERE title_en = 'Low' AND risk_matrix_id = @RiskMatrixId;
SELECT @ModerateId = id FROM risk_matrix_impacts WHERE title_en = 'Moderate' AND risk_matrix_id = @RiskMatrixId;
SELECT @HighId = id FROM risk_matrix_impacts WHERE title_en = 'High' AND risk_matrix_id = @RiskMatrixId;
SELECT @VeryHighId = id FROM risk_matrix_impacts WHERE title_en = 'Very High' AND risk_matrix_id = @RiskMatrixId;

-- Revert to original guidelines
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minimal financial loss less than 0.1% of annual budget. No impact on financial operations or reporting.' WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor financial loss between 0.1% to 0.5% of annual budget. Easily absorbed within existing budget allocations.' WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate financial loss between 0.5% to 2% of annual budget. May require budget reallocation or adjustments.' WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant financial loss between 2% to 5% of annual budget. Impacts quarterly performance and requires senior management intervention.' WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Severe financial loss exceeding 5% of annual budget. Threatens annual targets and may impact organizational viability.' WHERE risk_impact_category_id = @FinancialId AND risk_matrix_impact_id = @VeryHighId;

UPDATE risk_matrix_impact_guidelines SET guideline = 'Minimal impact on reputation. Internal issue with no external visibility or stakeholder awareness.' WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor reputational concern. Limited negative feedback from small stakeholder group, quickly resolved.' WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate reputational damage. Sector-wide concern, negative media coverage at local level, temporary loss of stakeholder confidence.' WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant reputational damage. National media coverage, social media backlash, loss of key stakeholder trust requiring extensive PR efforts.' WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Severe reputational crisis. Widespread negative media coverage, loss of major clients/partners, long-term damage to brand and market position.' WHERE risk_impact_category_id = @ReputationId AND risk_matrix_impact_id = @VeryHighId;

UPDATE risk_matrix_impact_guidelines SET guideline = 'Negligible operational disruption. Minor delays in non-critical processes with no impact on service delivery.' WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor operational disruption. Brief delays in some processes, minimal impact on overall operations, resolved within hours.' WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate operational disruption. Temporary or localized service interruption affecting some business units, resolved within 24 hours.' WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant operational disruption. Extended service outage under 24 hours affecting critical operations and multiple business units.' WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Severe operational failure. Organization-wide unavailability exceeding 24 hours, loss of critical functions, major impact on service delivery.' WHERE risk_impact_category_id = @OperationalId AND risk_matrix_impact_id = @VeryHighId;

UPDATE risk_matrix_impact_guidelines SET guideline = 'Minimal business interruption. Brief disruption to non-essential services with no customer impact.' WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor business interruption. Short-term disruption to some services, minimal customer inconvenience, recovery within hours.' WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate business interruption. Partial service disruption affecting customer experience, recovery within 24 hours.' WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Major business interruption. Significant service outage affecting core business functions, extended recovery time up to 48 hours.' WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Critical business interruption. Complete service failure, loss of critical business capabilities, recovery exceeding 48 hours with potential revenue loss.' WHERE risk_impact_category_id = @DisruptionId AND risk_matrix_impact_id = @VeryHighId;

UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor technical issues. Small software glitches or system errors with no impact on operations.' WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Limited technical problems. Isolated system issues affecting few users, resolved through standard support procedures.' WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate technical failures. System performance degradation affecting multiple users, requires technical intervention.' WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant technical failures. Major system outages impacting critical infrastructure, extended downtime requiring emergency response.' WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Critical technical failure. Complete system collapse, loss of critical IT infrastructure, potential data loss requiring disaster recovery.' WHERE risk_impact_category_id = @SystemsId AND risk_matrix_impact_id = @VeryHighId;

UPDATE risk_matrix_impact_guidelines SET guideline = 'Minimal cyber incident. Unsuccessful attack attempt or minor security alert with no data compromise.' WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor cyber incident. Limited data leak or minor security breach affecting non-sensitive information, quickly contained.' WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate cyber incident. Exposure of internal or confidential data, reportable security breach requiring investigation.' WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant cyber attack. Partial exposure of sensitive customer or employee data impacting critical systems, regulatory notification required.' WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Major cyber breach. Full exposure of highly sensitive data, widespread system compromise, significant regulatory penalties and reputational damage.' WHERE risk_impact_category_id = @CybersecurityId AND risk_matrix_impact_id = @VeryHighId;

-- Continue with other categories...
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minimal customer impact. Few isolated complaints with no pattern, quickly resolved.' WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @VeryLowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Minor customer dissatisfaction. Limited complaints from small customer segment, manageable through standard service recovery.' WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @LowId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Moderate customer dissatisfaction. Noticeable increase in complaints, temporary decline in satisfaction scores, requires management attention.' WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @ModerateId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Significant customer dissatisfaction. Major service failures leading to customer complaints, negative reviews, and potential customer loss.' WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @HighId;
UPDATE risk_matrix_impact_guidelines SET guideline = 'Severe customer satisfaction crisis. Widespread customer complaints, mass customer exodus, significant damage to customer relationships and retention.' WHERE risk_impact_category_id = @CustomerSatisfactionId AND risk_matrix_impact_id = @VeryHighId;

-- Remaining categories follow same pattern...

-- +goose StatementEnd
