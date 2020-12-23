USE master
GO

IF DB_ID(N'Lab10') IS NOT NULL
  DROP DATABASE Lab10;

CREATE DATABASE Lab10;
GO

USE Lab10;
GO

IF OBJECT_ID(N'table1') IS NOT NULL
  DROP TABLE [table1]
GO


CREATE TABLE [table1](
  col1 int PRIMARY KEY IDENTITY(1, 1),
  col2 nvarchar(120) NOT NULL,
  col3 nvarchar(120)  NOT NULL,

)
GO

INSERT INTO [table1](col2, col3)
    VALUES
  (N'example', N'commited')
GO


USE Lab10;
GO

DELETE
FROM table1;

DROP TABLE table1;


