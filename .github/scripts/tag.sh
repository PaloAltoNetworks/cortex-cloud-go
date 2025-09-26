#!/usr/bin/env bash

set -eo pipefail

# --- Help Function ---
usage() {
  echo "Usage: $0 <command> [options]"
  echo ""
  echo "A script to manage module tags and releases."
  echo ""
  echo "Commands:"
  echo "  create [modules_csv] [version]   Create new git tags for modules."
  echo "                                     - [modules_csv] (optional): Comma-separated list of modules. Defaults to all."
  echo "                                     - [version] (optional): Specific version (e.g., v1.2.3). Defaults to next patch."
  echo ""
  echo "  delete <tag_or_version>          Delete GitHub releases and associated git tags."
  echo "                                     - <tag_or_version>: A full tag (e.g., 'platform/v0.1.0') or just a version (e.g., 'v0.1.0') to delete across all modules."
  echo ""
  echo "Options:"
  echo "  -y, --yes, --auto-approve        Bypass confirmation prompts."
  exit 1
}

# --- Command Implementations ---

function cmd_create() {
    TARGET_MODULES_CSV=${ARGS[0]}
    TAG_VERSION=${ARGS[1]}

    # Find all modules (directories with a go.mod file)
    ALL_MODULES=$(find . -name go.mod -exec dirname {} \; | sed 's|^\./||' | tr '\n' ',' | sed 's/,$//')

    # If module list is empty, use all modules.
    if [ -z "$TARGET_MODULES_CSV" ]; then
        TARGET_MODULES_CSV=$ALL_MODULES
    fi

    # Convert CSV to an array to iterate over
    IFS=',' read -r -a TARGET_MODULES <<< "$TARGET_MODULES_CSV"

    echo "Target modules: ${TARGET_MODULES_CSV}"
    echo "---"

    TAGS_TO_CREATE=()

    for module in "${TARGET_MODULES[@]}"; do
        module=$(echo "$module" | xargs) # Trim whitespace
        if [ -z "$module" ]; then
            continue
        fi

        echo "Processing module: $module"

        if [ -n "$TAG_VERSION" ]; then
            if [[ ! "$TAG_VERSION" =~ ^v ]]; then
                TAG_VERSION="v$TAG_VERSION"
            fi
            NEW_VERSION=$TAG_VERSION
        else
            LATEST_TAG=$(git tag -l "${module}/v*" | sort -V | tail -n 1)

            if [ -z "$LATEST_TAG" ]; then
                NEW_VERSION="v0.0.1"
            else
                echo "Latest tag found: $LATEST_TAG"
                LATEST_TAG_VERSION=${LATEST_TAG#${module}/v}
                
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
}

function cmd_delete() {
    if [ ${#ARGS[@]} -eq 0 ]; then
        echo "Error: 'delete' command requires a tag or version argument."
        usage
    fi
    
    TAG_OR_VERSION=${ARGS[0]}
    TAGS_TO_DELETE=()

    # Check if the argument is a full tag (contains '/'') or just a version.
    if [[ "$TAG_OR_VERSION" == */* ]]; then
        TAGS_TO_DELETE+=("$TAG_OR_VERSION")
    else
        VERSION=$TAG_OR_VERSION
        if [[ ! "$VERSION" =~ ^v ]]; then
            VERSION="v$VERSION"
        fi
        echo "Searching for all tags with version '$VERSION'வுகளை"
        # Find all tags across all modules matching the version
        readarray -t found_tags < <(git tag -l "*/$VERSION")
        if [ ${#found_tags[@]} -gt 0 ]; then
            TAGS_TO_DELETE+=("${found_tags[@]}")
        fi
    fi

    if [ ${#TAGS_TO_DELETE[@]} -eq 0 ]; then
        echo "No matching tags found to delete for '$TAG_OR_VERSION'."
        exit 0
    fi

    echo "The following will be deleted (release and git tag):"
    for tag in "${TAGS_TO_DELETE[@]}"; do
        echo "  - $tag"
    done
    echo ""

    if [ "$AUTO_APPROVE" = false ]; then
        read -p "Proceed with deletion? This is irreversible. (y/n) " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted by user."
            exit 1
        fi
    fi

    for tag in "${TAGS_TO_DELETE[@]}"; do
        echo "---"
        echo "Processing deletion for: $tag"

        # First, try to delete the GitHub release. The --cleanup-tag flag will also delete the git tag.
        if gh release view "$tag" >/dev/null 2>&1; then
            echo "Found release for '$tag'. Deleting release and associated git tag..."
            gh release delete "$tag" --yes --cleanup-tag
            echo "Successfully deleted release and git tag '$tag'."
        else
            # If no release exists, delete the git tag manually from local and remote.
            echo "No GitHub release found for '$tag'. Deleting git tag only."
            
            # Delete remote tag
            if git ls-remote --tags origin | grep -q "refs/tags/$tag$"; then
                echo "Deleting remote git tag '$tag'..."
                git push --delete origin "$tag"
            else
                echo "No remote git tag found for '$tag'."
            fi

            # Delete local tag
            if git rev-parse "$tag" >/dev/null 2>&1; then
                echo "Deleting local git tag '$tag'..."
                git tag -d "$tag"
            else
                echo "No local git tag found for '$tag'."
            fi
        fi
    done
    
    echo "---"
    echo "Deletion complete."
}

#  To Create Tags (existing functionality):
#
#   1 # Calculate and create the next patch tags for all modules
#   2 .github/scripts/tag.sh create
#   3
#   4 # Create a specific version tag for the 'platform' module
#   5 .github/scripts/tag.sh create platform v1.5.0
#
#  To Delete Releases and Tags (new functionality):
#
#   1 # Delete a single, specific release and tag
#   2 .github/scripts/tag.sh delete platform/v1.5.0
#   3
#   4 # Find and delete all releases/tags for version v1.5.0 across all modules
#   5 .github/scripts/tag.sh delete v1.5.0
#   6
#   7 # You can bypass the confirmation prompt with --yes
#   8 .github/scripts/tag.sh delete v1.5.0 --yes

# --- Argument Parsing ---
if [ "$#" -eq 0 ]; then
  usage
fi

COMMAND=$1
shift

ARGS=()
AUTO_APPROVE=false
# Parse global flags first
while [[ "$#" -gt 0 ]]; do
  case $1 in
    -y|--yes|--auto-approve)
      AUTO_APPROVE=true
      shift
      ;;
    *)
      ARGS+=("$1")
      shift
      ;;
  esac
done


# --- Main Dispatcher ---
case $COMMAND in
  create)
    cmd_create
    ;;
  delete)
    cmd_delete
    ;;
  *)
    echo "Error: Unknown command '$COMMAND'"
    usage
    ;;
esac


