-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

-- Fetching Impact Categories
DECLARE @FinancialId INT, @ReputationId INT, @OperationalId INT, @BusinessInterruptionId INT;
DECLARE @TechnicalId INT, @CyberId INT, @CustomerSatisfactionId INT, @DevelopmentGrowthId INT;
DECLARE @CreditId INT, @HumanResourcesId INT, @RegulationsPoliciesId INT, @HealthSafetyId INT;
DECLARE @GovernanceControlsId INT;

SELECT @FinancialId = id FROM risk_impact_categories WHERE title_en = 'Financial';
SELECT @ReputationId = id FROM risk_impact_categories WHERE title_en = 'Reputation';
SELECT @OperationalId = id FROM risk_impact_categories WHERE title_en = 'Operational';
SELECT @BusinessInterruptionId = id FROM risk_impact_categories WHERE title_en = 'Business Interruption';
SELECT @TechnicalId = id FROM risk_impact_categories WHERE title_en = 'Technical';
SELECT @CyberId = id FROM risk_impact_categories WHERE title_en = 'Cyber';
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

-- Seeding Risk Matrix Impact Guidelines
INSERT INTO risk_matrix_impact_guidelines (risk_impact_category_id, risk_matrix_impact_id, guideline, status, created_by_id, created_at)
VALUES
    -- Financial Impact
    (@FinancialId, @VeryLowId, 'Minimal financial loss less than 0.1% of annual budget. No impact on financial operations or reporting.', 'active', 1, GETDATE()),
    (@FinancialId, @LowId, 'Minor financial loss between 0.1% to 0.5% of annual budget. Easily absorbed within existing budget allocations.', 'active', 1, GETDATE()),
    (@FinancialId, @ModerateId, 'Moderate financial loss between 0.5% to 2% of annual budget. May require budget reallocation or adjustments.', 'active', 1, GETDATE()),
    (@FinancialId, @HighId, 'Significant financial loss between 2% to 5% of annual budget. Impacts quarterly performance and requires senior management intervention.', 'active', 1, GETDATE()),
    (@FinancialId, @VeryHighId, 'Severe financial loss exceeding 5% of annual budget. Threatens annual targets and may impact organizational viability.', 'active', 1, GETDATE()),
    
    -- Reputation Impact
    (@ReputationId, @VeryLowId, 'Minimal impact on reputation. Internal issue with no external visibility or stakeholder awareness.', 'active', 1, GETDATE()),
    (@ReputationId, @LowId, 'Minor reputational concern. Limited negative feedback from small stakeholder group, quickly resolved.', 'active', 1, GETDATE()),
    (@ReputationId, @ModerateId, 'Moderate reputational damage. Sector-wide concern, negative media coverage at local level, temporary loss of stakeholder confidence.', 'active', 1, GETDATE()),
    (@ReputationId, @HighId, 'Significant reputational damage. National media coverage, social media backlash, loss of key stakeholder trust requiring extensive PR efforts.', 'active', 1, GETDATE()),
    (@ReputationId, @VeryHighId, 'Severe reputational crisis. Widespread negative media coverage, loss of major clients/partners, long-term damage to brand and market position.', 'active', 1, GETDATE()),
    
    -- Operational Impact
    (@OperationalId, @VeryLowId, 'Negligible operational disruption. Minor delays in non-critical processes with no impact on service delivery.', 'active', 1, GETDATE()),
    (@OperationalId, @LowId, 'Minor operational disruption. Brief delays in some processes, minimal impact on overall operations, resolved within hours.', 'active', 1, GETDATE()),
    (@OperationalId, @ModerateId, 'Moderate operational disruption. Temporary or localized service interruption affecting some business units, resolved within 24 hours.', 'active', 1, GETDATE()),
    (@OperationalId, @HighId, 'Significant operational disruption. Extended service outage under 24 hours affecting critical operations and multiple business units.', 'active', 1, GETDATE()),
    (@OperationalId, @VeryHighId, 'Severe operational failure. Organization-wide unavailability exceeding 24 hours, loss of critical functions, major impact on service delivery.', 'active', 1, GETDATE()),
    
    -- Business Interruption Impact
    (@BusinessInterruptionId, @VeryLowId, 'Minimal business interruption. Brief disruption to non-essential services with no customer impact.', 'active', 1, GETDATE()),
    (@BusinessInterruptionId, @LowId, 'Minor business interruption. Short-term disruption to some services, minimal customer inconvenience, recovery within hours.', 'active', 1, GETDATE()),
    (@BusinessInterruptionId, @ModerateId, 'Moderate business interruption. Partial service disruption affecting customer experience, recovery within 24 hours.', 'active', 1, GETDATE()),
    (@BusinessInterruptionId, @HighId, 'Major business interruption. Significant service outage affecting core business functions, extended recovery time up to 48 hours.', 'active', 1, GETDATE()),
    (@BusinessInterruptionId, @VeryHighId, 'Critical business interruption. Complete service failure, loss of critical business capabilities, recovery exceeding 48 hours with potential revenue loss.', 'active', 1, GETDATE()),
    
    -- Technical Impact
    (@TechnicalId, @VeryLowId, 'Minor technical issues. Small software glitches or system errors with no impact on operations.', 'active', 1, GETDATE()),
    (@TechnicalId, @LowId, 'Limited technical problems. Isolated system issues affecting few users, resolved through standard support procedures.', 'active', 1, GETDATE()),
    (@TechnicalId, @ModerateId, 'Moderate technical failures. System performance degradation affecting multiple users, requires technical intervention.', 'active', 1, GETDATE()),
    (@TechnicalId, @HighId, 'Significant technical failures. Major system outages impacting critical infrastructure, extended downtime requiring emergency response.', 'active', 1, GETDATE()),
    (@TechnicalId, @VeryHighId, 'Critical technical failure. Complete system collapse, loss of critical IT infrastructure, potential data loss requiring disaster recovery.', 'active', 1, GETDATE()),
    
    -- Cyber Impact
    (@CyberId, @VeryLowId, 'Minimal cyber incident. Unsuccessful attack attempt or minor security alert with no data compromise.', 'active', 1, GETDATE()),
    (@CyberId, @LowId, 'Minor cyber incident. Limited data leak or minor security breach affecting non-sensitive information, quickly contained.', 'active', 1, GETDATE()),
    (@CyberId, @ModerateId, 'Moderate cyber incident. Exposure of internal or confidential data, reportable security breach requiring investigation.', 'active', 1, GETDATE()),
    (@CyberId, @HighId, 'Significant cyber attack. Partial exposure of sensitive customer or employee data impacting critical systems, regulatory notification required.', 'active', 1, GETDATE()),
    (@CyberId, @VeryHighId, 'Major cyber breach. Full exposure of highly sensitive data, widespread system compromise, significant regulatory penalties and reputational damage.', 'active', 1, GETDATE()),
    
    -- Customer Satisfaction Impact
    (@CustomerSatisfactionId, @VeryLowId, 'Minimal customer impact. Few isolated complaints with no pattern, quickly resolved.', 'active', 1, GETDATE()),
    (@CustomerSatisfactionId, @LowId, 'Minor customer dissatisfaction. Limited complaints from small customer segment, manageable through standard service recovery.', 'active', 1, GETDATE()),
    (@CustomerSatisfactionId, @ModerateId, 'Moderate customer dissatisfaction. Noticeable increase in complaints, temporary decline in satisfaction scores, requires management attention.', 'active', 1, GETDATE()),
    (@CustomerSatisfactionId, @HighId, 'Significant customer dissatisfaction. Major service failures leading to customer complaints, negative reviews, and potential customer loss.', 'active', 1, GETDATE()),
    (@CustomerSatisfactionId, @VeryHighId, 'Severe customer satisfaction crisis. Widespread customer complaints, mass customer exodus, significant damage to customer relationships and retention.', 'active', 1, GETDATE()),
    
    -- Development and Growth Impact
    (@DevelopmentGrowthId, @VeryLowId, 'Negligible impact on growth. Minor delays in non-critical development initiatives with no strategic impact.', 'active', 1, GETDATE()),
    (@DevelopmentGrowthId, @LowId, 'Minor impact on development. Slight delays in growth projects, easily recovered with minimal resource adjustment.', 'active', 1, GETDATE()),
    (@DevelopmentGrowthId, @ModerateId, 'Moderate impact on growth. Delays in key development initiatives, requires reprioritization of growth strategies.', 'active', 1, GETDATE()),
    (@DevelopmentGrowthId, @HighId, 'Significant impact on growth. Major setbacks in strategic development projects, loss of market opportunities, impacts competitive position.', 'active', 1, GETDATE()),
    (@DevelopmentGrowthId, @VeryHighId, 'Critical impact on growth. Failure of major strategic initiatives, loss of market share, threatens long-term organizational growth and sustainability.', 'active', 1, GETDATE()),
    
    -- Credit Impact
    (@CreditId, @VeryLowId, 'Minimal credit risk. Isolated minor defaults with no impact on overall credit portfolio quality.', 'active', 1, GETDATE()),
    (@CreditId, @LowId, 'Minor credit deterioration. Limited increase in non-performing loans, manageable within existing provisions.', 'active', 1, GETDATE()),
    (@CreditId, @ModerateId, 'Moderate credit risk. Noticeable increase in defaults, requires additional provisioning and enhanced monitoring.', 'active', 1, GETDATE()),
    (@CreditId, @HighId, 'Significant credit losses. Major defaults affecting credit portfolio quality, substantial provisioning requirements, regulatory attention.', 'active', 1, GETDATE()),
    (@CreditId, @VeryHighId, 'Severe credit crisis. Widespread defaults, significant impairment of credit portfolio, threatens capital adequacy and regulatory compliance.', 'active', 1, GETDATE()),
    
    -- Human Resources Impact
    (@HumanResourcesId, @VeryLowId, 'Minimal HR impact. Minor employee concerns with no impact on morale or productivity.', 'active', 1, GETDATE()),
    (@HumanResourcesId, @LowId, 'Minor HR issues. Limited employee dissatisfaction or minor turnover in non-critical positions.', 'active', 1, GETDATE()),
    (@HumanResourcesId, @ModerateId, 'Moderate HR impact. Noticeable employee turnover, declining morale in some departments, requires HR intervention.', 'active', 1, GETDATE()),
    (@HumanResourcesId, @HighId, 'Significant HR crisis. High employee turnover including key personnel, widespread morale issues, difficulty in recruitment.', 'active', 1, GETDATE()),
    (@HumanResourcesId, @VeryHighId, 'Critical HR failure. Mass employee exodus, loss of critical talent, severe organizational culture issues threatening operations.', 'active', 1, GETDATE()),
    
    -- Regulations and Policies Impact
    (@RegulationsPoliciesId, @VeryLowId, 'Minimal compliance concern. Minor internal policy deviation with no regulatory implications.', 'active', 1, GETDATE()),
    (@RegulationsPoliciesId, @LowId, 'Minor compliance issue. Small breach of internal policies or minor audit finding, easily remediated.', 'active', 1, GETDATE()),
    (@RegulationsPoliciesId, @ModerateId, 'Moderate compliance breach. Reportable compliance incident, regulatory inquiry, requires corrective action plan.', 'active', 1, GETDATE()),
    (@RegulationsPoliciesId, @HighId, 'Significant regulatory violation. Formal regulatory warnings, substantial penalties, potential license restrictions.', 'active', 1, GETDATE()),
    (@RegulationsPoliciesId, @VeryHighId, 'Critical regulatory breach. Major violations leading to severe penalties, license suspension or revocation, potential criminal charges.', 'active', 1, GETDATE()),
    
    -- Health and Safety Impact
    (@HealthSafetyId, @VeryLowId, 'Minimal health and safety concern. Minor incident with no injuries, resolved through standard procedures.', 'active', 1, GETDATE()),
    (@HealthSafetyId, @LowId, 'Minor health and safety incident. Minor injuries requiring first aid, no lost time, no long-term effects.', 'active', 1, GETDATE()),
    (@HealthSafetyId, @ModerateId, 'Moderate health and safety incident. Injuries requiring medical treatment, temporary work restrictions, reportable incident.', 'active', 1, GETDATE()),
    (@HealthSafetyId, @HighId, 'Serious health and safety incident. Major injuries with long-term effects, lost time accidents, regulatory investigation.', 'active', 1, GETDATE()),
    (@HealthSafetyId, @VeryHighId, 'Critical health and safety failure. Multiple serious injuries or fatalities, major workplace disaster, criminal investigation and prosecution.', 'active', 1, GETDATE()),
    
    -- Governance and Internal Controls Impact
    (@GovernanceControlsId, @VeryLowId, 'Minimal control weakness. Minor gaps in internal controls with no material impact.', 'active', 1, GETDATE()),
    (@GovernanceControlsId, @LowId, 'Minor control deficiency. Limited control weaknesses identified in audit, easily remediated.', 'active', 1, GETDATE()),
    (@GovernanceControlsId, @ModerateId, 'Moderate control failure. Significant control deficiencies requiring management attention and remediation plan.', 'active', 1, GETDATE()),
    (@GovernanceControlsId, @HighId, 'Significant governance failure. Material control weaknesses, potential fraud exposure, board-level intervention required.', 'active', 1, GETDATE()),
    (@GovernanceControlsId, @VeryHighId, 'Critical governance breakdown. Complete failure of internal controls, fraud or misconduct, regulatory intervention, potential criminal liability.', 'active', 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM risk_matrix_impact_guidelines 
WHERE risk_matrix_impact_id IN (
    SELECT id FROM risk_matrix_impacts 
    WHERE risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'erm')
);
-- +goose StatementEnd
