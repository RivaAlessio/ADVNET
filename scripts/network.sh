#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh

function startNetwork() {
    local log_level=$1

    infoln "Starting the network"
    infoln $FABRIC_CFG_PATH

    FABRIC_LOG=$log_level COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION docker-compose -f ${DOCKER_COMPOSE_PATH} up -d 2>&1

    docker ps -a
    if [ $? -ne 0 ]; then
        fatalln "Unable to start network"
    fi
}

function stopNetwork() {
    infoln "Stopping the network"

    set -x
    FABRIC_LOG=INFO COMPOSE_PROJECT_NAME=$PROJECT_NAME PROJECT_NAME=$PROJECT_NAME IMAGE_TAG=$FABRIC_VERSION  docker-compose -f ${DOCKER_COMPOSE_PATH} down -v --rmi all 2>&1
    { set +x; } 2>/dev/null
}

function clearNetwork() {
    infoln "Cleaning the repository"

    stopNetwork

    # remove organizations
    rm -rf $ORGANIZATION_OUTPUTS

    # remove volumes
    rm -rf volumes

    # remove channels
    rm -rf channels
}


MODE=$1
LOG_LEVEL=$2 
if [ "$MODE" == "start" ]; then
  startNetwork $LOG_LEVEL
elif [ "$MODE" == "stop" ]; then
  stopNetwork
elif [ "$MODE" == "clear" ]; then
  clearNetwork  
fi