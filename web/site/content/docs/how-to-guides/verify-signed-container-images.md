---
title: "Verify signed container images"
type: docs
description: >
    Verify that container image is signed when creating a container from it in Kanto Container Management.
weight: 5
---

By following the steps below, you will sign a container image and push it to a local registry using a{{% refn "https://github.com/notaryproject/notation" %}}`notation`{{% /refn %}}. Then a notation trust policy and the Kanto Container Management service will be configured in a way that running containers from the signed image via kanto-cm CLI will be successful, while running containers from unsigned images will fail.

### Before you begin

To ensure that your edge device is capable to execute the steps in this guide, you need:

* If you don't have an installed and running Eclipse Kanto, follow {{% relrefn "install" %}} Install Eclipse Kanto {{% /relrefn %}}
* Installed {{% refn "https://notaryproject.dev/docs/user-guides/installation/cli/" %}} Notation CLI {{% /refn %}}
* Installed and running {{% refn "https://www.docker.com/products/docker-desktop/" %}} Docker {{% /refn %}}

### Create an image and push it to a local registry using docker and than sign it with notation

Create and run a local container registry:
```shell
sudo kanto-cm create --ports 5000:5000 --e REGISTRY_STORAGE_DELETE_ENABLED=true --name registry docker.io/library/registry:latest
sudo kanto-cm start -n registry
```

Build a dummy hello world image and push it to the registry:
```shell
cat <<EOF | sudo docker build -t localhost:5000/dummy-hello:signed -
FROM busybox:latest
CMD [ "echo", "Hello World" ]
EOF
sudo docker push localhost:5000/dummy-hello:signed
```

{{% tip %}}
When signing and verifying container images it is recommended to use the image digest instead of a tag as the digest is immutable.
{{% /tip %}}
Get the image digest and assign it to an environment variable to be used the next steps of the guide:
```shell
export IMAGE=$(sudo docker inspect --format='{{index .RepoDigests 0}}' localhost:5000/dummy-hello:signed)
echo $IMAGE
```

Generate a key-pair with notation and add it as the default signing key:
```shell
notation cert generate-test --default "kanto"
```

Sign the image and store the signature in the registry:
```shell
notation sign $IMAGE
```

### Configure notation truspolicy and container management verifier

Get the notation config directory and assign it to an environment variable to be used in the next steps of the guide:
```shell
export NOTATION_CONFIG=${XDG_CONFIG_HOME:-$HOME/.config}/notation
echo $NOTATION_CONFIG
``` 

Create a simple {{% refn "https://github.com/notaryproject/specifications/blob/main/specs/trust-store-trust-policy.md#trust-policy" %}} notation trustpolicy {{% /refn %}} as a `trustpolicy.json` file in the notation config directory: 
```shell
cat <<EOF | tee $NOTATION_CONFIG/trustpolicy.json
{
  "version": "1.0",
  "trustPolicies": [
    {
      "name": "kanto-images",
      "registryScopes": [ "*" ],
      "signatureVerification": {
        "level" : "strict"
      },
      "trustStores": [ "ca:kanto" ],
      "trustedIdentities": [ "*" ]
    }
  ]
}
EOF
```

Create a backup of the initial Kanto Container Management configuration that is found in `/etc/container-management/config.json`(the backup will be restored at the end of the guide):
```shell
sudo cp /etc/container-management/config.json /etc/container-management/config-backup.json
```

Configure the use of notation verifier, set its config directory, mark the local registry as an insecure one, and set the image expiry time to zero seconds, so the local cache of the images used in the how-to will be deleted upon container removal:
```shell
cat <<EOF | sudo tee /etc/container-management/config.json
{
  "log": {
    "log_file": "/var/log/container-management/container-management.log"
  },
  "containers": {
    "image_verifier_type": "notation",
    "image_verifier_config": {
      "configDir": "$NOTATION_CONFIG"
    },
    "insecure_registries": [ "localhost:5000" ],
    "image_expiry": "0s"
  }
}
EOF
```

Restart the Container Management service for the changes to take effect:
```shell
sudo systemctl restart container-management.service
```

### Verify

Create and run a container from the signed image. The container prints `Hello world` to the console:
```shell
sudo kanto-cm create --name dummy-hello --rp no --t $IMAGE
sudo kanto-cm start --name dummy-hello --a
```


Make sure that a docker hub hello-world image is not cached locally, by removing any containers with this image, and verify that creating containers from it fails, as the image is not signed, and the signature verification fails:
```shell
sudo kanto-cm remove -f $(sudo kanto-cm list --quiet --filter image=docker.io/library/hello-world:latest)
sudo kanto-cm create --name dockerhub-hello --rp no --t docker.io/library/hello-world:latest
```

### Clean up

Remove the created containers from the Kanto Container Management:
```shell
sudo kanto-cm remove -n dummy-hello
sudo kanto-cm remove -n registry -f
```

Restore the initial Kanto Container Management configuration and restart the service:
```shell
sudo mv -f /etc/container-management/config-backup.json /etc/container-management/config.json
sudo systemctl restart container-management.service
```

Remove the localy cached images from Docker:
```shell
sudo docker image rm localhost:5000/dummy-hello:signed registry:latest
```

Reset the notation configuration by removing the directory:
```shell
rm -r $NOTATION_CONFIG
```

Unset exported environment variables:
```shell
unset IMAGE NOTATION_CONFIG
