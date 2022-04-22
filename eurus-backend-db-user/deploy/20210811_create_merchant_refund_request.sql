-- Deploy eurus-backend-db-user:20210811_create_merchant_refund_request to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS merchant_refund_requests(
	id BIGSERIAL PRIMARY KEY,
	dest_address VARCHAR(255) NOT NULL,
	asset_name VARCHAR(20) NOT NULL,
	amount NUMERIC(75) NOT NULL,
	user_id BIGINT,
	purchase_trans_hash VARCHAR(255),
	refund_reason TEXT,
	operator_comment TEXT,
	status SMALLINT NOT NULL,
	merchant_id BIGINT NOT NULL,
	merchant_operator_id BIGINT,
	refund_trans_hash VARCHAR(255),
	created_date TIMESTAMPTZ NOT NULL,
	last_modified_date TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS merchant_refund_requests_merchant_id_idx ON merchant_refund_requests(merchant_id);
CREATE INDEX IF NOT EXISTS merchant_refund_requests_merchant_id_status ON merchant_refund_requests(merchant_id, status);
CREATE INDEX IF NOT EXISTS merchant_refund_requests_merchant_id_trans ON merchant_refund_requests(merchant_id, purchase_trans_hash);

COMMIT;
