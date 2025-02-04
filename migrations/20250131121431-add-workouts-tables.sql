
-- +migrate Up

CREATE TABLE muscle_groups (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE exercise_muscle_groups (
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    muscle_group_id UUID REFERENCES muscle_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, muscle_group_id)
);

CREATE TABLE workouts (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY NOT NULL,
    workout_id UUID REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    sets INT NOT NULL,
    repetitions INT,
    weight FLOAT,
    duration INT,
    rest_time INT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_muscle_groups_name ON muscle_groups(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_exercises_name ON exercises(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_workouts_user_id ON workouts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workout_exercises_workout_id ON workout_exercises(workout_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_workout_exercises_exercise_id ON workout_exercises(exercise_id) WHERE deleted_at IS NULL;

-- +migrate Down

DROP INDEX IF EXISTS idx_muscle_groups_name;
DROP INDEX IF EXISTS idx_exercises_name;
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workout_exercises_workout_id;
DROP INDEX IF EXISTS idx_workout_exercises_exercise_id;

DROP TABLE IF EXISTS workout_exercises;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS exercise_muscle_groups;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS muscle_groups;