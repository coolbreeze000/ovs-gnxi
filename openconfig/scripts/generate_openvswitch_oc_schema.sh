#!/bin/bash

# Source: https://github.com/google/link022/blob/master/openconfig/scripts/generate_wifi_oc_schema.sh

OVSGNXI=$HOME/go/src/ovs-gnxi

# Tools
YANG_CONVERTER=$OVSGNXI/vendor/github.com/openconfig/ygot/generator/generator.go

# Download OpenConfig models from https://github.com/openconfig/public
# Download ietf models from https://github.com/openconfig/yang/tree/master/standard/ietf/RFC
# Move downloaded models to a specific folder.
OC_FOLDER=$OVSGNXI/openconfig/models

# OpenConfig modules
YANG_MODELS=$OC_FOLDER/public/release/models
IETF_MODELS=$OC_FOLDER/yang/standard/ietf/RFC
OVS_TOP_MODULE=$OC_FOLDER/public/release/models/openflow/openconfig-openflow.yang
IGNORED_MODULES=openconfig-system,openconfig-extensions,openconfig-inet-types,openconfig-platform,openconfig-interfaces

# Output path
OUTPUT_PACKAGE_NAME=ocstruct
OUTPUT_FILE_PATH=$OVSGNXI/generated/$OUTPUT_PACKAGE_NAME/$OUTPUT_PACKAGE_NAME.go
mkdir -p $OVSGNXI/generated/$OUTPUT_PACKAGE_NAME

go run $YANG_CONVERTER \
-path=$YANG_MODELS,$IETF_MODELS \
-generate_fakeroot -fakeroot_name=device \
-package_name=ocstruct -compress_paths=false \
-exclude_modules=$IGNORED_MODULES \
-output_file=$OUTPUT_FILE_PATH \
$OVS_TOP_MODULE