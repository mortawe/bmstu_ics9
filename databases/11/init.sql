USE master
GO

IF DB_ID(N'music_archive_db') IS NOT NULL
ALTER DATABASE music_archive_db
    SET SINGLE_USER
    WITH ROLLBACK IMMEDIATE;
DROP DATABASE music_archive_db;
GO

CREATE DATABASE music_archive_db
GO

USE music_archive_db;

---------------

IF OBJECT_ID(N'artists') IS NOT NULL
    DROP TABLE artists;
GO
CREATE TABLE artists
(
    artist_id int IDENTITY PRIMARY KEY,
    name      varchar(30) NOT NULL,
    biography text,
    born_date date        NOT NULL,
);
IF OBJECT_ID(N'trigger_artist_update') IS NOT NULL
    DROP TRIGGER trigger_artist_update;
GO
CREATE TRIGGER trigger_artist_update
    ON dbo.artists
    INSTEAD OF UPDATE
    AS
BEGIN
    IF update(artist_id)
        THROW 500403, 'forbidden', 1
    UPDATE artists
    SET name      = inserted.name,
        born_date = inserted.born_date,
        biography = inserted.biography
    FROM inserted
    WHERE artists.artist_id = inserted.artist_id
END;
GO
IF OBJECT_ID(N'releases_seq') IS NOT NULL
    DROP SEQUENCE releases_seq;
GO
CREATE SEQUENCE releases_seq
    START WITH 1
    INCREMENT BY 1;
GO
IF OBJECT_ID(N'releases') IS NOT NULL
    DROP TABLE releases;
GO
CREATE TABLE releases
(
    release_id int DEFAULT (NEXT VALUE FOR releases_seq) PRIMARY KEY,
    genre      varchar(30),
    date       date NOT NULL,
    is_album   bit  NOT NULL,
);
IF OBJECT_ID(N'trigger_releases_update') IS NOT NULL
    DROP TRIGGER trigger_releases_update;
GO
CREATE TRIGGER trigger_releases_update
    ON dbo.releases
    INSTEAD OF UPDATE
    AS
BEGIN
    IF update(release_id)
        THROW 500403, 'forbidden id modify', 1
    IF update(is_album)
        THROW 500403, 'forbidden is album modify', 1
    UPDATE releases
    SET genre = inserted.genre,
        date  = inserted.date
    FROM inserted
    WHERE releases.release_id = inserted.release_id
END;
GO


IF OBJECT_ID(N'artist_release') IS NOT NULL
    DROP TABLE artist_release;
GO
CREATE TABLE artist_release
(
    artist_id  int NOT NULL,
    release_id int NOT NULL,
    CONSTRAINT artist_release__artist_fk FOREIGN KEY (artist_id)
        REFERENCES artists (artist_id)
        ON DELETE CASCADE,
    CONSTRAINT artist_release__release_fk FOREIGN KEY (release_id)
        REFERENCES releases (release_id)
        ON DELETE CASCADE,
    CONSTRAINT artist_release_unique UNIQUE (artist_id, release_id)
)

IF OBJECT_ID(N'albums') IS NOT NULL
    DROP TABLE albums;
GO
CREATE TABLE albums
(
    release_id    int             NOT NULL UNIQUE,
    name          varchar(30),
    tracks_number int DEFAULT (0) NOT NULL,
    duration      int DEFAULT (0) NOT NULL,
    CONSTRAINT album_release_fk FOREIGN KEY (release_id)
        REFERENCES releases (release_id)
        ON DELETE CASCADE,
);


IF OBJECT_ID(N'singles') IS NOT NULL
    DROP TABLE singles;
GO
CREATE TABLE singles
(
    release_id int NOT NULL UNIQUE,
    CONSTRAINT single_release_fk FOREIGN KEY (release_id)
        REFERENCES releases (release_id)
        ON DELETE CASCADE,
);

IF OBJECT_ID(N'Tracks') IS NOT NULL
    DROP TABLE tracks;
GO
CREATE TABLE tracks
(
    track_id   int IDENTITY PRIMARY KEY,
    release_id int NOT NULL,
    title      varchar(50),
    length     int NOT NULL DEFAULT (0),
    CONSTRAINT track_release_fk FOREIGN KEY (release_id)
        REFERENCES releases (release_id)
        ON DELETE CASCADE
);
GO
IF OBJECT_ID(N'trigger_tracks_update') IS NOT NULL
    DROP TRIGGER trigger_tracks_update;
GO
CREATE TRIGGER trigger_tracks_update
    ON dbo.tracks
    INSTEAD OF UPDATE
    AS
