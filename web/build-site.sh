#!/bin/bash
# Exit on error
set -e

if [[ -z "${DEPLOY_PRIME_URL}" ]]; then
  USE_BASE_URL_SITE=""
else
  echo -e "\033[0;32mNetlify DEPLOY_PRIME_URL detected, this seems to be a PR, deployment happening at ${DEPLOY_PRIME_URL}...\033[0m"
  USE_BASE_URL_SITE="--baseURL ${DEPLOY_PRIME_URL}"
fi

echo -e "\033[0;32mDeleting Hugo previously generated directories...\033[0m"
rm -rf site/public

echo -e "\033[0;32mSwitching to Docsy theme directory...\033[0m"
cd site && npm install

echo -e "\033[0;32mGenerating Hugo site for website...\033[0m"
hugo --minify ${USE_BASE_URL_SITE}

echo -e "\033[0;32mDone building site!\033[0m"
