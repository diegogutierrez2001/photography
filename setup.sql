DROP DATABASE IF EXISTS photography;
CREATE DATABASE photography;
use photography

DROP TABLE IF EXISTS photographs;
CREATE TABLE photographs (
	id 						INT NOT NULL AUTO_INCREMENT,
	name					VARCHAR(64) NOT NULL,
  title      		VARCHAR(128),
  description		VARCHAR(1024),
  PRIMARY KEY (`id`)
);

INSERT INTO photographs
  (name)
VALUES
	('bdayqueen.jpeg'),
  ('cozy.jpeg'),
  ('endofnight.jpeg'),
  ('group.jpeg'),
  ('inschock.jpeg'),
  ('ladiesmen.jpeg'),
  ('nips.jpeg'),
  ('notsponsored.jpeg'),
  ('point.jpeg'),
  ('renaissance.jpeg'),
  ('snug.jpeg'),
  ('star.jpeg'),
  ('studs.jpeg'),
  ('weldlove.jpeg'),
  ('wink.jpeg')
;

DROP TABLE IF EXISTS collections;
CREATE TABLE collections (
	id 						INT NOT NULL AUTO_INCREMENT,
  title      		VARCHAR(256) NOT NULL,
  description		VARCHAR(1024) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO collections
	(title, description)
VALUES
('Jenn\'s 22nd Birthday Party', 'We celebrated by going to a drag show in Boston'),
('Diego Gutierrez', 'Just pictures of me.');

DROP TABLE IF EXISTS entries;
CREATE TABLE entries (
	collectionID INT NOT NULL,
	photoID INT NOT NULL,
	FOREIGN KEY (collectionID) REFERENCES collections(id),
	FOREIGN KEY (photoID) REFERENCES photographs(id)
);

INSERT INTO entries
	(collectionID, photoID)
VALUES
	(1, 1),
	(1, 2),
	(1, 3),
	(1, 4),
	(1, 5),
	(1, 6),
	(1, 7),
	(1, 8),
	(1, 9),
	(1, 10),
	(1, 11),
	(1, 12),
	(1, 13),
	(1, 14),
	(1, 15),
	(2, 6),
	(2, 13),
	(2, 15)
;