BEGIN
    IF update(track_id)
        THROW 500403, 'forbidden id modify', 1
    UPDATE tracks
    SET title      = inserted.title,
        length     = inserted.length,
        release_id = inserted.track_id
    FROM inserted
    WHERE tracks.track_id = inserted.track_id
END;
GO

IF OBJECT_ID(N'trigger_tracks_insert') IS NOT NULL
    DROP TRIGGER trigger_tracks_insert;
GO
CREATE TRIGGER trigger_tracks_insert
    ON dbo.tracks
    AFTER INSERT
    AS
BEGIN
    UPDATE albums
    SET duration      = albums.duration + inserted.length,
        tracks_number = albums.tracks_number + (SELECT COUNT(*) FROM inserted)
    FROM inserted
    WHERE albums.release_id = inserted.release_id
END;
GO
-----------------------
IF OBJECT_ID(N'view_album_full_info') IS NOT NULL
    DROP VIEW view_album_full_info;
GO
CREATE VIEW view_album_full_info AS
SELECT a.name AS artist_name,
       a.born_date,
       genre,
       date,
       albums.name,
       tracks_number,
       duration
FROM releases
         JOIN albums ON releases.release_id = albums.release_id
         JOIN artist_release ar ON albums.release_id = ar.release_id
         JOIN artists a ON a.artist_id = ar.artist_id
WHERE releases.is_album = 1
GO
IF OBJECT_ID(N'trigger_view_album_full_info_insert') IS NOT NULL
    DROP TRIGGER trigger_view_album_full_info_insert;
GO
CREATE TRIGGER trigger_view_album_full_info_insert
    ON dbo.view_album_full_info
    INSTEAD OF INSERT
    AS
BEGIN
    DECLARE @inserted table
                      (
                          artist_id   int,
                          artist_name varchar(50),
                          born_date   date,
                          release_id  int DEFAULT (NEXT VALUE FOR releases_seq) NOT NULL,
                          genre       varchar(50),
                          date        date,
                          name        varchar(50)
                      );
    INSERT @inserted
    SELECT artist_id,
           artist_name,
           inserted.born_date,
           NEXT VALUE FOR releases_seq,
           genre,
           date,
           inserted.name
    FROM inserted
             JOIN artists
                  ON inserted.born_date = artists.born_date AND inserted.artist_name = artists.name;


    INSERT releases(release_id, genre, date, is_album)
    SELECT release_id, genre, date, 1
    FROM @inserted;

    INSERT artist_release(artist_id, release_id)
    SELECT artist_id, release_id
    FROM @inserted;

    INSERT albums(release_id, name)
    SELECT release_id, name
    FROM @inserted;
END;
GO

IF OBJECT_ID(N'view_single_track_full_info') IS NOT NULL
    DROP VIEW view_single_track_full_info;
GO
CREATE VIEW view_single_track_full_info AS
SELECT artists.name AS artist_name,
       artists.born_date,
       r2.genre,
       r2.date,
       t.title,
       t.length
FROM artists
         JOIN artist_release ar ON artists.artist_id = ar.artist_id
         JOIN releases r2 ON ar.release_id = r2.release_id
         JOIN singles s ON ar.release_id = s.release_id
         JOIN tracks t ON ar.release_id = t.release_id
WHERE r2.is_album = 0
GO
IF OBJECT_ID(N'tr_view_single_track_full_info') IS NOT NULL
    DROP TRIGGER tr_view_single_track_full_info;
GO
CREATE TRIGGER tr_view_single_track_full_info
    ON dbo.view_single_track_full_info
    INSTEAD OF INSERT
    AS
BEGIN
    DECLARE @inserted table
                      (
                          release_id   int NOT NULL,
                          artist_name  varchar(50),
                          track_title  varchar(50),
                          genre        varchar(50),
                          date         date,
                          track_length int,
                          artist_id    int NOT NULL
                      );

    INSERT @inserted(release_id, artist_name, track_title, genre, date, track_length, artist_id)
    SELECT NEXT VALUE FOR releases_seq,
           inserted.artist_name,
           inserted.title,
           inserted.genre,
           inserted.date,
           inserted.length,
           artists.artist_id
    FROM inserted
             JOIN artists ON artist_name = name AND artists.born_date = inserted.born_date;

    SELECT * FROM @inserted;
    INSERT INTO releases(release_id, genre, date, is_album)
    SELECT release_id, genre, date, 0
    FROM @inserted;

    INSERT INTO singles(release_id)
    SELECT release_id
    FROM @inserted;

    INSERT INTO tracks(release_id, title, length)
    SELECT release_id, track_title, track_length
    FROM @inserted;

    INSERT INTO artist_release(artist_id, release_id)
    SELECT artist_id, release_id
    FROM @inserted;
