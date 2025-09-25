#!/usr/bin/env bash
# Copyright (c) Palo Alto Networks, Inc.
# SPDX-License-Identifier: MPL-2.0


set -eo pipefail

# --- Argument Parsing ---
ARGS=()
AUTO_APPROVE=false
for arg in "$@"; do
  case $arg in
    -y|--yes|--auto-approve) 
      AUTO_APPROVE=true
      ;;
    *)
      ARGS+=("$arg")
      ;;
  esac
done

TARGET_MODULES_CSV=${ARGS[0]}
TAG_VERSION=${ARGS[1]}


# --- Defaults ---
# Find all modules (directories with a go.mod file)
ALL_MODULES=$(find . -name go.mod -exec dirname {} \; | sed 's|^\./||' | tr '\n' ',' | sed 's/,$//')

# If module list is empty, use all modules.
if [ -z "$TARGET_MODULES_CSV" ]; then
    TARGET_MODULES_CSV=$ALL_MODULES
fi


# --- Tag Calculation ---
# Convert CSV to an array to iterate over
IFS=',' read -r -a TARGET_MODULES <<< "$TARGET_MODULES_CSV"

echo "Target modules: ${TARGET_MODULES_CSV}"
echo "---"

TAGS_TO_CREATE=()

for module in "${TARGET_MODULES[@]}"; do
    # Trim whitespace
    module=$(echo "$module" | xargs)
    if [ -z "$module" ]; then
        continue
    fi

    echo "Processing module: $module"

    # If a specific version is provided, use it.
    # Otherwise, calculate the next patch version.
    if [ -n "$TAG_VERSION" ]; then
        NEW_VERSION=$TAG_VERSION
    else
        # Find the latest git tag for the module
        LATEST_TAG=$(git tag -l "${module}/v*" | sort -V | tail -n 1)

        if [ -z "$LATEST_TAG" ]; then
            # If no tag exists, start with v0.0.1
            NEW_VERSION="v0.0.1"
        else
            echo "Latest tag found: $LATEST_TAG"
            # Extract version
            LATEST_TAG_VERSION=$(echo "$LATEST_TAG" | sed -E "s|${module}/v||")
            
            # Increment the patch version
            MAJOR=$(echo "$LATEST_TAG_VERSION" | cut -d. -f1)
            MINOR=$(echo "$LATEST_TAG_VERSION" | cut -d. -f2)
            PATCH=$(echo "$LATEST_TAG_VERSION" | cut -d. -f3)
            
            if ! [[ "$PATCH" =~ ^[0-9]+$ ]]; then
                echo "Error: Could not parse patch version from '$LATEST_TAG'. Expected a number."
                exit 1
            fi

            NEW_PATCH=$((PATCH + 1))
            NEW_VERSION="v${MAJOR}.${MINOR}.${NEW_PATCH}"
        fi
    fi

    NEW_TAG="${module}/${NEW_VERSION}"
    echo "Calculated new tag: ${NEW_TAG}"
    TAGS_TO_CREATE+=("$NEW_TAG")
    
    echo "---"
done


# --- Approval and Tagging ---
if [ ${#TAGS_TO_CREATE[@]} -eq 0 ]; then
    echo "No tags to create."
    exit 0
fi

echo "The following tags will be created:"
for tag in "${TAGS_TO_CREATE[@]}"; do
    echo "  - $tag"
done
echo ""

if [ "$AUTO_APPROVE" = false ]; then
    read -p "Apply these tags? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted by user."
        exit 1
    fi
fi

echo "Applying tags..."
for tag in "${TAGS_TO_CREATE[@]}"; do
    git tag "$tag"
    echo "Created tag: $tag"
done

echo "---"
echo "Script finished successfully."
