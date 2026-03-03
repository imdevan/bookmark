# Context

This plan outlines improvements to the documentation system for the bookmark CLI tool. The goal is to create a clean, auto-generated documentation site that derives content from code comments, markdown files, and CLI structure.

# Definitions

- **Lander**: The landing/home page of the documentation site
- **Sidebar**: The navigation menu showing all documentation sections
- **Command docs**: Auto-generated documentation for CLI commands from godoc comments
- **Aesthetically pleasing**: Clean, modern design similar to https://devan.gg/prompter-cli/install/
- **just docs_* commands**: Justfile commands for generating and managing documentation

# v1.0.0

- [x] 1. Remove Current Lander
  - [x] 1.1 Remove existing landing page implementation
    notes: 
      root should link directly to docs
      same as first item sidebar (package_name from README)

- [ ] 2. Refactor just docs-generate script
  - [x] 2.1 Create sidebar layout structure
    - notes: 
      Sidebar should include: package_name (from README), install (from INSTALL.md), commands section, configuration, contributing
      the updates should happen within the docs-generate script. 
      and executed on just docs-dev and just docs-build
  - [x] 2.2 Generate package_name section from README.md
    - notes: Main package overview and introduction
  - [x] 2.3 Generate install section from INSTALL.md
    - notes: Installation instructions
      for this page specifically h3 elements should be converted to starlight tabs
      `import { Tabs, TabItem } from "@astrojs/starlight/components";`

  - [ ] 2.4 Auto-generate commands documentation
    - notes: Derive from cmd/bookmark/*.go files using godoc comments
  - [ ] 2.5 Generate root command documentation
    - notes: Document the main `bookmark` command from cmd/bookmark/root.go
  - [ ] 2.6 Generate subcommand documentation
    - notes: Auto-detect and document all subcommands from cmd/bookmark/[command].go files
  - [ ] 2.7 Generate configuration documentation
    - notes: Extract from domain.config godocs
  - [ ] 2.8 Generate contributing section
    - notes: Content from CONTRIBUTING.md if it exists

- [ ] 3. Implement Header Navigation
  - [ ] 3.1 Add package name to left side of header
    - notes: Display the package/project name prominently
  - [ ] 3.2 Add GitHub icon link to right side
    - notes: Link to repository
  - [ ] 3.3 Add dark/light mode toggle to right side
    - notes: Theme switcher for user preference
  - [ ] 3.4 Add search functionality to right side
    - notes: Enable searching across documentation

- [ ] 4. Update just docs-generate Command
  - [ ] 4.1 Ensure docs generation commands are helpful
    - notes: Commands should be intuitive and well-documented
  - [ ] 4.2 Ensure generated docs are aesthetically pleasing
    - notes: Follow design patterns from https://devan.gg/prompter-cli/install/
  - [ ] 4.3 Automate godoc comment extraction
    - notes: Parse Go files for documentation comments
  - [ ] 4.4 Support custom sidebar items
    - notes: Allow users to add additional documentation sections

- [ ] 5. Edit page should link back to appropriate go doc  or markdown source

# v2.0.0

- [ ] 1. Optional Lander Feature
  - [ ] 1.1 Implement configurable landing page
    - notes: Allow users to optionally enable a custom landing page
  - [ ] 1.2 Add lander configuration options
    - notes: Make lander opt-in with customization support

