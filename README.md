# Musannif

Server for Musannif, A sophisticated collaborative Markdown editor designed for teams, enabling seamless real-time collaboration and creativity.

## TODO

1. Init a basic logger - create wrapper over zerolog?
2. CI/CD: Publish new 'release' whenever a tagged commit is pushed
3. Mete out system architecture - how to approach horizontal scaling?
4. v1.0 File/'Note' management on disk & simple web-client for single user to create/edit their notes
5. Directory/Team management
6. Note sharing via URL - let guests view notes
7. ???

## Architecture

- Horizontal scaling of servers for geographical distribution
- Core idea: clients communicate diffs back and forth through a central server that resolves conflicts should any arise
