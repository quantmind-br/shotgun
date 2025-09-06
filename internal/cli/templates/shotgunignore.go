package templates

// ShotgunignoreTemplate contains the default .shotgunignore content
const ShotgunignoreTemplate = `# Shotgun ignore patterns
# This file specifies patterns for files and directories to exclude
# when scanning project files. Patterns follow .gitignore syntax.

# Build artifacts
build/
dist/
target/
out/
bin/
obj/
*.exe
*.dll
*.so
*.dylib

# Dependencies  
node_modules/
vendor/
.pnp/
.yarn/
Cargo.lock
package-lock.json
yarn.lock

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~
.project
.classpath
.settings/

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
Desktop.ini

# Logs and temporary files
*.log
logs/
tmp/
temp/
.tmp/

# Runtime and environment files
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
*.pid

# Version control
.git/
.svn/
.hg/
.bzr/

# Package managers
.npm/
.yarn-cache/
.pnpm-store/

# Go specific
go.sum
*.test
*.out
coverage.txt

# Language specific artifacts
__pycache__/
*.py[cod]
*$py.class
*.class
*.jar
target/
.gradle/
.mvn/
`

// IgnoreCategory represents a category of ignore patterns
type IgnoreCategory struct {
	Name        string
	Description string
	Patterns    []string
}

// ShotgunignoreCategories provides organized ignore patterns by category
var ShotgunignoreCategories = []IgnoreCategory{
	{
		Name:        "Build artifacts",
		Description: "Files and directories created during build processes",
		Patterns: []string{
			"build/",
			"dist/",
			"target/",
			"out/",
			"bin/",
			"obj/",
			"*.exe",
			"*.dll",
			"*.so",
			"*.dylib",
		},
	},
	{
		Name:        "Dependencies",
		Description: "Third-party dependencies and package manager files",
		Patterns: []string{
			"node_modules/",
			"vendor/",
			".pnp/",
			".yarn/",
			"Cargo.lock",
			"package-lock.json",
			"yarn.lock",
		},
	},
	{
		Name:        "IDE and editor files",
		Description: "Files created by IDEs and text editors",
		Patterns: []string{
			".vscode/",
			".idea/",
			"*.swp",
			"*.swo",
			"*~",
			".project",
			".classpath",
			".settings/",
		},
	},
	{
		Name:        "OS generated files",
		Description: "Files created by operating systems",
		Patterns: []string{
			".DS_Store",
			".DS_Store?",
			"._*",
			".Spotlight-V100",
			".Trashes",
			"ehthumbs.db",
			"Thumbs.db",
			"Desktop.ini",
		},
	},
	{
		Name:        "Logs and temporary files",
		Description: "Log files and temporary data",
		Patterns: []string{
			"*.log",
			"logs/",
			"tmp/",
			"temp/",
			".tmp/",
		},
	},
	{
		Name:        "Runtime and environment files",
		Description: "Runtime configuration and environment variables",
		Patterns: []string{
			".env",
			".env.local",
			".env.development.local",
			".env.test.local",
			".env.production.local",
			"*.pid",
		},
	},
	{
		Name:        "Version control",
		Description: "Version control system directories",
		Patterns: []string{
			".git/",
			".svn/",
			".hg/",
			".bzr/",
		},
	},
	{
		Name:        "Package managers",
		Description: "Package manager cache and data directories",
		Patterns: []string{
			".npm/",
			".yarn-cache/",
			".pnpm-store/",
		},
	},
	{
		Name:        "Go specific",
		Description: "Go language specific build artifacts and files",
		Patterns: []string{
			"go.sum",
			"*.test",
			"*.out",
			"coverage.txt",
		},
	},
	{
		Name:        "Language specific artifacts",
		Description: "Compiled files and artifacts from various programming languages",
		Patterns: []string{
			"__pycache__/",
			"*.py[cod]",
			"*$py.class",
			"*.class",
			"*.jar",
			"target/",
			".gradle/",
			".mvn/",
		},
	},
}
