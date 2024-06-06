CREATE DATABASE IF NOT EXISTS snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci;

USE snippetbox;

-- Not encoding inserts properly result in mojibaked test data
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET character_set_connection=utf8mb4;

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO 
    snippets (title, content, created, expires) 
VALUES 
    (
        'Żółw',
        'Żółw chciał pojechać koleją,\nLecz koleje nie tanieją.\nŻółwiowi szkoda pieniędzy:\n"Pójdę pieszo, będę prędzej."\n\n- Jan Brzechwa',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    ),
    (
        'Tygrys',
        '"Co słychać, panie tygrysie?"\n"A nic. Nudzi mi się."\n"Czy chciałby pan wyjść zza tych krat?"\n"Pewnie. Przynajmniej bym pana zjadł."\n\n- Jan Brzechwa',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
    ),
    (
        'Krokodyl',
        '"Skąd ty jesteś, krokodylu?"\n"Ja? Znad Nilu.\n"Wypuść mnie na kilka chwil,\nTo zawiozę cię nad Nil."\n\n- Jan Brzechwa',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 DAY)
    ),
    (
        'Papuga',
        '"Papużko, papużko,\nPowiedz mi coś na uszko."\n"Nic nie powiem, boś ty plotkarz,\nPowtórzysz każdemu, kogo spotkasz."\n\n- Jan Brzechwa',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

CREATE USER 'web'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'%';
ALTER USER 'web'@'%' IDENTIFIED BY 'pass';
FLUSH PRIVILEGES;
