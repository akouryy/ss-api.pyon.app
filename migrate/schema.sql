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
