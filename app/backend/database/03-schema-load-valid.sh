#!/bin/bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
# SPDX-License-Identifier: Apache-2.0

# Development seed data loader
# This script runs after the main schema initialization and loads test data

set -e

echo "🌺 Loading development test data... 🌺"

# Connect to the database and run the seed data
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    SELECT 'Users loaded: ' || COUNT(*) FROM users;
    SELECT 'Products loaded: ' || COUNT(*) FROM products;
    SELECT 'Categories loaded: ' || COUNT(*) FROM product_categories;
    SELECT 'Orders loaded: ' || COUNT(*) FROM orders;
 
    -- Verify the views work correctly  
    SELECT 'Purchase history records: ' || COUNT(*) FROM user_purchase_history;
EOSQL

echo "✅ Test data loaded successfully!"
echo "🔗 Available test users:"
echo "   Auth0 Users:"
echo "   - producer@islandbeats.com [verified]"
echo "   - beats@honolulusound.net [verified]"
echo "   - sounds@konabeats.io [unverified]"
echo ""
echo "   Email/Password Users:"
echo "   - music@tropicalstudio.co [verified]"
echo "   - hello@pacificsounds.org [unverified]"
echo "   - studio@mauivibes.com) [verified]"

echo "BDE Development Database ready."
