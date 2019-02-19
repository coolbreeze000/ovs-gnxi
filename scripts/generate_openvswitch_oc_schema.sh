#!/bin/bash

# Source: https://github.com/google/link022/blob/master/openconfig/scripts/generate_wifi_oc_schema.sh

OVSGNXI=$HOME/go/src/ovs-gnxi

# Tools
YANG_CONVERTER=$OVSGNXI/vendor/github.com/openconfig/ygot/generator/generator.go

# Download OpenConfig models from https://github.com/openconfig/public
# Download ietf models from https://github.com/openconfig/yang/tree/master/standard/ietf/RFC
# Move downloaded models to a specific folder.
MODEL_FOLDER=$OVSGNXI/yang

# OpenConfig modules
IGNORED_MODULES=ietf-interfaces
OC_MODELS=$MODEL_FOLDER/openconfig
IETF_MODELS=$MODEL_FOLDER/ietf

# Output path
OUTPUT_PACKAGE_NAME=ocstruct
OUTPUT_FILE_PATH=$OVSGNXI/generated/$OUTPUT_PACKAGE_NAME/$OUTPUT_PACKAGE_NAME.go
mkdir -p $OVSGNXI/generated/$OUTPUT_PACKAGE_NAME

go run $YANG_CONVERTER \
-path=yang \
-generate_fakeroot \
-compress_paths=true \
-package_name=ocstruct \
-exclude_modules=$IGNORED_MODULES \
-output_file=$OUTPUT_FILE_PATH \
$OC_MODELS/openconfig-interfaces.yang \
$OC_MODELS/openconfig-openflow.yang \
$OC_MODELS/openconfig-platform.yang \
$OC_MODELS/openconfig-system.yang \
