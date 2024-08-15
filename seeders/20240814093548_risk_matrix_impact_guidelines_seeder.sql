-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

-- Fetching Impact Categories
DECLARE @FinancialImpactCategoryId INT, @OperationalImpactCategoryId INT, @ReputationalImpactCategoryId INT, @ComplianceImpactCategoryId INT;
DECLARE @LegalImpactCategoryId INT, @HealthSafetyImpactCategoryId INT, @EnvironmentalImpactCategoryId INT, @CustomerImpactCategoryId INT;
DECLARE @TechnologicalImpactCategoryId INT, @StrategicImpactCategoryId INT;

SELECT 
    @FinancialImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Financial Impact';
SELECT 
    @OperationalImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Operational Impact';
SELECT 
    @ReputationalImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Reputational Impact';
SELECT 
    @ComplianceImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Compliance Impact';
SELECT 
    @LegalImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Legal Impact';
SELECT 
    @HealthSafetyImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Health and Safety Impact';
SELECT 
    @EnvironmentalImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Environmental Impact';
SELECT 
    @CustomerImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Customer Impact';
SELECT 
    @TechnologicalImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Technological Impact';
SELECT 
    @StrategicImpactCategoryId = id FROM risk_impact_categories WHERE title_en = 'Strategic Impact';

-- Fetching Matrix Impact Levels
DECLARE @MinorImpactId INT, @ModerateImpactId INT, @SignificantImpactId INT, @MajorImpactId INT, @CatastrophicImpactId INT;

SELECT 
    @MinorImpactId = id FROM risk_matrix_impacts WHERE title_en = 'Minor' AND risk_matrix_id = @RiskMatrixId;
SELECT 
    @ModerateImpactId = id FROM risk_matrix_impacts WHERE title_en = 'Moderate' AND risk_matrix_id = @RiskMatrixId;
SELECT 
    @SignificantImpactId = id FROM risk_matrix_impacts WHERE title_en = 'Significant' AND risk_matrix_id = @RiskMatrixId;
SELECT 
    @MajorImpactId = id FROM risk_matrix_impacts WHERE title_en = 'Major' AND risk_matrix_id = @RiskMatrixId;
SELECT 
    @CatastrophicImpactId = id FROM risk_matrix_impacts WHERE title_en = 'Catastrophic' AND risk_matrix_id = @RiskMatrixId;

