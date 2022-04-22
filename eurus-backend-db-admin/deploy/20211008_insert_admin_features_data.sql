-- Deploy eurus-backend-db-admin:20211008_insert_admin_features_data to pg

BEGIN;

-- XXX Add DDLs here.
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (1000, '餘額管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (2000, '用戶管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (3000, '商戶管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (4000, '訂單管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (5000, '賬單管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (6000, '統計管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (7000, '管理員管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (8000, '錢包管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (9000, '審批管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (10000, '配置管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (11000, '日誌管理', null, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (12000, '資訊管理', null, true);


INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (1, '用戶餘額列表', 1000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (2, '錢包餘額列表', 1000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (3, '用戶地址餘額列表', 1000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (4, '用戶列表（包括聯絡資料）', 2000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (5, '商戶列表', 3000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (6, '充值記錄', 4000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (7, '提現記錄', 4000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (8, '支付記錄', 4000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (9, '轉賬記錄', 4000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (10, '交易記錄', 4000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (11, '回調記錄', 4000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (12, '日結單', 5000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (13, '月結單', 5000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (14, '用戶統計', 6000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (15, '商戶統計', 6000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (16, '訂單統計', 6000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (17, '用戶地址統計', 6000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (18, '公鏈統計', 6000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (19, '私鏈統計', 6000, true);


INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (20, '角色列表', 7000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (21, '管理員列表', 7000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (22, '修改密碼', 7000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (23, '歸集錢包管理', 8000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (24, '出幣錢包管理', 8000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (25, 'KYC 審批', 9000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (26, '訂單審批', 9000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (27, '資產配置', 10000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (28, '費用配置', 10000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (29, '行情配置', 10000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (30, '報表時間範圍配置', 10000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (34, '用戶登錄日誌', 11000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (35, '後台操作日誌', 11000, true);

INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (36, '公告管理', 12000, true);
INSERT INTO admin_features (id, name, parent_feature_id, is_enabled) VALUES (37, '信息通知管理', 12000, true);

COMMIT;
