# Nitter Session Generator

This directory contains tools for generating Twitter/X OAuth sessions for use with Nitter.

## Files Overview

- `generate_session.py` - Python script for Twitter authentication with 2FA support
- `nitter.conf` - Nitter server configuration file
- `sessions.jsonl` - JSONL file storing OAuth tokens (SENSITIVE DATA)

## Purpose

This module handles Twitter authentication and session generation for Nitter, an alternative Twitter frontend. It uses Twitter's OAuth 1.0 flow to obtain access tokens.

## Key Components

### generate_session.py

**Authentication Flow:**
1. Obtains bearer token using Twitter consumer keys
2. Requests guest token
3. Initiates login flow with username
4. Handles password authentication
5. Supports 2FA (TOTP) and email verification
6. Extracts OAuth tokens from successful login
7. Appends tokens to sessions.jsonl

**Environment Variables Required:**
- `TWITTER_USERNAME` - Twitter account username
- `TWITTER_PASSWORD` - Twitter account password
- `TWITTER_OTP_SECRET` - TOTP secret for 2FA (if enabled)

**Usage:**
```bash
python generate_session.py <path-to-sessions-file>
```

**Dependencies:**
- requests
- pyotp
- cloudscraper (imported but not actively used)

### Security Considerations

1. **NEVER commit sessions.jsonl to version control** - contains active OAuth tokens
2. **Consumer keys are hardcoded** - these are public Twitter Android app keys
3. **Credentials from environment variables** - never hardcode passwords
4. **OAuth tokens grant account access** - treat as sensitive as passwords

## Code Guidelines

When working with this code:

1. **DO NOT modify authentication flow** without understanding Twitter's API requirements
2. **DO NOT log or expose OAuth tokens** in any output
3. **DO NOT commit files containing credentials**
4. **DO add error handling** for network failures
5. **DO validate environment variables** before starting auth flow
6. **DO respect rate limits** when generating multiple sessions

## Integration

The generated OAuth tokens in sessions.jsonl are used by Nitter to authenticate Twitter API requests, allowing it to fetch tweets and user data without requiring users to log in.
