[Back to Table of Contents](../documentation.md)


# Installation and Setup

## Install

### [Download the latest binary](https://github.com/vmare-tanzu-labs/yaml-overlay-tool/releases/latest)

### wget
Use wget to download the pre-compiled binaries:

```bash
wget https://github.com/vmware-tanzu-labs/yaml-overlay-tool/releases/download/${VERSION}/${BINARY}.tar.gz -O - |\
  tar xz && mv ${BINARY} /usr/bin/yot
```

For instance, VERSION=v0.3.1 and BINARY=yot_${VERSION}_linux_amd64

### Run with Docker

#### Oneshot use:

```bash
docker run --rm -v "${PWD}":/workdir ghcr.io/vmware-tanzu-labs/yot [flags]
```

#### Run commands interactively:

```bash
docker run --rm -it -v "${PWD}":/workdir --entrypoint sh ghcr.io/vmawre-tanzu-labs/yot
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yot() {
  docker run --rm -i -v "${PWD}":/workdir ghcr.io/vmware-tanzu-labs/yot "$@"
}
```


[Back to Table of Contents](../documentation.md)  
[Next Up: Command Line Usage and Overview](usage.md)