END;
GO

IF OBJECT_ID(N'view_track_album_info') IS NOT NULL
    DROP VIEW view_track_album_info;
GO
CREATE VIEW view_track_album_info AS
SELECT a.name      AS artist_name,
       born_date,
       genre,
       date,
       albums.name AS album_title,
       title       AS track_title,
       length
FROM releases
         JOIN albums ON releases.release_id = albums.release_id
         JOIN artist_release ar ON albums.release_id = ar.release_id
         JOIN artists a ON a.artist_id = ar.artist_id
         JOIN tracks t ON ar.release_id = t.release_id
WHERE releases.is_album = 1
GO

IF OBJECT_ID(N'trigger_view_track_album_info_insert') IS NOT NULL
    DROP TRIGGER trigger_view_track_album_info_insert;
GO
CREATE TRIGGER trigger_view_track_album_info_insert
    ON dbo.view_track_album_info
    INSTEAD OF INSERT
    AS
BEGIN
    DECLARE @inserted table
                      (
                          artist_id   int,
                          artist_name varchar(50),
                          born_date   date,
                          release_id  int NOT NULL,
                          genre       varchar(50),
                          date        date,
                          album_title varchar(50),
                          length      int,
                          track_title varchar(50)
                      );
    INSERT @inserted
    SELECT ar.artist_id,
           artist_name,
           inserted.born_date,
           albums.release_id,
           r2.genre,
           r2.date,
           inserted.album_title,
           inserted.length,
           inserted.track_title
    FROM inserted
             JOIN albums ON albums.name = inserted.album_title
             JOIN releases r2 ON albums.release_id = r2.release_id
             JOIN artist_release ar ON albums.release_id = ar.release_id
             JOIN artists a ON ar.artist_id = a.artist_id;

    INSERT tracks(release_id, title, length)
    SELECT release_id, track_title, length
    FROM @inserted
END;
GO

-----------------------

INSERT INTO artists (name, biography, born_date)
VALUES ('Grimes', 'Canadian musician, singer, songwriter, record producer, music video director, and visual artist.',
        '04.17.1988'),
       ('Mars Argo', 'American singer, songwriter, actress, photographer, and Internet personality.', '01.01.2009'),
       ('Phoebe Bridgers', 'American indie musician from Los Angeles, California.', '08.17.1994'),
       ('Thom Yorke', 'English musician and the main vocalist and songwriter of the rock band Radiohead.',
        '10.07.1968'),
       ('Björk', 'Icelandic singer, songwriter, record producer, actress, and DJ.', '11.21.1965'),
       ('Ashnikko', 'American singer, songwriter, and rapper.', '02.19.1996'),
       ('Toru Kitajima', 'Japanese musician and singer-songwriter.', '12.23.1982')


INSERT INTO view_album_full_info(artist_name, born_date, name, genre, date)
VALUES ('Grimes', '04.17.1988', 'Visions', 'pop, experimental', '01.31.2012'),
       ('Grimes', '04.17.1988', 'Art Angels', 'pop, dance', '11.06.2015'),
       ('Grimes', '04.17.1988', 'Miss Anthropocene', 'electronic', '02.21.2020'),
       ('Mars Argo', '01.01.2009', 'Technology Is a Dead Bird', 'rock', '11.06.2009'),
       ('Toru Kitajima', '12.23.1982', 'flowering', 'rock', '06.27.2012'),
       ('Toru Kitajima', '12.23.1982', 'Fantastic Magic', 'rock', '08.27.2014')


INSERT INTO view_single_track_full_info(artist_name, born_date, genre, date, title, length)
VALUES ('Grimes', '04.17.1988', 'rock, pop, experimental', '07.01.2012', 'Go', 3 * 60),
       ('Grimes', '04.17.1988', 'rock, pop', '12.10.2015', 'Flesh without Blood', 4 * 60),
       ('Toru Kitajima', '12.23.1982', 'rock, electronic', '01.01.2014', 'Unravel', 4 * 60),
       ('Toru Kitajima', '12.23.1982', 'rock, electronic', '01.08.2019', 'melt', 5 * 60)

INSERT INTO view_track_album_info(artist_name, born_date, album_title, track_title, length)
VALUES ('Grimes', '04.17.1988', 'Visions', 'track1', 10),
       ('Grimes', '04.17.1988', 'Visions', 'track2', 30),
       ('Grimes', '04.17.1988', 'Visions', 'track3', 40)

----------------------------
-- album nums order by artist's name
SELECT artist_name, COUNT(*) AS album_num
FROM view_album_full_info
GROUP BY artist_name
ORDER BY artist_name

