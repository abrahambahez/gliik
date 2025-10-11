# PRD-002: Tags and Language Support

## Overview

Add required `tags` and `lang` metadata fields to instructions for better organization, searchability, and internationalization support.

**Status:** Planning
**Version:** 0.2.0
**Depends on:** MVP 0.1.0

---

## Motivation

### Problems to Solve
1. **No categorization:** Users cannot organize instructions by topic/use case
2. **No language indication:** Multi-language instructions lack explicit language markers
3. **Poor terminal UX:** Current table format doesn't work well in narrow terminals
4. **Limited filtering:** Users rely on grep but tags aren't in output

### Goals
- Add structured categorization via tags
- Explicit language support via ISO 639-1 codes
- Improve list output for narrow terminals
- Enable grep-based filtering with visible tags

---

## Updated Metadata Schema

### `meta.yaml` Structure
```yaml
version: "1.0.0"
description: "Enhance resumes for job applications"
tags:
  - work
  - resume
lang: "en"
```

### Field Specifications

#### `tags` (required)
- **Type:** Array of strings
- **Validation:**
  - Must have at least one tag
  - Each tag: lowercase, alphanumeric + hyphens
  - Regex: `^[a-z0-9-]+$`
- **Examples:** `["coding", "review"]`, `["work", "cv", "job-search"]`

#### `lang` (required)
- **Type:** String
- **Validation:** ISO 639-1 two-letter code
- **Regex:** `^[a-z]{2}$`
- **Examples:** `"en"`, `"es"`, `"fr"`, `"de"`

### Migration Strategy
**Existing instructions without tags/lang:**
- Show warning when loading: `Warning: instruction 'X' missing required field 'tags' in meta.yaml`
- Show warning when loading: `Warning: instruction 'X' missing required field 'lang' in meta.yaml`
- Continue execution but mark as incomplete
- User must edit meta.yaml to add missing fields

---

## Updated Commands

### 1. `gliik add <name>`

#### New Flags
```bash
--tags, -t <tags>        Comma-separated tags (required)
--lang, -l <lang>        Language ISO code (required)
--description, -d <desc> Description (required)
```

#### Updated Behavior
- All three flags are now **required**
- Parse `--tags` comma-separated string into array
- Validate lang format matches `^[a-z]{2}$`
- Validate each tag matches `^[a-z0-9-]+$`
- Create meta.yaml with all fields populated

#### Examples
```bash
# Valid usage
gliik add cv_review -d "Review CVs" -t "work,hiring,review" -l "en"
gliik add resumen_mejorar -d "Mejorar resúmenes" -t "trabajo,cv" -l "es"

# Error cases
gliik add test -d "Test"
# Error: missing required flags: --tags, --lang

gliik add test -d "Test" -t "work" -l "EN"
# Error: lang must be lowercase two-letter ISO code (e.g., 'en')

gliik add test -d "Test" -t "Work,CV" -l "en"
# Error: tags must be lowercase alphanumeric with hyphens only
```

#### Error Messages
```
Error: missing required flags: --tags, --lang

Usage:
  gliik add <name> --description <text> --tags <tags> --lang <code>

Example:
  gliik add cv_review -d "Review CVs" -t "work,hiring" -l "en"
```

```
Error: invalid language code 'EN'

Language must be ISO 639-1 two-letter lowercase code
Examples: en, es, fr, de, pt
```

```
Error: invalid tag 'Work'

Tags must be lowercase alphanumeric with hyphens
Examples: work, code-review, job-search
```

---

### 2. `gliik list`

#### New Output Format
**Replace table format with compact line format:**

```
enhance_resume v1.0.0 [work, resume] - Enhance resumes for job applications
assess_paper v0.2.1 [academic, review] - Evaluate research papers for quality
code_review v1.1.0 [coding, review] - Review code for best practices
cv_mejorar v1.0.0 [trabajo, cv] - Mejorar currículums para aplicaciones
```

**Format specification:**
```
{name} v{version} [{tag1, tag2, ...}] - {description}
```

#### Benefits
- **Compact:** Works in narrow terminals (80+ columns)
- **Grep-friendly:** Tags visible for filtering
- **Sortable:** Natural alphabetical sorting
- **Readable:** Clear visual hierarchy

#### Filtering Examples
```bash
# Filter by tag
gliik list | grep "\[.*work.*\]"
gliik list | grep "\[.*resume.*\]"

# Filter by language (requires adding lang to output)
gliik list | grep "(es)"

# Filter by version
gliik list | grep "v1\."

# Multiple filters
gliik list | grep "\[.*coding.*\]" | grep "v1\."
```

#### Optional Enhancement
Consider adding language indicator:
```
enhance_resume v1.0.0 (en) [work, resume] - Enhance resumes for job applications
cv_mejorar v1.0.0 (es) [trabajo, cv] - Mejorar currículums para aplicaciones
```

---

## Implementation Plan

### Files to Modify

#### 1. `internal/instruction/instruction.go`
Update Meta struct:
```go
type Meta struct {
    Version     string   `yaml:"version"`
    Description string   `yaml:"description"`
    Tags        []string `yaml:"tags"`
    Lang        string   `yaml:"lang"`
}
```

