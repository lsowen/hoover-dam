INSERT INTO hoover_policy(name, policy, creation_date) VALUES
       ('FSFullAccess', '{"statement":[{"action":["fs:*"],"effect":"allow","resource":"*"}]}', NOW()),
       ('FSReadAll', '{"statement":[{"action":["fs:List*","fs:Read*"],"effect":"allow","resource":"*"}]}', NOW()),
       ('FSReadWriteAll', '{"statement":[{"action":["fs:Read*","fs:List*","fs:WriteObject","fs:DeleteObject","fs:RevertBranch","fs:CreateBranch","fs:CreateTag","fs:DeleteBranch","fs:DeleteTag","fs:CreateCommit","fs:CreateMetaRange"],"effect":"allow","resource":"*"}]}', NOW()),
       ('AuthFullAccess', '{"statement":[{"action":["auth:*"],"effect":"allow","resource":"*"}]}', NOW()),
       ('AuthManageOwnCredentials', '{"statement":[{"action":["auth:CreateCredentials","auth:DeleteCredentials","auth:ListCredentials","auth:ReadCredentials"],"effect":"allow","resource":"arn:lakefs:auth:::user/${user}"}]}', NOW()),
       ('RepoManagementFullAccess', '{"statement":[{"action":["ci:*"],"effect":"allow","resource":"*"},{"action":["retention:*"],"effect":"allow","resource":"*"}]}', NOW()),
       ('RepoManagementReadAll', '{"statement":[{"action":["ci:Read*"],"effect":"allow","resource":"*"},{"action":["retention:Get*"],"effect":"allow","resource":"*"}]}', NOW())
;

INSERT INTO hoover_group(name, creation_date) VALUES
  ('Admins', NOW()),
  ('SuperUsers', NOW()),
  ('Developers', NOW()),
  ('Viewers', NOW())
;

INSERT INTO hoover_group_policy (group_id, policy_id) VALUES 
  ((SELECT id FROM hoover_group WHERE name = 'Admins'), (SELECT id FROM hoover_policy WHERE name = 'FSFullAccess')),
  ((SELECT id FROM hoover_group WHERE name = 'Admins'), (SELECT id FROM hoover_policy WHERE name = 'AuthFullAccess')),
  ((SELECT id FROM hoover_group WHERE name = 'Admins'), (SELECT id FROM hoover_policy WHERE name = 'RepoManagementFullAccess')),

  ((SELECT id FROM hoover_group WHERE name = 'SuperUsers'), (SELECT id FROM hoover_policy WHERE name = 'FSFullAccess')),
  ((SELECT id FROM hoover_group WHERE name = 'SuperUsers'), (SELECT id FROM hoover_policy WHERE name = 'AuthManageOwnCredentials')),
  ((SELECT id FROM hoover_group WHERE name = 'SuperUsers'), (SELECT id FROM hoover_policy WHERE name = 'RepoManagementReadAll')),

  ((SELECT id FROM hoover_group WHERE name = 'Developers'), (SELECT id FROM hoover_policy WHERE name = 'AuthManageOwnCredentials')),
  ((SELECT id FROM hoover_group WHERE name = 'Developers'), (SELECT id FROM hoover_policy WHERE name = 'RepoManagementReadAll')),
  ((SELECT id FROM hoover_group WHERE name = 'Developers'), (SELECT id FROM hoover_policy WHERE name = 'FSReadWriteAll')),

  ((SELECT id FROM hoover_group WHERE name = 'Viewers'), (SELECT id FROM hoover_policy WHERE name = 'AuthManageOwnCredentials')),
  ((SELECT id FROM hoover_group WHERE name = 'Viewers'), (SELECT id FROM hoover_policy WHERE name = 'FSReadAll'))
;

---- create above / drop below ----

SELECT 1;
