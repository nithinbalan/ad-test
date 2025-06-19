package dtos

type DescriptionRewriteRequest struct {
	Content    string `json:"content" binding:"required"`
	FixGrammar bool   `json:"fix_grammer"`
	Rewrite    bool   `json:"rewrite"`
	AddArabic  bool   `json:"add_arabic"`
}

type LanguageTranslationRequest struct {
	Content  string `json:"content" binding:"required"`
	Language string `json:"language" binding:"required"`
}
