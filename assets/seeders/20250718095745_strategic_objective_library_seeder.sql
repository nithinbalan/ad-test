-- +goose Up
-- +goose StatementBegin

-- Governance & Compliance
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish Three Lines of Defense model with clear accountability frameworks', 'Governance & Compliance'),
('Implement GRC integrated platform for unified governance oversight', 'Governance & Compliance'),
('Develop board-level risk appetite and tolerance framework', 'Governance & Compliance'),
('Create regulatory compliance monitoring with automated breach detection', 'Governance & Compliance'),
('Establish ethics and conduct program with regular compliance attestation', 'Governance & Compliance'),
('Implement policy lifecycle management with governance approval workflows', 'Governance & Compliance'),
('Develop compliance risk assessment and control testing framework', 'Governance & Compliance'),
('Create governance reporting dashboard for executive and board oversight', 'Governance & Compliance'),
('Establish vendor governance and third-party risk management program', 'Governance & Compliance'),
('Implement segregation of duties and access control governance framework', 'Governance & Compliance');

-- Risk & Resilience
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish enterprise risk management framework aligned with ISO 31000', 'Risk & Resilience'),
('Implement risk register and heat map with automated escalation protocols', 'Risk & Resilience'),
('Develop operational resilience framework with impact tolerance thresholds', 'Risk & Resilience'),
('Create risk-based internal audit program with continuous monitoring', 'Risk & Resilience'),
('Establish key risk indicators (KRIs) with real-time monitoring capabilities', 'Risk & Resilience'),
('Implement crisis management governance with clear decision-making authority', 'Risk & Resilience'),
('Develop risk culture assessment and risk awareness training programs', 'Risk & Resilience'),
('Create integrated risk reporting for board and senior management', 'Risk & Resilience'),
('Establish business continuity governance with regular testing protocols', 'Risk & Resilience'),
('Implement risk treatment planning with control effectiveness monitoring', 'Risk & Resilience');

-- Operational Excellence
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish operational risk control framework with key control testing', 'Operational Excellence'),
('Implement process governance with documented procedures and controls', 'Operational Excellence'),
('Develop operational key performance indicators (KPIs) and risk metrics', 'Operational Excellence'),
('Create change management governance with risk impact assessment', 'Operational Excellence'),
('Establish operational compliance monitoring with exception reporting', 'Operational Excellence'),
('Implement service level agreement governance with performance monitoring', 'Operational Excellence'),
('Develop operational incident management with root cause analysis', 'Operational Excellence'),
('Create operational risk appetite framework with tolerance limits', 'Operational Excellence'),
('Establish operational audit trail and documentation standards', 'Operational Excellence'),
('Implement operational risk reporting with trend analysis capabilities', 'Operational Excellence');

-- Digital Transformation & Innovation
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish IT governance framework with risk-based technology decisions', 'Digital Transformation & Innovation'),
('Implement cybersecurity governance with board-level oversight', 'Digital Transformation & Innovation'),
('Develop data governance framework with privacy compliance controls', 'Digital Transformation & Innovation'),
('Create digital risk assessment methodology for emerging technologies', 'Digital Transformation & Innovation'),
('Establish cloud governance with security and compliance standards', 'Digital Transformation & Innovation'),
('Implement technology risk management with vendor assessment protocols', 'Digital Transformation & Innovation'),
('Develop digital compliance monitoring for regulatory requirements', 'Digital Transformation & Innovation'),
('Create innovation governance with ethical AI and responsible technology use', 'Digital Transformation & Innovation'),
('Establish digital audit capabilities with continuous compliance monitoring', 'Digital Transformation & Innovation'),
('Implement technology change governance with risk approval workflows', 'Digital Transformation & Innovation');

-- Customer & Stakeholder Engagement
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish customer data governance with privacy compliance framework', 'Customer & Stakeholder Engagement'),
('Implement stakeholder risk assessment and relationship governance', 'Customer & Stakeholder Engagement'),
('Develop customer complaint governance with regulatory reporting protocols', 'Customer & Stakeholder Engagement'),
('Create service delivery governance with compliance monitoring', 'Customer & Stakeholder Engagement'),
('Establish customer due diligence and know-your-customer (KYC) framework', 'Customer & Stakeholder Engagement'),
('Implement stakeholder engagement risk management with mitigation strategies', 'Customer & Stakeholder Engagement'),
('Develop customer risk profiling and segmentation governance', 'Customer & Stakeholder Engagement'),
('Create public consultation governance with regulatory compliance', 'Customer & Stakeholder Engagement'),
('Establish customer satisfaction monitoring with risk-based improvements', 'Customer & Stakeholder Engagement'),
('Implement stakeholder communication governance with approval workflows', 'Customer & Stakeholder Engagement');