-- rock singles
SELECT *
FROM view_single_track_full_info
WHERE genre LIKE '%rock%'

-- list of known artists with any release
SELECT DISTINCT name, born_date
FROM artists AS artitsts_list
         INNER JOIN artist_release ar ON artitsts_list.artist_id = ar.artist_id

DELETE releases
WHERE release_id IN (SELECT release_id FROM albums WHERE name = 'flowering')

------------------------------
IF OBJECT_ID(N'delete_artist_recursively') IS NOT NULL
    DROP PROCEDURE delete_artist_recursively;
GO
CREATE PROCEDURE delete_artist_recursively(@artist varchar(50))
AS
BEGIN
    DECLARE @to_delete_id table
                          (
                              release_id int
                          );
    INSERT @to_delete_id
    SELECT release_id
    FROM artist_release
             JOIN artists a ON a.artist_id = artist_release.artist_id
    WHERE a.name = @artist;

    DELETE albums WHERE release_id IN (SELECT release_id FROM @to_delete_id);
    DELETE singles WHERE release_id IN (SELECT release_id FROM @to_delete_id);
    DELETE releases WHERE release_id IN (SELECT release_id FROM @to_delete_id);
    DELETE artists WHERE name = @artist;
END;
GO

delete_artist_recursively 'Grimes'

UPDATE artists
SET name = N'Björk Guðmundsdóttir'
WHERE name = N'Björk'

SELECT artist_name, name, 1 AS is_album
FROM view_album_full_info
UNION ALL
SELECT artist_name, title AS name, 0 AS is_album
FROM view_single_track_full_info

IF OBJECT_ID(N'idx_artists_name_date') IS NOT NULL
    DROP INDEX idx_artists_name_date
        ON artists;
GO
CREATE UNIQUE INDEX idx_artists_name_date ON artists (name, born_date);

IF OBJECT_ID(N'idx_album_name_date') IS NOT NULL
    DROP INDEX idx_album_name_date
        ON albums;
GO
CREATE INDEX idx_album_name_date ON albums (name);


-- 4

SELECT DISTINCT name, born_date
FROM artists AS artitsts_list
         INNER JOIN artist_release ar ON artitsts_list.artist_id = ar.artist_id

SELECT * FROM artists INNER JOIN artist_release ar on artists.artist_id = ar.artist_id INNER JOIN albums a on ar.release_id = a.release_id
SELECT * FROM artists LEFT JOIN artist_release ar on artists.artist_id = ar.artist_id LEFT JOIN albums a on ar.release_id = a.release_id
SELECT * FROM artists RIGHT JOIN artist_release ar on artists.artist_id = ar.artist_id RIGHT JOIN albums a on ar.release_id = a.release_id
SELECT * FROM artists FULL OUTER JOIN artist_release ar on artists.artist_id = ar.artist_id FULL OUTER JOIN albums a on ar.release_id = a.release_id

SELECT * FROM view_album_full_info WHERE date BETWEEN '01.01.2000' AND '01.01.2010'
SELECT *
FROM view_single_track_full_info
WHERE genre LIKE '%rock%'

SELECT * FROM artists FULL OUTER JOIN artist_release ar on artists.artist_id = ar.artist_id FULL OUTER JOIN albums a on ar.release_id = a.release_id
WHERE a.release_id IS NOT NULL

SELECT * FROM artists WHERE EXISTS(SELECT * FROM artist_release WHERE artists.artist_id = artist_release.artist_id)

SELECT DATEDIFF(year,  date, CURRENT_TIMESTAMP) as year_passed, COUNT(*) AS album_num
FROM view_album_full_info
GROUP BY DATEDIFF(year,  date, CURRENT_TIMESTAMP)
HAVING COUNT(*) < 10

SELECT AVG(length) as avg, SUM(length) as sum, min(length) as min, max(length) as max FROM tracks
SELECT artist_name FROM view_single_track_full_info UNION ALL SELECT artist_name FROM view_album_full_info;
SELECT artist_name FROM view_single_track_full_info UNION ALL SELECT artist_name FROM view_album_full_info EXCEPT
SELECT vafi.artist_name FROM view_single_track_full_info JOIN view_album_full_info vafi ON view_single_track_full_info.artist_name = vafi.artist_name

SELECT artist_name FROM view_single_track_full_info UNION SELECT artist_name FROM view_album_full_info INTERSECT
SELECT vafi.artist_name FROM view_single_track_full_info INNER JOIN view_album_full_info vafi ON view_single_track_full_info.artist_name = vafi.artist_name
