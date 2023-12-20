---
title: "Install Eclipse Kanto"
type: docs
description: >
    Run Eclipse Kanto on your edge device.
weight: 1
---

### Before you begin

The `containerd` Debian package is required. You can install it manually or run:

```shell
curl -fsSL https://github.com/eclipse-kanto/kanto/raw/main/quickstart/install_ctrd.sh | sh
```

### Install Eclipse Kanto

Choose the Eclipse Kanto Debian package for your target device architecture from the ones available
at <a href="https://github.com/eclipse-kanto/kanto/releases" target="_blank">the project's GitHub Releases page</a>.
Download and install it via executing the following (adjusted to your package name):

```shell
wget https://github.com/eclipse-kanto/kanto/releases/download/v0.1.0-M3/kanto_0.1.0-M3_linux_x86_64.deb && \
sudo apt install ./kanto_0.1.0-M3_linux_x86_64.deb
```

### Verify

It's important to check if all the services provided by the Eclipse Kanto package are up and running successfully. You
can quickly do that via executing:

```shell
systemctl status \
suite-connector.service \
container-management.service \
software-update.service \
file-upload.service \
file-backup.service \
system-metrics.service
```

All listed services must be in an active running state.

### What's next

[Explore via Eclipse Hono]({{< relref "hono" >}})
