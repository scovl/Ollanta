package rules

import (
	"embed"

	ollantarules "github.com/scovl/ollanta/ollantarules"
)

// MetaFS holds the JSON metadata files for all Go rules in this package.
//
//go:embed *.json
var MetaFS embed.FS

func init() {
	ollantarules.MustRegister(MetaFS, "*.json",
		NamingConventions,
		TodoComment,
		UselessIfElse,
		BadTmp,
		MathRandom,
		MD5UsedAsPassword,
		// Wave 2 — Go
		BindAll,
		MissingSSLMinVersion,
		WeakCrypto,
		DecompressionBomb,
		FilepathCleanMisuse,
		LoopPointer,
		// Wave 3 — Go
		CookieMissingHttponly,
		CookieMissingSecure,
		TemplateHTMLDoesNotEscape,
		UnsafeUsage,
		ZipTraversal,
		// Wave 4 — Go
		SwitchNoDefault,
	)
}
