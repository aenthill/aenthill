#!/bin/sh

set -e

print_pheromones()
{
    printf "\n\n-- PHEROMONES --\n\n"
    env
}

print_manifest()
{
    MANIFEST_PATH="${PHEROMONE_CONTAINER_PROJECT_DIR}/aenthill.json"
    printf "\n\n-- MANIFEST --\n\n"
    printf "Manifest should be located at ${MANIFEST_PATH}\n"
    printf "Does it exist? "
    if [ -f ${MANIFEST_PATH} ]; then
        printf "YES\n\n"
        cat ${MANIFEST_PATH}
    else
        printf "NO"
    fi
}

test_register()
{
    printf "\n\nTesting register...\n"
    aenthill register aenthill/cassandra FOO -m FOO=BAR
    printf "\nRegister done!"
    print_pheromones
    print_manifest
}

EVENT="$1"
PAYLOAD="$2"

printf "Received event ${EVENT} with payload ${PAYLOAD}"
print_pheromones
print_manifest

if [ ${EVENT} = "TEST_REGISTER" ]; then
    test_register
fi

printf "\nBye!\n"