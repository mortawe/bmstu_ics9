USE Lab10;
SELECT @@spid AS Session1Id;
GO
-- 1. dirty read:
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
-- SET TRANSACTION ISOLATION LEVEL READ COMMITED;
GO

BEGIN TRAN;

SELECT * FROM table1;

COMMIT TRAN;
GO

WAITFOR DELAY '00:00:03';

SELECT * FROM table1;
GO

--------------------
-- 2. non-repeatable read

SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
-- SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
GO

BEGIN TRAN;

SELECT * FROM table1;
GO

WAITFOR DELAY '00:00:03';

SELECT * FROM table1;
GO

COMMIT TRAN;


--------------------

-- 3. phantom read

SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
-- SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
GO

BEGIN TRAN;

SELECT * FROM table1;
GO

WAITFOR DELAY '00:00:03';

SELECT * FROM table1;
GO

COMMIT TRAN;


--------------------

-- 4. snapshot
ALTER DATABASE Lab10
    SET ALLOW_SNAPSHOT_ISOLATION ON;
GO

SET TRANSACTION ISOLATION LEVEL SNAPSHOT;
GO

BEGIN TRAN;

SELECT * FROM table1;
GO

WAITFOR DELAY '00:00:03';

SELECT * FROM table1;
GO

COMMIT TRAN;
