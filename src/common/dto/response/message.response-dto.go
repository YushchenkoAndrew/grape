package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func prettier(field validator.FieldError) string {
	name := strings.ToLower(field.Field())
	param := strings.ToLower(field.Param())
	value := fmt.Sprintf("%v", field.Value())

	switch field.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", name)
	case "alpha":
		return fmt.Sprintf("'%s' must contain only alphabetic characters", name)
	case "alphanum", "alphanumunicode":
		return fmt.Sprintf("'%s' must contain only alphanumeric characters", name)
	case "alphaunicode":
		return fmt.Sprintf("'%s' must contain only alphabetic Unicode characters", name)
	case "ascii":
		return fmt.Sprintf("'%s' must contain only ASCII characters", name)
	case "boolean":
		return fmt.Sprintf("'%s' must be a boolean", name)
	case "contains":
		return fmt.Sprintf("'%s' must contain the substring '%s'", name, param)
	case "containsany":
		return fmt.Sprintf("'%s' must contain at least one of the characters '%s'", name, param)
	case "containsrune":
		return fmt.Sprintf("'%s' must contain the rune '%s'", name, param)
	case "endsnotwith":
		return fmt.Sprintf("'%s' must not end with '%s'", name, param)
	case "endswith":
		return fmt.Sprintf("'%s' must end with '%s'", name, param)
	case "excludes":
		return fmt.Sprintf("'%s' must not contain the substring '%s'", name, param)
	case "excludesall":
		return fmt.Sprintf("'%s' must not contain any of the characters '%s'", name, param)
	case "excludesrune":
		return fmt.Sprintf("'%s' must not contain the rune '%s'", name, param)
	case "lowercase":
		return fmt.Sprintf("'%s' must be lowercase", name)
	case "uppercase":
		return fmt.Sprintf("'%s' must be uppercase", name)
	case "multibyte":
		return fmt.Sprintf("'%s' must contain multi-byte characters", name)
	case "number", "numeric":
		return fmt.Sprintf("'%s' must be a valid number", name)
	case "printascii":
		return fmt.Sprintf("'%s' must contain only printable ASCII characters", name)
	case "startsnotwith":
		return fmt.Sprintf("'%s' must not start with '%s'", name, param)
	case "startswith":
		return fmt.Sprintf("'%s' must start with '%s'", name, param)
	case "min":
		return fmt.Sprintf("'%s' must be at least '%s' characters long", name, param)
	case "max":
		return fmt.Sprintf("'%s' must be at most '%s' characters long", name, param)
	case "email":
		return fmt.Sprintf("'%s' must be a valid email address", name)
	case "eqcsfield", "eqfield", "eq":
		return fmt.Sprintf("'%s' must be equal to '%s'", name, param)
	case "fieldcontains":
		return fmt.Sprintf("'%s' must contain the indicated characters", name)
	case "fieldexcludes":
		return fmt.Sprintf("'%s' must not contain the indicated characters", name)
	case "gtcsfield", "gtfield", "gt":
		return fmt.Sprintf("'%s' must be greater than '%s'", name, param)
	case "gtecsfield", "gtefield", "gte":
		return fmt.Sprintf("'%s' must be greater than or equal to '%s'", name, param)
	case "ltcsfield", "ltfield", "lt":
		return fmt.Sprintf("'%s' must be less than '%s'", name, param)
	case "ltecsfield", "ltefield", "lte":
		return fmt.Sprintf("'%s' must be less than or equal to '%s'", name, param)
	case "necsfield", "nefield", "ne":
		return fmt.Sprintf("'%s' must not be equal to '%s'", name, param)
	case "eq_ignore_case":
		return fmt.Sprintf("'%s' must be equal to '%s' (case insensitive)", name, param)
	case "ne_ignore_case":
		return fmt.Sprintf("'%s' must not be equal to '%s' (case insensitive)", name, param)
	case "dir":
		return fmt.Sprintf("'%s' must be an existing directory", name)
	case "dirpath":
		return fmt.Sprintf("'%s' must be a valid directory path", name)
	case "file":
		return fmt.Sprintf("'%s' must be an existing file", name)
	case "filepath":
		return fmt.Sprintf("'%s' must be a valid file path", name)
	case "image":
		return fmt.Sprintf("'%s' must be a valid image", name)
	case "isdefault":
		return fmt.Sprintf("'%s' must be the default value", name)
	case "len":
		return fmt.Sprintf("'%s' must have a length of '%s'", name, param)
	case "oneof":
		return fmt.Sprintf("'%s' must be one of [%s]", name, param)
	case "required_if":
		return fmt.Sprintf("'%s' is required if '%s' is '%s'", name, param, value)
	case "required_unless":
		return fmt.Sprintf("'%s' is required unless '%s' is '%s'", name, param, value)
	case "required_with":
		return fmt.Sprintf("'%s' is required when '%s' is present", name, param)
	case "required_with_all":
		return fmt.Sprintf("'%s' is required when all of '%s' are present", name, param)
	case "required_without":
		return fmt.Sprintf("'%s' is required when '%s' is not present", name, param)
	case "required_without_all":
		return fmt.Sprintf("'%s' is required when none of '%s' are present", name, param)
	case "excluded_if":
		return fmt.Sprintf("'%s' must not be included if '%s' is '%s'", name, param, value)
	case "excluded_unless":
		return fmt.Sprintf("'%s' must not be included unless '%s' is '%s'", name, param, value)
	case "excluded_with":
		return fmt.Sprintf("'%s' must not be included when '%s' is present", name, param)
	case "excluded_with_all":
		return fmt.Sprintf("'%s' must not be included when all of '%s' are present", name, param)
	case "excluded_without":
		return fmt.Sprintf("'%s' must not be included when '%s' is not present", name, param)
	case "excluded_without_all":
		return fmt.Sprintf("'%s' must not be included when none of '%s' are present", name, param)
	case "unique":
		return fmt.Sprintf("'%s' must be unique", name)
	case "base64":
		return fmt.Sprintf("'%s' must be a valid base64 string", name)
	case "base64url":
		return fmt.Sprintf("'%s' must be a valid base64url string", name)
	case "base64rawurl":
		return fmt.Sprintf("'%s' must be a valid base64rawurl string", name)
	case "bic":
		return fmt.Sprintf("'%s' must be a valid Business Identifier Code (ISO 9362)", name)
	case "bcp47_language_tag":
		return fmt.Sprintf("'%s' must be a valid language tag (BCP 47)", name)
	case "btc_addr":
		return fmt.Sprintf("'%s' must be a valid Bitcoin Address", name)
	case "btc_addr_bech32":
		return fmt.Sprintf("'%s' must be a valid Bitcoin Bech32 Address (segwit)", name)
	case "credit_card":
		return fmt.Sprintf("'%s' must be a valid Credit Card Number", name)
	case "mongodb":
		return fmt.Sprintf("'%s' must be a valid MongoDB ObjectID", name)
	case "cron":
		return fmt.Sprintf("'%s' must be a valid Cron expression", name)
	case "spicedb":
		return fmt.Sprintf("'%s' must be a valid SpiceDb ObjectID/Permission/Type", name)
	case "datetime":
		return fmt.Sprintf("'%s' must be a valid Datetime", name)
	case "e164":
		return fmt.Sprintf("'%s' must be a valid e164 formatted phone number", name)
	case "eth_addr":
		return fmt.Sprintf("'%s' must be a valid Ethereum Address", name)
	case "hexadecimal":
		return fmt.Sprintf("'%s' must be a valid Hexadecimal string", name)
	case "hexcolor":
		return fmt.Sprintf("'%s' must be a valid Hexcolor string", name)
	case "hsl":
		return fmt.Sprintf("'%s' must be a valid HSL string", name)
	case "hsla":
		return fmt.Sprintf("'%s' must be a valid HSLA string", name)
	case "html":
		return fmt.Sprintf("'%s' must contain valid HTML tags", name)
	case "html_encoded":
		return fmt.Sprintf("'%s' must be a valid HTML encoded string", name)
	case "isbn":
		return fmt.Sprintf("'%s' must be a valid International Standard Book Number", name)
	case "isbn10":
		return fmt.Sprintf("'%s' must be a valid International Standard Book Number 10", name)
	case "isbn13":
		return fmt.Sprintf("'%s' must be a valid International Standard Book Number 13", name)
	case "issn":
		return fmt.Sprintf("'%s' must be a valid International Standard Serial Number", name)
	case "iso3166_1_alpha2":
		return fmt.Sprintf("'%s' must be a valid Two-letter country code (ISO 3166-1 alpha-2)", name)
	case "iso3166_1_alpha3":
		return fmt.Sprintf("'%s' must be a valid Three-letter country code (ISO 3166-1 alpha-3)", name)
	case "iso3166_1_alpha_numeric":
		return fmt.Sprintf("'%s' must be a valid Numeric country code (ISO 3166-1 numeric)", name)
	case "iso3166_2":
		return fmt.Sprintf("'%s' must be a valid Country subdivision code (ISO 3166-2)", name)
	case "iso4217":
		return fmt.Sprintf("'%s' must be a valid Currency code (ISO 4217)", name)
	case "json":
		return fmt.Sprintf("'%s' must be a valid JSON", name)
	case "jwt":
		return fmt.Sprintf("'%s' must be a valid JSON Web Token (JWT)", name)
	case "latitude":
		return fmt.Sprintf("'%s' must be a valid Latitude", name)
	case "longitude":
		return fmt.Sprintf("'%s' must be a valid Longitude", name)
	case "luhn_checksum":
		return fmt.Sprintf("'%s' must pass Luhn Algorithm Checksum", name)
	case "postcode_iso3166_alpha2":
		return fmt.Sprintf("'%s' must be a valid Postcode for ISO 3166 alpha-2 country", name)
	case "postcode_iso3166_alpha2_field":
		return fmt.Sprintf("'%s' must be a valid Postcode for ISO 3166 alpha-2 country", name)
	case "rgb":
		return fmt.Sprintf("'%s' must be a valid RGB string", name)
	case "rgba":
		return fmt.Sprintf("'%s' must be a valid RGBA string", name)
	case "ssn":
		return fmt.Sprintf("'%s' must be a valid Social Security Number (SSN)", name)
	case "timezone":
		return fmt.Sprintf("'%s' must be a valid Timezone", name)
	case "uuid":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID)", name)
	case "uuid3":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v3", name)
	case "uuid3_rfc4122":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v3 RFC4122", name)
	case "uuid4":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v4", name)
	case "uuid4_rfc4122":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v4 RFC4122", name)
	case "uuid5":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v5", name)
	case "uuid5_rfc4122":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) v5 RFC4122", name)
	case "uuid_rfc4122":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Identifier (UUID) RFC4122", name)
	case "md4":
		return fmt.Sprintf("'%s' must be a valid MD4 hash", name)
	case "md5":
		return fmt.Sprintf("'%s' must be a valid MD5 hash", name)
	case "sha256":
		return fmt.Sprintf("'%s' must be a valid SHA256 hash", name)
	case "sha384":
		return fmt.Sprintf("'%s' must be a valid SHA384 hash", name)
	case "sha512":
		return fmt.Sprintf("'%s' must be a valid SHA512 hash", name)
	case "ripemd128":
		return fmt.Sprintf("'%s' must be a valid RIPEMD-128 hash", name)
	case "ripemd160":
		return fmt.Sprintf("'%s' must be a valid RIPEMD-160 hash", name)
	case "tiger128":
		return fmt.Sprintf("'%s' must be a valid TIGER128 hash", name)
	case "tiger160":
		return fmt.Sprintf("'%s' must be a valid TIGER160 hash", name)
	case "tiger192":
		return fmt.Sprintf("'%s' must be a valid TIGER192 hash", name)
	case "semver":
		return fmt.Sprintf("'%s' must be a valid Semantic Versioning 2.0.0", name)
	case "ulid":
		return fmt.Sprintf("'%s' must be a valid Universally Unique Lexicographically Sortable Identifier (ULID)", name)
	case "cve":
		return fmt.Sprintf("'%s' must be a valid Common Vulnerabilities and Exposures Identifier (CVE id)", name)
	default:
		return fmt.Sprintf("'%s' is invalid", name)
	}
}