#### 2. `internal/instruction/create.go`
- Update `Create()` signature: add `tags []string, lang string` parameters
- Add validation functions:
  - `ValidateLanguageCode(lang string) error`
  - `ValidateTags(tags []string) error`
- Update Meta initialization with new fields

#### 3. `cmd/add.go`
- Add required flags: `--tags`, `--lang`
- Make `--description` explicitly required
- Parse comma-separated tags into slice
- Validate inputs before calling Create()
- Update error messages

#### 4. `cmd/list.go`
- Replace tabwriter with simple line format
- Update output template:
  ```go
  fmt.Printf("%s v%s [%s] - %s\n",
      inst.Name,
      inst.Meta.Version,
      strings.Join(inst.Meta.Tags, ", "),
      inst.Meta.Description)
  ```
- Sort by name (already implemented)

#### 5. `internal/instruction/load.go` (if exists)
- Add validation warnings for missing fields
- Check if `Tags` or `Lang` are empty/nil
- Print warnings to stderr but continue execution

### Validation Functions

#### Language Validation
```go
func ValidateLanguageCode(lang string) error {
    validLangRegex := regexp.MustCompile(`^[a-z]{2}$`)
    if !validLangRegex.MatchString(lang) {
        return fmt.Errorf("invalid language code '%s': must be ISO 639-1 two-letter lowercase code (e.g., 'en', 'es', 'fr')", lang)
    }
    return nil
}
```

#### Tag Validation
```go
func ValidateTags(tags []string) error {
    if len(tags) == 0 {
        return fmt.Errorf("at least one tag is required")
    }

    validTagRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
    for _, tag := range tags {
        if !validTagRegex.MatchString(tag) {
            return fmt.Errorf("invalid tag '%s': must be lowercase alphanumeric with hyphens only", tag)
        }
    }
    return nil
}
```

---

## Migration & Backward Compatibility

### Handling Existing Instructions

When loading instructions without new fields:
1. Check if `Tags` is nil or empty
2. Check if `Lang` is empty
3. Print warnings to stderr:
   ```
   Warning: instruction 'old_name' missing required field 'tags' in meta.yaml
   Warning: instruction 'old_name' missing required field 'lang' in meta.yaml
   ```
4. Continue execution (degraded mode)
5. User must manually edit meta.yaml to fix

### Update Workflow
Users must manually update existing instructions:
```bash
# Option 1: Edit meta.yaml directly
vim ~/.gliik/instructions/old_name/meta.yaml

# Option 2: Re-create instruction (not recommended for complex ones)
gliik remove old_name -f
gliik add old_name -d "Description" -t "tag1,tag2" -l "en"
```

---

## Testing Checklist

- [ ] Create instruction with all required fields
- [ ] Create instruction missing --tags (should error)
- [ ] Create instruction missing --lang (should error)
- [ ] Create instruction with invalid lang code (should error)
- [ ] Create instruction with invalid tag format (should error)
- [ ] List instructions shows new format correctly
- [ ] List with tags containing commas formats properly
- [ ] Load legacy instruction without tags/lang shows warnings
- [ ] Grep filtering works with new format
- [ ] YAML serialization/deserialization works correctly

---

## Examples

### Creating Instructions
```bash
# English coding instruction
gliik add code_review \
  -d "Review code for best practices" \
  -t "coding,review,quality" \
  -l "en"

# Spanish work instruction
gliik add cv_mejorar \
  -d "Mejorar currículums para aplicaciones" \
  -t "trabajo,cv,empleo" \
  -l "es"

# Multi-domain instruction
gliik add api_doc_gen \
  -d "Generate API documentation" \
  -t "coding,documentation,api" \
  -l "en"
```

### Filtering Instructions
```bash
# All work-related
gliik list | grep "\[.*work.*\]"

# All Spanish instructions
gliik list | grep "(es)"

# Coding and review
gliik list | grep "\[.*coding.*\]" | grep "\[.*review.*\]"

# Version 1.x only
gliik list | grep "v1\."
```

### List Output Examples
```
code_review v1.0.0 [coding, review, quality] - Review code for best practices
cv_mejorar v0.1.0 [trabajo, cv, empleo] - Mejorar currículums para aplicaciones
api_doc_gen v1.2.0 [coding, documentation, api] - Generate API documentation
enhance_resume v2.0.0 [work, resume, job-search] - Enhance resumes for applications
```

---

## Out of Scope

- Tag autocomplete/suggestions
- Predefined tag taxonomy
- Language validation against full ISO 639-1 list (just regex check)
- Tag-based filtering built into `gliik list` (use grep)
- Multi-language support in single instruction
- Tag hierarchies or categories
- Instruction search command (use grep)

---

## Success Criteria

- [x] All new instructions require tags and lang
- [x] Validation enforces correct formats
- [x] List output works in narrow terminals (80+ cols)
- [x] Tags visible for grep-based filtering
- [x] Warnings shown for legacy instructions
- [x] Backward compatible with existing workflows
- [x] Clear error messages guide users
