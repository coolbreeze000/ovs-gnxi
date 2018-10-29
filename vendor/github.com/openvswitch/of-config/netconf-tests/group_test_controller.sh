#!/bin/bash
. ./config

./reset.sh
./get-config.sh > start_state
./run_edit_config.sh create_controller.xml running
check_startup_nonchange start_state
./get-config.sh | tee before_change
./run_edit_config.sh change_controller.xml running
check_startup_nonchange before_change
./get-config.sh
./run_edit_config.sh remove_controller.xml running
check_startup_difference start_state

rm -f start_state before_change

