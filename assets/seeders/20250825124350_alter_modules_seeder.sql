-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET route = '/settings/organization'
WHERE identifier = 'organization_setting';

UPDATE modules
SET route = '/settings/compliance'
WHERE identifier = 'compliance_setting';


UPDATE modules
SET route = '/settings/risk-management'
WHERE identifier = 'risk_setting';


UPDATE modules
SET route = '/settings/audit'
WHERE identifier = 'audit_setting';


UPDATE modules
SET route = '/settings/document'
WHERE identifier = 'document_setting';


UPDATE modules
SET route = '/settings/incident'
WHERE identifier = 'incident_setting';


UPDATE modules
SET route = '/settings/review-meeting'
WHERE identifier = 'review_meeting_setting';


UPDATE modules
SET route = '/settings/audit-insight'
WHERE identifier = 'audit_insight_setting';


UPDATE modules
SET route = '/settings/business-continuity'
WHERE identifier = 'business_continuity_setting';


UPDATE modules
SET route = '/settings/general'
WHERE identifier = 'general_setting';


UPDATE modules
SET route = '/settings/masters'
WHERE identifier = 'common_setting';


UPDATE modules
SET route = '/settings/asset'
WHERE identifier = 'asset_setting';


UPDATE modules
SET route = '/settings/log'
WHERE identifier = 'log';

UPDATE modules
SET route = '/settings/strategy'
WHERE identifier = 'strategy_setting';

UPDATE modules
SET route = '/settings/organization-context'
WHERE identifier = 'organizational_context_setting';

UPDATE modules
SET route = '/settings/tprm'
WHERE identifier = 'tprm_setting';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET route = '/setting/organization-setting'
WHERE identifier = 'organization_setting';

UPDATE modules
SET route = '/setting/compliance-setting'
WHERE identifier = 'compliance_setting';


UPDATE modules
SET route = '/setting/risk-setting'
WHERE identifier = 'risk_setting';


UPDATE modules
SET route = '/setting/audit-setting'
WHERE identifier = 'audit_setting';


UPDATE modules
SET route = '/setting/document-setting'
WHERE identifier = 'document_setting';


UPDATE modules
SET route = '/setting/incident-setting'
WHERE identifier = 'incident_setting';


UPDATE modules
SET route = '/setting/review-meeting-setting'
WHERE identifier = 'review_meeting_setting';


UPDATE modules
SET route = '/setting/audit-insight-setting'
WHERE identifier = 'audit_insight_setting';


UPDATE modules
SET route = '/setting/business-continuity-setting'
WHERE identifier = 'business_continuity_setting';


UPDATE modules
SET route = '/setting/general-setting'
WHERE identifier = 'general_setting';


UPDATE modules
SET route = '/setting/master-setting'
WHERE identifier = 'common_setting';


UPDATE modules
SET route = '/setting/asset-setting'
WHERE identifier = 'asset_setting';


UPDATE modules
SET route = '/setting/log'
WHERE identifier = 'log';

UPDATE modules
SET route = '/setting/strategy-setting'
WHERE identifier = 'strategy_setting';

UPDATE modules
SET route = '/setting/organizational-context-setting'
WHERE identifier = 'organizational_context_setting';

UPDATE modules
SET route = '/setting/tprm-setting'
WHERE identifier = 'tprm_setting';
-- +goose StatementEnd
