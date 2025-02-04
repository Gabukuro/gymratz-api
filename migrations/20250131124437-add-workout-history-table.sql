
-- +migrate Up

CREATE TABLE workout_history (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT
);

CREATE TABLE exercise_progress (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE,
    sets INT NOT NULL,
    repetitions INT,
    weight FLOAT,
    duration INT,
    notes TEXT
);

CREATE INDEX idx_workout_history_user_id ON workout_history(user_id);
CREATE INDEX idx_exercise_progress_user_id ON exercise_progress(user_id);
CREATE INDEX idx_exercise_progress_exercise_id ON exercise_progress(exercise_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_workout_history_user_id;
DROP INDEX IF EXISTS idx_exercise_progress_user_id;
DROP INDEX IF EXISTS idx_exercise_progress_exercise_id;
DROP TABLE IF EXISTS workout_history;
DROP TABLE IF EXISTS exercise_progress;
