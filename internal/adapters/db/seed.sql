-- Extensiones
CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- opcional

-- -------------------------
-- Dificultad / tamaño tablero
-- -------------------------
CREATE TABLE IF NOT EXISTS difficulty (
                            id                SERIAL PRIMARY KEY,
                            name              VARCHAR(64) NOT NULL UNIQUE,   -- easy/medium/hard
                            number_of_blocks  INT NOT NULL,
                            CHECK (number_of_blocks % 2 = 1)                 -- hueco central
    );

-- -------------------------
-- Sesión (una ejecución del juego)
-- -------------------------
CREATE TABLE IF NOT EXISTS sessions (
                          id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          player_id    UUID,                 -- opcional
                          device       TEXT,                 -- Quest, etc.
                          is_finished  BOOLEAN NOT NULL DEFAULT FALSE,
                          started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
                          ended_at     TIMESTAMPTZ
);

-- -------------------------
-- Partida (nivel) dentro de la sesión
-- -------------------------
CREATE TABLE IF NOT EXISTS matches (
                         id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         session_id     UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
                         difficulty_id  INT  NOT NULL REFERENCES difficulty(id),
                         level_n        INT  NOT NULL,                -- 3,4,5...
                         is_active      BOOLEAN NOT NULL DEFAULT TRUE,
                         started_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
                         ended_at       TIMESTAMPTZ,
                         outcome        VARCHAR(16),                  -- win/lose/aborted
                         meta           JSONB
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_matches_one_active_per_session
    ON matches (session_id)
    WHERE is_active;


CREATE INDEX idx_matches_session   ON matches(session_id);
CREATE INDEX idx_matches_active    ON matches(session_id, is_active);

-- -------------------------
-- Movimientos (append-only)
-- -------------------------
CREATE TABLE moves (
                       id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       match_id      UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
                       seq           INT  NOT NULL,                         -- 1..N
                       occurred_at   TIMESTAMPTZ NOT NULL DEFAULT now(),    -- tiempo absoluto
                       elapsed_ms    INT  NOT NULL,                         -- delta vs. movimiento anterior
                       from_idx      INT  NOT NULL,
                       to_idx        INT  NOT NULL,
                       move_kind     SMALLINT NOT NULL,                     -- 1=paso, 2=salto
                       frog_side     SMALLINT NOT NULL,                     -- 1=izq, 2=der
                       is_correct    BOOLEAN  NOT NULL DEFAULT TRUE,
                       interruption  BOOLEAN  NOT NULL DEFAULT FALSE,

    -- (opcional, útil para auditoría / replays)
                       board_before  JSONB,
                       board_after   JSONB,

    -- Métricas por movimiento (para promediar en KPIs)
                       branching_factor  INT,            -- nº opciones válidas vistas
                       buclicidad        DOUBLE PRECISION, -- 0..1 (ciclos/retrocesos)

                       UNIQUE (match_id, seq)
);

CREATE INDEX idx_moves_match_seq   ON moves(match_id, seq);
CREATE INDEX idx_moves_match_time  ON moves(match_id, occurred_at);

-- -------------------------
-- KPIs por partida (se rellenan solos por trigger)
-- -------------------------
CREATE TABLE match_stats (
                             match_id           UUID PRIMARY KEY REFERENCES matches(id) ON DELETE CASCADE,
                             total_moves        INT    NOT NULL DEFAULT 0,
                             errors             INT    NOT NULL DEFAULT 0,
                             avg_time_ms        INT    NOT NULL DEFAULT 0,
                             buclicidad_avg     DOUBLE PRECISION NOT NULL DEFAULT 0,
                             branch_factor_avg  DOUBLE PRECISION NOT NULL DEFAULT 0,
                             computed_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ========== Funciones y triggers ==========
-- Recalcular y upsert de KPIs por partida
CREATE OR REPLACE FUNCTION _recompute_match_stats(p_match UUID)
RETURNS VOID AS $$
BEGIN
INSERT INTO match_stats AS ms (match_id, total_moves, errors, avg_time_ms, buclicidad_avg, branch_factor_avg, computed_at)
SELECT
    p_match,
    COUNT(*)                                        AS total_moves,
    SUM(CASE WHEN is_correct = FALSE THEN 1 ELSE 0 END) AS errors,
    COALESCE(ROUND(AVG(elapsed_ms))::INT,0)         AS avg_time_ms,
    COALESCE(AVG(buclicidad), 0)                    AS buclicidad_avg,
    COALESCE(AVG(branching_factor), 0)              AS branch_factor_avg,
    now()
FROM moves
WHERE match_id = p_match
    ON CONFLICT (match_id) DO UPDATE
                                  SET total_moves       = EXCLUDED.total_moves,
                                  errors            = EXCLUDED.errors,
                                  avg_time_ms       = EXCLUDED.avg_time_ms,
                                  buclicidad_avg    = EXCLUDED.buclicidad_avg,
                                  branch_factor_avg = EXCLUDED.branch_factor_avg,
                                  computed_at       = EXCLUDED.computed_at;
END;
$$ LANGUAGE plpgsql;

-- Trigger: al insertar/actualizar/eliminar movimientos, recalcular KPIs
CREATE OR REPLACE FUNCTION _tg_moves_update_stats()
RETURNS TRIGGER AS $$
DECLARE
v_match UUID;
BEGIN
  IF (TG_OP = 'INSERT') THEN
    v_match := NEW.match_id;
  ELSIF (TG_OP = 'UPDATE') THEN
    v_match := NEW.match_id;
ELSE
    v_match := OLD.match_id;
END IF;

  PERFORM _recompute_match_stats(v_match);
RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tg_moves_update_stats_iud ON moves;
CREATE TRIGGER tg_moves_update_stats_iud
    AFTER INSERT OR UPDATE OR DELETE ON moves
    FOR EACH ROW EXECUTE FUNCTION _tg_moves_update_stats();

-- Cerrar partida: marca fin, outcome y asegura recálculo final
CREATE OR REPLACE FUNCTION close_match(p_match UUID, p_outcome VARCHAR DEFAULT 'win')
RETURNS VOID AS $$
BEGIN
UPDATE matches
SET is_active = FALSE,
    ended_at  = now(),
    outcome   = p_outcome
WHERE id = p_match;

PERFORM _recompute_match_stats(p_match);
END;
$$ LANGUAGE plpgsql;

-- Abrir nueva partida en una sesión (falla si ya hay activa)
CREATE OR REPLACE FUNCTION open_match(p_session UUID, p_difficulty INT, p_level INT)
RETURNS UUID AS $$
DECLARE
v_id UUID;
BEGIN
INSERT INTO matches (session_id, difficulty_id, level_n)
VALUES (p_session, p_difficulty, p_level)
    RETURNING id INTO v_id;
RETURN v_id;
END;
$$ LANGUAGE plpgsql;

INSERT INTO difficulty (name, number_of_blocks) VALUES ('easy', 7);
INSERT INTO difficulty (name, number_of_blocks) VALUES ('medium', 9);
INSERT INTO difficulty (name, number_of_blocks) VALUES ('hard', 11);
