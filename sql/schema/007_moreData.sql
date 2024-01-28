-- ++ goose Up
ALTER TABLE visits
ADD COLUMN browser TEXT NOT NULL default 'unknown',
ADD COLUMN device TEXT NOT NULL default 'unknown',
ADD COLUMN os TEXT NOT NULL default 'unknown';

-- goose down
ALTER TABLE visits 
DROP COLUMN browser,
DROP COLUMN device,
DROP COLUMN os;
