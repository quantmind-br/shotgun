# Database Schema

N/A - Shotgun CLI uses the file system for all persistence. Configuration and templates are stored as files:

```
~/.config/shotgun-cli/
├── config.toml          # User configuration
├── templates/           # User-defined templates
│   ├── my_template.toml
│   └── team_template.toml
└── history.json         # Session history (future feature)
```
