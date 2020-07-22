#!/usr/bin/env bash

echo "*** test-realm-values-cli ****"
echo "ATLAS_PUBLIC_KEY=${ATLAS_PUBLIC_KEY}"
echo "ATLAS_PRIVATE_KEY=${ATLAS_PRIVATE_KEY}"
echo "ATLAS_GROUP_ID=${ATLAS_GROUP_ID}"
echo ""
TEST_APP=$( curl -s https://frightanic.com/goodies_content/docker-names.php | tr '_' '-' )
# List apps
echo "Running tests: TEST_APP=${TEST_APP}"
echo "*** "
echo "Listing apps -----------------------------------"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \

# Create a new app

echo "Creating test app: ${TEST_APP} -----------------------------------"
APP=$(./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--create-app \
--value "{\"name\": \"${TEST_APP}\" }")

echo "${APP}"
APP_ID=$(echo "${APP}" | jq -r '._id')
echo "${APP_ID}"

# List apps
echo "Listing apps -----------------------------------"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \



# List values
echo "Listing values for app=${APP_NAME} (${APP_ID})--"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}"

# Create new value
echo "Create new value for app=${APP_NAME} (${APP_ID})--"
key="key-123"
value='{ "a" : 1, "test" : [0, 55, "Hello World!"]}'
VALUE=$(./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}" \
--key "${key}" \
--value "${value}")

echo "${VALUE}"
VALUE_ID=$(echo "${VALUE}" | jq -r '._id')
echo "${VALUE_ID}"

# Get the new value back
echo "Get value for key=${key} VALUE_ID=${VALUE_ID} for app=${APP_NAME} (${APP_ID})--"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}" \
--key "${VALUE_ID}"

# List values
echo "Listing values for app=${APP_NAME} (${APP_ID})--"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}"

# Delete value
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}" \
--key "${VALUE_ID}"
--delete-key

# List values
echo "Listing values for app=${APP_NAME} (${APP_ID})--"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}"

# Delete app

echo "Deleting app=${APP_NAME} (${APP_ID})--"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \
--appid "${APP_ID}" \
--delete-app

# List apps
echo "Listing apps -----------------------------------"
./realm-values-cli \
--publicApiKey "${ATLAS_PUBLIC_KEY}" \
--privateApiKey "${ATLAS_PRIVATE_KEY}" \
--groupid "${ATLAS_GROUP_ID}" \


