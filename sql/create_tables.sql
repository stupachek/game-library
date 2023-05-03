CREATE TABLE IF NOT EXISTS users(
        id UUID PRIMARY KEY ,
    	email VARCHAR NOT NULL UNIQUE,
        userName VARCHAR NOT NULL,
		badgeColor VARCHAR,
		role VARCHAR NOT NULL,
        hashedPassword VARCHAR NOT NULL,
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO users(id, email, username, role, hashedPassword, badgeColor) values('699a5565-4c7e-4d18-be3e-ea04eb4f5e4d', 'admin@a.a', 'admin', 'admin', '$argon2id$v=19$m=65536,t=3,p=4$iZAGQSDhbte1+l0oF+rD/g$QWVeHaFYiR8iUA1BvK9+Pkua9EV3K/6y6CMTqSSet4Y', '' ) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS genres(
        id UUID PRIMARY KEY,
    	name VARCHAR NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS platforms(
        id UUID PRIMARY KEY,
    	name VARCHAR NOT NULL UNIQUE
    );

CREATE TABLE IF NOT EXISTS games(
        id UUID PRIMARY KEY,
		publisherId UUID,
    	title VARCHAR NOT NULL UNIQUE,
        description VARCHAR ,
		imageLink VARCHAR,
		ageRestriction INT,
		releaseYear INT,
		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS genresOnGames(
        id SERIAL PRIMARY KEY ,
		gameId UUID REFERENCES games(id),
		genreId UUID REFERENCES genres(id)
    );

CREATE TABLE IF NOT EXISTS platformsOnGames(
        id SERIAL PRIMARY KEY ,
		gameId UUID REFERENCES games(id),
		platformId UUID REFERENCES platforms(id)
    );