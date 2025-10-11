CREATE TABLE users
(
    id INT SERIAL PRIMARY KEY,
    email VARCHAR(254) UNIQUE NOT NULL,
    password CHAR(16) NOT NULL,
    bday DATE NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE user_configs
(
    user_id INT REFERENCES users(id),
    sex CHAR(1) NOT NULL CHECK (sex IN ('M', 'W')),
    height SMALLINT NOT NULL CHECK (height > 100),
    weight SMALLINT NOT NULL CHECK (weight > 0),
    cal_goal SMALLINT NOT NULL CHECK (cal_goal > 0),
    activity SMALLINT NOT NULL CHECK (activity BETWEEN 1 AND 5),
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE meals
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    cal SMALLINT NOT NULL,
    protein NUMERIC(6,2) NOT NULL,
    carbs NUMERIC(6,2) NOT NULL,
    fats NUMERIC(6,2) NOT NULL,
    created_at TIMESTAMP NOT NULL
);