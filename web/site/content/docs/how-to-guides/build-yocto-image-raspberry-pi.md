---
title: "Build Yocto Image for Raspberry Pi"
type: docs
description: >
    This page provides all the instructions needed to build yocto image for raspberry pi with kanto recipes
weight: 9
---

Building a Yocto Image for a Raspberry Pi involves several steps, including setting up the environment, download the necessary layers, configuring the build and compiling the image.

### Before you begin

* Install the required packages required for the build on a Ubuntu/Debian system
  ```shell
  sudo apt-get update
  sudo apt install -y gawk wget git-core diffstat unzip texinfo gcc-multilib \
  build-essential chrpath socat cpio python3 python3-pip python3-pexpect \
  xz-utils debianutils iputils-ping libsdl1.2-dev xterm zstd liblz4-tool \
  ```

### Clone the Yocto/Poky Repository

* Verify the Yocto Version of Kanto availabe in meta-kanto layer, {{% relrefn "hono" %}} https://github.com/eclipse-kanto/meta-kanto {{% /relrefn %}}, this example provides information for kirkstone branch.

* Create a Source folder 
  ```shell
  mkdir sources
  cd sources
  ``` 

* If it is kirkstone branch then, clone poky for kirkstone version in source directory
  
  ```shell
  git clone https://github.com/yoctoproject/poky.git
  cd poky
  git checkout kirkstone
  ```
  Note : Change branch based on meta-kanto layer.
  
### Add Meta Layers to the source directory

* Based on the yocto version, clone the below meta layers in the source directory
  `meta-raspberrypi`
  `meta-openembedded`
  `meta-virtualization`
  
* Clone meta-raspberry pi layer to sources directory
  ```shell
  cd ..
  git clone git://git.yoctoproject.org/meta-raspberrypi
  cd meta-raspberrypi
  git checkout kirkstone
  ```
  
* Clone meta-openembedded layer to sources directory
  ```shell
  cd ..
  git clone https://github.com/openembedded/meta-openembedded.git
  cd meta-openembedded
  git checkout kirkstone
  ```
  
* Clone meta-virtualization layer to sources directory
  ```shell
  cd ..
  git clone https://git.yoctoproject.org/git/meta-virtualization
  cd meta-virtualization
  git checkout kirkstone
  ```

### Add meta-kanto layer to the sources directory

* Clone meta-kanto layer to the sources directory
  ```shell
  cd ..
  git clone https://github.com/eclipse-kanto/meta-kanto.git
  cd meta-kanto
  git checkout kirkstone
  ```

  Note : Make sure all the layers cloned are of same yocto version ( kirkstone in this case)

### Create Build Directory

* After cloning all the required meta layers, move out of source directory to create build directory 
  
  ```shell
  cd ../..
  source sources/poky/oe-init-build-env
  ```
* Run the below command to view the layers present in `bblayers.conf` file

  ```shell
  bitbake-layers show-layers
  ```
  Note : Resolve any dependendencies if occured while running bitbake command.
   
### Configure bblayer.conf file

* Add all the layers to the bblayers.conf file with below command
  ```shell
  bitbake-layers add-layer /home/path/to/meta/layer/directory
  ```
* Add layers to bblayers.conf file
  
  ```shell
  bitbake-layers add-layers /path/to/meta/layer
  ```
  while adding layers bitbake might have dependencies, add the dependent layers.
  
  Note : Provide path for the meta-layer that has been cloned in the previous steps.
  
  Example,
  
  To add meta-kanto layer to bblayer.conf file which is kept at `/home/yocto/sources/meta-kanto`
  
  `bitbake-layers add-layers /home/yocto/sources/meta-kanto`
  
* After adding all the required layers in bblayer.conf file, verify again by running the below command
  
  ```shell
  bitbake-layers show-layers
  ```

### Configure local.conf file

* Open local.conf file which is placed at the below location in build directory
  
  ```shell
  vi conf/local.conf
  ```

* Change the machine variable in local.conf file to raspberry pi machine
  ```json
  MACHINE ??= "raspberrypi4" 
  ```

  Note: Check the sources/meta-raspberrypi/conf/machine for the availabe machines for raspberry pi.
  
* Add required variables in local.conf file as provided in the link below

  {{% relrefn "hono" %}} https://github.com/eclipse-kanto/meta-kanto {{% /relrefn %}} readme.md file
  
  ```json
    # Add the required DISTRO_FEATURES
    DISTRO_FEATURES:append = " virtualization systemd"

    # Configure the kernel modules required to be included
    MACHINE_ESSENTIAL_EXTRA_RRECOMMENDS += "kernel-modules"

    # System initialization manager setup
    VIRTUAL-RUNTIME_init_manager = "systemd"
    DISTRO_FEATURES_BACKFILL_CONSIDERED = "sysvinit"
    VIRTUAL-RUNTIME_initscripts = "systemd-compat-units"

    # Add the Eclipse Kanto components
    IMAGE_INSTALL:append = " suite-connector"
  ``` 
  
  Note : Add IMAGE_INSTALL:append with other kanto components to build all the kanto recepies.
  
  ```json
    IMAGE_INSTALL:append = " aws-connector"
    IMAGE_INSTALL:append = " azure-connector"
    IMAGE_INSTALL:append = " software-updates"
    IMAGE_INSTALL:append = " file-upload"
    IMAGE_INSTALL:append = " file-backup"
    IMAGE_INSTALL:append = " update-manager"
    IMAGE_INSTALL:append = " container-management"
    IMAGE_INSTALL:append = " local-digital-twins"
  ```
  
* Run bitbake `target-name' availabe as shown below,

  ```json
    Common targets are:
    core-image-minimal
    core-image-full-cmdline
    core-image-sato
    core-image-weston
    meta-toolchain
    meta-ide-support
  ```

  Note : Build issues if something arises needs to be resolved and again 'bitbake target-name' to be run.

### Build Image Repository

* After the successful build, the image will be availabe at the below location.

  ```json
  build/tmp/deploy/images/`machine_name`/
  ```
