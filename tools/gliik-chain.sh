#!/usr/bin/env bash
#
# gliik-chain - Execute Gliik instructions with intelligent context selection
#
# Usage:
#   gliik-chain <context-file> <instruction> [flags...]
#
# Example:
#   gliik-chain context.md analyze_metrics --output report.md
#
# This tool chains two Gliik instructions:
# 1. select_context: Determines which linked files are relevant
# 2. <instruction>: Executes with selected context files

set -euo pipefail

if [ $# -lt 2 ]; then
    cat << 'EOF'
Usage: gliik-chain <context-file> <instruction> [flags...]

Executes a Gliik instruction with intelligent context selection.

Arguments:
  context-file    Path to context file with [link: ...] references
  instruction     Name of Gliik instruction to execute
  [flags...]      Additional flags passed to the instruction

Example:
  gliik-chain context.md analyze_metrics --output report.md

Process:
  1. Runs 'select_context' to determine relevant linked files
  2. Combines context file with selected files
  3. Executes specified instruction with combined context

Requirements:
  - gliik must be installed and in PATH
  - 'select_context' instruction must exist
  - Context file must contain [link: path] references
EOF
    exit 1
fi

CONTEXT_FILE="$1"
INSTRUCTION="$2"
shift 2

# Validate context file exists
if [ ! -f "$CONTEXT_FILE" ]; then
    echo "Error: Context file not found: $CONTEXT_FILE" >&2
    exit 1
fi

# Get directory of context file for resolving relative paths
CONTEXT_DIR=$(dirname "$CONTEXT_FILE")

# Execute select_context to get list of needed files
echo "→ Selecting relevant context files..." >&2
SELECTED_FILES=$(cat "$CONTEXT_FILE" | gliik run select_context --intent "$INSTRUCTION" "$@" 2>/dev/null)

if [ -z "$SELECTED_FILES" ]; then
    echo "→ No additional context files selected" >&2
    cat "$CONTEXT_FILE" | gliik run "$INSTRUCTION" "$@"
    exit 0
fi

# Display selected files
echo "→ Selected context files:" >&2
echo "$SELECTED_FILES" | sed 's/^/  - /' >&2

# Resolve paths relative to context file directory and validate
RESOLVED_FILES=""
while IFS= read -r file; do
    # Skip empty lines
    [ -z "$file" ] && continue
    
    # Resolve relative path
    if [[ "$file" = /* ]]; then
        # Absolute path
        FULL_PATH="$file"
    else
        # Relative path - resolve relative to context file directory
        FULL_PATH="$CONTEXT_DIR/$file"
    fi
    
    if [ ! -f "$FULL_PATH" ]; then
        echo "Warning: Selected file not found: $FULL_PATH" >&2
        continue
    fi
    
    RESOLVED_FILES="$RESOLVED_FILES $FULL_PATH"
done <<< "$SELECTED_FILES"

# Execute instruction with combined context
echo "→ Executing instruction: $INSTRUCTION" >&2
if [ -n "$RESOLVED_FILES" ]; then
    cat "$CONTEXT_FILE" $RESOLVED_FILES | gliik run "$INSTRUCTION" "$@"
else
    # No valid files found, use just context file
    cat "$CONTEXT_FILE" | gliik run "$INSTRUCTION" "$@"
fi
