#!/bin/bash

# Railway Environment Check Script
# Run this locally with: railway run ./scripts/check-railway-env.sh

echo "🔍 Checking Railway Environment Configuration..."
echo ""

# Check critical environment variables
check_var() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo "❌ $var_name: NOT SET"
        return 1
    else
        # Mask sensitive values
        if [[ $var_name == *"SECRET"* ]] || [[ $var_name == *"PASSWORD"* ]] || [[ $var_name == *"KEY"* ]]; then
            echo "✅ $var_name: ********** (set)"
        else
            # Show first 50 chars for URLs
            if [ ${#var_value} -gt 50 ]; then
                echo "✅ $var_name: ${var_value:0:50}... (${#var_value} chars)"
            else
                echo "✅ $var_name: $var_value"
            fi
        fi
        return 0
    fi
}

echo "📊 Critical Variables:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
check_var "DATABASE_URL"
check_var "JWT_SECRET"
check_var "ENV"
check_var "PORT"
echo ""

echo "📦 Optional Services:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
check_var "REDIS_URL"
check_var "AWS_ACCESS_KEY_ID"
check_var "OPENAI_API_KEY"
check_var "FCM_PROJECT_ID"
echo ""

# Check DATABASE_URL format
if [ ! -z "$DATABASE_URL" ]; then
    echo "🔗 Database URL Analysis:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Check if it contains localhost
    if [[ $DATABASE_URL == *"localhost"* ]] || [[ $DATABASE_URL == *"127.0.0.1"* ]]; then
        echo "⚠️  WARNING: DATABASE_URL contains localhost!"
        echo "   This will NOT work on Railway."
        echo "   Make sure you're using Railway's PostgreSQL service."
    else
        echo "✅ DATABASE_URL does not contain localhost"
    fi
    
    # Check SSL mode
    if [[ $DATABASE_URL == *"sslmode=require"* ]]; then
        echo "✅ SSL mode is set to 'require'"
    elif [[ $DATABASE_URL == *"sslmode=disable"* ]]; then
        echo "⚠️  WARNING: SSL mode is 'disable' - this may cause issues"
        echo "   Add ?sslmode=require to your DATABASE_URL"
    else
        echo "⚠️  WARNING: No SSL mode specified"
        echo "   Add ?sslmode=require to your DATABASE_URL"
    fi
    
    # Check if it's a Railway URL
    if [[ $DATABASE_URL == *"railway.app"* ]]; then
        echo "✅ Using Railway PostgreSQL service"
    fi
    echo ""
fi

# Check JWT_SECRET strength
if [ ! -z "$JWT_SECRET" ]; then
    echo "🔐 JWT Secret Analysis:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    jwt_length=${#JWT_SECRET}
    if [ $jwt_length -lt 32 ]; then
        echo "⚠️  WARNING: JWT_SECRET is only $jwt_length characters"
        echo "   Recommended minimum: 32 characters"
    else
        echo "✅ JWT_SECRET length: $jwt_length characters (good)"
    fi
    echo ""
fi

# Summary
echo "📋 Summary:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

errors=0
warnings=0

# Count critical missing vars
[ -z "$DATABASE_URL" ] && ((errors++))
[ -z "$JWT_SECRET" ] && ((errors++))
[ -z "$ENV" ] && ((warnings++))

# Check for localhost in DATABASE_URL
if [[ $DATABASE_URL == *"localhost"* ]] || [[ $DATABASE_URL == *"127.0.0.1"* ]]; then
    ((errors++))
fi

# Check SSL mode
if [[ ! $DATABASE_URL == *"sslmode=require"* ]]; then
    ((warnings++))
fi

if [ $errors -eq 0 ] && [ $warnings -eq 0 ]; then
    echo "✅ All checks passed! Your configuration looks good."
elif [ $errors -eq 0 ]; then
    echo "⚠️  Configuration is OK but has $warnings warning(s)"
    echo "   Review warnings above for optimization"
else
    echo "❌ Configuration has $errors critical error(s)"
    echo "   Fix errors above before deploying"
    exit 1
fi

echo ""
echo "💡 Tips:"
echo "   - Run this script with: railway run ./scripts/check-railway-env.sh"
echo "   - Make sure PostgreSQL service is added in Railway dashboard"
echo "   - DATABASE_URL should reference Railway's Postgres service"
echo "   - Add ?sslmode=require to DATABASE_URL for production"
