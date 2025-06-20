-- +goose Up
-- +goose StatementBegin

-- 1. Remove `.png` from all icon values
UPDATE knowledge_center_resources
SET icon = REPLACE(icon, '.png', ''), updated_at = GETDATE()
WHERE icon LIKE '%.png';

-- 2. Set type and item_count with updated_at for relevant resources

-- Risk Library
UPDATE knowledge_center_resources
SET 
  type = 'risk',
  item_count = (SELECT COUNT(*) FROM risk_library),
  updated_at = GETDATE()
WHERE title_en = 'Risk Library';

-- KRI Library
UPDATE knowledge_center_resources
SET 
  type = 'kri',
  item_count = (SELECT COUNT(*) FROM kri_library),
  updated_at = GETDATE()
WHERE title_en = 'Key Risk Indicators Library';

-- KPI Library
UPDATE knowledge_center_resources
SET type = 'kpi', updated_at = GETDATE()
WHERE title_en = 'KPI Library';

-- Compliance
UPDATE knowledge_center_resources SET type = 'sama', updated_at = GETDATE() WHERE title_en = 'Saudi Arabian Monetary Authority (SAMA)';
UPDATE knowledge_center_resources SET type = 'nca', updated_at = GETDATE() WHERE title_en = 'National Cybersecurity Authority (NCA)';
UPDATE knowledge_center_resources SET type = 'iso', updated_at = GETDATE() WHERE title_en = 'ISO/IEC 27001';

-- Document
UPDATE knowledge_center_resources SET type = 'policy_template', updated_at = GETDATE() WHERE title_en = 'Policy Templates Library';
UPDATE knowledge_center_resources SET type = 'standard', updated_at = GETDATE() WHERE title_en = 'Standards Library';

-- Audit
UPDATE knowledge_center_resources SET type = 'audit_checklist', updated_at = GETDATE() WHERE title_en = 'Audit Checklist Library';

-- Strategy
UPDATE knowledge_center_resources SET type = 'strategy', updated_at = GETDATE() WHERE title_en = 'Strategic Objectives Library';

-- Asset
UPDATE knowledge_center_resources SET type = 'asset', updated_at = GETDATE() WHERE title_en = 'Asset Library';
UPDATE knowledge_center_resources SET type = 'vulnerability', updated_at = GETDATE() WHERE title_en = 'Vulnerabilities Library';
UPDATE knowledge_center_resources SET type = 'threat', updated_at = GETDATE() WHERE title_en = 'Threats Library';

-- Knowledge Center Topics
UPDATE knowledge_center_resources SET type = 'cyber_security', updated_at = GETDATE() WHERE title_en = 'Cyber Security';
UPDATE knowledge_center_resources SET type = 'ai_generative_ai', updated_at = GETDATE() WHERE title_en = 'AI & Generative AI';
UPDATE knowledge_center_resources SET type = 'compliance_frameworks', updated_at = GETDATE() WHERE title_en = 'Compliance & Frameworks';
UPDATE knowledge_center_resources SET type = 'metrics_utilities', updated_at = GETDATE() WHERE title_en = 'Metrics & Utilities';
UPDATE knowledge_center_resources SET type = 'erm', updated_at = GETDATE() WHERE title_en = 'Enterprise Risk Management (ERM)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert icon changes
UPDATE knowledge_center_resources
SET icon = icon + '.png', updated_at = GETDATE()
WHERE icon NOT LIKE '%.png';

-- Revert type and item_count
UPDATE knowledge_center_resources
SET type = NULL,
    item_count = 0,
    updated_at = GETDATE();

-- +goose StatementEnd
