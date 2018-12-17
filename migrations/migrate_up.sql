-- SQL in this section is executed when migration is applied.

-- name: create-user-accounts
CREATE TABLE IF NOT EXISTS user_accounts
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    username        VARCHAR(255)   NOT NULL,
    password        VARCHAR(255)   NOT NULL,
    is_admin        BOOLEAN        DEFAULT TRUE,
    created_at      DATETIME       DEFAULT NOW(),
    last_login_at   DATETIME       NULL,
    updated_at      DATETIME       NULL,      
    PRIMARY KEY(id)
);

-- name: create-user-accounts-username-index
CREATE UNIQUE INDEX user_accounts_username_index ON user_accounts(username);


-- name: create-mobile-users
CREATE TABLE IF NOT EXISTS mobile_users
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    msisdn          VARCHAR(255)   NOT NULL,
    created_at      DATETIME       DEFAULT NOW(),
    last_login_at   DATETIME       NULL,     
    PRIMARY KEY(id)
);

-- name: create-mobile-users-msisdn-index
CREATE UNIQUE INDEX mobile_users_msisdn_index ON mobile_users(msisdn);


-- name: create-mobile-user-tokens
CREATE TABLE IF NOT EXISTS mobile_user_tokens
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    user_id         INT            NOT NULL,
    platform        VARCHAR(255)   NOT NULL,
    token           VARCHAR(255)   NOT NULL,
    created_at      DATETIME       DEFAULT NOW(),
    updated_at      DATETIME       NULL,     
    PRIMARY KEY(id),
    CONSTRAINT fk_mobile_user_tokens_user_id  FOREIGN KEY  (user_id) REFERENCES mobile_users(id)
);

-- name: create-mobile-users-token-index
CREATE UNIQUE INDEX mobile_users_token_index ON mobile_user_tokens(token);

-- name: create-contacts
CREATE TABLE IF NOT EXISTS mobile_user_contacts
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    msisdn          VARCHAR(255)   NOT NULL,
    fullname        VARCHAR(255)   NOT NULL,
    user_id         INT            NOT NULL,
    created_at      DATETIME       DEFAULT NOW(),   
    PRIMARY KEY(id),
    CONSTRAINT fk_mobile_user_contacts_user_id  FOREIGN KEY  (user_id) REFERENCES mobile_users(id)
);


-- name: create-mobile-user-profiles
CREATE TABLE IF NOT EXISTS mobile_user_profiles
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    user_id         INT            NOT NULL,
    title           VARCHAR(255)   NULL,
    fullname        VARCHAR(255)   NULL,
    street          VARCHAR(255)   NULL,
    city            VARCHAR(255)   NULL,
    post_code       VARCHAR(255)   NULL,
    geo_lng         VARCHAR(255)   NULL,
    geo_lat         VARCHAR(255)   NULL,
    created_at      DATETIME       DEFAULT NOW(),
    updated_at      DATETIME       NULL,     
    PRIMARY KEY(id),
    CONSTRAINT fk_mobile_user_profiles_user_id  FOREIGN KEY  (user_id) REFERENCES mobile_users(id)
);

-- name: create-mobile-user-alerts
CREATE TABLE IF NOT EXISTS mobile_user_alerts
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    user_id         INT            NOT NULL,
    geo_lng         VARCHAR(255)   NULL,
    geo_lat         VARCHAR(255)   NULL,
    created_at      DATETIME       DEFAULT NOW(),   
    PRIMARY KEY(id),
    CONSTRAINT fk_mobile_user_alerts_user_id  FOREIGN KEY  (user_id) REFERENCES mobile_users(id)
);

