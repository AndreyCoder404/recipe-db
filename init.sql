CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    protein FLOAT,
    fat FLOAT,
    carbs FLOAT,
    calories INT,
    instructions TEXT
);

INSERT INTO recipes (name, protein, fat, carbs, calories, instructions)
VALUES ('Test Recipe', 10, 5, 20, 200, 'Cook for 10 minutes');