-- Financial & Performance Management
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish financial risk governance with board-level oversight', 'Financial & Performance Management'),
('Implement financial compliance monitoring with regulatory reporting', 'Financial & Performance Management'),
('Develop budget governance framework with risk-based allocations', 'Financial & Performance Management'),
('Create financial controls testing with audit trail documentation', 'Financial & Performance Management'),
('Establish financial risk appetite framework with tolerance metrics', 'Financial & Performance Management'),
('Implement financial performance governance with variance analysis', 'Financial & Performance Management'),
('Develop credit risk management with portfolio monitoring', 'Financial & Performance Management'),
('Create investment governance with risk assessment protocols', 'Financial & Performance Management'),
('Establish procurement governance with vendor risk evaluation', 'Financial & Performance Management'),
('Implement financial audit readiness with continuous compliance monitoring', 'Financial & Performance Management');

-- Workforce & Human Capital Development
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish HR governance framework with compliance monitoring', 'Workforce & Human Capital Development'),
('Implement workforce risk assessment with mitigation strategies', 'Workforce & Human Capital Development'),
('Develop talent risk management with succession planning governance', 'Workforce & Human Capital Development'),
('Create conduct risk framework with behavioral monitoring', 'Workforce & Human Capital Development'),
('Establish employment compliance monitoring with regulatory updates', 'Workforce & Human Capital Development'),
('Implement competency governance with risk-based training requirements', 'Workforce & Human Capital Development'),
('Develop workplace safety governance with incident reporting protocols', 'Workforce & Human Capital Development'),
('Create performance management governance with objective assessment', 'Workforce & Human Capital Development'),
('Establish diversity and inclusion governance with compliance tracking', 'Workforce & Human Capital Development'),
('Implement HR audit capabilities with documentation standards', 'Workforce & Human Capital Development');

-- Sustainability & Environmental Stewardship
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish environmental governance framework with board oversight', 'Sustainability & Environmental Stewardship'),
('Implement environmental risk assessment with impact monitoring', 'Sustainability & Environmental Stewardship'),
('Develop climate risk governance with scenario planning capabilities', 'Sustainability & Environmental Stewardship'),
('Create environmental compliance monitoring with regulatory reporting', 'Sustainability & Environmental Stewardship'),
('Establish sustainability risk appetite with performance indicators', 'Sustainability & Environmental Stewardship'),
('Implement environmental audit program with corrective action tracking', 'Sustainability & Environmental Stewardship'),
('Develop ESG governance framework with stakeholder reporting', 'Sustainability & Environmental Stewardship'),
('Create environmental incident management with root cause analysis', 'Sustainability & Environmental Stewardship'),
('Establish green procurement governance with supplier risk assessment', 'Sustainability & Environmental Stewardship'),
('Implement carbon risk management with reduction target governance', 'Sustainability & Environmental Stewardship');

-- Strategic Planning & National Alignment
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish strategic governance framework with risk-informed planning', 'Strategic Planning & National Alignment'),
('Implement policy governance with compliance impact assessment', 'Strategic Planning & National Alignment'),
('Develop strategic risk assessment with national alignment monitoring', 'Strategic Planning & National Alignment'),
('Create planning governance with stakeholder consultation protocols', 'Strategic Planning & National Alignment'),
('Establish strategic compliance monitoring with regulatory alignment', 'Strategic Planning & National Alignment'),
('Implement strategic performance governance with outcome tracking', 'Strategic Planning & National Alignment'),
('Develop strategic change governance with risk approval processes', 'Strategic Planning & National Alignment'),
('Create strategic communication governance with approval workflows', 'Strategic Planning & National Alignment'),
('Establish strategic audit capabilities with effectiveness assessment', 'Strategic Planning & National Alignment'),
('Implement strategic risk reporting with board-level escalation', 'Strategic Planning & National Alignment');

-- Knowledge & Information Management
INSERT INTO strategic_objective_library (title, category) VALUES
('Establish information governance framework with data classification', 'Knowledge & Information Management'),
('Implement information security governance with risk-based controls', 'Knowledge & Information Management'),
('Develop knowledge risk assessment with intellectual property protection', 'Knowledge & Information Management'),
('Create information compliance monitoring with privacy regulations', 'Knowledge & Information Management'),
('Establish data governance with quality and integrity controls', 'Knowledge & Information Management'),
('Implement information audit program with access control testing', 'Knowledge & Information Management'),
('Develop information risk management with breach response protocols', 'Knowledge & Information Management'),
('Create records management governance with retention compliance', 'Knowledge & Information Management'),
('Establish information sharing governance with confidentiality controls', 'Knowledge & Information Management'),
('Implement knowledge management risk assessment with continuity planning', 'Knowledge & Information Management');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM strategic_objective_library;
-- +goose StatementEnd
