#!/bin/bash

# Generate secure JWT secrets for Kodia Framework
# This script generates two random 32-character secrets for JWT signing

set -e

echo "🔐 Generating secure JWT secrets for Kodia Framework..."
echo ""

# Generate access secret
ACCESS_SECRET=$(openssl rand -base64 32 | tr -d '\n')
echo "✅ Access Secret (32 chars): $ACCESS_SECRET"

# Generate refresh secret
REFRESH_SECRET=$(openssl rand -base64 32 | tr -d '\n')
echo "✅ Refresh Secret (32 chars): $REFRESH_SECRET"

echo ""
echo "📋 Add these to your .env file:"
echo ""
echo "APP_JWT_ACCESS_SECRET=$ACCESS_SECRET"
echo "APP_JWT_REFRESH_SECRET=$REFRESH_SECRET"
echo ""
echo "⚠️  Keep these secrets SECURE and NEVER commit them to git!"
echo "📝 Add .env to .gitignore if not already done"
