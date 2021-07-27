[Back to Table of contents](../index.md)


# Installation and setup

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Manual install](#manual-install)
  - [Download the latest binary](#download-the-latest-binaryhttpsgithubcomvmare-tanzu-labsyaml-overlay-toolreleaseslatest)
  - [wget](#wget)
- [MacOS / Linux via Homebrew install](#macos-linux-via-homebrew-install)
- [Linux snap install](#linux-snap-install)
- [Docker image pull](#docker-image-pull)
  - [One-shot container use](#one-shot-container-use)
  - [Run container commands interactively](#run-container-commands-interactively)
- [Go install](#go-install)

<!-- /code_chunk_output -->


## Manual install

### [Download the latest binary](https://github.com/vmare-tanzu-labs/yaml-overlay-tool/releases/latest)

### wget
Use wget to download the pre-compiled binaries:

```bash
wget https://github.com/vmware-tanzu-labs/yaml-overlay-tool/releases/download/${VERSION}/${BINARY}.tar.gz -O - |\
  tar xz && mv ${BINARY} /usr/bin/yot
```

For instance, VERSION=v0.3.1 and BINARY=yot_${VERSION}_linux_amd64

## MacOS / Linux via Homebrew install

Using [Homebrew](https://brew.sh/)  

```bash
brew tap vmware-tanzu-labs/tap
brew install yot
```

## Linux snap install

```bash
snap install yaml-overlay-tool
```

>**NOTE**: `yot` installs with [_strict confinement_](https://docs.snapcraft.io/snap-confinement/6233) in snap, this means it doesn't have direct access to root files. To read root files you can:

```bash
sudo cat /etc/myfile.yaml | yot -i -
```

And to write to a root file you can either use [sponge](https://linux.die.net/man/1/sponge):

```bash
sudo cat /etc/myfile.yaml | yot -s -i - | sudo sponge /etc/myfile.yaml
```

or write to a temporary file:

```bash
sudo cat /etc/myfile.yaml | yot -s -i  | sudo tee /etc/myfiletmp.yaml
sudo mv /etc/myfiletmp.yaml /etc/myfile.yaml
```

## Docker image pull

```bash
docker pull ghcr.io/vmawre-tanzu-labs/yot
```

### One-shot container use

```bash
docker run --rm -v "${PWD}":/workdir ghcr.io/vmware-tanzu-labs/yot [flags]
```


### Run container commands interactively

```bash
docker run --rm -it -v "${PWD}":/workdir --entrypoint sh ghcr.io/vmawre-tanzu-labs/yot
```

It can be useful to have a bash function to avoid typing the whole docker command:

```bash
yot() {
  docker run --rm -i -v "${PWD}":/workdir ghcr.io/vmware-tanzu-labs/yot "$@"
}
```


## Go install

```bash
GO111MODULE=on go get github.com/vmware-tanzu-labs/yaml-overlay-tool/cmd/yot
```


[Back to Table of contents](../index.md)  
[Next Up: Configuration file](configFile.md)