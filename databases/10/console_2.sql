USE Lab10;
SELECT @@spid AS Session1Id;
GO
-- SELECT @@spid AS Session2Id;
-- GO

--------------------

-- 1. dirty read

BEGIN TRANSACTION;

UPDATE [table1]
SET col3 = N'uncommited'
WHERE col2 = N'example';

WAITFOR DELAY '00:00:03';

ROLLBACK;


--------------------

-- 2. non-repeatable read

BEGIN TRAN;

DELETE [table1]
WHERE col2 = N'example';

COMMIT TRAN;


--------------------

-- 3. phantom read

BEGIN TRAN;

INSERT INTO [table1](col2, col3)
VALUES
(N'phantom1', N'uncommited1'),
(N'phantom2', N'uncommited2');

COMMIT TRAN;