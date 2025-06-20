-- +goose Up
-- +goose StatementBegin

INSERT INTO compliance_libraries (
    icon,
    title_en,
    title_ar,
    description_en,
    description_ar,
    issued_by,
    classification,
    document_classification,
    disclaimer,
    created_at,
    updated_at
)
VALUES (
    'nca-icon',
    'Essential Cybersecurity Controls (ECC – 2: 2024)',
    'ضوابط الأمن السيبراني الأساسية (ECC – 2: 2024)',
    'The ECC-2:2024, published by the National Cybersecurity Authority (NCA) of Saudi Arabia, outlines mandatory baseline cybersecurity requirements for national entities to protect information systems and digital infrastructure.',
    'يحدد ECC-2:2024، الذي نشرته الهيئة الوطنية للأمن السيبراني (NCA) في المملكة العربية السعودية، متطلبات أساسية إلزامية للأمن السيبراني للجهات الوطنية لحماية أنظمتها الرقمية.',
    'National Cybersecurity Authority (NCA), KSA',
    'Public (TLP: White)',
    'Public',
    'Disclaimer: The following controls will be governed by and implemented in accordance with the laws of the Kingdom of Saudi Arabia, and must be subject to the exclusive jurisdiction of the courts of the Kingdom of Saudi Arabia. Therefore, the Arabic version will be the binding language for all matters relating to the meaning or interoperation of this document.',
    GETDATE(),
    GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE compliance_libraries;

-- +goose StatementEnd