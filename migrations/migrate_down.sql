-- SQL in this section is executed when migration is rolled back.

-- name: remove-user-accounts
DROP TABLE IF EXISTS user_accounts;

-- name: remove-mobile_users
DROP TABLE IF EXISTS mobile_users;

-- name: remove-mobile-user-tokens
DROP TABLE IF EXISTS mobile_user_tokens;

-- name: remove-mobile_user_contacts
DROP TABLE IF EXISTS mobile_user_contacts;

-- name: remove-mobile_user_profiles
DROP TABLE IF EXISTS mobile_user_profiles;

-- name: remove-mobile_user_alerts
DROP TABLE IF EXISTS mobile_user_alerts;