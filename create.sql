DROP TABLE IF EXISTS Mappings;

CREATE TABLE Mappings (
  slug VARCHAR(100) PRIMARY KEY,
  url VARCHAR(1000)
);

INSERT INTO Mappings
VALUES ('/cs144', 'https://suclass.stanford.edu/courses/course-v1:Engineering+CS144+Fall2017/about');
INSERT INTO Mappings
VALUES ('/cs107', 'http://web.stanford.edu/class/cs107/');