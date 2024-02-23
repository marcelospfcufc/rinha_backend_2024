#!/usr/bin/bash

# Use este script para executar testes locais

RESULTS_WORKSPACE="/load-test/user-files/results"
GATLING_BIN_DIR=/opt/gatling/bin
GATLING_WORKSPACE="/load-test/user-files"

runGatling() {
    sh /opt/gatling/bin/gatling.sh -rm local -s RinhaBackendCrebitosSimulation \
        -rd "Rinha de Backend - 2024/Q1: Crébito" \
        -rf $RESULTS_WORKSPACE \
        -sf "$GATLING_WORKSPACE/simulations"
}

startTest() {
    for i in {1..20}; do
        # 2 requests to wake the 2 api instances up :)
        curl --fail http://localhost:8081/clientes/1/extrato && \
        echo "" && \
        curl --fail http://localhost:8081/clientes/1/extrato && \
        echo "" && \
        runGatling && \
        break || sleep 2;
    done
}

startTest