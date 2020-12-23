use Lab10;

SELECT
    resource_type,
    resource_subtype,
    resource_description,
    request_mode,
    request_type,
    request_status,
    request_owner_type,
    request_session_id
FROM sys.dm_tran_locks;
GO

WAITFOR DELAY '00:00:03';
GO

SELECT
    resource_type,
    resource_subtype,
    resource_description,
    request_mode,
    request_type,
    request_status,
    request_owner_type,
    request_session_id
FROM sys.dm_tran_locks;
GO