-- Seeding Risk Matrix Impact Guidelines
INSERT INTO risk_matrix_impact_guidelines (risk_impact_category_id, risk_matrix_impact_id, guideline, status, created_by_id, created_at)
VALUES
    -- Financial Impact
    (@FinancialImpactCategoryId, @MinorImpactId, 'Minimal financial loss that does not affect the organization''s overall financial health (e.g., less than 0.5% of the annual budget).', 'active', 1, GETDATE()),
    (@FinancialImpactCategoryId, @ModerateImpactId, 'Manageable financial loss that may require budget adjustments (e.g., 0.5% to 1% of the annual budget).', 'active', 1, GETDATE()),
    (@FinancialImpactCategoryId, @SignificantImpactId, 'Considerable financial loss impacting quarterly performance (e.g., 1% to 5% of the annual budget).', 'active', 1, GETDATE()),
    (@FinancialImpactCategoryId, @MajorImpactId, 'Severe financial loss affecting annual performance, possibly leading to layoffs or major budget cuts (e.g., 5% to 10% of the annual budget).', 'active', 1, GETDATE()),
    (@FinancialImpactCategoryId, @CatastrophicImpactId, 'Extreme financial loss threatening the organization''s viability (e.g., more than 10% of the annual budget).', 'active', 1, GETDATE()),
    
    -- Operational Impact
    (@OperationalImpactCategoryId, @MinorImpactId, 'Negligible disruption to operations with no significant effect on overall performance (e.g., minor delays in non-critical processes).', 'active', 1, GETDATE()),
    (@OperationalImpactCategoryId, @ModerateImpactId, 'Noticeable disruption requiring additional resources or minor adjustments (e.g., moderate delays in critical processes).', 'active', 1, GETDATE()),
    (@OperationalImpactCategoryId, @SignificantImpactId, 'Substantial disruption affecting key operations and requiring significant management attention (e.g., a temporary halt in key production lines).', 'active', 1, GETDATE()),
    (@OperationalImpactCategoryId, @MajorImpactId, 'Severe disruption leading to prolonged downtime of critical operations (e.g., failure of major systems affecting overall operations).', 'active', 1, GETDATE()),
    (@OperationalImpactCategoryId, @CatastrophicImpactId, 'Complete operational failure with widespread impact, potentially leading to a shutdown of the organization (e.g., destruction of a primary facility).', 'active', 1, GETDATE()),
    
    -- Reputational Impact
    (@ReputationalImpactCategoryId, @MinorImpactId, 'Insignificant impact on reputation, unlikely to be noticed outside the organization (e.g., minor internal issue resolved quickly).', 'active', 1, GETDATE()),
    (@ReputationalImpactCategoryId, @ModerateImpactId, 'Some negative perception within certain stakeholders but manageable (e.g., temporary customer dissatisfaction).', 'active', 1, GETDATE()),
    (@ReputationalImpactCategoryId, @SignificantImpactId, 'Negative public perception leading to a decline in customer trust or stakeholder confidence (e.g., negative media coverage).', 'active', 1, GETDATE()),
    (@ReputationalImpactCategoryId, @MajorImpactId, 'Significant reputational damage requiring extensive public relations efforts (e.g., a major scandal affecting public trust).', 'active', 1, GETDATE()),
    (@ReputationalImpactCategoryId, @CatastrophicImpactId, 'Long-term reputational damage threatening the organization''s existence (e.g., widespread scandal leading to loss of major clients and partners).', 'active', 1, GETDATE()),

    -- Compliance Impact
    (@ComplianceImpactCategoryId, @MinorImpactId, 'Minor non-compliance with minimal regulatory consequences (e.g., a small breach of internal policies).', 'active', 1, GETDATE()),
    (@ComplianceImpactCategoryId, @ModerateImpactId, 'Noticeable non-compliance leading to warnings or minor penalties (e.g., minor regulatory violations).', 'active', 1, GETDATE()),
    (@ComplianceImpactCategoryId, @SignificantImpactId, 'Significant non-compliance resulting in substantial penalties or legal action (e.g., failure to meet key regulatory requirements).', 'active', 1, GETDATE()),
    (@ComplianceImpactCategoryId, @MajorImpactId, 'Severe non-compliance leading to legal action, large fines, or sanctions (e.g., major violations of industry regulations).', 'active', 1, GETDATE()),
    (@ComplianceImpactCategoryId, @CatastrophicImpactId, 'Complete failure to comply, leading to criminal charges, loss of licenses, or business closure (e.g., criminal regulatory violations).', 'active', 1, GETDATE()),

    -- Legal Impact
    (@LegalImpactCategoryId, @MinorImpactId, 'Minimal legal implications with no significant effect on the organization (e.g., minor contractual disputes).', 'active', 1, GETDATE()),
    (@LegalImpactCategoryId, @ModerateImpactId, 'Manageable legal issues requiring attention but not posing a serious threat (e.g., moderate legal claims or disputes).', 'active', 1, GETDATE()),
    (@LegalImpactCategoryId, @SignificantImpactId, 'Significant legal challenges with potential for substantial penalties or settlements (e.g., major lawsuits).', 'active', 1, GETDATE()),
    (@LegalImpactCategoryId, @MajorImpactId, 'Severe legal consequences, including the possibility of large settlements, fines, or business restrictions (e.g., class action lawsuits).', 'active', 1, GETDATE()),
    (@LegalImpactCategoryId, @CatastrophicImpactId, 'Extreme legal consequences threatening the organization''s survival, such as criminal charges against key executives (e.g., major criminal cases or government intervention).', 'active', 1, GETDATE()),

    -- Health and Safety Impact
    (@HealthSafetyImpactCategoryId, @MinorImpactId, 'Minor injuries or health impacts with no long-term effects (e.g., minor workplace incidents).', 'active', 1, GETDATE()),
    (@HealthSafetyImpactCategoryId, @ModerateImpactId, 'Noticeable injuries or health impacts requiring medical attention but not life-threatening (e.g., moderate workplace accidents).', 'active', 1, GETDATE()),
    (@HealthSafetyImpactCategoryId, @SignificantImpactId, 'Serious injuries or health impacts with long-term effects on employees or the public (e.g., severe workplace accidents).', 'active', 1, GETDATE()),
    (@HealthSafetyImpactCategoryId, @MajorImpactId, 'Major health and safety incidents causing multiple injuries or fatalities (e.g., significant industrial accidents).', 'active', 1, GETDATE()),
    (@HealthSafetyImpactCategoryId, @CatastrophicImpactId, 'Catastrophic health and safety failures leading to multiple fatalities and widespread impact (e.g., major disaster at a facility).', 'active', 1, GETDATE()),

    -- Environmental Impact
    (@EnvironmentalImpactCategoryId, @MinorImpactId, 'Minimal environmental impact with no significant long-term effects (e.g., minor spill contained quickly).', 'active', 1, GETDATE()),
    (@EnvironmentalImpactCategoryId, @ModerateImpactId, 'Noticeable environmental impact requiring cleanup or mitigation (e.g., moderate pollution incidents).', 'active', 1, GETDATE()),
    (@EnvironmentalImpactCategoryId, @SignificantImpactId, 'Significant environmental damage affecting local ecosystems or communities (e.g., major spill or contamination).', 'active', 1, GETDATE()),
    (@EnvironmentalImpactCategoryId, @MajorImpactId, 'Severe environmental impact with long-term consequences, requiring extensive remediation efforts (e.g., large-scale deforestation or pollution).', 'active', 1, GETDATE()),
    (@EnvironmentalImpactCategoryId, @CatastrophicImpactId, 'Catastrophic environmental damage with widespread and possibly irreversible effects (e.g., environmental disaster leading to loss of biodiversity).', 'active', 1, GETDATE()),

    -- Customer Impact
    (@CustomerImpactCategoryId, @MinorImpactId, 'Minimal impact on customers, with no significant effect on customer satisfaction (e.g., minor service disruptions).', 'active', 1, GETDATE()),
    (@CustomerImpactCategoryId, @ModerateImpactId, 'Noticeable impact leading to temporary dissatisfaction but recoverable (e.g., moderate delays in service delivery).', 'active', 1, GETDATE()),
    (@CustomerImpactCategoryId, @SignificantImpactId, 'Significant impact on customer experience, leading to a loss of customers or negative feedback (e.g., major service outages).', 'active', 1, GETDATE()),
    (@CustomerImpactCategoryId, @MajorImpactId, 'Severe impact leading to a significant loss of customer base or long-term damage to customer relationships (e.g., failure to deliver critical services).', 'active', 1, GETDATE()),
    (@CustomerImpactCategoryId, @CatastrophicImpactId, 'Catastrophic impact causing widespread loss of customers and potentially threatening the organization''s future (e.g., major product recall leading to customer exodus).', 'active', 1, GETDATE()),

    -- Technological Impact
    (@TechnologicalImpactCategoryId, @MinorImpactId, 'Minor technology issues with no significant impact on operations (e.g., minor software glitches).', 'active', 1, GETDATE()),
    (@TechnologicalImpactCategoryId, @ModerateImpactId, 'Noticeable technology issues requiring intervention but manageable (e.g., moderate system downtimes).', 'active', 1, GETDATE()),
    (@TechnologicalImpactCategoryId, @SignificantImpactId, 'Significant technology failures affecting critical systems or processes (e.g., major system outages).', 'active', 1, GETDATE()),
    (@TechnologicalImpactCategoryId, @MajorImpactId, 'Severe technology failures causing prolonged downtime or data loss (e.g., cyberattacks leading to data breaches).', 'active', 1, GETDATE()),
    (@TechnologicalImpactCategoryId, @CatastrophicImpactId, 'Complete technology failure with widespread and irreversible impact, potentially leading to business closure (e.g., destruction of IT infrastructure).', 'active', 1, GETDATE()),

    -- Strategic Impact
    (@StrategicImpactCategoryId, @MinorImpactId, 'Minor impact on strategic objectives, easily adjusted with no long-term consequences (e.g., slight delays in strategic initiatives).', 'active', 1, GETDATE()),
    (@StrategicImpactCategoryId, @ModerateImpactId, 'Noticeable impact requiring adjustments to strategic plans but still achievable (e.g., moderate changes to key strategies).', 'active', 1, GETDATE()),
    (@StrategicImpactCategoryId, @SignificantImpactId, 'Significant impact necessitating major revisions to strategic objectives (e.g., failure to meet key strategic goals).', 'active', 1, GETDATE()),
    (@StrategicImpactCategoryId, @MajorImpactId, 'Severe impact threatening the achievement of critical strategic objectives (e.g., loss of a major market or customer base).', 'active', 1, GETDATE()),
    (@StrategicImpactCategoryId, @CatastrophicImpactId, 'Catastrophic impact leading to the failure of the organization’s strategic vision and possible dissolution (e.g., complete strategic collapse due to market failure).', 'active', 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_matrix_impact_guidelines;
-- +goose StatementEnd
