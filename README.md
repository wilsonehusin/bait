# Bait

Bait is a thing-doer based on webhook receiver endpoint.

Initial use case is to refresh cache / mirror in Nginx deployment when upstream has updates.

```yaml
paths:
  - request: "/my-awesome-site.com"
    workdir: "/www"
    command: "wget my-awesome-site.com"
```

## Should you use it?

This software is "it works on my machine" quality in the foreseeable future. Highly recommended to only run in internal network.

## Trivia

You'd put the bait on the hook, wouldn't ya?
