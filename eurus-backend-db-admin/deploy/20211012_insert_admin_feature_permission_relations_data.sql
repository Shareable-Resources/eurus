-- Deploy eurus-backend-db-admin:20211012_insert_admin_feature_permission_relations_data to pg

BEGIN;

-- XXX Add DDLs here.
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (1, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (2, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (3, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (4, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (4, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (4, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (5, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (5, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (5, 3);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (6, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (7, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (8, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (9, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (10, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (11, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (12, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (13, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (14, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (15, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (16, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (17, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (18, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (19, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (20, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (20, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (20, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (20, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (21, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (21, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (21, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (21, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (22, 3);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (23, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (24, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (24, 3);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (25, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (25, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (25, 5);


INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (26, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (26, 5);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (27, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (27, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (27, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (27, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (28, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (28, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (28, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (28, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (29, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (29, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (29, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (29, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (30, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (34, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (35, 1);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (36, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (36, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (36, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (36, 4);

INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (37, 1);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (37, 2);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (37, 3);
INSERT INTO admin_feature_permission_relations (feature_id, permission_id) VALUES (37, 4);

COMMIT;
