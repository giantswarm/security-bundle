# Automated Chart Updates (simplified)

A high-level view of how an upstream bump flows from Renovate into a testing PR on
`security-bundle`. For the full version with every workflow, see
[automated-chart-updates.md](./automated-chart-updates.md).

```mermaid
sequenceDiagram
    actor Renovate
    participant App as App repo automation
    participant SB as security-bundle

    Renovate->>App: pushes vendir.yml update to renovate/vendir/{APP}
    App->>App: pushes to main#update-chart
    App->>App: calls sync-from-upstream
    App->>App: updates schema / changelog / values
    App->>SB: dispatches
    SB->>SB: creates testing PR
    SB->>SB: updates PR with testing version
    SB->>SB: runs e2e testing
```
