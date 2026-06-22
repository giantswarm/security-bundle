# Automated Chart Updates (simplified)

A high-level view of how an upstream bump flows from Renovate into a testing PR on
`security-bundle`. For the full version with every workflow, see
[automated-chart-updates.md](./automated-chart-updates.md).

```mermaid
flowchart TD
    renovate(["👤 Renovate"])

    renovate -->|"pushes vendir.yml update"| branch["renovate/vendir/{APP}"]

    subgraph apprepo["App repo automation"]
        branch -->|"pushes to"| update["main#update-chart"]
        update -->|"calls"| sync["sync-from-upstream"]
        sync --> changes["updates schema / changelog / values"]
    end

    changes -->|"dispatches to security-bundle"| pr["creates testing PR"]

    subgraph sb["security-bundle"]
        pr --> version["updates PR with testing version"]
        version --> e2e["runs e2e testing"]
    end
```
