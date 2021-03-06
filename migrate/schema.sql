ALTER DATABASE ________ CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;
--> SET NAMES utf8mb4 COLLATE utf8mb4_bin;
--> SET collation_server = utf8mb4_bin;
--> SHOW VARIABLES LIKE '%character%'; */
--> SHOW VARIABLES LIKE '%collation%'; */

CREATE TABLE books (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(300) NOT NULL,
  created_at DATETIME NOT NULL
)
  /*CHARACTER SET utf8mb4
  COLLATE utf8mb4_bin*/;

CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  nickname VARCHAR(30) NOT NULL UNIQUE,
  password VARCHAR(70) NOT NULL
);

CREATE TABLE authors (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  url VARCHAR(300) NOT NULL,
  created_at DATETIME NOT NULL
);

CREATE TABLE book_authors (
  book_id INT NOT NULL,
  author_id INT NOT NULL,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE,
  UNIQUE (book_id, author_id)
);

CREATE TABLE episodes (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(300) NOT NULL,
  book_id INT NOT NULL,
  indexnum INT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
  UNIQUE (book_id, title),
  UNIQUE (book_id, indexnum)
);

CREATE TABLE sections (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  content MEDIUMTEXT NOT NULL,
  episode_id INT NOT NULL,
  indexnum INT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (episode_id) REFERENCES episodes(id) ON DELETE CASCADE,
  UNIQUE (episode_id, indexnum)
) ROW_FORMAT=DYNAMIC;
