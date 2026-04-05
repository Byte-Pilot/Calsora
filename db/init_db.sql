CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    bday DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE user_profile
(
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    sex CHAR(1) NOT NULL CHECK (sex IN ('M', 'W')),
    height SMALLINT NOT NULL CHECK (height > 100),
    weight SMALLINT NOT NULL CHECK (weight > 0),
    activity SMALLINT NOT NULL CHECK (activity BETWEEN 1 AND 5),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE user_goal
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    type    VARCHAR(16) NOT NULL CHECK (type IN ('lose', 'maintain', 'gain')),
    target_weight SMALLINT NOT NULL CHECK (target_weight > 0),
    weekly_rate NUMERIC(3,1) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE nutrition_target
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    cal SMALLINT NOT NULL,
    protein NUMERIC(6,2) NOT NULL,
    carbs NUMERIC(6,2) NOT NULL,
    fats NUMERIC(6,2) NOT NULL,
    is_custom BOOL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE meals
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE meal_items (
    id SERIAL PRIMARY KEY,
    meal_id INT REFERENCES meals(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    grams SMALLINT  NOT NULL,
    cal SMALLINT NOT NULL,
    protein NUMERIC(6,2) NOT NULL,
    carbs NUMERIC(6,2) NOT NULL,
    fats NUMERIC(6,2) NOT NULL,
    confidence NUMERIC(3,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    plan VARCHAR(16) DEFAULT 'free',
    status VARCHAR(16) DEFAULT 'active',
    started_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_id_meals ON meals(user_id);
CREATE INDEX idx_meals_user_data ON meals(user_id, created_at);
CREATE INDEX idx_meal_id_meal_items ON meal_items(meal_id);
CREATE INDEX idx_user_id_profile ON user_profile(user_id);
CREATE INDEX idx_user_goal_latest ON user_goal(user_id, created_at DESC);
CREATE INDEX idx_nutrition_target_latest ON nutrition_target(user_id, created_at DESC);
CREATE INDEX idx_user_subscriptions ON subscriptions(user_id, status, expires_at DESC);
